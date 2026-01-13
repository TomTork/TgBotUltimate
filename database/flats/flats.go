package flats

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"context"
)

func GetFlats(ctx context.Context, db *Database.DB, whereQuery string) ([]Database.Query, error) {
	rows, err := db.Query(ctx, queries.FlatsQuery+" "+whereQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	qs := make([]Database.Query, 0)
	for rows.Next() {
		var query Database.Query
		err = rows.Scan(
			&query.ProjectName,
			&query.City,
			&query.District,
			&query.Address,
			&query.AddressOffice,
			&query.BuildingName,
			&query.FlatNumber,
			&query.LivingSquare,
			&query.TotalSquare,
			&query.RoomsAmount,
			&query.Floor,
			&query.Cost,
			&query.FlatImg,
			&query.FloorImg,
			&query.Path,
			&query.Status,
			&query.PlaceType,
			&query.Infos,
			&query.Tags,
		)
		if err != nil {
			return nil, err
		}
		qs = append(qs, query)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return qs, nil
}
