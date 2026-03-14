package Database

import "time"

type ChatMessage struct {
	TgId    uint64
	Message string
	Parameters
}

type Message struct {
	Id        uint64
	CreatedAt time.Time
	ChatMessage
}

type Parameters struct {
	ProjectName    string
	BuildingLiter  string
	FloorMin       string
	FloorMax       string
	RoomsAmountMin string
	RoomsAmountMax string
	SquareMin      string
	SquareMax      string
	CostMin        string
	CostMax        string
}

type indexParameter struct {
	Data  string
	Index int
}

type IndexesParameters struct {
	ProjectName    indexParameter
	BuildingLiter  indexParameter
	FloorMin       indexParameter
	FloorMax       indexParameter
	RoomsAmountMin indexParameter
	RoomsAmountMax indexParameter
	SquareMin      indexParameter
	SquareMax      indexParameter
	CostMin        indexParameter
	CostMax        indexParameter
}
