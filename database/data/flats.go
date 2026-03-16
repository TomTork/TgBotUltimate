package data

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/processing"
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Sync"
	"context"
)

func GetFlatsByParameters(ctx context.Context, db *Database.DB, user *Database.User) ([]Database.Query, error) {
	summarize, err := processing.Summarize(ctx, db, uint64(*user.TgId))
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(ctx, queries.FlatsQuery+processing.Converter(summarize, user))
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
			&query.PhoneNumber,
			&query.Site,
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
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return flats, nil
}

func GetFlatByCode(ctx context.Context, db *Database.DB, code string) (*Database.IFlat, error) {
	flat := Database.IFlat{}
	err := db.QueryRow(ctx, queries.GetS("flats", "code", code)).Scan(
		&flat.Id,
		&flat.Code,
		&flat.BuildingCode,
		&flat.SectionCode,
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

func CreateFlat(ctx context.Context, db *Database.DB, flat Sync.Flat) error {
	existsApartment, _ := GetFlatByCode(ctx, db, *flat.Code)
	if existsApartment == nil {
		_, err := db.Exec(
			ctx,
			queries.Create(
				"flats",
				queries.ApartmentsFields,
				queries.ApartmentsValues(flat),
			),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateFlat(ctx context.Context, db *Database.DB, flat Sync.Flat) error {
	_, err := db.Exec(
		ctx,
		queries.UpdateS(
			"flats",
			"code",
			*flat.Code,
			queries.ApartmentsFields,
			queries.ApartmentsValues(flat),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
