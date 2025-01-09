package model

type MainPlan struct {
	WorkcellID int    `json:"workcell"`
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	Hours      int    `json:"hours"`
	ActionBy   string `json:"action"`
	UserGroup  int    `json:"userGroup"`
}

type BodyPlanOB struct {
	Factory   int    `json:"factory"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
	Hours     int    `json:"hours"`
	ActionBy  string `json:"action"`
	UserGroup int    `json:"userGroup"`
}

type ResultMainPlan struct {
	ID_PLAN       int
	ID_FACTORY    int
	ID_WORK_CELL  int
	NAME_WORKCELL string
	FACTORY_NAME  string
	NAME_UGROUP   string
	ID_UGROUP     int
	CREATED_AT    string
	MONTH         int
	YEAR          int
	HOURS         float64
	UPDATED_AT    interface{}
	FNAME         string
}

type ResultPlanOB struct {
	ID_OB_PLAN   int
	ID_FACTORY   int
	FACTORY_NAME string
	NAME_UGROUP  string
	ID_UGROUP    int
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
type ResultPlanByFactory struct {
	YEAR       int
	MONTH      int
	ID_FACTORY int
	SUM_HOURS  float64
}

type PlanByWorkcell struct {
	REQUEST_NO   string
	SUM_HOURS    float64
	ID_WORK_CELL int
}
