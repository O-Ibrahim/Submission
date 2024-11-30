package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"takehome/pkg/db"
)

type runBodyT struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type handleRunResponse struct {
	ID string `json:"id"`
}

// handleCreateJob creates a new job and runs it
func (a *App) handleCreateJob(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	runBody := new(runBodyT)
	err = json.Unmarshal(body, runBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	job, err := NewJob(a.JobHub, runBody.Command, runBody.Args...)
	if err != nil {
		log.Println("error creating job:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := a.Store.CreateJob(job.toModel()); err != nil {
		log.Println("error creating job:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	status := job.Status
	a.SendJobToHook(job.ID, string(status))
	a.JobHub.AddJob(job)
	go func() {
		job.Status = Running
		err := a.Store.UpdateJob(job.toModel())
		if err != nil {
			log.Println("error updating job:", err)
			return
		}
		a.JobHub.UpdateJob(job)
		err = job.Run()
		if err != nil {
			log.Println("error running job:", err)
			a.JobHub.RemoveJob(job.ID)
		}
		err = a.Store.UpdateJob(job.toModel())
		if err != nil {
			log.Println("error updating job:", err)
			return
		}
		a.SendJobToHook(job.ID, string(job.Status))
	}()

	writeJson(w, http.StatusCreated, handleRunResponse{
		ID: job.ID,
	})

}

// jobStatusResponse is the response for the status endpoint
type jobStatusResponse struct {
	ID     string `json:"id"`
	Status Status `json:"status"`
}

// handleGetStatus returns the status of a job
func (a *App) handleGetStatus(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		log.Println("no id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jobFromDB, err := a.Store.GetJobByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("job not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			log.Println("error getting job:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	writeJson(w, http.StatusOK, jobStatusResponse{
		ID:     jobFromDB.ID,
		Status: Status(jobFromDB.Status),
	})
}

// handleKillJob kills a job based on ID in the path value
func (a *App) handleKillJob(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		log.Println("no id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	job := a.JobHub.GetJob(id)
	err := job.Kill()
	if err != nil {
		log.Println("error killing job:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := a.Store.UpdateJob(job.toModel()); err != nil {
		log.Println("error updating job:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	go a.SendJobToHook(job.ID, string(job.Status))
	writeJson(w, http.StatusOK, jobStatusResponse{
		ID:     job.ID,
		Status: job.Status,
	})
}

// handleGetLogs returns the logs of a job based on ID in the path value, lines query string is optional and it provides the latest n lines of the log
func (a *App) handleGetLogs(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	latestLines := 0
	if req.URL.Query().Get("lines") != "" {
		lines, err := strconv.Atoi(req.URL.Query().Get("lines"))
		if err != nil {
			log.Println("invalid lines parameter")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if lines < 0 {
			log.Println("invalid lines parameter")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		latestLines = lines
	}

	if id == "" {
		log.Println("no id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if latestLines == 0 {
		file, err := os.Open(fmt.Sprintf("%s.log", id))
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()
		io.Copy(w, file)
		return
	}

	cmd := exec.Command("tail", "-n", strconv.Itoa(latestLines), fmt.Sprintf("%s.log", id))

	output, err := cmd.Output()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(output)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// getAllJobsResponse is the response for the get all jobs endpoint
type getAllJobsResponse struct {
	Jobs []*db.Job `json:"jobs"`
}

// handleGetAllJobs returns all jobs
func (a *App) handleGetAllJobs(w http.ResponseWriter, req *http.Request) {
	//TODO add pagination
	jobs, err := a.Store.GetJobs()
	if err != nil && err != sql.ErrNoRows {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJson(w, http.StatusOK, getAllJobsResponse{
		Jobs: jobs,
	})

}

func (a *App) handleGetJobByID(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		log.Println("no id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	job, err := a.Store.GetJobByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("job not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			log.Println("error getting job:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	writeJson(w, http.StatusOK, job)
}

// SendJobToHook sends the job and status to the hook url, request of type GET and data in query strings
func (a *App) SendJobToHook(jobID string, status string) {
	if a.Config.HookUrl == "" {
		return
	}
	url := fmt.Sprintf("%s?jobId=%s&status=%s", a.Config.HookUrl, jobID, status)
	log.Println("sending hook request to: ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("error creating request:", err)
		return
	}
	hc := http.Client{}
	_, err = hc.Do(req)
	if err != nil {
		log.Println("error sending request:", err)
		return
	}
}
