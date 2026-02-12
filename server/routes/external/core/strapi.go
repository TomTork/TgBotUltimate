package core

import (
	"TgBotUltimate/database"
	data2 "TgBotUltimate/database/data"
	"TgBotUltimate/server/routes/helper"
	"TgBotUltimate/types/Sync/SyncStrapi"
	"context"
	"fmt"
	"os"
)

func Strapi(ctx context.Context) string {
	db, err := database.NewDatabase(ctx)
	if err != nil {
		return err.Error()
	}
	var data SyncStrapi.Strapi
	if err := helper.Get(ctx, fmt.Sprintf("%s/external/bot", os.Getenv("URL_STRAPI")), nil, &data); err != nil {
		return err.Error()
	}
	for _, project := range data.Projects {
		if err := data2.CreateProject(ctx, db, project); err != nil {
			return err.Error()
		} else if err := data2.UpdateProject(ctx, db, project); err != nil {
			return err.Error()
		}
	}
	for _, building := range data.Buildings {
		if err := data2.CreateBuilding(ctx, db, building); err != nil {
			return err.Error()
		} else if err := data2.UpdateBuilding(ctx, db, building); err != nil {
			return err.Error()
		}
	}
	for _, section := range data.Sections {
		if err := data2.CreateSection(ctx, db, section); err != nil {
			return err.Error()
		} else if err := data2.UpdateSection(ctx, db, section); err != nil {
			return err.Error()
		}
	}
	for _, flat := range data.Flats {
		if err := data2.CreateFlat(ctx, db, flat); err != nil {
			return err.Error()
		} else if err := data2.UpdateFlat(ctx, db, flat); err != nil {
			return err.Error()
		}
	}
	return "strapi"
}
