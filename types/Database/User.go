package Database

type User struct {
	TgId        *int64
	UserName    *string
	FirstName   *string
	LastName    *string
	PhoneNumber *string
	Email       *string
	UOffset     *int
	ExpertSystem
}

type ExpertSystem struct {
	ExProjectName    *string
	ExBuildingLiter  *string
	ExFloorMin       *string
	ExFloorMax       *string
	ExRoomsAmountMin *string
	ExRoomsAmountMax *string
	ExSquareMin      *string
	ExSquareMax      *string
	ExCostMin        *string
	ExCostMax        *string
}
