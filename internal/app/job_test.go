package app

import (
	"bytes"
	"io"
	"os"
	"takehome/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JobCreation(t *testing.T) {
	jobHub := NewJobHub()

	job, err := NewJob(jobHub, "echo", "hello")
	assert.Nil(t, err)
	assert.NotNil(t, job)
	assert.Equal(t, job.Status, New)
	assert.Equal(t, job.Command, "echo")
	assert.Equal(t, job.Args, []string{"hello"})
	assert.NotNil(t, job.Logfile)
	// cleanup
	os.Remove(job.Logfile.Name())
}

func Test_JobRun(t *testing.T) {
	jobHub := NewJobHub()

	job, err := NewJob(jobHub, "echo", "hello")
	assert.Nil(t, err)
	assert.NotNil(t, job)

	err = job.Run()
	assert.Nil(t, err)
	assert.Equal(t, job.Status, Finished)
	//Read log file to make sure that the command was executed
	logFile, _ := os.Open(job.Logfile.Name())
	var buff bytes.Buffer

	io.Copy(&buff, logFile)
	assert.Equal(t, "hello\n", buff.String())
	// cleanup
	os.Remove(job.Logfile.Name())
}

func Test_JobRunWithFaultyCommand(t *testing.T) {
	jobHub := NewJobHub()

	job, err := NewJob(jobHub, "ecasdasdasdho", "hasdasdello")
	assert.Nil(t, err)
	assert.NotNil(t, job)

	err = job.Run()
	assert.NotNil(t, err)
	// cleanup
	os.Remove(job.Logfile.Name())
}

func Test_CreateLogFile(t *testing.T) {
	file, err := createLogFile("test")
	assert.Nil(t, err)
	assert.NotNil(t, file)
	//cleanup
	os.Remove(file.Name())
}

// func Test_KillJob(t *testing.T) {
// 	jobHub := NewJobHub()

// 	job, err := NewJob(jobHub, "sleep", "3")
// 	assert.Nil(t, err)
// 	assert.NotNil(t, job)
// 	go job.Run()
// 	time.Sleep(1 * time.Second)
// 	err = job.Kill()
// 	assert.Nil(t, err)
// 	assert.Equal(t, Killed, job.Status)
// }

func Test_ToMode(t *testing.T) {
	jobHub := NewJobHub()

	job, err := NewJob(jobHub, "echo", "hello")
	assert.Nil(t, err)
	assert.NotNil(t, job)
	model := job.toModel()
	assert.Equal(t, model.ID, job.ID)
	assert.Equal(t, model.Command, job.Command)
	assert.Equal(t, model.Status, string(job.Status))
	assert.Equal(t, model.Logfile, job.Logfile.Name())

	//cleanup
	os.Remove(job.Logfile.Name())

}

func Test_FromModel(t *testing.T) {
	jobhub := NewJobHub()
	jobModel := &db.Job{
		ID:      "test",
		Command: "echo",
		Args:    "hello",
		Status:  "finished",
		Logfile: "test.log",
	}
	var job Job
	err := job.fromModel(jobModel, jobhub)
	assert.Nil(t, err)
	assert.Equal(t, job.ID, jobModel.ID)
	assert.Equal(t, job.Command, jobModel.Command)
	assert.Equal(t, job.Args, []string{"hello"})
	assert.Equal(t, job.Status, Finished)
	assert.Equal(t, job.Logfile.Name(), jobModel.Logfile)
	//cleanup
	os.Remove(job.Logfile.Name())
}

func Test_GetModeLogFile(t *testing.T) {
	jobhub := NewJobHub()
	job, err := NewJob(jobhub, "echo", "hello")
	assert.Nil(t, err)
	assert.NotNil(t, job)

	file, err := getModelLogFile(job.Logfile.Name())
	assert.Nil(t, err)
	assert.NotNil(t, file)
	//cleanup
	os.Remove(job.Logfile.Name())
}
