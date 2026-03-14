package processing

import (
	messages2 "TgBotUltimate/database/messages"
	"TgBotUltimate/database/users"
	"TgBotUltimate/types/Database"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Summarize(ctx context.Context, db *Database.DB, id uint64) (*Database.Parameters, error) {
	var parameters Database.IndexesParameters
	var finalParameters Database.Parameters
	messages, err := messages2.GetMessagesByTgId(ctx, db, id)
	if err != nil {
		return nil, err
	}
	user, err := users.GetUserById(ctx, db, int64(id))
	if err != nil {
		return nil, err
	}
	for index, message := range messages {
		if message.ProjectName != "" {
			parameters.ProjectName.Data = message.ProjectName
			parameters.ProjectName.Index = index
		}
		if message.BuildingLiter != "" {
			parameters.BuildingLiter.Data = message.BuildingLiter
			parameters.BuildingLiter.Index = index
		}
		if message.FloorMin != "" {
			parameters.FloorMin.Data = message.FloorMin
			parameters.FloorMin.Index = index
		}
		if message.FloorMax != "" {
			parameters.FloorMax.Data = message.FloorMax
			parameters.FloorMax.Index = index
		}
		if message.RoomsAmountMin != "" {
			parameters.RoomsAmountMin.Data = message.RoomsAmountMin
			parameters.RoomsAmountMin.Index = index
		}
		if message.RoomsAmountMax != "" {
			parameters.RoomsAmountMax.Data = message.RoomsAmountMax
			parameters.RoomsAmountMax.Index = index
		}
		if message.SquareMin != "" {
			parameters.SquareMin.Data = message.SquareMin
			parameters.SquareMin.Index = index
		}
		if message.SquareMax != "" {
			parameters.SquareMax.Data = message.SquareMax
			parameters.SquareMax.Index = index
		}
		if message.CostMin != "" {
			parameters.CostMin.Data = message.CostMin
			parameters.CostMin.Index = index
		}
		if message.CostMax != "" {
			parameters.CostMax.Data = message.CostMax
			parameters.CostMax.Index = index
		}
	}
	limit, err := strconv.Atoi(os.Getenv("MESSAGE_HISTORY_COUNT"))
	if *user.ExProjectName != "" && parameters.ProjectName.Index < limit/2 {
		finalParameters.ProjectName = *user.ExProjectName
	} else {
		finalParameters.ProjectName = parameters.ProjectName.Data
	}
	if *user.ExBuildingLiter != "" && parameters.BuildingLiter.Index < limit/2 {
		finalParameters.BuildingLiter = *user.ExBuildingLiter
	} else {
		finalParameters.BuildingLiter = parameters.BuildingLiter.Data
	}
	if *user.ExFloorMin != "" && parameters.FloorMin.Index < limit/2 {
		finalParameters.FloorMin = *user.ExFloorMin
	} else {
		finalParameters.FloorMin = parameters.FloorMin.Data
	}
	if *user.ExFloorMax != "" && parameters.FloorMax.Index < limit/2 {
		finalParameters.FloorMax = *user.ExFloorMax
	} else {
		finalParameters.FloorMax = parameters.FloorMax.Data
	}
	if *user.ExRoomsAmountMin != "" && parameters.RoomsAmountMin.Index < limit/2 {
		finalParameters.RoomsAmountMin = *user.ExRoomsAmountMin
	} else {
		finalParameters.RoomsAmountMin = parameters.RoomsAmountMin.Data
	}
	if *user.ExRoomsAmountMax != "" && parameters.RoomsAmountMax.Index < limit/2 {
		finalParameters.RoomsAmountMax = *user.ExRoomsAmountMax
	} else {
		finalParameters.RoomsAmountMax = parameters.RoomsAmountMax.Data
	}
	if *user.ExSquareMin != "" && parameters.SquareMin.Index < limit/2 {
		finalParameters.SquareMin = *user.ExSquareMin
	} else {
		finalParameters.SquareMin = parameters.SquareMin.Data
	}
	if *user.ExSquareMax != "" && parameters.SquareMax.Index < limit/2 {
		finalParameters.SquareMax = *user.ExSquareMax
	} else {
		finalParameters.SquareMax = parameters.SquareMax.Data
	}
	if *user.ExCostMin != "" && parameters.CostMin.Index < limit/2 {
		finalParameters.CostMin = *user.ExCostMin
	} else {
		finalParameters.CostMin = parameters.CostMin.Data
	}
	if *user.ExCostMax != "" && parameters.CostMax.Index < limit/2 {
		finalParameters.CostMax = *user.ExCostMax
	} else {
		finalParameters.CostMax = parameters.CostMax.Data
	}
	return &finalParameters, nil
}

func Converter(parameters *Database.Parameters, user *Database.User) string {
	result := make([]string, 0)
	result = append(result, "f.status = 0")
	if parameters.ProjectName != "" {
		result = append(result, fmt.Sprintf("p.name LIKE '%%%s%%'", parameters.ProjectName))
	}
	if parameters.BuildingLiter != "" {
		result = append(result, fmt.Sprintf("b.liter LIKE '%%%s%%'", parameters.BuildingLiter))
	}
	if parameters.FloorMin != "" {
		result = append(result, fmt.Sprintf("f.floor >= %s", parameters.FloorMin))
	}
	if parameters.FloorMax != "" {
		result = append(result, fmt.Sprintf("f.floor <= %s", parameters.FloorMax))
	}
	if parameters.RoomsAmountMin != "" {
		result = append(result, fmt.Sprintf("f.rooms_amount >= %s", parameters.RoomsAmountMin))
	}
	if parameters.RoomsAmountMax != "" {
		result = append(result, fmt.Sprintf("f.rooms_amount <= %s", parameters.RoomsAmountMax))
	}
	if parameters.SquareMin != "" {
		result = append(result, fmt.Sprintf("f.total_square >= %s", parameters.SquareMin))
	}
	if parameters.SquareMax != "" {
		result = append(result, fmt.Sprintf("f.total_square <= %s", parameters.SquareMax))
	}
	if parameters.CostMin != "" {
		result = append(result, fmt.Sprintf("f.cost >= %s", parameters.CostMin))
	}
	if parameters.CostMax != "" {
		result = append(result, fmt.Sprintf("f.cost <= %s", parameters.CostMax))
	}
	return " WHERE " + strings.Join(result, " AND ") + " " + fmt.Sprintf("OFFSET %d", *user.UOffset*3) + " LIMIT 3"
}

func ShowFlat(flat Database.Query) (string, string, string) {
	result := make([]string, 0)
	fullCost := os.Getenv("FULL_COST") == "true"
	if flat.ProjectName != nil {
		result = append(result, fmt.Sprintf("Проект: %s", *flat.ProjectName))
	}
	if flat.City != nil {
		result = append(result, fmt.Sprintf("Город: %s", *flat.City))
	}
	if flat.District != nil {
		result = append(result, fmt.Sprintf("Район: %s", *flat.District))
	}
	if flat.AddressOffice != nil {
		result = append(result, fmt.Sprintf("Адрес офиса: %s", *flat.AddressOffice))
	}
	if flat.BuildingAddress != nil {
		result = append(result, fmt.Sprintf("Адрес здания: %s", *flat.BuildingAddress))
	}
	if flat.BuildingName != nil {
		result = append(result, fmt.Sprintf("Здание: %s", *flat.BuildingName))
	}
	if flat.FlatNumber != nil {
		result = append(result, fmt.Sprintf("№ квартиры: %d", *flat.FlatNumber))
	}
	if flat.RoomsAmount != nil {
		result = append(result, fmt.Sprintf("Количество комнат: %d", *flat.RoomsAmount))
	}
	if flat.Floor != nil {
		result = append(result, fmt.Sprintf("Этаж: %d", *flat.Floor))
	}
	if flat.TotalSquare != nil {
		result = append(result, fmt.Sprintf("Общая площадь: %.2f", *flat.TotalSquare))
	}
	if flat.LivingSquare != nil {
		result = append(result, fmt.Sprintf("Жилая площадь: %.2f", *flat.LivingSquare))
	}
	if flat.Cost != nil && fullCost {
		result = append(result, fmt.Sprintf("Цена: %.0f", *flat.Cost))
	} else if flat.Cost != nil && !fullCost && flat.TotalSquare != nil {
		result = append(result, fmt.Sprintf("Цена: %.0f", *flat.Cost**flat.TotalSquare))
	}
	var flatImg, floorImg string
	if flat.FlatImg != nil {
		flatImg = *flat.FlatImg
	} else {
		flatImg = ""
	}
	if flat.FloorImg != nil {
		floorImg = *flat.FloorImg
	} else {
		floorImg = ""
	}
	return strings.Join(result, "\n"), flatImg, floorImg
}
