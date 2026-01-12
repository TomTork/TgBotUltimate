package users

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
)

func GetUserById(db *Database.DB, id uint64) (*Database.User, error) {
	user := Database.User{}
	err := db.QueryRow(db.Context, queries.Get("users", "tg_id", id)).Scan(
		&user.TgId,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(db *Database.DB, user Database.User) error {
	existsUser, _ := GetUserById(db, user.TgId)
	if existsUser == nil {
		err := db.QueryRow(
			db.Context,
			queries.Create(
				"users",
				queries.UsersFields,
				queries.UsersValues(user),
			),
		).Scan()
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateUser(db *Database.DB, user Database.User) error {
	err := db.QueryRow(
		db.Context,
		queries.Update(
			"users",
			"tg_id",
			user.TgId,
			queries.UsersFields,
			queries.UsersValues(user),
		),
	).Scan()
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(db *Database.DB, id uint64) (bool, error) {
	err := db.QueryRow(db.Context, queries.Delete("users", "tg_id", id)).Scan()
	if err != nil {
		return false, err
	}
	return true, nil
}
