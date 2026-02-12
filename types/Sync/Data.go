package Sync

type Project struct {
	Code *string `json:"code"`
	Name *string `json:"name"`
}

type Building struct {
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Liter       *string `json:"liter"`
	ProjectCode *string `json:"project_code"`
}

type Section struct {
	Code         *string `json:"code"`
	SectionNum   *int    `json:"section_num"`
	SectionLiter *string `json:"section_liter"`
	BuildingCode *string `json:"building_code"`
}

type Flat struct {
	Id           *int     `json:"id"`
	Code         *string  `json:"code"`
	FlatNumber   *string  `json:"flat_number"`
	RoomsAmount  *uint8   `json:"rooms_amount"`
	Floor        *int     `json:"floor"`
	TotalSquare  *float32 `json:"total_square"`
	LivingSquare *float32 `json:"living_square"`
	Cost         *float32 `json:"cost"`
	FlatImg      *string  `json:"flat_img"`
	FloorImg     *string  `json:"floor_img"`
	Path         *string  `json:"path"`
	Status       *string  `json:"stature"`
	PlaceType    *string  `json:"place_type"`
	SectionCode  *string  `json:"section_code"`
	BuildingCode *string  `json:"building_code"`
}
