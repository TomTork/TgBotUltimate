package Neuro

type Response struct {
	ProjectName    string `json:"project_name"`
	BuildingLiter  string `json:"building_liter"`
	FloorMin       string `json:"floor_min"`
	FloorMax       string `json:"floor_max"`
	RoomsAmountMin string `json:"rooms_amount_min"`
	RoomsAmountMax string `json:"rooms_amount_max"`
	SquareMin      string `json:"square_min"`
	SquareMax      string `json:"square_max"`
	CostMin        string `json:"cost_min"`
	CostMax        string `json:"cost_max"`
}
