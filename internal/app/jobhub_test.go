package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateJobHub(t *testing.T) {
	jobHub := NewJobHub()
	assert.NotNil(t, jobHub)
	assert.NotNil(t, jobHub.Jobs)
	assert.Len(t, jobHub.Jobs, 0)
}

func Test_AddJob(t *testing.T) {
	jobHub := NewJobHub()
	job := &Job{
		ID: "1",
	}
	jobHub.AddJob(job)
	assert.Len(t, jobHub.Jobs, 1)
}

func Test_GetJob(t *testing.T) {
	jobHub := NewJobHub()
	job := &Job{
		ID: "1",
	}
	jobHub.AddJob(job)
	assert.Equal(t, jobHub.GetJob("1"), job)
}

func Test_RemoveJob(t *testing.T) {
	jobHub := NewJobHub()
	job := &Job{
		ID: "1",
	}
	jobHub.AddJob(job)
	jobHub.RemoveJob("1")
	assert.Len(t, jobHub.Jobs, 0)
}

func Test_UpdateJob(t *testing.T)	{
	jobHub := NewJobHub()
	job := &Job{
		ID: "1",
	}
	jobHub.AddJob(job)
	job.Status = Running
	jobHub.UpdateJob(job)
	assert.Equal(t, jobHub.GetJob("1").Status, Running)
}

func Test_Shutdown(t *testing.T) {
	jobHub := NewJobHub()
	job := &Job{
		ID: "1",
	}
	jobHub.AddJob(job)
	job.Status = Running
	jobHub.Shutdown(context.Background())
	assert.Len(t, jobHub.Jobs, 0)
}