package model

type RoleEncepyt struct {
	EMPLOYEE_CODE string
	NAME_ROLE     string
	FACTORY_NAME  string
	NAME_GROUP    string
	FULLNAME      string
	MAIL          string
	REAL_DEPT     string
	ID_FACTORY    int
	ID_GROUP_DEPT int
}

type UserEncepyt struct {
	EmployeeCode string
	Department   string
	Fullname     string
	Email        string
	Role         []RoleEncepyt
}

type UsersInfo struct {
	EMPLOYEE_CODE string
	FULLNAME      string
	FACTORY_NAME  string
	NAME_ROLE     string
	NAME_GROUP    string
	POSITION      string
}

type Approver struct {
	CODE_APPROVER string
	NAME_APPROVER string
	ID_GROUP_DEPT int
	NAME_GROUP    string
	ROLE          int
	UHR_Position  string
	STEP          int
	ID_APPROVER   int
}

type ApproverRequest struct {
	EmpCode   string `json:"empCode"`
	Name      string `json:"name"`
	GroupID   int    `json:"groupId"`
	RoleID    int    `json:"roleId"`
	Step      int    `json:"step"`
	CreatedBy string `json:"createBy"`
	Factory   int    `json:"factory"`
}

type ApproverRequestUpdate struct {
	EmpCode  string `json:"empCode"`
	Name     string `json:"name"`
	GroupID  int    `json:"groupId"`
	RoleID   int    `json:"roleId"`
	Step     int    `json:"step"`
	UpdateBy string `json:"updateBy"`
	Active   string `json:"active"`
}

type ResultCountReqByYearMonth struct {
	AMOUNT_REQ int
	YEAR_RQ    int
	MONTH_RQ   int
}

type OptionMenuByYear struct {
	AMOUNT_REQ int
	YEAR_RQ    int
}

type OptionMenuMonth struct {
	AMOUNT_REQ int
	MONTH_RQ   int
}

type ListUserByRequest struct {
	EMPLOYEE_CODE string
	FULLNAME      string
	POSITION      interface{}
}
