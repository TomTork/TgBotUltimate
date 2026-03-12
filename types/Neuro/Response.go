package Neuro

import (
	"bytes"
	"encoding/json"
)

type Value string

func (v *Value) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if bytes.Equal(trimmed, []byte("null")) || bytes.Equal(trimmed, []byte("<UNK>")) {
		*v = ""
		return nil
	}

	var str string
	if err := json.Unmarshal(trimmed, &str); err == nil {
		*v = Value(str)
		return nil
	}

	var number json.Number
	if err := json.Unmarshal(trimmed, &number); err == nil {
		*v = Value(number.String())
		return nil
	}

	return nil
}

type Response struct {
	ProjectName    Value `json:"project_name"`
	BuildingLiter  Value `json:"building_liter"`
	FloorMin       Value `json:"floor_min"`
	FloorMax       Value `json:"floor_max"`
	RoomsAmountMin Value `json:"rooms_amount_min"`
	RoomsAmountMax Value `json:"rooms_amount_max"`
	SquareMin      Value `json:"square_min"`
	SquareMax      Value `json:"square_max"`
	CostMin        Value `json:"cost_min"`
	CostMax        Value `json:"cost_max"`
}
