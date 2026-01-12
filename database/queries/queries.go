package queries

import (
	"TgBotUltimate/database/queries/helper"
	"fmt"
	"log"
	"strings"
)

const CreateUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	tg_id BIGINT PRIMARY KEY,
	username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
	phone_number VARCHAR(12),
	email VARCHAR(255)
);
`

const CreateMessagesTable = `
CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	tg_id BIGINT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	message TEXT NOT NULL,
	CONSTRAINT fk_user
		FOREIGN KEY (tg_id)
		REFERENCES users(tg_id)
		ON DELETE CASCADE
		ON UPDATE CASCADE
);
`

func Get(tableName string, idName string, id uint64) string {
	return fmt.Sprintf(`
		SELECT * FROM %s WHERE %s = %d;
	`, tableName, idName, id)
}

func GetOneByMinValue(tableName string, idName string, minValue string) string {
	log.Println(fmt.Sprintf(`
		SELECT DISTINCT ON (%s)
    		*
		FROM %s
		ORDER BY %s, %s ASC LIMIT 1;
	`, idName, tableName, idName, minValue))
	return fmt.Sprintf(`
		SELECT DISTINCT ON (%s)
    		*
		FROM %s
		ORDER BY %s, %s ASC LIMIT 1;
	`, idName, tableName, idName, minValue)
}

func Create(tableName string, fields []string, values []interface{}) string {
	return fmt.Sprintf(`
		INSERT INTO %s (%s) VALUES (%s);
	`, tableName, strings.Join(fields, ","), helper.ConvertValuesToSQLCreate(values))
}

func Update(tableName string, idName string, id uint64, fields []string, values []interface{}) string {
	return fmt.Sprintf(`
		UPDATE %s
		SET %s
		WHERE %s = %d;
	`, tableName, helper.ConvertValuesToSQLUpdate(fields, values), idName, id)
}

func Delete(tableName string, idName string, id uint64) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE %s = %d;`, tableName, idName, id)
}

func Count(tableName string, idName string, id uint64) string {
	log.Println(fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s = %d;`, tableName, idName, id))
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s = %d;`, tableName, idName, id)
}
