package model

type ListActual struct {
	Date         string  `JSON:"date"`
	EmployeeCode string  `JSON:"employeeCode"`
	Start        string  `JSON:"start"`
	End          string  `JSON:"end"`
	Overtime1    float64 `JSON:"overtime1"`
	Overtime15   float64 `JSON:"overtime15"`
	Overtime2    float64 `JSON:"overtime2"`
	Overtime3    float64 `JSON:"overtime3"`
	Shift        string  `JSON:"shift"`
	Total        float64 `JSON:"total"`
}

type BodyRequestSaveActual struct {
	ActionBy string       `JSON:"action"`
	Overtime []ListActual `JSON:"overtime"`
}

type AllActualOvertime struct {
	Id            int
	EMPLOYEE_CODE string
	SCAN_IN       string
	SCAN_OUT      string
	OT_DATE       string
	SHIFT         string
	OT1_HOURS     float64
	OT15_HOURS    float64
	OT2_HOURS     float64
	OT3_HOURS     float64
	TOTAL_HOURS   float64
	UPDATED_AT    interface{}
	CREATED_BY    interface{}
	UPDATED_BY    interface{}
	FACTORY_NAME  string
	NAME_UGROUP   string
	NAME_UTYPE    string
}

type SummaryActualComparePlan struct {
	MONTH_NO      int
	MONTH_NAME    string
	MONTH         int
	SUM_OT_ACTUAL float64
	SUM_OT_PLANWC float64
	SUM_OT_PLANOB float64
}

type SummaryActualByFactory struct {
	ID_FACTORY   int
	FACTORY_NAME string
	SUM_ACTUAL   float64
	SUM_PLAN     float64
	SUM_PLAN_OB  float64
}
type CountActualOvertime struct {
	COUNT_OT int
}

type SummaryActualByDuration struct {
	DATE_OT   string
	SUM_TOTAL float64
	DAY_OT    int
	COUNT_OT  int
}

type OvertimeActual struct {
	EMPLOYEE_CODE  string
	OT_DATE        string
	SCAN_IN        string
	SCAN_OUT       string
	HOURS          float64
	FACTORY_NAME   string
	NAME_UGROUP    string
	UHR_Department string
	HOURS_AMOUNT   string
	NAME_UTYPE     string
	ID_FACTORY     interface{}
	ID_UTYPE       interface{}
	ID_UGROUP      interface{}
	ID_TYPE_OT     interface{}
}
type CalActualByFac struct {
	SUN_HOURS float64
}
