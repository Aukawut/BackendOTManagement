package model

type MainPlan struct {
	FactoryID int `json:"factory"`
	Month     int `json:"month"`
	Year      int `json:"year"`
	Hours     int `json:"hours"`
}

type ResultMainPlan struct {
	ID_PLAN      int
	ID_FACTORY   int
	FACTORY_NAME string
	CREATED_AT   string
	MONTH        int
	YEAR         int
	HOURS        float64
	UPDATED_AT   interface{}
	FNAME        string
}

type BodyGetMainPlan struct {
	Year    int `JSON:"year"`
	Factory int `JSON:"factory"`
	Month   int `JSON:"month"`
}
