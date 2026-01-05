package platform

import (
	"TgBotUltimate/platform/telegram"
	"context"
	"golang.org/x/sync/errgroup"
	"os"
)

func Platform(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return telegram.Telegram(ctx, os.Getenv("TELEGRAM_TOKEN"))
	})
	//g.Go(func() error {
	//	// ...
	//})

	return g.Wait()
}
