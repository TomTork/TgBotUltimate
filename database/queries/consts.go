package queries

import "TgBotUltimate/types/Database"

var UsersFields = []string{"tg_id", "name", "phone_number", "email"}
var UsersValues = func(user *Database.User) []interface{} {
	return []interface{}{user.TgId, user.Name, user.PhoneNumber, user.Email}
}
