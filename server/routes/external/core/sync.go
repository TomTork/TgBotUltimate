package core

import (
	"TgBotUltimate/database"
	data2 "TgBotUltimate/database/data"
	"TgBotUltimate/server/routes"
	"TgBotUltimate/types/Sync/Sync1C"
	"context"
	"fmt"
	"os"
)

func Feed(ctx context.Context) string {
	db, err := database.NewDatabase(ctx)
	if err != nil {
		return err.Error()
	}
	var __projects []Sync1C.Project
	if err := routes.Get(ctx, fmt.Sprintf("%s/get_projects", os.Getenv("URL_1C")), nil, &__projects); err != nil {
		return err.Error()
	}
	for _, __project := range __projects {
		var data Sync1C.Data
		if err := routes.Get(ctx, fmt.Sprintf("%s/get_all?project_id=%s", os.Getenv("URL_1C"), __project.Uid), nil, &data); err != nil {
			return err.Error()
		}
		for _, project := range data.Projects {
			_project := Sync1C.TypeProject{ProjectId: project.ProjectId, ProjectName: project.ProjectName}
			if err := data2.CreateProject(ctx, db, _project); err != nil {
				return err.Error()
			} else if err := data2.UpdateProject(ctx, db, _project); err != nil {
				return err.Error()
			}
			for _, house := range project.Houses {
				for _, building := range house.Buildings {
					_building := Sync1C.TTypeBuilding{BuildingId: building.BuildingId, BuildingName: building.BuildingName, ProjectCode: project.ProjectId}
					if err := data2.CreateBuilding(ctx, db, _building); err != nil {
						return err.Error()
					} else if err := data2.UpdateBuilding(ctx, db, _building); err != nil {
						return err.Error()
					}
					for _, apartment := range building.Apartments {
						_flat := Sync1C.TTypeApartment{ApartmentId: apartment.ApartmentId, Floor: apartment.Floor, Number: apartment.Number, NumberOld: apartment.NumberOld, NumberForSort: apartment.NumberForSort, Type: apartment.Type, TypeAlias: apartment.TypeAlias, Status: apartment.Status, StatusText: apartment.StatusText, StatusColor: apartment.StatusColor, RoomsAmount: apartment.RoomsAmount, Tags: apartment.Tags, TotalSquare: apartment.TotalSquare, LivingSquare: apartment.LivingSquare, BltSquare: apartment.BltSquare, PriceKvM: apartment.PriceKvM, PriceTotal: apartment.PriceTotal, SalerInn: apartment.SalerInn, SalerName: apartment.SalerName, DateSale: apartment.DateSale, DogovorStatusText: apartment.DogovorStatusText, Pokupatel: apartment.Pokupatel, DogovorNumber: apartment.DogovorNumber, FlatPlanImg: apartment.FlatPlanImg, FloorPlanImg: apartment.FloorPlanImg, BuildingCode: building.BuildingId}
						if err := data2.CreateFlat(ctx, db, _flat); err != nil {
							return err.Error()
						} else if err := data2.UpdateFlat(ctx, db, _flat); err != nil {
							return err.Error()
						}
					}
				}
			}
		}
	}
	return "feed"
}
