package Database

type Info struct {
	Type string
	Name string
}

type Tag struct {
	Name string
}

type FlatFilter struct {
	ProjectName  *string
	City         *string
	District     *string
	BuildingName *string
	FlatNumber   *uint32
	LivingSquare *float32
	TotalSquare  *float32
	RoomsAmount  *uint8
	Floor        *uint16
	Cost         *float32
	PlaceType    *string
	Offset       *uint16
	Limit        *uint16
}

type Flat struct {
	ProjectName  string
	City         string
	District     string
	BuildingName string
	FlatNumber   uint32
	LivingSquare float32
	TotalSquare  float32
	RoomsAmount  uint8
	Floor        uint16
	Cost         float32
	PlaceType    string
	Offset       uint16
	Limit        uint16
}

type Query struct {
	Flat
	Address       string
	AddressOffice string
	FlatImg       string
	FloorImg      string
	Path          string
	Status        uint8

	Infos []Info
	Tags  []Tag
}
