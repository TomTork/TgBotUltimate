package messages

import (
	"TgBotUltimate/database/queries"
	"TgBotUltimate/database/users"
	"TgBotUltimate/types/Database"
	"context"
	"os"
	"strconv"
)

func GetMessagesByTgId(ctx context.Context, db *Database.DB, id uint64) ([]Database.Message, error) {
	rows, err := db.Query(ctx, queries.GetSort("messages", "tg_id", id, "created_at", "ASC"))
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
			&msg.ProjectName,
			&msg.BuildingLiter,
			&msg.FloorMin,
			&msg.FloorMax,
			&msg.RoomsAmountMin,
			&msg.RoomsAmountMax,
			&msg.SquareMin,
			&msg.SquareMax,
			&msg.CostMin,
			&msg.CostMax,
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

func getCountMessagesByTgId(ctx context.Context, db *Database.DB, id uint64) (uint8, error) {
	var count uint8
	err := db.QueryRow(ctx, queries.Count("messages", "tg_id", id)).Scan(&count)
	if err != nil {
		return 0, nil
	}
	return count, nil
}

func ParametersHasBeenChanged(ctx context.Context, db *Database.DB, newMessage Database.ChatMessage) (bool, error) {
	messages, err := GetMessagesByTgId(ctx, db, newMessage.TgId)
	if err != nil {
		return false, err
	}
	if len(messages) == 0 {
		return true, nil
	}
	lastMessage := messages[len(messages)-1]
	if lastMessage.ProjectName != newMessage.ProjectName ||
		lastMessage.BuildingLiter != newMessage.BuildingLiter ||
		lastMessage.FloorMin != newMessage.FloorMin ||
		lastMessage.FloorMax != newMessage.FloorMax ||
		lastMessage.RoomsAmountMin != newMessage.RoomsAmountMin ||
		lastMessage.RoomsAmountMax != newMessage.RoomsAmountMax ||
		lastMessage.SquareMin != newMessage.SquareMin ||
		lastMessage.SquareMax != newMessage.SquareMax ||
		lastMessage.CostMin != newMessage.CostMin ||
		lastMessage.CostMax != newMessage.CostMax {
		return true, nil
	}
	return false, nil
}

func CreateMessage(ctx context.Context, db *Database.DB, message Database.ChatMessage) error {
	message = removeUNK(message)
	parametersChanged, err := ParametersHasBeenChanged(ctx, db, message)
	if err != nil {
		return err
	} else if parametersChanged {
		_ = users.DropUserOffset(ctx, db, int64(message.TgId))
	}
	count, err := getCountMessagesByTgId(ctx, db, message.TgId)
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(os.Getenv("MESSAGE_HISTORY_COUNT"))
	if err != nil {
		return err
	}
	if !(count < uint8(limit)) {
		__message := Database.Message{}
		err = db.QueryRow(ctx, queries.GetOneByMinValue("messages", "tg_id", "created_at")).Scan(
			&__message.Id,
			&__message.TgId,
			&__message.CreatedAt,
			&__message.Message,
			&__message.ProjectName,
			&__message.BuildingLiter,
			&__message.FloorMin,
			&__message.FloorMax,
			&__message.RoomsAmountMin,
			&__message.RoomsAmountMax,
			&__message.SquareMin,
			&__message.SquareMax,
			&__message.CostMin,
			&__message.CostMax,
		)
		if err != nil {
			return err
		}
		db.QueryRow(ctx, queries.Delete("messages", "id", __message.Id))
	}
	db.QueryRow(ctx, queries.Create("messages", queries.MessagesFields, queries.MessagesValues(message)))
	_ = users.IncreaseUserOffset(ctx, db, int64(message.TgId))
	return nil
}

func removeUNK(message Database.ChatMessage) Database.ChatMessage {
	if message.ProjectName == "<UNK>" {
		message.ProjectName = ""
	}
	if message.BuildingLiter == "<UNK>" {
		message.BuildingLiter = ""
	}
	if message.FloorMin == "<UNK>" {
		message.FloorMin = ""
	}
	if message.FloorMax == "<UNK>" {
		message.FloorMax = ""
	}
	if message.RoomsAmountMin == "<UNK>" {
		message.RoomsAmountMin = ""
	}
	if message.RoomsAmountMax == "<UNK>" {
		message.RoomsAmountMax = ""
	}
	if message.SquareMin == "<UNK>" {
		message.SquareMin = ""
	}
	if message.SquareMax == "<UNK>" {
		message.SquareMax = ""
	}
	if message.CostMin == "<UNK>" {
		message.CostMin = ""
	}
	if message.CostMax == "<UNK>" {
		message.CostMax = ""
	}
	return message
}
