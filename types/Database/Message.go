package Database

import "time"

type ChatMessage struct {
	TgId           uint64
	Message        string
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

type Message struct {
	Id        uint64
	CreatedAt time.Time
	ChatMessage
}
