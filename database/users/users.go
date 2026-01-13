package users

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"context"
)

func GetUserById(ctx context.Context, db *Database.DB, id uint64) (*Database.User, error) {
	user := Database.User{}
	err := db.QueryRow(ctx, queries.Get("users", "tg_id", id)).Scan(
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

func CreateUser(ctx context.Context, db *Database.DB, user Database.User) error {
	existsUser, _ := GetUserById(ctx, db, user.TgId)
	if existsUser == nil {
		err := db.QueryRow(
			ctx,
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

func UpdateUser(ctx context.Context, db *Database.DB, user Database.User) error {
	err := db.QueryRow(
		ctx,
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

func DeleteUser(ctx context.Context, db *Database.DB, id uint64) (bool, error) {
	err := db.QueryRow(ctx, queries.Delete("users", "tg_id", id)).Scan()
	if err != nil {
		return false, err
	}
	return true, nil
}
