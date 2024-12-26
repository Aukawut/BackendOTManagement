package handler

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

type ApproverList struct {
	REQUEST_NO string
	REV        int
	MAIL       string
	FULLNAME   string
}
type MailReturn struct {
	Email    string
	FULLNAME string
}

func CheckSendEmailRequestor(rev int, requestNo string) model.MailReturnRequestor {

	var mailResponse []model.MailReturnRequestor
	var mailResult model.MailReturnRequestor

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
		mailResult.MAIL = "N/A"
		return mailResult
	}
	// 4 Request -> Cancel

	step, errStep := db.Query(`SELECT h.REQUEST_NO,REV,h.CREATED_BY as REQUESTOR,ISNULL(hr.AD_Mail,'N/A') as MAIL,CONCAT(hr.UHR_FirstName_th,' ',hr.UHR_LastName_th) as FULLNAME FROM TBL_REQUESTS_HISTORY h
LEFT JOIN  V_AllUserPSTH hr ON h.CREATED_BY COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS
LEFT JOIN TBL_REQUESTS r ON r.REQUEST_NO = h.REQUEST_NO
WHERE h.REQUEST_NO = @req AND REV = @rev AND r.REQ_STATUS <> 4`,
		sql.Named("req", requestNo),
		sql.Named("rev", rev))

	if errStep != nil {
		fmt.Print("Query error : ", errStep.Error())
		mailResult.MAIL = "N/A"
		return mailResult
	}

	defer step.Close()

	for step.Next() {
		var mail model.MailReturnRequestor

		errScan := step.Scan(&mail.REQUEST_NO, &mail.REV, &mail.REQUESTOR, &mail.MAIL, &mail.FULLNAME)

		if errScan != nil {
			fmt.Println(errScan.Error())

		} else {
			mailResponse = append(mailResponse, mail)
		}

	}
	if len(mailResponse) > 0 {
		mailResult.REQUEST_NO = mailResponse[0].REQUEST_NO
		mailResult.REV = mailResponse[0].REV
		mailResult.REQUESTOR = mailResponse[0].REQUESTOR
		mailResult.MAIL = mailResponse[0].MAIL
		mailResult.FULLNAME = mailResponse[0].FULLNAME
	} else {
		mailResult.MAIL = "N/A"
		fmt.Println("Email Requestor isn't Found.")
	}

	return mailResult

}

func CheckSendEmail(rev int, requestNo string) MailReturn {

	var lastedApproved []ApproverList
	var mailResponse MailReturn

	connString := config.LoadDatabaseConfig()

	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
		mailResponse.Email = "N/A"
		return mailResponse
	}

	step, errStep := db.Query(`SELECT aa.REQUEST_NO,aa.REV,ISNULL(hr.AD_Mail,'N/A') as [AD_Mail],CONCAT(hr.UHR_FirstName_th,' ',hr.UHR_LastName_th) as FULLNAME FROM (
  SELECT h.REQUEST_NO,h.REV,a.CODE_APPROVER,h.ID_GROUP_DEPT,h.ID_FACTORY, al.ID_STATUS_APV FROM TBL_APPROVAL  al
  LEFT JOIN TBL_REQUESTS_HISTORY h ON h.REQUEST_NO = al.REQUEST_NO AND h.REV = al.REV
  LEFT JOIN TBL_APPROVERS a ON a.ID_GROUP_DEPT = h.ID_GROUP_DEPT AND a.ID_FACTORY = h.ID_FACTORY AND al.STEP = a.STEP
  WHERE h.REQUEST_NO = @req AND h.REV = @rev  AND ID_STATUS_APV = 1 ) aa
  LEFT JOIN V_AllUserPSTH hr ON aa.CODE_APPROVER COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS
  LEFT JOIN TBL_REQUESTS r ON aa.REQUEST_NO = r.REQUEST_NO WHERE r.REQ_STATUS = 1`, sql.Named("req", requestNo), sql.Named("rev", rev))

	if errStep != nil {
		fmt.Print("Query error : ", errStep.Error())
		mailResponse.Email = "N/A"
		return mailResponse
	}

	defer step.Close()

	for step.Next() {
		var approver ApproverList

		errScan := step.Scan(&approver.REQUEST_NO, &approver.REV, &approver.MAIL, &approver.FULLNAME)

		if errScan != nil {
			fmt.Println(errScan.Error())

		} else {
			lastedApproved = append(lastedApproved, approver)
		}
	}

	if len(lastedApproved) > 0 {
		fmt.Println(lastedApproved)
		mailResponse.Email = lastedApproved[0].MAIL
		mailResponse.FULLNAME = lastedApproved[0].FULLNAME
		return mailResponse

	} else {
		mailResponse.Email = "N/A"
		return mailResponse
	}

}

func TestApp(c *fiber.Ctx) error {

	rev, _ := strconv.Atoi(c.Params("rev"))
	name, _ := os.Hostname()
	CheckSendEmail(rev, c.Params("requestNo"))

	return c.JSON(fiber.Map{"err": false, "msg": "Hello", "container": name})
}
