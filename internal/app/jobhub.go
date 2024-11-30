package app

import (
	"context"
	"log"
	"sync"
)

// JobStatus represents a hub that holds all active jobs
type JobHub struct {
	jobsLock sync.RWMutex
	Jobs     map[string]*Job
}

func NewJobHub() *JobHub {
	return &JobHub{
		Jobs: make(map[string]*Job),
	}
}

func (h *JobHub) AddJob(j *Job) {
	h.jobsLock.Lock()
	defer h.jobsLock.Unlock()
	h.Jobs[j.ID] = j
}

func (h *JobHub) GetJob(id string) *Job {
	h.jobsLock.RLock()
	defer h.jobsLock.RUnlock()
	return h.Jobs[id]
}

func (h *JobHub) RemoveJob(id string) {
	h.jobsLock.Lock()
	defer h.jobsLock.Unlock()
	delete(h.Jobs, id)
}

func (h *JobHub) UpdateJob(j *Job) {
	h.jobsLock.Lock()
	defer h.jobsLock.Unlock()
	h.Jobs[j.ID] = j
}
func (h *JobHub) Shutdown(ctx context.Context) {
	for _, j := range h.Jobs {
		if j.Status == Running {
			//kill the job
			log.Println("killing job", j.ID)
			//kill the job
			err := j.Kill()
			if err != nil {
				log.Println("error killing job", j.ID)
			}
		}
		h.RemoveJob(j.ID)
	}
	//itterate over all jobs and stop them
}
