package db

type Job struct {
	ID        string `json:"id"`
	Command   string `json:"command"`
	Args      string `json:"args"`
	Status    string `json:"status"`
	Logfile   string `json:"logfile"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
