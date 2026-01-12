package queries

import "TgBotUltimate/types/Database"

var UsersFields = []string{"tg_id", "username", "first_name", "last_name", "phone_number", "email"}
var UsersValues = func(user Database.User) []interface{} {
	return []interface{}{user.TgId, user.UserName, user.FirstName, user.LastName, user.PhoneNumber, user.Email}
}

var MessagesFields = []string{"tg_id", "message"}
var MessagesValues = func(message Database.ChatMessage) []interface{} {
	return []interface{}{message.TgId, message.Message}
}
