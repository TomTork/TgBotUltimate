package helper

import (
	"TgBotUltimate/types/Database"
	"fmt"
	"os"
	"strconv"
)

func CreateQueryForSearchFlats(data Database.FlatFilter) string {
	var where string = ""
	LimitMax, _ := strconv.Atoi(os.Getenv("LIMIT_MAX"))
	if data.Limit != nil && *data.Limit > uint16(LimitMax) {
		*data.Limit = uint16(LimitMax)
	} else if data.Limit == nil {
		*data.Limit = 3
	}
	if data.Offset == nil {
		*data.Offset = 0
	}
	if data.ProjectName != nil {
		where += fmt.Sprintf(" AND p.name LIKE '%%%s%%'", *data.ProjectName)
	}
	if data.City != nil {
		where += fmt.Sprintf(" AND p.city LIKE '%%%s%%'", *data.City)
	}
	if data.District != nil {
		where += fmt.Sprintf(" AND p.district LIKE '%%%s%%'", *data.District)
	}
	if data.BuildingName != nil {
		where += fmt.Sprintf(" AND b.name LIKE '%%%s%%'", *data.BuildingName)
	}
	if data.FlatNumber != nil {
		where += fmt.Sprintf(" AND f.flat_number = %d", *data.FlatNumber)
	}
	if data.LivingSquare != nil {
		where += fmt.Sprintf(" AND f.living_square = %f", *data.LivingSquare)
	}
	if data.TotalSquare != nil {
		where += fmt.Sprintf(" AND f.total_square = %f", *data.TotalSquare)
	}
	if data.RoomsAmount != nil {
		where += fmt.Sprintf(" AND f.rooms_amount = %d", *data.RoomsAmount)
	}
	if data.Floor != nil {
		where += fmt.Sprintf(" AND f.floor = %d", *data.Floor)
	}
	if data.Cost != nil {
		where += fmt.Sprintf(" AND f.cost = %f", *data.Cost)
	}
	if data.PlaceType != nil {
		if *data.PlaceType != "Квартира" || *data.PlaceType != "Подсобное" || *data.PlaceType != "Торговля/коммерция" {
			where += fmt.Sprintf(" AND f.place_type = 'Квартира'", *data.PlaceType)
		} else {
			where += fmt.Sprintf(" AND f.place_type = '%s'", *data.PlaceType)
		}
	}
	return where + ";"
}
