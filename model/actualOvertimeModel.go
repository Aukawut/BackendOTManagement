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
	EMPLOYEE_CODE interface{}
	SCAN_IN       interface{}
	SCAN_OUT      interface{}
	OT_DATE       interface{}
	SHIFT         interface{}
	OT1_HOURS     interface{}
	OT15_HOURS    interface{}
	OT2_HOURS     interface{}
	OT3_HOURS     interface{}
	TOTAL_HOURS   interface{}
	UPDATED_AT    interface{}
	CREATED_BY    interface{}
	UPDATED_BY    interface{}
	FACTORY_NAME  interface{}
	NAME_UGROUP   interface{}
	NAME_UTYPE    interface{}
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
	EMPLOYEE_CODE  interface{}
	OT_DATE        interface{}
	SCAN_IN        interface{}
	SCAN_OUT       interface{}
	HOURS          float64
	FACTORY_NAME   interface{}
	NAME_UGROUP    interface{}
	UHR_Department interface{}
	HOURS_AMOUNT   interface{}
	NAME_UTYPE     interface{}
	ID_FACTORY     interface{}
	ID_UTYPE       interface{}
	ID_UGROUP      interface{}
	ID_TYPE_OT     interface{}
}

type OvertimeActualByFac struct {
	FACTORY_NAME  string
	ID_FACTORY    int
	INLINE_HOURS  float64
	OFFLINE_HOURS float64
}

type OvertimeActualByType struct {
	ID_TYPE_OT   int
	HOURS_AMOUNT string
	SUM_HOURS    float64
}
type OvertimeActualByDate struct {
	INLINE_HOURS  float64
	OFFLINE_HOURS float64
	DATE_OT       string
}

type CalActualByFac struct {
	SUM_HOURS float64
}
type CalActualByWorkcell struct {
	SUM_HOURS   float64
	WORKCELL_ID int
}

type RequestList struct {
	REQUEST_NO   string
	REV          string
	FACTORY_NAME string
	ID_TYPE_OT   int
	HOURS_AMOUNT string
	PERSON       int
	DURATION     float64
	HOURS_TOTAL  float64
}

type OldRequestDetail struct {
	REQUEST_NO    string
	REV           string
	ID_FACTORY    int
	ID_GROUP_DEPT int
	ID_WORK_CELL  string
	START_DATE    string
	END_DATE      string
	ID_TYPE_OT    int
	ID_WORKGRP    int
	REMARK        interface{}
}

type ActualCompareWorkgroup struct {
	OT_DATE       interface{}
	EMPLOYEE_CODE interface{}
	NAME_WORKGRP  interface{}
	NAME_UGROUP   interface{}
	NAME_WORKCELL interface{}
	HOURS         interface{}
	HOURS_AMOUNT  interface{}
}

type ActualGroupByWorkcell struct {
	SUM_HOURS     float64
	NAME_WORKCELL string
}
type ActualGroupByWorkgroup struct {
	SUM_HOURS    float64
	NAME_WORKGRP string
}
