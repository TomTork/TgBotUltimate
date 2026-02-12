package data

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Sync"
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
		&project.AddressOffice,
	)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func CreateProject(ctx context.Context, db *Database.DB, project Sync.Project) error {
	existsProject, _ := GetProjectByCode(ctx, db, *project.Code)
	if existsProject == nil {
		_, err := db.Exec(
			ctx,
			queries.Create(
				"projects",
				queries.ProjectsFields,
				queries.ProjectsValues(project),
			),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateProject(ctx context.Context, db *Database.DB, project Sync.Project) error {
	_, err := db.Exec(
		ctx,
		queries.UpdateS(
			"projects",
			"code",
			*project.Code,
			queries.ProjectsFields,
			queries.ProjectsValues(project),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
