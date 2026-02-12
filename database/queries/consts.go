package queries

import (
	"TgBotUltimate/database/queries/helper"
	"TgBotUltimate/types/Database"
	"TgBotUltimate/types/Sync"
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
