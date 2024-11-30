package store

import (
	"database/sql"
	"log"
	"takehome/pkg/db"
)

type Store interface {
	CreateJob(j *db.Job) error
	UpdateJob(j *db.Job) error
	GetJobByID(id string) (*db.Job, error)
	GetJobs() ([]*db.Job, error)
}

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	log.Println("creating new store")
	return &store{
		db: db,
	}
}
