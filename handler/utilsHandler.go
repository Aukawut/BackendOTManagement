package handler

import (
	"database/sql"
	"fmt"
	"strconv"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"github.com/gofiber/fiber/v2"
)

type ApproverList struct {
	REQUEST_NO string
	REV        int
	MAIL       string
}
type MailReturn struct {
	Email string
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

	step, errStep := db.Query(`SELECT aa.REQUEST_NO,aa.REV,ISNULL(hr.AD_Mail,'N/A') as [AD_Mail] FROM (
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

		errScan := step.Scan(&approver.REQUEST_NO, &approver.REV, &approver.MAIL)

		if errScan != nil {
			fmt.Println(errScan.Error())

		} else {
			lastedApproved = append(lastedApproved, approver)
		}
	}

	if len(lastedApproved) > 0 {
		fmt.Println(lastedApproved)
		mailResponse.Email = lastedApproved[0].MAIL
		return mailResponse

	} else {
		mailResponse.Email = "N/A"
		return mailResponse
	}

}

func TestApp(c *fiber.Ctx) error {

	rev, _ := strconv.Atoi(c.Params("rev"))

	CheckSendEmail(rev, c.Params("requestNo"))

	return c.JSON(fiber.Map{"err": false, "msg": "123"})
}
