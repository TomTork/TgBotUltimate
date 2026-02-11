package main

import (
	db "TgBotUltimate/database"
	"TgBotUltimate/platform"
	"TgBotUltimate/processing/neuro"
	"TgBotUltimate/server"
	cron_tasks "TgBotUltimate/server/cron-tasks"
	"context"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	_ = godotenv.Load()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	var statements []string = strings.Split(os.Getenv("STATEMENTS"), ",")
	for _, statement := range statements {
		switch statement {
		case "neuro":
			g.Go(func() error {
				return neuro.InitPython()
			})
		case "platform":
			g.Go(func() error {
				return platform.Platform(ctx)
			})
			break
		case "server":
			g.Go(func() error {
				return server.RunHTTP()
			})
			break
		case "database":
			g.Go(func() error {
				_, err := db.NewDatabase(ctx)
				return err
			})
			break
		case "cron":
			g.Go(func() error {
				return cron_tasks.CronTasks(ctx)
			})
			break
		}
	}
	if err := g.Wait(); err != nil {
		log.Println("Error:", err)
	} else {
		log.Println("Shutdown gracefully")
	}
}
