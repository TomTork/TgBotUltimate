package queries

import (
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Sync/Sync1C"
)

var UsersFields = []string{"tg_id", "username", "first_name", "last_name", "phone_number", "email"}
var UsersValues = func(user Database.User) []interface{} {
	return []interface{}{user.TgId, user.UserName, user.FirstName, user.LastName, user.PhoneNumber, user.Email}
}

var MessagesFields = []string{"tg_id", "message"}
var MessagesValues = func(message Database.ChatMessage) []interface{} {
	return []interface{}{message.TgId, message.Message}
}

var ProjectsFields = []string{"code", "name"}
var ProjectsValues = func(project Sync1C.TypeProject) []interface{} {
	return []interface{}{project.ProjectId, project.ProjectName}
}

var BuildingsFields = []string{"code", "name", "project_code"}
var BuildingsValues = func(building Sync1C.TTypeBuilding) []interface{} {
	return []interface{}{building.BuildingId, building.BuildingName, building.ProjectCode}
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
var ApartmentsValues = func(apartment Sync1C.TTypeApartment) []interface{} {
	return []interface{}{
		apartment.ApartmentId,
		apartment.BuildingCode,
		apartment.Number,
		apartment.RoomsAmount,
		apartment.Floor,
		apartment.TotalSquare,
		apartment.LivingSquare,
		apartment.PriceTotal,
		apartment.FlatPlanImg,
		apartment.FloorPlanImg,
		apartment.Status,
		apartment.Type,
	}
}

var TagsFields = []string{
	"code",
	"flat_code",
	"name",
}
var TagsValues = func(tag Database.ITag) []interface{} {
	return []interface{}{tag.Code, tag.FlatCode, tag.Name}
}
