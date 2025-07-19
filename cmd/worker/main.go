package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/riichi-mahjong-dev/backend-riichi/configs"
	"github.com/riichi-mahjong-dev/backend-riichi/database"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/jobs"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"gorm.io/gorm"
)

const maxWorkers = 5

func main() {
	env := configs.LoadEnv()
	dbConfig := env.LoadDatabaseConfig()
	db, err := database.ConnectDatabase(dbConfig)

	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
		return
	}

	handlers := map[string]func(context.Context, *gorm.DB, models.Job) error{
		"calculate_mmr": jobs.HandleMMR,
	}

	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				log.Println("Shutting down worker loop")
				return
			default:
				job, err := jobs.FetchPendingJob(db.Conn)
				log.Printf("fetch job %v\n", job)
				if err != nil {
					time.Sleep(5 * time.Second)
					continue
				}

				if job == nil {
					time.Sleep(5 * time.Second)
					continue
				}

				sem <- struct{}{}
				wg.Add(1)

				go func(job *models.Job) {
					defer wg.Done()
					defer func() { <-sem }()

					handler, ok := handlers[job.JobType]
					if !ok {
						log.Printf("Unknown job type: %s\n", job.JobType)
						jobs.MarkJobFailed(db.Conn, job.ID, "unknown job type")
						return
					}

					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()

					if err := handler(ctx, db.Conn, *job); err != nil {
						log.Printf("job %d failed: %v\n", job.ID, err)
						jobs.MarkJobFailed(db.Conn, job.ID, err.Error())
					} else {
						jobs.MarkJobDone(db.Conn, job.ID)
					}
				}(job)
			}
		}
	}()

	<-stopChan
	log.Println("Shutdown signal received")

	close(done)

	wg.Wait()
	log.Println("All jobs completed. Exiting")
}
