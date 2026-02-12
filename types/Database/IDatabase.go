package Database

type IDatabase interface {
	GetFlats()
}

type IProject struct {
	Id            *uint64
	Code          *string
	Name          *string
	City          *string
	District      *string
	AddressOffice *string
}

type IProjectInfo struct {
	Id          *uint64
	Code        *string
	ProjectCode *string
	Type        *string
	Name        *string
}

type IBuilding struct {
	Id              *uint64
	ProjectCode     *string
	Code            *string
	Name            *string
	Liter           *string
	DeliveryDate    *string
	BuildingAddress *string
}

type ISection struct {
	Id           *uint64
	Code         *string
	BuildingCode *string
	SectionNum   *string
	SectionLiter *string
}

type IFlat struct {
	Id           *uint64
	Code         *string
	BuildingCode *string
	SectionCode  *string
	FlatNumber   *uint32
	RoomsAmount  *uint8
	Floor        *uint8
	TotalSquare  *float32
	LivingSquare *float32
	Cost         *float32
	FlatImg      *string
	FloorImg     *string
	Path         *string
	Status       *uint8
	PlaceType    *string
}

type ITag struct {
	Id       *uint64
	Code     *string
	FlatCode *string
	Name     *string
}
