package Sync1C

type TypeProject struct {
	ProjectId   string `json:"project_id"`
	ProjectName string `json:"project_name"`
}

type TypeHouse struct {
	HouseId   string `json:"house_id"`
	HouseName string `json:"house_name"`
}

type TypeBuilding struct {
	BuildingId   string `json:"building_id"`
	BuildingName string `json:"building_name"`
}

type TTypeBuilding struct {
	BuildingId   string `json:"building_id"`
	BuildingName string `json:"building_name"`
	ProjectCode  string
}

type TypeApartment struct {
	ApartmentId   string `json:"apartment_id"`
	Floor         int    `json:"floor"`
	Number        string `json:"number"`
	NumberOld     string `json:"number_old"`
	NumberForSort string `json:"number_for_sort"`
	Type          string `json:"type"`
	TypeAlias     string `json:"type_alias"`
	Status        string `json:"status"`
	StatusText    string `json:"status_text"`
	StatusColor   struct {
		Red   uint16 `json:"red"`
		Green uint16 `json:"green"`
		Blue  uint16 `json:"blue"`
		Web   string `json:"web"`
	} `json:"status_color"`
	RoomsAmount uint8 `json:"rooms_amount"`
	Tags        []struct {
		TagUid  string `json:"tag_uid"`
		TagName string `json:"tag_name"`
		TagType string `json:"tag_type"`
	} `json:"tags"`
	TotalSquare       float32 `json:"total_square"`
	LivingSquare      float32 `json:"living_square"`
	BltSquare         float32 `json:"blt_square"`
	PriceKvM          float32 `json:"price_kv_m"`
	PriceTotal        float32 `json:"price_total"`
	SalerInn          string  `json:"saler_inn"`
	SalerName         string  `json:"saler_name"`
	DateSale          string  `json:"date_sale"`
	DogovorStatusText string  `json:"dogovor_status_text"`
	Pokupatel         string  `json:"pokupatel"`
	DogovorNumber     string  `json:"dogovor_number"`
	FlatPlanImg       string  `json:"flat_plan_img"`
	FloorPlanImg      string  `json:"floor_plan_img"`
}

type TTypeApartment struct {
	ApartmentId   string `json:"apartment_id"`
	Floor         int    `json:"floor"`
	Number        string `json:"number"`
	NumberOld     string `json:"number_old"`
	NumberForSort string `json:"number_for_sort"`
	Type          string `json:"type"`
	TypeAlias     string `json:"type_alias"`
	Status        string `json:"status"`
	StatusText    string `json:"status_text"`
	StatusColor   struct {
		Red   uint16 `json:"red"`
		Green uint16 `json:"green"`
		Blue  uint16 `json:"blue"`
		Web   string `json:"web"`
	} `json:"status_color"`
	RoomsAmount uint8 `json:"rooms_amount"`
	Tags        []struct {
		TagUid  string `json:"tag_uid"`
		TagName string `json:"tag_name"`
		TagType string `json:"tag_type"`
	} `json:"tags"`
	TotalSquare       float32 `json:"total_square"`
	LivingSquare      float32 `json:"living_square"`
	BltSquare         float32 `json:"blt_square"`
	PriceKvM          float32 `json:"price_kv_m"`
	PriceTotal        float32 `json:"price_total"`
	SalerInn          string  `json:"saler_inn"`
	SalerName         string  `json:"saler_name"`
	DateSale          string  `json:"date_sale"`
	DogovorStatusText string  `json:"dogovor_status_text"`
	Pokupatel         string  `json:"pokupatel"`
	DogovorNumber     string  `json:"dogovor_number"`
	FlatPlanImg       string  `json:"flat_plan_img"`
	FloorPlanImg      string  `json:"floor_plan_img"`
	BuildingCode      string
}

type Data struct {
	ActDate  string `json:"act_date"`
	Projects []struct {
		TypeProject
		Houses []struct {
			TypeHouse
			Buildings []struct {
				TypeBuilding
				Apartments []TypeApartment `json:"apartments"`
			} `json:"buildings"`
		} `json:"houses"`
	} `json:"data"`
}
