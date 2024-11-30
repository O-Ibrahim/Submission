package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"takehome/config"
	"takehome/pkg/db"
	"takehome/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// App is the main application struct
type App struct {
	DB     *sql.DB
	Store  store.Store
	JobHub *JobHub
	Config *config.Config
}

// NewApp creates a new App struct
func NewApp() (*App, error) {
	log.Println("creating new app")
	db, err := db.GetDBDriver()
	if err != nil {
		return nil, err
	}
	store := store.NewStore(db)
	hub := NewJobHub()
	config := config.NewConfig()
	return &App{
		DB:     db,
		JobHub: hub,
		Config: config,
		Store:  store,
	}, nil
}

// Run starts the app
func (a *App) Run() {
	a.killUnfinishedJobs()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))

	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/jobs", func(jobsRoute chi.Router) {
		jobsRoute.Use(a.AuthMiddleware)
		jobsRoute.Post("/", a.handleCreateJob)
		jobsRoute.Get("/", a.handleGetAllJobs)
		jobsRoute.Route("/{id}", func(idRoute chi.Router) {
			idRoute.Get("/", a.handleGetJobByID)
			idRoute.Get("/status", a.handleGetStatus)
			idRoute.Get("/logs", a.handleGetLogs)
			idRoute.Get("/kill", a.handleKillJob)
		})
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.Config.Port),
		Handler: r,
	}

	go func() {
		log.Println("starting server on port", a.Config.Port)
		log.Fatal(server.ListenAndServe())
	}()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	log.Println("shutting down server")
	//Graceful shutdown - Cleanup
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownRelease()
	a.JobHub.Shutdown(shutdownCtx)
	a.DB.Close()
	server.Shutdown(shutdownCtx)
}

// killUnfinishedJobs kills all jobs that are still running when the app starts
func (a *App) killUnfinishedJobs() {
	jobs, err := a.Store.GetJobs()
	if err != nil {
		log.Println("error getting all jobs:", err)
		return
	}
	for _, job := range jobs {
		if job.Status == "running" {
			job.Status = "killed"
			err := a.Store.UpdateJob(job)
			if err != nil {
				log.Println("error updating job:", err)
			}
			go a.SendJobToHook(job.ID, job.Status)
		}
	}
}
