package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"takehome/pkg/db"
	"time"

	"math/rand"
)

type Status string

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	Running  Status = "running"
	Finished Status = "finished"
	Error    Status = "error"
	New      Status = "new"
	Killed   Status = "killed"
)

var JobStatuses = map[string]Status{
	"running":  Running,
	"finished": Finished,
	"error":    Error,
	"new":      New,
	"killed":   Killed,
}

type Job struct {
	ID      string
	PID     int
	Logfile *os.File
	Command string
	Args    []string
	Status  Status
	Jobhub  *JobHub
}

// NewJob creates a new job with a unique id, also initializes log file
func NewJob(Jobhub *JobHub, command string, args ...string) (*Job, error) {
	id := generateID()
	file, err := createLogFile(id)

	if err != nil {
		log.Println("Error creating log file:", err)
		return nil, err
	}

	return &Job{
		ID:      id,
		Command: command,
		Args:    args,
		Status:  New,
		Logfile: file,
		Jobhub:  Jobhub,
	}, nil
}

func generateID() string {
	allCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, 8)
	for i := range id {
		id[i] = allCharacters[rand.Intn(len(allCharacters))]
	}
	return string(id)
}

// Run runs a job and sets the status accordingly
func (j *Job) Run() error {
	cmd := exec.Command(j.Command, j.Args...)
	cmd.Stdout = j.Logfile
	cmd.Stderr = j.Logfile
	defer j.Logfile.Close()
	err := cmd.Start()
	if err != nil {
		log.Println("Error starting command:", err)
		return err
	}
	j.PID = cmd.Process.Pid
	if err := cmd.Wait(); err != nil {
		log.Println("Error waiting for command:", err)
		if !strings.Contains(err.Error(), "signal: killed") {
			j.Status = Error
		}
		return err
	}
	j.Status = Finished
	return nil
}

// Kill kills a job
func (j *Job) Kill() error {
	if j.Status != Running {
		return fmt.Errorf("job is not running")
	}
	process, err := os.FindProcess(j.PID)
	if err != nil {
		log.Println("Error finding process:", err)
		return err
	}
	err = process.Kill()
	if err != nil {
		log.Println("Error killing process:", err)
		return err
	}
	j.Jobhub.RemoveJob(j.ID)
	j.Status = Killed
	return nil
}

// CreateLogFile creates a log file for a job
func createLogFile(id string) (*os.File, error) {
	log.Println("creating log file for job", id)
	return os.Create(getLogFileName(id))
}

// ToModel converts a job to a db.Job
func (j *Job) toModel() *db.Job {
	return &db.Job{
		ID:      j.ID,
		Command: j.Command,
		Args:    strings.Join(j.Args, ","),
		Status:  string(j.Status),
		Logfile: getLogFileName(j.ID),
	}
}

// FromModel converts a db.Job to a job
func (j *Job) fromModel(e *db.Job, hub *JobHub) error {
	file, err := getModelLogFile(e.ID)
	if err != nil {
		return err
	}
	j.ID = e.ID
	j.Command = e.Command
	j.Args = strings.Split(e.Args, ",")
	j.Status = JobStatuses[e.Status]
	j.Logfile = file
	j.Jobhub = hub
	return nil
}

// GetLogFileName returns the name of the log file for a job
func getLogFileName(id string) string {
	return fmt.Sprintf("%s.log", id)
}

// GetModelLogFile returns the log file for a job
func getModelLogFile(id string) (*os.File, error) {
	file, err := os.Open(getLogFileName(id))
	if err != nil {
		log.Println("Error opening log file:", err)
		f, err := createLogFile(id)
		if err != nil {
			log.Println("Error creating log file:", err)
			return nil, err
		}
		return f, nil
	}
	return file, nil
}
