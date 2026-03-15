package queries

import (
	"TgBotUltimate/database/queries/helper"
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Sync"
)

var UsersFields = []string{"tg_id", "username", "first_name", "last_name", "phone_number", "email"}
var UsersValues = func(user Database.User) []interface{} {
	var username, firstname, lastname, phoneNumber, email string
	if user.UserName != nil {
		username = *user.UserName
	} else {
		username = ""
	}
	if user.FirstName != nil {
		firstname = *user.FirstName
	} else {
		firstname = ""
	}
	if user.LastName != nil {
		lastname = *user.LastName
	} else {
		lastname = ""
	}
	if user.PhoneNumber != nil {
		phoneNumber = *user.PhoneNumber
	} else {
		phoneNumber = ""
	}
	if user.Email != nil {
		email = *user.Email
	} else {
		email = ""
	}
	return []interface{}{*user.TgId, username, firstname, lastname, phoneNumber, email}
}

var UserExpertSystem = []string{"ex_project_name", "ex_building_liter", "ex_floor_min", "ex_floor_max", "ex_rooms_amount_min", "ex_rooms_amount_max", "ex_square_min", "ex_square_max", "ex_cost_min", "ex_cost_max"}
var UserExpertSystemValues = func(system Database.ExpertSystem) []interface{} {
	return []interface{}{system.ExProjectName, system.ExBuildingLiter, system.ExFloorMin, system.ExFloorMax, system.ExRoomsAmountMin, system.ExRoomsAmountMax, system.ExSquareMin, system.ExSquareMax, system.ExCostMin, system.ExCostMax}
}

var UserExpertSystemAnswersFields = []string{"user_tg_id", "question_id", "variant_index"}
var UserExpertSystemAnswersValues = func(answer Database.ExpertSystemAnswer) []interface{} {
	return []interface{}{answer.UserTgID, answer.QuestionID, answer.VariantIndex}
}

var MessagesFields = []string{"tg_id", "message", "project_name", "building_liter", "floor_min", "floor_max", "rooms_amount_min", "rooms_amount_max", "square_min", "square_max", "cost_min", "cost_max"}
var MessagesValues = func(message Database.ChatMessage) []interface{} {
	if message.ProjectName == "<UNK>" {
		message.ProjectName = ""
	}
	if message.BuildingLiter == "<UNK>" {
		message.BuildingLiter = ""
	}
	if message.FloorMin == "<UNK>" {
		message.FloorMin = ""
	}
	if message.FloorMax == "<UNK>" {
		message.FloorMax = ""
	}
	if message.RoomsAmountMin == "<UNK>" {
		message.RoomsAmountMin = ""
	}
	if message.RoomsAmountMax == "<UNK>" {
		message.RoomsAmountMax = ""
	}
	if message.SquareMin == "<UNK>" {
		message.SquareMin = ""
	}
	if message.SquareMax == "<UNK>" {
		message.SquareMax = ""
	}
	if message.CostMin == "<UNK>" {
		message.CostMin = ""
	}
	if message.CostMax == "<UNK>" {
		message.CostMax = ""
	}
	return []interface{}{
		message.TgId,
		message.Message,
		message.ProjectName[0:min(255, len(message.ProjectName))],
		message.BuildingLiter[0:min(7, len(message.BuildingLiter))],
		message.FloorMin[0:min(7, len(message.FloorMin))],
		message.FloorMax[0:min(7, len(message.FloorMax))],
		message.RoomsAmountMin[0:min(7, len(message.RoomsAmountMin))],
		message.RoomsAmountMax[0:min(7, len(message.RoomsAmountMax))],
		message.SquareMin[0:min(7, len(message.SquareMin))],
		message.SquareMax[0:min(7, len(message.SquareMax))],
		message.CostMin[0:min(7, len(message.CostMin))],
		message.CostMax[0:min(7, len(message.CostMax))],
	}
}

var ProjectsFields = []string{"code", "name"}
var ProjectsValues = func(project Sync.Project) []interface{} {
	return []interface{}{*project.Code, *project.Name}
}

var BuildingsFields = []string{"code", "name", "project_code", "liter"}
var BuildingsValues = func(building Sync.Building) []interface{} {
	return []interface{}{*building.Code, *building.Name, *building.ProjectCode, *building.Liter}
}

var SectionsFields = []string{"code", "building_code", "section_num", "section_liter"}
var SectionsValues = func(section Sync.Section) []interface{} {
	return []interface{}{*section.Code, *section.BuildingCode, *section.SectionNum, helper.SafeNil(section.SectionLiter)}
}

var ApartmentsFields = []string{
	"code",
	"building_code",
	"flat_number",
	"rooms_amount",
	"floor",
	"total_square",
	"living_square",
	"cost",
	"flat_img",
	"floor_img",
	"status",
	"place_type",
}
var ApartmentsValues = func(apartment Sync.Flat) []interface{} {
	return []interface{}{
		*apartment.Code,
		*apartment.BuildingCode,
		helper.SafeNil(apartment.FlatNumber),
		helper.SafeNil(apartment.RoomsAmount),
		helper.SafeNil(apartment.Floor),
		helper.SafeNil(apartment.TotalSquare),
		helper.SafeNil(apartment.LivingSquare),
		helper.SafeNil(apartment.Cost),
		helper.SafeNil(apartment.FlatImg),
		helper.SafeNil(apartment.FloorImg),
		helper.SafeNil(apartment.Status),
		helper.SafeNil(apartment.PlaceType),
	}
}

var TagsFields = []string{
	"code",
	"flat_code",
	"name",
}
var TagsValues = func(tag Database.ITag) []interface{} {
	return []interface{}{*tag.Code, *tag.FlatCode, *tag.Name}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
