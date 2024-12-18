package model

type UserRequest struct {
	EmpCode string `json:"empCode"`
}

type RequestOvertimeBody struct {
	OvertimeDateStart string        `json:"overtimeDateStart"`
	OvertimeDateEnd   string        `json:"overtimeDateEnd"`
	OvertimeType      int           `json:"overtimeType"`
	GroupDept         int           `json:"group"`
	Factory           int           `json:"factory"`
	Remark            string        `json:"remark"`
	ActionBy          string        `json:"actionBy"`
	Users             []UserRequest `json:"users"`
	GroupWorkCell     int           `json:"groupworkcell"`
	WorkCell          int           `json:"workcell"`
}

type RewriteRequestOvertimeBody struct {
	OvertimeDate  string        `json:"overtimeDate"`
	OvertimeType  int           `json:"overtimeType"`
	GroupDept     int           `json:"group"`
	Factory       int           `json:"factory"`
	Remark        string        `json:"remark"`
	ActionBy      string        `json:"actionBy"`
	Start         string        `json:"start"`
	End           string        `json:"end"`
	Users         []UserRequest `json:"users"`
	GroupWorkCell int           `json:"groupworkcell"`
	WorkCell      int           `json:"workcell"`
	RequestNo     string        `json:"requestNo"`
}

type CountRequest struct {
	AMOUNT      int
	NAME_STATUS string
}

type ResultCheckApproved struct {
	REQUEST_NO     string
	APPROVED_COUNT int
	REQ_STATUS     int
	STEP           int
	CREATED_BY     string
}

type ResultRequestByUser struct {
	REQUEST_NO  string
	REQ_STATUS  int
	REV         int
	NAME_STATUS string
}

type ApproverPendingAll struct {
	REQUEST_NO    string
	REV           int
	NEXT_APPROVER int
	ID_GROUP_DEPT int
	ID_FACTORY    int
	NAME_GROUP    string
	FACTORY_NAME  string
	CODE_APPROVER interface{}
	NAME_APPROVER string
}

type ResultCountApproveByEmpCode struct {
	NAME_STATUS   string
	ID_STATUS_APV int
	AMOUNT        int
}

type ResultListRequestByEmpIdAndStatus struct {
	REQUEST_NO    string
	CODE_APPROVER interface{}
	REV           int
	FACTORY_NAME  string
	NAME_GROUP    string
	ID_FACTORY    int
	ID_GROUP_DEPT int
	COUNT_USER    int
	DURATION      int
	HOURS_AMOUNT  float64
	SUM_MINUTE    int
	MINUTE_TOTAL  int
}

type SummaryRequestLastRev struct {
	REQUEST_NO    string
	REV           int
	NAME_STATUS   string
	FULLNAME      string
	OT_TYPE       float64
	FACTORY_NAME  string
	NAME_WORKGRP  string
	NAME_WORKCELL string
	USERS_AMOUNT  int
	SUM_MINUTE    int
	START_DATE    string
	END_DATE      string
	ID_FACTORY    int
	SUM_PLAN      float64
	SUM_PLAN_OB   float64
	ID_WORK_CELL  int
}

type RequestCommentApprover struct {
	REQUEST_NO    string
	ID_STATUS_APV interface{}
	CODE_APPROVER interface{}
	CREATED_AT    interface{}
	UPDATED_AT    interface{}
	NAME_STATUS   interface{}
	REMARK        interface{}
	DEPARTMENT    interface{}
	POSITION      interface{}
	FULLNAME      interface{}
}

type BodyApproveRequest struct {
	Status   int    `json:"status"`
	ActionBy string `json:"actionBy"`
	Remark   string `json:"remark"`
}

type ResponseApproverStepByReq struct {
	APPROVER string
	STEP     int
}
