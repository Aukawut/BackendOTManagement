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
