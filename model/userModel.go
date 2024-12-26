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
	FACTORY_NAME  string
	MAIL          string
	ID_FACTORY    int
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
	AMOUNT_REQ interface{}
	YEAR_RQ    interface{}
	MONTH_RQ   interface{}
}

type OptionMenuByYear struct {
	AMOUNT_REQ interface{}
	YEAR_RQ    interface{}
}

type OptionMenuMonth struct {
	AMOUNT_REQ interface{}
	MONTH_RQ   interface{}
}

type ListUserByRequest struct {
	EMPLOYEE_CODE string
	FULLNAME      string
	POSITION      interface{}
}

type UserType struct {
	ID_UTYPE   int
	NAME_UTYPE string
	CREATED_AT interface{}
	UPDATED_AT interface{}
}

type UserGroup struct {
	ID_UGROUP   int
	NAME_UGROUP string
}

type EmployeeCode struct {
	EMPLOYEE_CODE string
}
type Employee struct {
	EMPLOYEE_CODE string
	ID_FACTORY    int
	GROUP_ID      int
	PREFIX        string
	FNAME_TH      string
	LNAME_TH      string
	FNAME_EN      string
	LNAME_EN      string
	ID_ROLE       int
	TYPE_ID       int
	ID_UGROUP     int
	FACTORY_NAME  string
	NAME_ROLE     string
	NAME_GROUP    string
	NAME_UGROUP   string
	NAME_UTYPE    string
	CREATED_AT    interface{}
	UPDATED_AT    interface{}
	UPDATED_BY    interface{}
	CREATED_BY    interface{}
}

type BodyEmployee struct {
	EmployeeCode string `json:"code"`
	Factory      int    `json:"factory"`
	Group        int    `json:"group"`
	Prefix       string `json:"prefix"`
	FnameTH      string `json:"fnameTH"`
	LnameTH      string `json:"lnameTH"`
	FnameEN      string `json:"fnameEN"`
	LnameEN      string `json:"lnameEN"`
	Role         int    `json:"role"`
	Ugroup       int    `json:"ugroup"`
	Type         int    `json:"type"`
	ActionBy     string `json:"actionBy"`
}

type BodyUpdateEmployee struct {
	Factory  int    `json:"factory"`
	Group    int    `json:"group"`
	Role     int    `json:"role"`
	Ugroup   int    `json:"ugroup"`
	Type     int    `json:"type"`
	ActionBy string `json:"actionBy"`
}

type ResultEmployeeByCode struct {
	Prefix       string
	EmployeeCode string
	FnameTH      string
	LnameTH      string
	FnameEN      string
	LnameEN      string
}

type UserBodyMail struct {
	EMPLOYEE_CODE string
	FULLNAME      string
	DEPARTMENT    interface{}
}

type ApproverCheck struct {
	CODE_APPROVER string
	ID_FACTORY    int
	ID_GROUP_DEPT int
	FACTORY_NAME  string
	NAME_GROUP    string
}
