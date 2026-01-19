package cron_tasks

import (
	"github.com/go-co-op/gocron/v2"
	"log"
)

func CronTasks() error {
	s, _ := gocron.NewScheduler()

	_, err := s.NewJob(
		gocron.CronJob(
			"0 0 3 * * *",
			true,
		),
		gocron.NewTask(
			func() {
				log.Println("running...")
			},
		),
	)
	if err != nil {
		_ = s.Shutdown()
		return err
	}
	s.Start()
	return nil
}
