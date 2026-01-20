package cron_tasks

import (
	"TgBotUltimate/server/routes/external/core"
	"context"
	"github.com/go-co-op/gocron/v2"
)

func CronTasks(ctx context.Context) error {
	s, _ := gocron.NewScheduler()

	_, err := s.NewJob(
		gocron.CronJob(
			"0 0 3 * * *",
			true,
		),
		gocron.NewTask(
			func() {
				core.Feed(ctx)
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
