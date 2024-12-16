package model

type Permission struct {
	ID_PERMISSION int
	EMPLOYEE_CODE string
	ID_FACTORY    int
	ID_ROLE       int
	ID_GROUP_DEPT int
	CREATED_AT    interface{}
	UPDATED_AT    interface{}
}

type BodyPermission struct {
	EmployeeCode string `json:"code"`
	Factory      int    `json:"factory"`
	Role         int    `json:"role"`
	GroupDept    int    `json:"groupDept"`
}
