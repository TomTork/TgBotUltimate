package messages

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/types/Database"
	"os"
	"strconv"
)

func GetMessagesByTgId(db *Database.DB, id uint64) ([]Database.Message, error) {
	rows, err := db.Query(db.Context, queries.Get("messages", "tg_id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]Database.Message, 0)
	for rows.Next() {
		var msg Database.Message
		err = rows.Scan(
			&msg.Id,
			&msg.TgId,
			&msg.CreatedAt,
			&msg.Message,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func getCountMessagesByTgId(db *Database.DB, id uint64) (uint8, error) {
	var count uint8
	err := db.QueryRow(db.Context, queries.Count("messages", "tg_id", id)).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func CreateMessage(db *Database.DB, message Database.ChatMessage) error {
	count, err := getCountMessagesByTgId(db, message.TgId)
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(os.Getenv("MESSAGE_HISTORY_COUNT"))
	if err != nil {
		return err
	}
	if !(count < uint8(limit)) {
		__message := Database.Message{}
		err = db.QueryRow(db.Context, queries.GetOneByMinValue("messages", "tg_id", "created_at")).Scan(
			&__message.Id,
			&__message.TgId,
			&__message.CreatedAt,
			&__message.Message,
		)
		if err != nil {
			return err
		}
		db.QueryRow(db.Context, queries.Delete("messages", "id", __message.Id))
	}
	db.QueryRow(db.Context, queries.Create("messages", queries.MessagesFields, queries.MessagesValues(message)))
	return nil
}
