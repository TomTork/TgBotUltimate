package database

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func NewDatabase(ctx context.Context) (*Database.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	pool.QueryRow(ctx, queries.CreateUsersTable)
	pool.QueryRow(ctx, queries.CreateMessagesTable)

	return &Database.DB{
		Pool:    pool,
		Context: ctx,
	}, nil
}

func Close(db *Database.DB) {
	db.Pool.Close()
}
