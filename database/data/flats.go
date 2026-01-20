package data

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/database/queries/helper"
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Sync/Sync1C"
	"context"
)

func GetFlats(ctx context.Context, db *Database.DB, data Database.FlatFilter) ([]Database.Query, error) {
	rows, err := db.Query(ctx, queries.FlatsQuery+helper.CreateQueryForSearchFlats(data))
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

func GetFlatByCode(ctx context.Context, db *Database.DB, code string) (*Database.IFlat, error) {
	flat := Database.IFlat{}
	err := db.QueryRow(ctx, queries.GetS("flats", "code", code)).Scan(
		&flat.Id,
		&flat.Code,
		&flat.BuildingCode,
		&flat.FlatNumber,
		&flat.RoomsAmount,
		&flat.Floor,
		&flat.TotalSquare,
		&flat.LivingSquare,
		&flat.Cost,
		&flat.FlatImg,
		&flat.FloorImg,
		&flat.Path,
		&flat.Status,
		&flat.PlaceType,
	)
	if err != nil {
		return nil, err
	}
	return &flat, nil
}

func CreateFlat(ctx context.Context, db *Database.DB, flat Sync1C.TTypeApartment) error {
	existsApartment, _ := GetFlatByCode(ctx, db, flat.ApartmentId)
	if existsApartment != nil {
		err := db.QueryRow(
			ctx,
			queries.Create(
				"flats",
				queries.ApartmentsFields,
				queries.ApartmentsValues(flat),
			),
		).Scan()
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateFlat(ctx context.Context, db *Database.DB, flat Sync1C.TTypeApartment) error {
	err := db.QueryRow(
		ctx,
		queries.UpdateS(
			"flats",
			"code",
			flat.ApartmentId,
			queries.ApartmentsFields,
			queries.ApartmentsValues(flat),
		),
	).Scan()
	if err != nil {
		return err
	}
	return nil
}
