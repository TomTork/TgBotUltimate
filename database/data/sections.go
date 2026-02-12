package data

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Sync"
	"context"
)

func GetSectionsByCode(ctx context.Context, db *Database.DB, code string) (*Database.ISection, error) {
	section := Database.ISection{}
	err := db.QueryRow(ctx, queries.GetS("sections", "code", code)).Scan(
		&section.Id,
		&section.Code,
		&section.BuildingCode,
		&section.SectionNum,
		&section.SectionLiter,
	)
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func CreateSection(ctx context.Context, db *Database.DB, section Sync.Section) error {
	existsSection, _ := GetSectionsByCode(ctx, db, *section.Code)
	if existsSection == nil {
		_, err := db.Exec(
			ctx,
			queries.Create(
				"sections",
				queries.SectionsFields,
				queries.SectionsValues(section),
			),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateSection(ctx context.Context, db *Database.DB, section Sync.Section) error {
	_, err := db.Exec(
		ctx,
		queries.UpdateS(
			"sections",
			"code",
			*section.Code,
			queries.SectionsFields,
			queries.SectionsValues(section),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
