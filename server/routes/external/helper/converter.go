package helper

import (
	"TgBotUltimate/types/Sync"
	"TgBotUltimate/types/Sync/Sync1C"
)

func ConvertProjectToType1C(project Sync1C.TypeProject) Sync.Project {
	return Sync.Project{
		Code: &project.ProjectId,
		Name: &project.ProjectName,
	}
}

func ConvertBuildingToType1C(building Sync1C.TTypeBuilding) Sync.Building {
	return Sync.Building{
		Code:        &building.BuildingId,
		Name:        &building.BuildingName,
		ProjectCode: &building.ProjectCode,
	}
}

func ConvertApartmentToType1C(apartment Sync1C.TTypeApartment) Sync.Flat {
	return Sync.Flat{
		Code:         &apartment.ApartmentId,
		Floor:        &apartment.Floor,
		FlatNumber:   &apartment.Number,
		PlaceType:    &apartment.Type,
		Status:       &apartment.Status,
		RoomsAmount:  &apartment.RoomsAmount,
		TotalSquare:  &apartment.TotalSquare,
		LivingSquare: &apartment.LivingSquare,
		Cost:         &apartment.PriceKvM,
		FlatImg:      &apartment.FlatPlanImg,
		FloorImg:     &apartment.FloorPlanImg,
		BuildingCode: &apartment.BuildingCode,
	}
}
