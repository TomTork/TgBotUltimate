package data

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Sync/Sync1C"
	"context"
)

func GetBuildingByCode(ctx context.Context, db *Database.DB, code string) (*Database.IBuilding, error) {
	building := Database.IBuilding{}
	err := db.QueryRow(ctx, queries.GetS("buildings", "code", code)).Scan(
		&building.Id,
		&building.ProjectCode,
		&building.Code,
		&building.Name,
		&building.Liter,
		&building.SectionNum,
		&building.SectionLiter,
	)
	if err != nil {
		return nil, err
	}
	return &building, nil
}

func CreateBuilding(ctx context.Context, db *Database.DB, building Sync1C.TTypeBuilding) error {
	existsBuilding, _ := GetBuildingByCode(ctx, db, building.BuildingId)
	if existsBuilding != nil {
		err := db.QueryRow(
			ctx,
			queries.Create(
				"buildings",
				queries.BuildingsFields,
				queries.BuildingsValues(building),
			),
		).Scan()
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateBuilding(ctx context.Context, db *Database.DB, building Sync1C.TTypeBuilding) error {
	err := db.QueryRow(
		ctx,
		queries.UpdateS(
			"buildings",
			"code",
			building.BuildingId,
			queries.BuildingsFields,
			queries.BuildingsValues(building),
		),
	).Scan()
	if err != nil {
		return err
	}
	return nil
}
