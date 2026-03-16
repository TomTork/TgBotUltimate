package favorites

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

func IsFavorite(ctx context.Context, db *Database.DB, userID int64, flatCode string) (bool, error) {
	var favoriteID uint64
	err := db.QueryRow(
		ctx,
		`SELECT id FROM user_favorite_flats WHERE user_tg_id = $1 AND flat_code = $2`,
		userID,
		flatCode,
	).Scan(&favoriteID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func AddFavorite(ctx context.Context, db *Database.DB, userID int64, flatCode string) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO user_favorite_flats (user_tg_id, flat_code) VALUES ($1, $2) ON CONFLICT (user_tg_id, flat_code) DO NOTHING`,
		userID,
		flatCode,
	)
	return err
}

func RemoveFavorite(ctx context.Context, db *Database.DB, userID int64, flatCode string) error {
	_, err := db.Exec(
		ctx,
		`DELETE FROM user_favorite_flats WHERE user_tg_id = $1 AND flat_code = $2`,
		userID,
		flatCode,
	)
	return err
}

func GetFavoriteFlatsByUser(ctx context.Context, db *Database.DB, userID int64) ([]Database.Query, error) {
	rows, err := db.Query(
		ctx,
		queries.FlatsQuery+` INNER JOIN user_favorite_flats uff ON uff.flat_code = f.code WHERE uff.user_tg_id = $1 ORDER BY uff.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	flats := make([]Database.Query, 0)
	for rows.Next() {
		var query Database.Query
		err = rows.Scan(
			&query.FlatCode,
			&query.ProjectName,
			&query.City,
			&query.District,
			&query.AddressOffice,
			&query.BuildingAddress,
			&query.BuildingName,
			&query.FlatNumber,
			&query.RoomsAmount,
			&query.Floor,
			&query.TotalSquare,
			&query.LivingSquare,
			&query.Cost,
			&query.FlatImg,
			&query.FloorImg,
			&query.Path,
			&query.PlaceType,
		)
		if err != nil {
			return nil, err
		}
		flats = append(flats, query)
	}

	return flats, rows.Err()
}
