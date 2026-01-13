package platform

import (
	"TgBotUltimate/database"
	"TgBotUltimate/platform/telegram"
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
)

func Platform(ctx context.Context) error {
	db, err := database.NewDatabase(ctx)
	if err != nil {
		log.Println("Failed to connect to database")
	}
	defer database.Close(db)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return telegram.Telegram(ctx, os.Getenv("TELEGRAM_TOKEN"), db)
	})
	//g.Go(func() error {
	//	// ...
	//})

	return g.Wait()
}
