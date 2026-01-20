package data

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Sync/Sync1C"
	"context"
)

func GetProjectByCode(ctx context.Context, db *Database.DB, code string) (*Database.IProject, error) {
	project := Database.IProject{}
	err := db.QueryRow(ctx, queries.GetS("projects", "code", code)).Scan(
		&project.Id,
		&project.Code,
		&project.Name,
		&project.City,
		&project.District,
		&project.Address,
		&project.AddressOffice,
	)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func CreateProject(ctx context.Context, db *Database.DB, project Sync1C.TypeProject) error {
	existsProject, _ := GetProjectByCode(ctx, db, project.ProjectId)
	if existsProject != nil {
		err := db.QueryRow(
			ctx,
			queries.Create(
				"projects",
				queries.ProjectsFields,
				queries.ProjectsValues(project),
			),
		).Scan()
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateProject(ctx context.Context, db *Database.DB, project Sync1C.TypeProject) error {
	err := db.QueryRow(
		ctx,
		queries.UpdateS(
			"projects",
			"code",
			project.ProjectId,
			queries.ProjectsFields,
			queries.ProjectsValues(project),
		),
	).Scan()
	if err != nil {
		return err
	}
	return nil
}
