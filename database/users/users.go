package users

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"context"
	"fmt"
	"log"
	"strings"
)

func SaveAllUsersDataToFile(ctx context.Context, db *Database.DB) ([]byte, error) {
	rows, err := db.Query(ctx, queries.GetAll("users"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users = make([]Database.User, 0)
	for rows.Next() {
		var user Database.User
		err = rows.Scan(
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
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	var __users = make([]string, len(users))
	for _, user := range users {
		__users = append(__users, fmt.Sprintf("%d;%s;%s;%s;%s;%s", user.TgId, user.UserName, user.FirstName, user.LastName, user.PhoneNumber, user.Email))
	}
	return []byte(strings.Join(__users, "\n")), nil
}

func GetUserById(ctx context.Context, db *Database.DB, id int64) (*Database.User, error) {
	user := Database.User{}
	log.Println("req", queries.Get("users", "tg_id", uint64(id)))
	err := db.QueryRow(ctx, queries.Get("users", "tg_id", uint64(id))).Scan(
		&user.TgId,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Email,
		&user.ExProjectName,
		&user.ExBuildingLiter,
		&user.ExFloorMin,
		&user.ExFloorMax,
		&user.ExRoomsAmountMin,
		&user.ExRoomsAmountMax,
		&user.ExSquareMin,
		&user.ExSquareMax,
		&user.ExCostMin,
		&user.ExCostMax,
		&user.UOffset,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(ctx context.Context, db *Database.DB, user Database.User) error {
	existsUser, _ := GetUserById(ctx, db, *user.TgId)
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

func SetExpertSystemFields(ctx context.Context, db *Database.DB, id int64, system Database.ExpertSystem) error {
	err := DropUserOffset(ctx, db, id)
	if err != nil {
		return err
	}
	err = db.QueryRow(
		ctx,
		queries.Update(
			"users",
			"tg_id",
			uint64(id),
			queries.UserExpertSystem,
			queries.UserExpertSystemValues(system),
		),
	).Scan()
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(ctx context.Context, db *Database.DB, user Database.User) error {
	err := db.QueryRow(
		ctx,
		queries.Update(
			"users",
			"tg_id",
			uint64(*user.TgId),
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

func DropUserOffset(ctx context.Context, db *Database.DB, id int64) error {
	err := db.QueryRow(
		ctx,
		fmt.Sprintf("UPDATE users SET uoffset = 0 WHERE tg_id = %d", id),
	).Scan()
	if err != nil {
		return err
	}
	return nil
}

func IncreaseUserOffset(ctx context.Context, db *Database.DB, id int64) error {
	user, err := GetUserById(ctx, db, id)
	log.Println("my user", *user.TgId, *user.UOffset)
	if err != nil {
		log.Println("get user", err)
		return err
	}
	log.Println("User offset increased:", *user.UOffset, *user.UOffset+1)
	err = db.QueryRow(
		ctx,
		fmt.Sprintf("UPDATE users SET uoffset = %d WHERE tg_id = %d", *user.UOffset+1, id),
	).Scan()
	log.Println("error:", err, fmt.Sprintf("UPDATE users SET uoffset = %d WHERE tg_id = %d", *user.UOffset+1, id))
	user, err = GetUserById(ctx, db, id)
	log.Println("new user", *user.TgId, *user.UOffset)
	if err != nil {
		return err
	}
	return nil
}
