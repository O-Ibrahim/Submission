package store

import (
	"takehome/pkg/db"
	"time"
)

func (s *store) CreateJob(j *db.Job) error {
	query := "INSERT INTO jobs (id,command, args, status, logfile,created_at,updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	j.CreatedAt = time.Now().Unix()
	j.UpdatedAt = time.Now().Unix()
	_, err := s.db.Exec(query, j.ID, j.Command, j.Args, j.Status, j.Logfile, j.CreatedAt, j.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) UpdateJob(j *db.Job) error {
	query := "UPDATE jobs SET status = ?, updated_at = ? WHERE id = ?"
	j.UpdatedAt = time.Now().Unix()
	_, err := s.db.Exec(query, j.Status, j.UpdatedAt, j.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) GetJobByID(id string) (*db.Job, error) {
	query := "SELECT * FROM jobs WHERE id = ?"
	row := s.db.QueryRow(query, id)
	j := &db.Job{}
	err := row.Scan(&j.ID, &j.Command, &j.Args, &j.Status, &j.Logfile, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (s *store) GetJobs() ([]*db.Job, error) {
	query := "SELECT * FROM jobs"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	jobs := []*db.Job{}
	for rows.Next() {
		j := &db.Job{}
		err := rows.Scan(&j.ID, &j.Command, &j.Args, &j.Status, &j.Logfile, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}
