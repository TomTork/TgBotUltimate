package users

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/database/queries/helper"
	"TgBotUltimate/types/Database"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
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
		__users = append(__users, fmt.Sprintf(
			"%d;%s;%s;%s;%s;%s",
			derefInt64(user.TgId),
			derefString(user.UserName),
			derefString(user.FirstName),
			derefString(user.LastName),
			derefString(user.PhoneNumber),
			derefString(user.Email),
		))
	}
	return []byte(strings.Join(__users, "\n")), nil
}

func GetUserById(ctx context.Context, db *Database.DB, id int64) (*Database.User, error) {
	user := Database.User{}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(ctx context.Context, db *Database.DB, user Database.User) error {
	existsUser, _ := GetUserById(ctx, db, *user.TgId)
	if existsUser == nil {
		_, err := db.Exec(
			ctx,
			queries.Create(
				"users",
				queries.UsersFields,
				queries.UsersValues(user),
			),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetExpertSystemFields(ctx context.Context, db *Database.DB, id int64, system Database.ExpertSystem) error {
	fields, values := nonEmptyExpertSystemFields(system)
	if len(fields) == 0 {
		return nil
	}

	changed, err := expertSystemFieldsChanged(ctx, db, id, system)
	if err != nil {
		return err
	}
	if changed {
		err = DropUserOffset(ctx, db, id)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec(
		ctx,
		queries.Update(
			"users",
			"tg_id",
			uint64(id),
			fields,
			values,
		),
	)
	if err != nil {
		return err
	}
	return nil
}

func ReplaceExpertSystemFields(ctx context.Context, db *Database.DB, id int64, system Database.ExpertSystem) error {
	changed, err := expertSystemReplaceChanged(ctx, db, id, system)
	if err != nil {
		return err
	}
	if changed {
		err = DropUserOffset(ctx, db, id)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec(
		ctx,
		queries.Update(
			"users",
			"tg_id",
			uint64(id),
			queries.UserExpertSystem,
			normalizedExpertSystemValues(system),
		),
	)
	return err
}

func ResetExpertSystemFields(ctx context.Context, db *Database.DB, id int64) error {
	changed, err := expertSystemResetChanged(ctx, db, id)
	if err != nil {
		return err
	}
	if changed {
		err = DropUserOffset(ctx, db, id)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec(
		ctx,
		queries.DropExpertSystemFields(id),
	)
	return err
}

func UpdateUser(ctx context.Context, db *Database.DB, user Database.User) error {
	_, err := db.Exec(
		ctx,
		queries.Update(
			"users",
			"tg_id",
			uint64(*user.TgId),
			queries.UsersFields,
			queries.UsersValues(user),
		),
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(ctx context.Context, db *Database.DB, id uint64) (bool, error) {
	_, err := db.Exec(ctx, queries.Delete("users", "tg_id", id))
	if err != nil {
		return false, err
	}
	return true, nil
}

func DropUserOffset(ctx context.Context, db *Database.DB, id int64) error {
	_, err := db.Exec(
		ctx,
		fmt.Sprintf("UPDATE users SET uoffset = 0 WHERE tg_id = %d", id),
	)
	if err != nil {
		return err
	}
	return nil
}

func IncreaseUserOffset(ctx context.Context, db *Database.DB, id int64) error {
	user, err := GetUserById(ctx, db, id)
	if err != nil {
		return err
	}
	if user == nil || user.UOffset == nil {
		return nil
	}
	_, err = db.Exec(
		ctx,
		fmt.Sprintf("UPDATE users SET uoffset = %d WHERE tg_id = %d", *user.UOffset+1, id),
	)
	if err != nil {
		return err
	}
	return nil
}

func nonEmptyExpertSystemFields(system Database.ExpertSystem) ([]string, []interface{}) {
	allFields := queries.UserExpertSystem
	allValues := queries.UserExpertSystemValues(system)

	fields := make([]string, 0, len(allFields))
	values := make([]interface{}, 0, len(allValues))

	for i, value := range allValues {
		if !hasExpertValue(value) {
			continue
		}

		fields = append(fields, allFields[i])
		values = append(values, helper.SafeNil(value))
	}

	return fields, values
}

func hasExpertValue(value interface{}) bool {
	str, ok := helper.SafeNil(value).(string)
	if !ok {
		return false
	}

	return strings.TrimSpace(str) != ""
}

func expertSystemFieldsChanged(ctx context.Context, db *Database.DB, id int64, incoming Database.ExpertSystem) (bool, error) {
	user, err := GetUserById(ctx, db, id)
	if err != nil {
		return false, err
	}
	if user == nil {
		return true, nil
	}

	currentValues := queries.UserExpertSystemValues(user.ExpertSystem)
	incomingValues := queries.UserExpertSystemValues(incoming)

	for i := range incomingValues {
		incomingValue, ok := helper.SafeNil(incomingValues[i]).(string)
		if !ok || strings.TrimSpace(incomingValue) == "" {
			continue
		}

		currentValue, _ := helper.SafeNil(currentValues[i]).(string)
		if currentValue != incomingValue {
			return true, nil
		}
	}

	return false, nil
}

func expertSystemResetChanged(ctx context.Context, db *Database.DB, id int64) (bool, error) {
	user, err := GetUserById(ctx, db, id)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}

	for _, value := range queries.UserExpertSystemValues(user.ExpertSystem) {
		currentValue, _ := helper.SafeNil(value).(string)
		if strings.TrimSpace(currentValue) != "" {
			return true, nil
		}
	}

	return false, nil
}

func expertSystemReplaceChanged(ctx context.Context, db *Database.DB, id int64, incoming Database.ExpertSystem) (bool, error) {
	user, err := GetUserById(ctx, db, id)
	if err != nil {
		return false, err
	}
	if user == nil {
		return true, nil
	}

	currentValues := normalizedExpertSystemValues(user.ExpertSystem)
	incomingValues := normalizedExpertSystemValues(incoming)

	for i := range incomingValues {
		currentValue, _ := currentValues[i].(string)
		incomingValue, _ := incomingValues[i].(string)
		if currentValue != incomingValue {
			return true, nil
		}
	}

	return false, nil
}

func normalizedExpertSystemValues(system Database.ExpertSystem) []interface{} {
	rawValues := queries.UserExpertSystemValues(system)
	values := make([]interface{}, 0, len(rawValues))

	for _, value := range rawValues {
		normalized, ok := helper.SafeNil(value).(string)
		if !ok {
			normalized = ""
		}
		values = append(values, normalized)
	}

	return values
}

func stringPtr(value string) *string {
	return &value
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func derefInt64(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}
