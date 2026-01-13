package Database

type Info struct {
	Type string
	Name string
}

type Tag struct {
	Name string
}

type Query struct {
	ProjectName   string
	City          string
	District      string
	Address       string
	AddressOffice string
	BuildingName  string
	FlatNumber    uint32
	LivingSquare  float32
	TotalSquare   float32
	RoomsAmount   uint8
	Floor         uint16
	Cost          float32
	FlatImg       string
	FloorImg      string
	Path          string
	Status        uint8
	PlaceType     string

	Infos []Info
	Tags  []Tag
}
