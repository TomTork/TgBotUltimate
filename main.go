package main

import (
	db "TgBotUltimate/database"
	"TgBotUltimate/platform"
	"TgBotUltimate/server"
	cron_tasks "TgBotUltimate/server/cron-tasks"
	"context"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	_ = godotenv.Load()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return platform.Platform(ctx)
	})
	g.Go(func() error {
		return server.RunHTTP()
	})
	g.Go(func() error {
		_, err := db.NewDatabase(ctx)
		return err
	})
	g.Go(func() error {
		return cron_tasks.CronTasks(ctx)
	})
	//g.Go(func() error {
	//	// ...
	//})
	if err := g.Wait(); err != nil {
		log.Println("Error:", err)
	} else {
		log.Println("Shutdown gracefully")
	}
}
