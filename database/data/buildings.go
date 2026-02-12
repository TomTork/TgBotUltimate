package data

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Sync"
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
		&building.DeliveryDate,
		&building.BuildingAddress,
	)
	if err != nil {
		return nil, err
	}
	return &building, nil
}

func CreateBuilding(ctx context.Context, db *Database.DB, building Sync.Building) error {
	existsBuilding, _ := GetBuildingByCode(ctx, db, *building.Code)
	if existsBuilding == nil {
		_, err := db.Exec(
			ctx,
			queries.Create(
				"buildings",
				queries.BuildingsFields,
				queries.BuildingsValues(building),
			),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateBuilding(ctx context.Context, db *Database.DB, building Sync.Building) error {
	_, err := db.Exec(
		ctx,
		queries.UpdateS(
			"buildings",
			"code",
			*building.Code,
			queries.BuildingsFields,
			queries.BuildingsValues(building),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
