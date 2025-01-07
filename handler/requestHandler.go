package handler

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"strconv"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetRunningNoRequest() string {

	today := time.Now()

	prefix := "RQ" + today.Format("20060102") // RQ20241202

	strConfig := config.LoadDatabaseConfig()

	var runningNo string

	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	// Execute SELECT query
	rows, err := db.Query(`SELECT TOP 1 RIGHT([REQUEST_NO],4) as [REQUEST_NO] FROM TBL_REQUESTS WHERE LEFT([REQUEST_NO],10) = @prefix ORDER BY [REQUEST_NO] DESC`, sql.Named("prefix", prefix))
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var transectionNo string
		err := rows.Scan(&transectionNo)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			runningNo = ""
		} else {
			runningNo = transectionNo
		}

	}

	// มี Transection ของวันนี้
	if runningNo == "" {
		query := fmt.Sprintf(`SELECT '%s' + FORMAT(1, '0000') AS [REQUEST_NO]`, prefix)

		results, errQuery := db.Query(query)

		if errQuery != nil {

			log.Fatalf("Query execution failed: %v", err)
		}

		defer results.Close()

		for results.Next() {
			var transec string
			if errScan := results.Scan(&transec); errScan != nil {
				log.Fatalf("Failed to scan row: %v", errScan)
			}
			runningNo = transec
		}

		if errResult := results.Err(); errResult != nil {
			log.Fatalf("Row iteration error: %v", errResult)
			runningNo = ""
		}

	} else {
		// ไม่มี Transection ของวันนี้
		lastTrans, err := strconv.Atoi(runningNo) // แปลง String เป็น Int

		if err != nil {

			log.Fatalf("Row iteration error: %v", err)
		}

		nextTrans := lastTrans + 1

		query := fmt.Sprintf(`SELECT '%s' + FORMAT(%d, '0000') AS [REQUEST_NO]`, prefix, nextTrans)

		results, errQuery := db.Query(query)

		if errQuery != nil {

			log.Fatalf("Query execution failed: %v", err)
		}

		defer results.Close()

		for results.Next() {
			var transec string
			if errScan := results.Scan(&transec); errScan != nil {
				log.Fatalf("Failed to scan row: %v", errScan)
			}
			runningNo = transec
		}

		if errResult := results.Err(); errResult != nil {
			log.Fatalf("Row iteration error: %v", errResult)
			runningNo = ""
		}

	}

	return runningNo

}

func RequestOvertime(c *fiber.Ctx) error {
	var req model.RequestOvertimeBody

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
	}

	users := req.Users

	if len(users) > 0 {
		// Get Transection Request No
		running := GetRunningNoRequest()

		strConfig := config.LoadDatabaseConfig()

		db, err := sql.Open("sqlserver", strConfig)
		if err != nil {
			fmt.Println("Error creating connection: " + err.Error())
		}
		defer db.Close()

		// Test connection
		err = db.Ping()
		if err != nil {
			fmt.Println("Error connecting to the database: " + err.Error())
		}

		// <-- Execute Insert Request -->
		_, errInsert := db.Exec(`INSERT INTO [dbo].[TBL_REQUESTS] ([REQUEST_NO],[ID_GROUP_DEPT],[REMARK],[REQ_STATUS],[ID_TYPE_OT],[ID_FACTORY],[CREATED_AT],[CREATED_BY],[START_DATE],[END_DATE],[ID_WORKGRP],[ID_WORK_CELL]) VALUES 
			(@reqNo,@group,@remark,@status,@type,@factory,GETDATE(),@actionBy,@dateStart,@dateEnd,@groupWorkcell,@workcell)`,
			sql.Named("reqNo", running),
			sql.Named("group", req.GroupDept),
			sql.Named("dateStart", req.OvertimeDateStart),
			sql.Named("dateEnd", req.OvertimeDateEnd),
			sql.Named("remark", req.Remark),
			sql.Named("status", 1), // Pending
			sql.Named("type", req.OvertimeType),
			sql.Named("factory", req.Factory),
			sql.Named("actionBy", req.ActionBy),
			sql.Named("groupWorkcell", req.GroupWorkCell),
			sql.Named("workcell", req.WorkCell),
		)

		if errInsert != nil {
			fmt.Println("Query failed: " + errInsert.Error())
		}

		_, errorInsertHistory := db.Exec(`INSERT INTO [dbo].[TBL_REQUESTS_HISTORY] ([REQUEST_NO],[ID_GROUP_DEPT],[REMARK],[REQ_STATUS],[ID_TYPE_OT],[ID_FACTORY],[CREATED_AT],[CREATED_BY],[REV],[START_DATE],[END_DATE],[ID_WORKGRP],[ID_WORK_CELL]) 
		VALUES (@reqNo,@group,@remark,@status,@type,@factory,GETDATE(),@actionBy,@rev,@dateStart,@dateEnd,@groupWorkcell,@workcell)`,
			sql.Named("reqNo", running),
			sql.Named("group", req.GroupDept),
			sql.Named("dateStart", req.OvertimeDateStart),
			sql.Named("dateEnd", req.OvertimeDateEnd),
			sql.Named("remark", req.Remark),
			sql.Named("status", 1), // Pending
			sql.Named("type", req.OvertimeType),
			sql.Named("factory", req.Factory),
			sql.Named("actionBy", req.ActionBy),
			sql.Named("rev", 1), // Rev : 1
			sql.Named("groupWorkcell", req.GroupWorkCell),
			sql.Named("workcell", req.WorkCell),
		)
		if errorInsertHistory != nil {

			fmt.Println("Query failed: " + errorInsertHistory.Error())

		}

		results, errResult := db.Query(`SELECT STEP FROM TBL_OT_TYPE WHERE [ID_TYPE_OT] = @id`, sql.Named("id", req.OvertimeType))

		if errResult != nil {
			fmt.Println("Query failed: " + errResult.Error())
		}
		defer results.Close()

		var step int

		for results.Next() {
			if err := results.Scan(&step); err != nil {
				log.Fatalf("Failed to scan row: %v", err)
			}

		}

		if errResult := results.Err(); errResult != nil {
			log.Fatalf("Row iteration error: %v", errResult)
			step = 0
		}

		if step > 0 {
			// Insert Step Approval
			for i := 0; i < step; i++ {

				_, errorInsertApproval := db.Exec(`INSERT INTO [dbo].[TBL_APPROVAL] 
			([REQUEST_NO],[STEP],[REV],[CREATED_AT]) VALUES (@reqNo,@step,@rev,GETDATE())`,
					sql.Named("reqNo", running),
					sql.Named("step", i+1),
					sql.Named("rev", 1),
				)
				if errorInsertApproval != nil {
					fmt.Println("Query failed: " + errResult.Error())
				}

			}
			_, errorUpdate := db.Exec(`UPDATE [dbo].[TBL_APPROVAL] SET [ID_STATUS_APV] = 1 
			WHERE REQUEST_NO = @reqNo AND [STEP] = 1 AND [REV] = 1`, sql.Named("reqNo", running))
			if errorUpdate != nil {
				fmt.Println("Query failed: " + errResult.Error())
			}
		}

		for j := 0; j < len(users); j++ {
			_, errorInsertUser := db.Exec(`INSERT INTO [dbo].[TBL_USERS_REQ] 
			([EMPLOYEE_CODE],[REQUEST_NO],[CREATED_AT],[CREATED_BY],[REV]) VALUES (@code,@reqNo,GETDATE(),@actionBy,@rev)`,
				sql.Named("code", users[j].EmpCode),
				sql.Named("reqNo", running),
				sql.Named("actionBy", req.ActionBy),
				sql.Named("rev", 1),
			)
			if errorInsertUser != nil {
				fmt.Println("Query failed: " + errResult.Error())
			}
		}

		mail := CheckSendEmail(1, running)

		if mail.Email != "N/A" {

			SendEMailToApprover(running, 1, mail.Email, mail.FULLNAME)
		} else {

			fmt.Println("E-mail address isn't found.")
		}

		return c.JSON(fiber.Map{
			"err":    false,
			"msg":    "Requested successfully!",
			"status": "Ok",
		})

	} else {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Users is required!!",
		})
	}

}

func CountRequestByEmpCode(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()
	var code = c.Params("code")
	var requestCount []model.CountRequest

	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	results, err := db.Query(`SELECT SUM(CASE WHEN REQUEST_NO IS NOT NULL THEN 1 ELSE 0 END) as AMOUNT,NAME_STATUS FROM [dbo].[TBL_REQ_STATUS] s LEFT JOIN (SELECT * FROM TBL_REQUESTS WHERE CREATED_BY = @code) r ON s.ID_STATUS = r.REQ_STATUS 
							     GROUP BY NAME_STATUS`, sql.Named("code", code))

	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer results.Close()

	for results.Next() {
		var countReq model.CountRequest

		errScan := results.Scan(&countReq.AMOUNT, &countReq.NAME_STATUS)

		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			requestCount = append(requestCount, countReq)
		}

	}

	if len(requestCount) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": requestCount,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": requestCount,
		})
	}

}

func CancelRequestByReqNo(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()
	var requestNo = c.Params("requestNo")
	var rev = c.Params("rev")
	fmt.Println("testing..")
	user := c.Locals("user").(jwt.MapClaims)

	var resultsCheckApproved []model.ResultCheckApproved

	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	results, err := db.Query(`SELECT REQUEST_NO,SUM( CASE WHEN CODE_APPROVER IS NOT NULL AND REQ_STATUS NOT IN(2,3,4,5) THEN  1 ELSE 0 END) as APPROVED_COUNT,
REQ_STATUS,STEP FROM (
	SELECT a.REQUEST_NO,r.CREATED_BY,REQ_STATUS,a.CODE_APPROVER,t.STEP FROM [dbo].[TBL_REQUESTS_HISTORY] r
	LEFT JOIN TBL_APPROVAL a ON r.REQUEST_NO = a.REQUEST_NO AND r.REV = a.REV
	LEFT JOIN TBL_OT_TYPE t ON r.ID_TYPE_OT = t.ID_TYPE_OT
	WHERE  r.REQUEST_NO = @requestNo AND r.REV = @rev )m GROUP BY m.REQUEST_NO,REQ_STATUS,STEP`, sql.Named("requestNo", requestNo), sql.Named("rev", rev))

	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer results.Close()

	for results.Next() {
		var approved model.ResultCheckApproved

		errScan := results.Scan(&approved.REQUEST_NO, &approved.APPROVED_COUNT, &approved.REQ_STATUS, &approved.STEP)

		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			resultsCheckApproved = append(resultsCheckApproved, approved)
		}

	}

	if len(resultsCheckApproved) > 0 {
		status := resultsCheckApproved[0].REQ_STATUS
		approvedCount := resultsCheckApproved[0].APPROVED_COUNT
		step := resultsCheckApproved[0].STEP

		if approvedCount > 0 && approvedCount < step {
			// return to clinet job approved by some approver.
			return c.JSON(fiber.Map{
				"err": true,
				"msg": "Request approved by some approver!",
			})
		} else if approvedCount > 0 && (approvedCount == step || approvedCount > step) {
			// return to clinet job approved by all approver.
			return c.JSON(fiber.Map{
				"err": true,
				"msg": "Request approved by all approver!",
			})
		} else if status != 1 && status == 2 {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": "Can't cancel Request Approved.",
			})

		} else if status != 1 && status == 3 {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": "Can't cancel Request No: Approved.",
			})

		} else if status != 1 && status == 4 {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": "Can't cancel Request No: Canceled.",
			})

		} else {

			// Update Status to Cencel
			stmtCancel := `UPDATE [TBL_REQUESTS] 
                SET REQ_STATUS = 4, UPDATED_AT = GETDATE(), UPDATED_BY = @actionBy 
                WHERE REQUEST_NO = @requestNo`

			stmtCancelHistory := `UPDATE [TBL_REQUESTS_HISTORY] 
                SET REQ_STATUS = 4, UPDATED_AT = GETDATE(), UPDATED_BY = @actionBy 
                WHERE REQUEST_NO = @requestNo AND [REV] = @rev`

			_, errUpdate := db.Exec(stmtCancel,
				sql.Named("requestNo", requestNo),
				sql.Named("actionBy", user["employee_code"]))

			_, errUpdateHistory := db.Exec(stmtCancelHistory,
				sql.Named("requestNo", requestNo),
				sql.Named("rev", rev),
				sql.Named("actionBy", user["employee_code"]),
			)

			if errUpdate != nil {
				return c.JSON(fiber.Map{
					"err": true,
					"msg": errUpdate.Error(), // Use errUpdate instead of err
				})
			}
			if errUpdateHistory != nil {
				return c.JSON(fiber.Map{
					"err": true,
					"msg": errUpdateHistory.Error(), // Use errUpdate instead of err
				})
			}

			// Return success message
			return c.JSON(fiber.Map{
				"err":    false,
				"msg":    "Request canceled successfully.",
				"status": "Ok",
			})

		}

	} else {

		return c.JSON(fiber.Map{
			"err": false,
			"msg": "Request isn't found!",
		})

	}

}

func ReviseRequestOvertime(c *fiber.Ctx) error {
	var req model.ReviseRequestOvertimeBody
	var lastRevNo int
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
	}

	users := req.Users

	if len(users) > 0 {

		strConfig := config.LoadDatabaseConfig()

		db, err := sql.Open("sqlserver", strConfig)
		if err != nil {
			fmt.Println("Error creating connection: " + err.Error())
		}
		defer db.Close()

		// Test connection
		err = db.Ping()
		if err != nil {
			fmt.Println("Error connecting to the database: " + err.Error())
		}

		lastRev, errLast := db.Query(`SELECT TOP 1 REV FROM TBL_REQUESTS_HISTORY WHERE REQUEST_NO = @reqNo GROUP BY REQUEST_NO,REV ORDER BY REV DESC`,
			sql.Named("reqNo", req.RequestNo))

		if errLast != nil {
			return c.JSON(fiber.Map{"err": true, "msg": "Last Rev. isn't found."})
		}

		for lastRev.Next() {
			var last int
			errScan := lastRev.Scan(&last)

			if errScan != nil {
				fmt.Println("Row scan failed: " + errScan.Error())

			} else {
				lastRevNo = last
			}

		}
		// Not Lasted Revise
		if lastRevNo != req.Rev {
			return c.JSON(fiber.Map{"err": true, "msg": "Rev. isn't corrected!"})
		}

		_, errorInsertHistory := db.Exec(`INSERT INTO [dbo].[TBL_REQUESTS_HISTORY] ([REQUEST_NO],[ID_GROUP_DEPT],
		[REMARK],[REQ_STATUS],[ID_TYPE_OT],[ID_FACTORY],[CREATED_AT],[CREATED_BY],[REV],[START_DATE],[END_DATE],[ID_WORKGRP],[ID_WORK_CELL])
		VALUES (@reqNo,@group,@remark,@status,@type,@factory,GETDATE(),@actionBy,@rev,@start,@end,@groupWorkcell,@workcell)`,
			sql.Named("reqNo", req.RequestNo),
			sql.Named("group", req.GroupDept),
			sql.Named("remark", req.Remark),
			sql.Named("status", 1), // Pending
			sql.Named("type", req.OvertimeType),
			sql.Named("factory", req.Factory),
			sql.Named("actionBy", req.ActionBy),
			sql.Named("rev", lastRevNo+1), // Rev + 1
			sql.Named("start", req.Start),
			sql.Named("end", req.End),
			sql.Named("groupWorkcell", req.GroupWorkCell),
			sql.Named("workcell", req.WorkCell),
		)
		if errorInsertHistory != nil {

			fmt.Println("Query failed: " + errorInsertHistory.Error())

		}

		results, errResult := db.Query(`SELECT STEP FROM TBL_OT_TYPE WHERE [ID_TYPE_OT] = @id`, sql.Named("id", req.OvertimeType))

		if errResult != nil {
			fmt.Println("Query failed: " + errResult.Error())
		}
		defer results.Close()

		var step int

		for results.Next() {
			if err := results.Scan(&step); err != nil {
				log.Fatalf("Failed to scan row: %v", err)
			}

		}

		if errResult := results.Err(); errResult != nil {
			log.Fatalf("Row iteration error: %v", errResult)
			step = 0
		}

		if step > 0 {
			// Insert Step Approval
			for i := 0; i < step; i++ {

				_, errorInsertApproval := db.Exec(`INSERT INTO [dbo].[TBL_APPROVAL]
			([REQUEST_NO],[STEP],[REV],[CREATED_AT]) VALUES (@reqNo,@step,@rev,GETDATE())`,
					sql.Named("reqNo", req.RequestNo),
					sql.Named("step", i+1),
					sql.Named("rev", lastRevNo+1),
				)
				if errorInsertApproval != nil {
					fmt.Println("Query failed: " + errResult.Error())
				}

			}
			_, errorUpdate := db.Exec(`UPDATE [dbo].[TBL_APPROVAL] SET [ID_STATUS_APV] = 1 
			WHERE REQUEST_NO = @reqNo AND [STEP] = 1 AND [REV] = @rev`, sql.Named("reqNo", req.RequestNo),
				sql.Named("rev", lastRevNo+1))
			if errorUpdate != nil {
				fmt.Println("Query failed: " + errResult.Error())
			}

		}

		for j := 0; j < len(users); j++ {
			fmt.Println(users[j].EmpCode)
			_, errorInsertUser := db.Exec(`INSERT INTO [dbo].[TBL_USERS_REQ]
			([EMPLOYEE_CODE],[REQUEST_NO],[CREATED_AT],[CREATED_BY],[REV]) VALUES (@code,@reqNo,GETDATE(),@actionBy,@rev)`,
				sql.Named("code", users[j].EmpCode),
				sql.Named("reqNo", req.RequestNo),
				sql.Named("actionBy", req.ActionBy),
				sql.Named("rev", lastRevNo+1),
			)
			if errorInsertUser != nil {
				fmt.Println("Query failed: " + errResult.Error())
			}
		}

		// Update main Request

		// <-- Execute Insert Request -->
		_, errInsert := db.Exec(`UPDATE [dbo].[TBL_REQUESTS] SET [ID_GROUP_DEPT] = @group,
		[REMARK] = @remark,[REQ_STATUS] = @status,
		[ID_TYPE_OT] = @type,[ID_FACTORY] = @factory,[UPDATED_AT] = GETDATE(),[UPDATED_BY] = @actionBy,
		[START_DATE] = @start,[END_DATE] = @end,[ID_WORKGRP] = @groupWorkcell,[ID_WORK_CELL] = @workcell WHERE [REQUEST_NO] = @reqNo`,
			sql.Named("group", req.GroupDept),
			sql.Named("remark", req.Remark),
			sql.Named("status", 1), // Pending
			sql.Named("type", req.OvertimeType),
			sql.Named("factory", req.Factory),
			sql.Named("actionBy", req.ActionBy),
			sql.Named("start", req.Start),
			sql.Named("end", req.End),
			sql.Named("groupWorkcell", req.GroupWorkCell),
			sql.Named("workcell", req.WorkCell),
			sql.Named("reqNo", req.RequestNo),
		)

		if errInsert != nil {
			fmt.Println("Query failed: " + errInsert.Error())
		}

		_, errorUpdate := db.Exec(`UPDATE [dbo].[TBL_APPROVAL] SET [ID_STATUS_APV] = 1 
			WHERE REQUEST_NO = @reqNo AND [STEP] = 1 AND [REV] = @rev`,
			sql.Named("reqNo", req.RequestNo),
			sql.Named("rev", lastRevNo+1))

		if errorUpdate != nil {
			fmt.Println("Query failed: " + errResult.Error())
		}

		mail := CheckSendEmail(lastRevNo+1, req.RequestNo)

		if mail.Email != "N/A" {

			SendEMailToApprover(req.RequestNo, lastRevNo+1, mail.Email, mail.FULLNAME)
		} else {

			fmt.Println("E-mail address isn't found.")
		}

		return c.JSON(fiber.Map{
			"err":    false,
			"msg":    "Revised successfully!",
			"status": "Ok",
		})

	} else {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Users is required!!",
		})
	}

}

func CountRequestByYear(c *fiber.Ctx) error {
	info := []model.ResultCountReqByYearMonth{}

	var year = c.Params("year")
	var code = c.Params("code")

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
	}

	queryUser := `SELECT COUNT(*) as AMOUNT_REQ,YEAR(m.DATE_OT) as YEAR_RQ,MONTH(m.DATE_OT) as MONTH_RQ FROM(
			SELECT a.*,b.EMPLOYEE_CODE,r.REQ_STATUS,CONVERT(DATE,r.START_DATE) as DATE_OT FROM (
			SELECT REQUEST_NO,MAX(REV) as REV FROM TBL_USERS_REQ GROUP BY REQUEST_NO ) a
			LEFT JOIN TBL_USERS_REQ b ON  a.REQUEST_NO = b.REQUEST_NO AND a.REV = b.REV 
			LEFT JOIN (
				SELECT * FROM TBL_REQUESTS_HISTORY WHERE 
					REV IN (SELECT MAX(REV) FROM TBL_REQUESTS_HISTORY GROUP BY REQUEST_NO)  AND
					CREATED_BY = @code AND REQ_STATUS = 2
			) r ON a.REQUEST_NO = r.REQUEST_NO ) m
			WHERE YEAR(m.DATE_OT) = @year
			GROUP BY YEAR(m.DATE_OT),MONTH(m.DATE_OT) ORDER BY MONTH(m.DATE_OT) DESC`

	results, errorQueryser := db.Query(queryUser, sql.Named("code", code), sql.Named("year", year)) //Query

	if errorQueryser != nil {
		fmt.Println("Query failed: " + errorQueryser.Error())
	}

	for results.Next() {
		var result model.ResultCountReqByYearMonth

		errScan := results.Scan(&result.AMOUNT_REQ, &result.YEAR_RQ, &result.MONTH_RQ) // Scan เก็บข้อมูลใน Struct
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}

func GetYearMenu(c *fiber.Ctx) error {
	info := []model.OptionMenuByYear{}

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
	}

	query := `SELECT COUNT(*) as AMOUNT_REQ,YEAR(m.DATE_OT) as YEAR_RQ FROM(
			SELECT a.*,b.EMPLOYEE_CODE,r.REQ_STATUS,CONVERT(DATE,r.START_DATE) as [DATE_OT]  FROM (
			SELECT REQUEST_NO,MAX(REV) as REV FROM TBL_USERS_REQ GROUP BY REQUEST_NO ) a
			LEFT JOIN TBL_USERS_REQ b ON  a.REQUEST_NO = b.REQUEST_NO AND a.REV = b.REV 
			LEFT JOIN TBL_REQUESTS r ON a.REQUEST_NO = r.REQUEST_NO ) m
			GROUP BY YEAR(m.DATE_OT) ORDER BY YEAR(m.DATE_OT) DESC`

	results, errorQueryser := db.Query(query) //Query

	if errorQueryser != nil {
		fmt.Println("Query failed: " + errorQueryser.Error())
	}

	for results.Next() {
		var result model.OptionMenuByYear

		errScan := results.Scan(&result.AMOUNT_REQ, &result.YEAR_RQ) // Scan เก็บข้อมูลใน Struct
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}

func GetMonthMenu(c *fiber.Ctx) error {
	info := []model.OptionMenuMonth{}
	year := c.Params("year")
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
	}

	query := `SELECT COUNT(*) as AMOUNT_REQ,MONTH(m.DATE_OT) as MONTH_RQ FROM(
			SELECT a.*,b.EMPLOYEE_CODE,r.REQ_STATUS,CONVERT(DATE,r.START_DATE) AS DATE_OT FROM (
			SELECT REQUEST_NO,MAX(REV) as REV FROM TBL_USERS_REQ GROUP BY REQUEST_NO ) a
			LEFT JOIN TBL_USERS_REQ b ON  a.REQUEST_NO = b.REQUEST_NO AND a.REV = b.REV 
			LEFT JOIN TBL_REQUESTS r ON a.REQUEST_NO = r.REQUEST_NO ) m
			WHERE YEAR(m.DATE_OT) = @year
			GROUP BY MONTH(m.DATE_OT),YEAR(m.DATE_OT) 
			ORDER BY MONTH(m.DATE_OT) DESC`

	results, errorQueryser := db.Query(query, sql.Named("year", year)) //Query

	if errorQueryser != nil {
		fmt.Println("Query failed: " + errorQueryser.Error())
	}

	for results.Next() {
		var result model.OptionMenuMonth

		errScan := results.Scan(&result.AMOUNT_REQ, &result.MONTH_RQ) // Scan เก็บข้อมูลใน Struct
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}

func GetRequestByEmpCode(c *fiber.Ctx) error {
	info := []model.ResultRequestByUser{}
	empCode := c.Params("code")
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
	}

	query := `SELECT h.REQUEST_NO,h.REV,s.NAME_STATUS,p.PERSON,START_DATE,END_DATE,
			CAST(DATEDIFF(MINUTE,START_DATE,END_DATE)/60  as decimal(18,2)) as DURATION,
			f.FACTORY_NAME,wc.NAME_WORKCELL
			FROM TBL_REQUESTS_HISTORY h
			LEFT JOIN TBL_REQ_STATUS s ON h.REQ_STATUS = s.ID_STATUS
			LEFT JOIN (
				SELECT COUNT(*) as PERSON, u.REQUEST_NO,u.REV 
			FROM TBL_USERS_REQ u GROUP BY u.REQUEST_NO,u.REV
			) p ON h.REQUEST_NO = p.REQUEST_NO AND h.REV = p.REV
			LEFT JOIN TBL_FACTORY f ON h.ID_FACTORY = f.ID_FACTORY
			LEFT JOIN TBL_WORKCELL wc ON h.ID_WORK_CELL = wc.ID_WORK_CELL
			WHERE h.CREATED_BY = @code
			`

	results, errorQuery := db.Query(query, sql.Named("code", empCode)) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.ResultRequestByUser

		errScan := results.Scan(&result.REQUEST_NO, &result.REV, &result.NAME_STATUS, &result.PERSON,
			&result.START_DATE, &result.END_DATE, &result.DURATION, &result.FACTORY_NAME, &result.NAME_WORKCELL,
		) // Scan เก็บข้อมูลใน Struct
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}

func GetApproverPending(c *fiber.Ctx) error {
	info := []model.ApproverPendingAll{}

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
	}

	query := `WITH CTE_REQUEST AS (
	SELECT a.*,a.CURRENT_APPROVE + 1 as NEXT_APPROVER  FROM (
	SELECT REQUEST_NO,MAX(REV) as [REV],SUM(CASE WHEN CODE_APPROVER IS NOT NULL THEN 1 ELSE 0 END) as CURRENT_APPROVE
	FROM TBL_APPROVAL GROUP BY REQUEST_NO ) a  )
	SELECT r.REQUEST_NO,r.REV,r.NEXT_APPROVER,rs.ID_GROUP_DEPT,rs.ID_FACTORY,g.NAME_GROUP,f.FACTORY_NAME,ap.CODE_APPROVER,ap.NAME_APPROVER FROM CTE_REQUEST r 
		LEFT JOIN TBL_REQUESTS rs ON r.REQUEST_NO = rs.REQUEST_NO
		LEFT JOIN TBL_FACTORY f ON rs.ID_FACTORY = f.ID_FACTORY 
		LEFT JOIN TBL_GROUP_DEPT g ON rs.ID_GROUP_DEPT = g.ID_GROUP_DEPT
		LEFT JOIN TBL_APPROVERS ap ON g.ID_GROUP_DEPT = ap.ID_GROUP_DEPT AND f.ID_FACTORY = ap.ID_FACTORY 
		AND r.NEXT_APPROVER = ap.STEP
		ORDER BY r.REQUEST_NO`

	results, errorQuery := db.Query(query) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.ApproverPendingAll

		errScan := results.Scan(&result.REQUEST_NO, &result.REV,
			&result.NEXT_APPROVER, &result.ID_GROUP_DEPT,
			&result.ID_FACTORY, &result.NAME_GROUP, &result.FACTORY_NAME, &result.CODE_APPROVER, &result.NAME_APPROVER) // Scan เก็บข้อมูลใน Struct
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}
func ApproveRequestByNo(c *fiber.Ctx) error {

	var req model.BodyApproveRequest
	var currentStepApprover []model.ResponseApproverStepByReq

	requestNo := c.Params("requestNo")

	rev := c.Params("rev")
	iRev, _ := strconv.Atoi(rev)

	var completeStep int
	var approved int
	var jobStatus int

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
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
	}

	if req.Status == 1 {
		jobStatus = 1
	} else if req.Status == 2 {
		jobStatus = 5
	} else if req.Status == 3 {
		jobStatus = 2
	} else if req.Status == 4 {
		jobStatus = 3
	} else {
		return c.JSON(fiber.Map{"err": true, "msg": "[Status] Approval Process don't have."})
	}

	// Check Step User
	selectApprover, errorSelect := db.Query(`SELECT a.CODE_APPROVER as APPROVER,a.STEP FROM TBL_REQUESTS_HISTORY h LEFT JOIN 
TBL_APPROVERS a ON a.ID_GROUP_DEPT = h.ID_GROUP_DEPT AND a.ID_FACTORY = h.ID_FACTORY
WHERE a.CODE_APPROVER = @code AND h.REQUEST_NO = @requestNo AND REV = @rev  AND 
STEP NOT IN (SELECT STEP FROM TBL_APPROVAL WHERE REQUEST_NO = @requestNo AND REV = rev AND ID_STATUS_APV = 3)`,
		sql.Named("code", req.ActionBy),
		sql.Named("requestNo", requestNo),
		sql.Named("rev", rev))

	if errorSelect != nil {
		fmt.Println("Query failed: " + errorSelect.Error())
	}

	for selectApprover.Next() {
		var stepCurrent model.ResponseApproverStepByReq
		errScan := selectApprover.Scan(&stepCurrent.APPROVER, &stepCurrent.STEP)
		if errScan != nil {
			fmt.Println("Error scanning approver step:", errScan.Error())
		} else {
			currentStepApprover = append(currentStepApprover, stepCurrent)
		}
	}

	if len(currentStepApprover) > 0 {

		// Update Status Approve

		stmtUpdateApproval := `UPDATE [dbo].[TBL_APPROVAL] SET 
		[CODE_APPROVER] = @code,[ID_STATUS_APV] = @status,[REMARK] = @remark,[UPDATED_AT] = GETDATE() ,[SPECIAL_OT] = @special
		WHERE [REQUEST_NO] = @requestNo AND [REV] = @rev AND [STEP] = @step`

		_, errUpdate := db.Exec(stmtUpdateApproval,
			sql.Named("code", req.ActionBy),
			sql.Named("status", req.Status),
			sql.Named("remark", req.Remark),
			sql.Named("requestNo", requestNo),
			sql.Named("rev", rev),
			sql.Named("step", currentStepApprover[0].STEP),
			sql.Named("special", req.Special),
		)

		if errUpdate != nil {
			fmt.Println("Query failed: " + errorSelect.Error())
		}

		query := `SELECT t.STEP FROM TBL_REQUESTS_HISTORY h LEFT JOIN TBL_OT_TYPE t ON h.ID_TYPE_OT = t.ID_TYPE_OT WHERE REQUEST_NO = @requestNo AND REV = @rev 
	AND REQ_STATUS = 1`

		stepComplete, errorStep := db.Query(query, sql.Named("requestNo", requestNo), sql.Named("rev", rev))

		if errorStep != nil {
			fmt.Println("Query failed: " + errorStep.Error())
		}

		for stepComplete.Next() {
			var step int
			errScan := stepComplete.Scan(&step)
			if errScan != nil {
				fmt.Println(errScan.Error())
			} else {
				completeStep = step
			}

		}

		queryApproved, errQueryApproved := db.Query(`SELECT COUNT(*) as [COUNT] FROM TBL_APPROVAL al WHERE al.REQUEST_NO = @requestNo
    AND REV = @rev AND ID_STATUS_APV = 3`, sql.Named("requestNo", requestNo), sql.Named("rev", rev))

		if errQueryApproved != nil {
			fmt.Println("Query failed: " + errQueryApproved.Error())
		}
		for queryApproved.Next() {
			var step int
			errScan := queryApproved.Scan(&step)
			if errScan != nil {
				fmt.Println(errScan.Error())
			} else {
				approved = step
			}

		}
		// 3 is Done (Status Process Approve)
		if req.Status != 3 {
			fmt.Println("Update Status Request and Send Mail to Requestor")

			// Update History Table
			_, errUpdate := db.Exec(`UPDATE TBL_REQUESTS_HISTORY SET REQ_STATUS = @status,
			UPDATED_AT = GETDATE(),UPDATED_BY = @code,[SPECIAL_OT] = @special WHERE REQUEST_NO = @requestNo AND REV = @rev`,
				sql.Named("code", req.ActionBy),
				sql.Named("status", jobStatus),
				sql.Named("requestNo", requestNo),
				sql.Named("rev", rev),
				sql.Named("special", req.Special),
			)

			if errUpdate != nil {
				fmt.Println("Query failed: " + errUpdate.Error())
			}

			// Update Main Request Table
			_, errUpdateMainReq := db.Exec(`UPDATE TBL_REQUESTS SET REQ_STATUS = @status,
			UPDATED_AT = GETDATE(),UPDATED_BY = @code,[SPECIAL_OT] = @special WHERE REQUEST_NO = @requestNo`,
				sql.Named("code", req.ActionBy),
				sql.Named("status", jobStatus),
				sql.Named("requestNo", requestNo),
				sql.Named("special", req.Special),
			)

			if errUpdateMainReq != nil {
				fmt.Println("Query failed: " + errUpdateMainReq.Error())
			}
			var mailRequestor model.MailReturnRequestor
			//SendMail to Requestor

			intRev, errIntConvert := strconv.Atoi(rev)
			if errIntConvert != nil {
				fmt.Println(errIntConvert)
			} else {

				// Check Mail Requestor
				mailRequestor = CheckSendEmailRequestor(intRev, requestNo)

				status := ""

				if req.Status == 2 {
					status = "Reject"
				} else if req.Status == 3 {
					status = "อนุมัติ"
				} else if req.Status == 4 {
					status = "ไม่อนุมัติ"
				}

				if mailRequestor.MAIL != "N/A" {

					SendEMailToRequestor(requestNo, iRev, mailRequestor.MAIL, mailRequestor.FULLNAME, status)
				}
			}

			return c.JSON(fiber.Map{"err": false, "msg": "Updated successfully!", "status": "Ok"})

		}

		fmt.Println("Request is completeStep steps", completeStep)
		fmt.Println("approved", approved)

		// 3 is Done and  All Approved
		if (completeStep == approved) && req.Status == 3 {

			// Update Status and Send Mail to Requestor
			fmt.Println("Update Status Request and Send Mail to Requestor")

			// Update History Table
			_, errUpdate := db.Exec(`UPDATE TBL_REQUESTS_HISTORY SET REMARK = @remark ,REQ_STATUS = @status,
			UPDATED_AT = GETDATE(),UPDATED_BY = @code WHERE REQUEST_NO = @requestNo AND REV = @rev`,
				sql.Named("code", req.ActionBy),
				sql.Named("status", jobStatus),
				sql.Named("remark", req.Remark),
				sql.Named("requestNo", requestNo),
				sql.Named("rev", rev),
			)

			if errUpdate != nil {
				fmt.Println("Query failed: " + errUpdate.Error())
			}

			stmtUpdateApprove := `UPDATE TBL_REQUESTS SET REMARK = @remark ,REQ_STATUS = @status,
			UPDATED_AT = GETDATE(),UPDATED_BY = @code WHERE REQUEST_NO = @requestNo`
			// Update Main Request Table
			_, errUpdateMainReq := db.Exec(stmtUpdateApprove,
				sql.Named("code", req.ActionBy),
				sql.Named("status", jobStatus),
				sql.Named("remark", req.Remark),
				sql.Named("requestNo", requestNo),
			)

			if errUpdateMainReq != nil {
				fmt.Println("Query failed: " + errUpdateMainReq.Error())
			}
			var mailRequestor model.MailReturnRequestor
			//SendMail to Requestor

			intRev, errIntConvert := strconv.Atoi(rev)
			if errIntConvert != nil {
				fmt.Println(errIntConvert)
			} else {

				// Check Mail Requestor
				mailRequestor = CheckSendEmailRequestor(intRev, requestNo)

				if mailRequestor.MAIL != "N/A" {

					SendEMailToRequestor(requestNo, iRev, mailRequestor.MAIL, mailRequestor.FULLNAME, "อนุมัติ")
				}
			}

			fmt.Println(mailRequestor)

		}

		fmt.Println("completeStep", completeStep)
		fmt.Println("approved", approved)
		fmt.Println(completeStep)

		// 3 is Done (Status Process Approve)
		if (completeStep > approved) && req.Status == 3 {
			fmt.Println("Update Approval Status and Send mail to Next")
			fmt.Println("currentStepApprover", currentStepApprover[0].STEP+1)
			if (currentStepApprover[0].STEP + 1) <= completeStep {

				// Update Pending to Next Approver
				stmtUpdateApprove := `UPDATE [dbo].[TBL_APPROVAL] SET [ID_STATUS_APV] = 1 ,[UPDATED_AT] = GETDATE() 
				WHERE [REQUEST_NO] = @requestNo AND [REV] = @rev AND [STEP] = @step`
				// Update Main Request Table
				_, errUpdateStatus := db.Exec(stmtUpdateApprove,

					sql.Named("step", currentStepApprover[0].STEP+1),
					sql.Named("requestNo", requestNo),
					sql.Named("rev", iRev),
				)

				if errUpdateStatus != nil {
					fmt.Println(errUpdateStatus.Error())
				} else {
					// mail
					mailApprover := CheckSendEmail(iRev, requestNo)

					fmt.Println(mailApprover)
					if mailApprover.Email != "N/A" {
						SendEMailToApprover(requestNo, iRev, mailApprover.Email, mailApprover.FULLNAME)
					}
				}

			}

		}

		return c.JSON(fiber.Map{"err": false, "msg": "Updated successfully!", "status": "Ok"})
	} else {
		return c.JSON(fiber.Map{"err": true, "msg": "Permission is denined!"})
	}
}

func CountApproveStatusByCode(c *fiber.Ctx) error {
	info := []model.ResultCountApproveByEmpCode{}
	empCode := c.Params("code")
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
	}

	query := `SELECT asp.NAME_STATUS,asp.ID_STATUS_APV,ISNULL(app.AMOUNT,0) as [AMOUNT] 
						FROM TBL_APPROVE_STATUS asp LEFT JOIN (
							SELECT m.CODE_APPROVER,m.ID_STATUS_APV,m.NAME_STATUS,COUNT(*) as [AMOUNT] FROM (
							SELECT 
							  ap.REQUEST_NO,ap.REV,ISNULL(ap.ID_STATUS_APV,0) as ID_STATUS_APV,apv.CODE_APPROVER,
							  ISNULL(s.NAME_STATUS,'N/A') as [NAME_STATUS] FROM TBL_APPROVAL ap LEFT JOIN 
							  TBL_REQUESTS_HISTORY h ON ap.REQUEST_NO = h.REQUEST_NO
							  AND ap.REV = h.REV
							  LEFT JOIN TBL_APPROVERS  apv 
							  ON h.ID_GROUP_DEPT = apv.ID_GROUP_DEPT
							  AND h.ID_FACTORY = apv.ID_FACTORY
							  LEFT JOIN TBL_APPROVE_STATUS s ON ap.ID_STATUS_APV = s.ID_STATUS_APV
							  ) m 
							  WHERE m.CODE_APPROVER = @code
							  GROUP BY m.CODE_APPROVER,m.ID_STATUS_APV,m.NAME_STATUS ) app 
						ON asp.ID_STATUS_APV = app.ID_STATUS_APV ORDER BY asp.ID_STATUS_APV ASC`

	results, errorQuery := db.Query(query, sql.Named("code", empCode)) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.ResultCountApproveByEmpCode

		errScan := results.Scan(&result.NAME_STATUS, &result.ID_STATUS_APV, &result.AMOUNT) // Scan เก็บข้อมูลใน Struct
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}

func GetRequestListByCodeAndStatus(c *fiber.Ctx) error {
	info := []model.ResultListRequestByEmpIdAndStatus{}
	empCode := c.Params("code")
	status := c.Params("status")
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
	}

	query := `SELECT REQUEST_NO,CODE_APPROVER,REV,FACTORY_NAME,NAME_GROUP,ID_FACTORY,ID_GROUP_DEPT,COUNT_USER,DURATION,HOURS_AMOUNT,SUM_MINUTE,MINUTE_TOTAL  
				FROM [dbo].[Func_GetLists_Status_And_Empc] (@code, @status) a
				ORDER BY REQUEST_NO DESC`

	results, errorQuery := db.Query(query, sql.Named("code", empCode), sql.Named("status", status)) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.ResultListRequestByEmpIdAndStatus

		errScan := results.Scan(
			&result.REQUEST_NO,
			&result.CODE_APPROVER,
			&result.REV,
			&result.FACTORY_NAME,
			&result.NAME_GROUP,
			&result.ID_FACTORY,
			&result.ID_GROUP_DEPT,
			&result.COUNT_USER,
			&result.DURATION,
			&result.HOURS_AMOUNT,
			&result.SUM_MINUTE,
			&result.MINUTE_TOTAL,
		) // Scan เก็บข้อมูลใน Struct

		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}

func SummaryLastRevRequestAll(c *fiber.Ctx) error {
	info := []model.SummaryRequestLastRev{}

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
	}

	query := `SELECT r.REQUEST_NO,h.REV ,s.NAME_STATUS,hr.UHR_FullName_th as FULLNAME,tot.HOURS_AMOUNT as OT_TYPE,f.FACTORY_NAME,wg.NAME_WORKGRP,
	wc.NAME_WORKCELL,us.USERS as USERS_AMOUNT,mi.SUM_MINUTE,r.START_DATE,r.END_DATE FROM  [dbo].[TBL_REQUESTS] r LEFT JOIN 
	(SELECT MAX(REV) as REV,REQUEST_NO FROM TBL_REQUESTS_HISTORY GROUP BY REQUEST_NO) h ON
	r.REQUEST_NO = h.REQUEST_NO
	LEFT JOIN TBL_REQ_STATUS s ON r.REQ_STATUS = s.ID_STATUS
	LEFT JOIN V_AllUserPSTH hr ON r.CREATED_BY COLLATE thai_CI_AS = hr.UHR_EmpCode COLLATE thai_CI_AS 
	LEFT JOIN TBL_OT_TYPE tot ON tot.ID_TYPE_OT = r.ID_TYPE_OT
	LEFT JOIN TBL_FACTORY f ON r.ID_FACTORY = f.ID_FACTORY
	LEFT JOIN TBL_WORK_GROUP wg ON r.ID_WORKGRP = wg.ID_WORKGRP
	LEFT JOIN TBL_WORKCELL wc ON r.ID_WORK_CELL = wc.ID_WORK_CELL
	LEFT JOIN (SELECT COUNT(*) as USERS,REQUEST_NO,REV 
	    FROM TBL_USERS_REQ GROUP BY REQUEST_NO,REV) us ON r.REQUEST_NO = us.REQUEST_NO AND h.REV = us.REV
	LEFT JOIN   (SELECT uq.REQUEST_NO,uq.REV,SUM(m.MINUTE_DIFF) as SUM_MINUTE FROM TBL_USERS_REQ uq  
	LEFT JOIN (SELECT REQUEST_NO,REV,DATEDIFF(MINUTE,START_DATE,END_DATE) as MINUTE_DIFF FROM TBL_REQUESTS_HISTORY hh) m 
   ON uq.REQUEST_NO = m.REQUEST_NO AND uq.REV = m.REV GROUP BY uq.REQUEST_NO,uq.REV) mi ON r.REQUEST_NO = mi.REQUEST_NO AND h.REV = mi.REV
   ORDER BY r.REQUEST_NO DESC`

	results, errorQuery := db.Query(query) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.SummaryRequestLastRev

		errScan := results.Scan(
			&result.REQUEST_NO,
			&result.REV,
			&result.NAME_STATUS,
			&result.FULLNAME,
			&result.OT_TYPE,
			&result.FACTORY_NAME,
			&result.NAME_WORKGRP,
			&result.NAME_WORKCELL,
			&result.USERS_AMOUNT,
			&result.SUM_MINUTE,
			&result.START_DATE,
			&result.END_DATE,
		) // Scan เก็บข้อมูลใน Struct

		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}

func SummaryLastRevRequestAllByReqNo(c *fiber.Ctx) error {
	info := []model.SummaryRequestLastRev{}
	reqNo := c.Params("reqNo")
	rev := c.Params("rev")

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
	}

	query := `SELECT r.REQUEST_NO,h.REV ,s.NAME_STATUS,hr.UHR_FullName_th as FULLNAME,tot.HOURS_AMOUNT as OT_TYPE,f.FACTORY_NAME,wg.NAME_WORKGRP,
	wc.NAME_WORKCELL,us.USERS as USERS_AMOUNT,mi.SUM_MINUTE,r.START_DATE,r.END_DATE,f.ID_FACTORY,ISNULL(pln.SUM_PLAN,0) as [SUM_PLAN],ISNULL(pln.SUM_PLAN_OB,0) as [SUM_PLAN_OB],wc.ID_WORK_CELL,h.REMARK FROM  [dbo].[TBL_REQUESTS] r LEFT JOIN 
	(SELECT REV,REQUEST_NO,REMARK FROM TBL_REQUESTS_HISTORY) h ON
	r.REQUEST_NO = h.REQUEST_NO
	LEFT JOIN TBL_REQ_STATUS s ON r.REQ_STATUS = s.ID_STATUS
	LEFT JOIN V_AllUserPSTH hr ON r.CREATED_BY COLLATE thai_CI_AS = hr.UHR_EmpCode COLLATE thai_CI_AS 
	LEFT JOIN TBL_OT_TYPE tot ON tot.ID_TYPE_OT = r.ID_TYPE_OT
	LEFT JOIN TBL_FACTORY f ON r.ID_FACTORY = f.ID_FACTORY
	LEFT JOIN TBL_WORK_GROUP wg ON r.ID_WORKGRP = wg.ID_WORKGRP
	LEFT JOIN TBL_WORKCELL wc ON r.ID_WORK_CELL = wc.ID_WORK_CELL
	LEFT JOIN (SELECT COUNT(*) as USERS,REQUEST_NO,REV 
	    FROM TBL_USERS_REQ GROUP BY REQUEST_NO,REV) us ON r.REQUEST_NO = us.REQUEST_NO AND h.REV = us.REV
	LEFT JOIN   (SELECT uq.REQUEST_NO,uq.REV,SUM(m.MINUTE_DIFF) as SUM_MINUTE FROM TBL_USERS_REQ uq  
	LEFT JOIN (SELECT REQUEST_NO,REV,DATEDIFF(MINUTE,START_DATE,END_DATE) as MINUTE_DIFF FROM TBL_REQUESTS_HISTORY hh) m 
   ON uq.REQUEST_NO = m.REQUEST_NO AND uq.REV = m.REV GROUP BY uq.REQUEST_NO,uq.REV) mi ON r.REQUEST_NO = mi.REQUEST_NO AND h.REV = mi.REV
	LEFT JOIN (
	
   SELECT ISNULL(pwc.ID_FACTORY,pp.ID_FACTORY) as ID_FACTORY,ISNULL(pp.SUM_PLAN_OB,0) as [SUM_PLAN_OB],ISNULL(pwc.SUM_PLAN,0) as [SUM_PLAN],
   ISNULL(pp.[YEAR],pwc.[YEAR]) as [YEAR],ISNULL(pp.[MONTH],pwc.[MONTH]) as [MONTH] FROM (
   SELECT SUM(HOURS) as SUM_PLAN,f.ID_FACTORY,YEAR,MONTH FROM TBL_PLAN_OVERTIME  po
   LEFT JOIN TBL_WORKCELL wc ON wc.ID_WORK_CELL = po.ID_WORK_CELL
   LEFT JOIN TBL_FACTORY f ON f.ID_FACTORY = wc.ID_FACTORY
   GROUP BY YEAR,MONTH,f.ID_FACTORY ) pwc
   FULL JOIN (
   SELECT SUM(HOURS) as SUM_PLAN_OB,pob.ID_FACTORY,pob.YEAR,pob.MONTH FROM TBL_PLAN_OB pob GROUP BY pob.ID_FACTORY,pob.YEAR,pob.MONTH )
	pp ON pwc.ID_FACTORY = pp.ID_FACTORY AND pwc.YEAR = pp.YEAR AND pwc.MONTH = pp.MONTH
	)pln ON r.ID_FACTORY = pln.ID_FACTORY AND YEAR(r.START_DATE) = pln.[YEAR] AND MONTH(r.START_DATE) = pln.[MONTH]

   WHERE h.REQUEST_NO = @reqNo AND h.REV = @rev
   ORDER BY r.REQUEST_NO DESC`

	results, errorQuery := db.Query(query, sql.Named("reqNo", reqNo), sql.Named("rev", rev)) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.SummaryRequestLastRev

		errScan := results.Scan(
			&result.REQUEST_NO,
			&result.REV,
			&result.NAME_STATUS,
			&result.FULLNAME,
			&result.OT_TYPE,
			&result.FACTORY_NAME,
			&result.NAME_WORKGRP,
			&result.NAME_WORKCELL,
			&result.USERS_AMOUNT,
			&result.SUM_MINUTE,
			&result.START_DATE,
			&result.END_DATE,
			&result.ID_FACTORY,
			&result.SUM_PLAN,
			&result.SUM_PLAN_OB,
			&result.ID_WORK_CELL,
			&result.REMARK,
		) // Scan เก็บข้อมูลใน Struct

		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}

func GetApproverCommentByRequestNo(c *fiber.Ctx) error {
	info := []model.RequestCommentApprover{}
	reqNo := c.Params("requestNo")
	rev := c.Params("rev")

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
	}

	query := `SELECT a.REQUEST_NO,a.ID_STATUS_APV,a.CODE_APPROVER,a.CREATED_AT,a.UPDATED_AT,s.NAME_STATUS,
a.REMARK,hr.UHR_Department as DEPARTMENT,hr.UHR_Position as POSITION,hr.UHR_FullName_th as FULLNAME FROM TBL_APPROVAL  a
LEFT JOIN [dbo].[TBL_APPROVE_STATUS] s ON a.ID_STATUS_APV = s.ID_STATUS_APV
LEFT JOIN V_AllUserPSTH hr ON a.CODE_APPROVER COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS
WHERE REQUEST_NO = @reqNo AND REV = @rev 
ORDER BY a.UPDATED_AT DESC`

	results, errorQuery := db.Query(query, sql.Named("reqNo", reqNo), sql.Named("rev", rev)) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.RequestCommentApprover

		errScan := results.Scan(
			&result.REQUEST_NO,
			&result.ID_STATUS_APV,
			&result.CODE_APPROVER,
			&result.CREATED_AT,
			&result.UPDATED_AT,
			&result.NAME_STATUS,
			&result.REMARK,
			&result.DEPARTMENT,
			&result.POSITION,
			&result.FULLNAME,
		) // Scan เก็บข้อมูลใน Struct

		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}

func GetRequestListByStatusApprove(c *fiber.Ctx) error {
	var actualAll []model.RequestList
	code := c.Params("code")
	status := c.Params("status")
	fmt.Println(code)
	strConfig := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	stmt := `SELECT m.REQUEST_NO,m.REV,m.FACTORY_NAME,m.ID_TYPE_OT,CONCAT('OT',m.HOURS_AMOUNT) as HOURS_AMOUNT,p.PERSON,du.DURATION,p.PERSON * du.DURATION as HOURS_TOTAL FROM (
  SELECT MAX(a.REV) as REV,a.REQUEST_NO,t.ID_TYPE_OT,t.HOURS_AMOUNT,f.FACTORY_NAME FROM TBL_APPROVAL a
  LEFT JOIN TBL_REQUESTS_HISTORY h ON a.REQUEST_NO = h.REQUEST_NO AND a.REV = h.REV
  LEFT JOIN TBL_GROUP_DEPT d ON h.ID_GROUP_DEPT = d.ID_GROUP_DEPT
  LEFT JOIN TBL_FACTORY f ON f.ID_FACTORY = h.ID_FACTORY
  LEFT JOIN TBL_APPROVERS ap ON ap.ID_GROUP_DEPT = d.ID_GROUP_DEPT AND ap.ID_FACTORY = f.ID_FACTORY
  LEFT JOIN TBL_OT_TYPE t ON h.ID_TYPE_OT = t.ID_TYPE_OT
 
  WHERE ap.CODE_APPROVER = @code AND a.ID_STATUS_APV = @status
   GROUP BY a.REQUEST_NO,t.ID_TYPE_OT,t.HOURS_AMOUNT,f.FACTORY_NAME ) m 
    LEFT JOIN (SELECT COUNT(*) as PERSON,REQUEST_NO,REV FROM TBL_USERS_REQ u GROUP BY REQUEST_NO,REV) p 
	ON m.REQUEST_NO = p.REQUEST_NO AND m.REV = p.REV
	LEFT JOIN (SELECT REQUEST_NO,REV,CAST(DATEDIFF(MINUTE,START_DATE,END_DATE)/60 as decimal(18,2)) as DURATION 
	FROM TBL_REQUESTS_HISTORY) du ON m.REQUEST_NO = du.REQUEST_NO AND m.REV = du.REV ORDER BY REQUEST_NO DESC`

	rows, errSelect := db.Query(stmt, sql.Named("code", code), sql.Named("status", status))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.RequestList

		errorScan := rows.Scan(
			&actual.REQUEST_NO,
			&actual.REV,
			&actual.FACTORY_NAME,
			&actual.ID_TYPE_OT,
			&actual.HOURS_AMOUNT,
			&actual.PERSON,
			&actual.DURATION,
			&actual.HOURS_TOTAL,
		)

		if errorScan != nil {
			fmt.Println("Error Scan : ", errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			actualAll = append(actualAll, actual)
		}
	}

	if len(actualAll) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": actualAll,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": actualAll,
			"msg":     "Not Found",
		})
	}

}

func GetRequestListByStatusApproveAndCode(c *fiber.Ctx) error {
	var actualAll []model.RequestList
	code := c.Params("code")
	status := c.Params("status")

	strConfig := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	stmt := `SELECT  mt.*,mt.PERSON * mt.DURATION as HOURS_TOTAL FROM (
  SELECT DISTINCT h.REQUEST_NO,h.REV,f.FACTORY_NAME,ot.ID_TYPE_OT,CONCAT('OT',ot.HOURS_AMOUNT) as [HOURS_AMOUNT],p.PERSON,CAST(DATEDIFF(MINUTE,h.START_DATE,h.END_DATE) / 60  as decimal(18,2))as DURATION FROM TBL_APPROVAL a 
  LEFT JOIN TBL_REQUESTS_HISTORY h ON h.REQUEST_NO = a.REQUEST_NO AND a.REV = h.REV
  LEFT JOIN TBL_OT_TYPE ot ON h.ID_TYPE_OT = ot.ID_TYPE_OT
  LEFT JOIN TBL_FACTORY f ON h.ID_FACTORY = f.ID_FACTORY
  LEFT JOIN (SELECT COUNT(*) as PERSON,REQUEST_NO,REV FROM TBL_USERS_REQ u GROUP BY REQUEST_NO,REV) p 
  ON h.REQUEST_NO = p.REQUEST_NO AND h.REV = p.REV
  WHERE a.CODE_APPROVER = @code AND h.REQUEST_NO IS NOT NULL
  AND a.ID_STATUS_APV = @status)mt`

	rows, errSelect := db.Query(stmt, sql.Named("code", code), sql.Named("status", status))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.RequestList

		errorScan := rows.Scan(
			&actual.REQUEST_NO,
			&actual.REV,
			&actual.FACTORY_NAME,
			&actual.ID_TYPE_OT,
			&actual.HOURS_AMOUNT,
			&actual.PERSON,
			&actual.DURATION,
			&actual.HOURS_TOTAL,
		)

		if errorScan != nil {
			fmt.Println("Error Scan : ", errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			actualAll = append(actualAll, actual)
		}
	}

	if len(actualAll) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": actualAll,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": actualAll,
			"msg":     "Not Found",
		})
	}

}

func GetUserRequestListByStatusPendApprove(c *fiber.Ctx) error {
	var actualAll []model.RequestList
	code := c.Params("code")
	status := c.Params("status")

	strConfig := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	stmt := `SELECT m.REQUEST_NO,m.REV,m.FACTORY_NAME,m.ID_TYPE_OT,CONCAT('OT',m.HOURS_AMOUNT) as HOURS_AMOUNT,p.PERSON,du.DURATION,p.PERSON * du.DURATION as HOURS_TOTAL FROM (
  SELECT MAX(a.REV) as REV,a.REQUEST_NO,t.ID_TYPE_OT,t.HOURS_AMOUNT,f.FACTORY_NAME FROM TBL_APPROVAL a
  LEFT JOIN TBL_REQUESTS_HISTORY h ON a.REQUEST_NO = h.REQUEST_NO AND a.REV = h.REV
  LEFT JOIN TBL_GROUP_DEPT d ON h.ID_GROUP_DEPT = d.ID_GROUP_DEPT
  LEFT JOIN TBL_FACTORY f ON f.ID_FACTORY = h.ID_FACTORY
  LEFT JOIN TBL_APPROVERS ap ON ap.ID_GROUP_DEPT = d.ID_GROUP_DEPT AND ap.ID_FACTORY = f.ID_FACTORY
  LEFT JOIN TBL_OT_TYPE t ON h.ID_TYPE_OT = t.ID_TYPE_OT
 
  WHERE h.CREATED_BY = @code AND a.ID_STATUS_APV = @status
   GROUP BY a.REQUEST_NO,t.ID_TYPE_OT,t.HOURS_AMOUNT,f.FACTORY_NAME ) m 
    LEFT JOIN (SELECT COUNT(*) as PERSON,REQUEST_NO,REV FROM TBL_USERS_REQ u GROUP BY REQUEST_NO,REV) p 
	ON m.REQUEST_NO = p.REQUEST_NO AND m.REV = p.REV
	LEFT JOIN (SELECT REQUEST_NO,REV,CAST(DATEDIFF(MINUTE,START_DATE,END_DATE)/60 as decimal(18,2)) as DURATION 
	FROM TBL_REQUESTS_HISTORY) du ON m.REQUEST_NO = du.REQUEST_NO AND m.REV = du.REV`

	rows, errSelect := db.Query(stmt, sql.Named("code", code), sql.Named("status", status))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.RequestList

		errorScan := rows.Scan(
			&actual.REQUEST_NO,
			&actual.REV,
			&actual.FACTORY_NAME,
			&actual.ID_TYPE_OT,
			&actual.HOURS_AMOUNT,
			&actual.PERSON,
			&actual.DURATION,
			&actual.HOURS_TOTAL,
		)

		if errorScan != nil {
			fmt.Println("Error Scan : ", errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			actualAll = append(actualAll, actual)
		}
	}

	if len(actualAll) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": actualAll,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": actualAll,
			"msg":     "Not Found",
		})
	}

}

func GetUserRequestListByStatusApprove(c *fiber.Ctx) error {
	var actualAll []model.RequestList
	code := c.Params("code")
	status := c.Params("status")

	strConfig := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	stmt := `SELECT  ms.REQUEST_NO,ms.REV,ar.FACTORY_NAME,ar.ID_TYPE_OT,ar.HOURS_AMOUNT,ar.PERSON,ar.DURATION,ar.HOURS_TOTAL FROM (
SELECT MAX(ht.REV) as REV,
ht.REQUEST_NO FROM TBL_REQUESTS_HISTORY ht GROUP BY ht.REQUEST_NO ) ms
LEFT JOIN (
	SELECT mt.*,mt.PERSON * mt.DURATION as HOURS_TOTAL FROM (
		SELECT h.REQUEST_NO,h.REV,f.FACTORY_NAME,t.ID_TYPE_OT,
		CONCAT('OT',t.HOURS_AMOUNT) as HOURS_AMOUNT,p.PERSON,
		CAST(DATEDIFF(MINUTE,h.START_DATE,h.END_DATE)/60 as decimal(18,2)) as DURATION  
		FROM TBL_REQUESTS_HISTORY h
		LEFT JOIN TBL_FACTORY f ON h.ID_FACTORY = f.ID_FACTORY
		LEFT JOIN TBL_OT_TYPE t ON h.ID_TYPE_OT = t.ID_TYPE_OT
	
		LEFT JOIN (
		SELECT COUNT(*) as PERSON,REQUEST_NO,REV FROM TBL_USERS_REQ u GROUP BY REQUEST_NO,REV
		) p 
		ON h.REQUEST_NO = p.REQUEST_NO AND h.REV = p.REV 
		WHERE h.CREATED_BY = @code AND h.REQ_STATUS = @status
		) mt 
	) ar ON ms.REQUEST_NO = ar.REQUEST_NO AND ms.REV = ar.REV WHERE FACTORY_NAME IS NOT NULL ORDER BY REQUEST_NO DESC`

	rows, errSelect := db.Query(stmt, sql.Named("code", code), sql.Named("status", status))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.RequestList

		errorScan := rows.Scan(
			&actual.REQUEST_NO,
			&actual.REV,
			&actual.FACTORY_NAME,
			&actual.ID_TYPE_OT,
			&actual.HOURS_AMOUNT,
			&actual.PERSON,
			&actual.DURATION,
			&actual.HOURS_TOTAL,
		)

		if errorScan != nil {
			fmt.Println("Error Scan : ", errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			actualAll = append(actualAll, actual)
		}
	}

	if len(actualAll) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": actualAll,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": actualAll,
			"msg":     "Not Found",
		})
	}

}

func CountRequestByStatusAndCode(c *fiber.Ctx) error {
	var amountRequest []model.CountRequestByStatus
	code := c.Params("code")
	status := c.Params("status")

	strConfig := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	stmt := `SELECT COUNT(*) as STATUS_PENDING FROM TBL_REQUESTS WHERE REQ_STATUS = @status AND CREATED_BY = @code`

	rows, errSelect := db.Query(stmt, sql.Named("code", code), sql.Named("status", status))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.CountRequestByStatus

		errorScan := rows.Scan(
			&actual.STATUS_PENDING,
		)

		if errorScan != nil {
			fmt.Println("Error Scan : ", errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			amountRequest = append(amountRequest, actual)
		}
	}

	if len(amountRequest) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": amountRequest,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": amountRequest,
			"msg":     "Not Found",
		})
	}

}

func CountRequestStatusAndApproveByCode(c *fiber.Ctx) error {
	var amountRequest []model.CountRequestByStatus
	code := c.Params("code")
	status := c.Params("status")

	strConfig := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	stmt := `SELECT  COUNT(*) as [STATUS_PENDING] FROM TBL_APPROVAL  a
LEFT JOIN TBL_REQUESTS_HISTORY h ON a.REQUEST_NO = h.REQUEST_NO AND a.REV = h.REV
LEFT JOIN TBL_APPROVERS ap ON h.ID_FACTORY = ap.ID_FACTORY AND h.ID_GROUP_DEPT = ap.ID_GROUP_DEPT
WHERE ap.CODE_APPROVER = @code AND a.ID_STATUS_APV = @status`

	rows, errSelect := db.Query(stmt, sql.Named("code", code), sql.Named("status", status))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.CountRequestByStatus

		errorScan := rows.Scan(
			&actual.STATUS_PENDING,
		)

		if errorScan != nil {
			fmt.Println("Error Scan : ", errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			amountRequest = append(amountRequest, actual)
		}
	}
	fmt.Println(amountRequest)

	if len(amountRequest) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": amountRequest,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": amountRequest,
			"msg":     "Not Found",
		})
	}

}

func GetDetailOldRequestByStatus(c *fiber.Ctx) error {
	var details []model.OldRequestDetail

	requestNo := c.Params("requestNo")
	status := c.Params("status")
	rev := c.Params("rev")

	strConfig := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	stmt := `SELECT REQUEST_NO,REV,ID_FACTORY,ID_GROUP_DEPT,ID_WORK_CELL,START_DATE,END_DATE,ID_TYPE_OT,ID_WORKGRP,REMARK FROM TBL_REQUESTS_HISTORY 
WHERE REQUEST_NO = @requestNo AND REV = @rev AND REQ_STATUS = @status`

	rows, errSelect := db.Query(stmt, sql.Named("requestNo", requestNo), sql.Named("rev", rev), sql.Named("status", status))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var info model.OldRequestDetail

		errorScan := rows.Scan(
			&info.REQUEST_NO,
			&info.REV,
			&info.ID_FACTORY,
			&info.ID_GROUP_DEPT,
			&info.ID_WORK_CELL,
			&info.START_DATE,
			&info.END_DATE,
			&info.ID_TYPE_OT,
			&info.ID_WORKGRP,
			&info.REMARK,
		)

		if errorScan != nil {
			fmt.Println("Error Scan : ", errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			details = append(details, info)
		}
	}

	if len(details) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": details,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": details,
			"msg":     "Not Found",
		})
	}

}

func GetRequestHistoryByFactory(c *fiber.Ctx) error {
	var details []model.HistoryRequest

	factory := c.Params("factory")
	start := c.Params("start")
	end := c.Params("end")

	strConfig := config.LoadDatabaseConfig()
	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}

	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	condition := ``

	stmt := fmt.Sprintf(`SELECT [REQUEST_NO]
	  ,[REV]
	  ,[START_DATE]
	  ,[END_DATE]
	  ,[FACTORY_NAME]
	  ,[NAME_WORKCELL]
	  ,[NAME_WORKGRP]
	  ,[ID_FACTORY]
	  ,[ID_GROUP_DEPT]
	  ,[NAME_GROUP]
	  ,[ID_TYPE_OT]
	  ,[HOURS_AMOUNT]
	  ,[NAME_STATUS]
	  ,[REQUESTOR]
	  ,[REQUESTOR_NAME]
	  ,[PERSON]
	  ,[FINAL_STEP]
	  ,[STATUS_DESC]
	  ,[DURATION]
	  ,[TOTAL_DURATION]
	  ,[PENDING_APPROVER]
	  ,[PENDING_CODE]
  FROM [DB_OT_MANAGEMENT].[dbo].[V_History_Requests] WHERE CONVERT(DATE,START_DATE) BETWEEN '%s' AND '%s'
  %s ORDER BY REQUEST_NO,REV DESC`, start, end, condition)

	if factory != "all" {
		condition += ` AND ID_FACTORY  = @factory`
	}

	rows, errSelect := db.Query(stmt, sql.Named("factory", factory))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var info model.HistoryRequest

		errorScan := rows.Scan(
			&info.REQUEST_NO,
			&info.REV,
			&info.START_DATE,
			&info.END_DATE,
			&info.FACTORY_NAME,
			&info.NAME_WORKCELL,
			&info.NAME_WORKGRP,
			&info.ID_FACTORY,
			&info.ID_GROUP_DEPT,
			&info.NAME_GROUP,
			&info.ID_TYPE_OT,
			&info.HOURS_AMOUNT,
			&info.NAME_STATUS,
			&info.REQUESTOR,
			&info.REQUESTOR_NAME,
			&info.PERSON,
			&info.FINAL_STEP,
			&info.STATUS_DESC,
			&info.DURATION,
			&info.TOTAL_DURATION,
			&info.PENDING_APPROVER,
			&info.PENDING_CODE,
		)

		if errorScan != nil {
			fmt.Println("Error Scan : ", errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			details = append(details, info)
		}
	}

	if len(details) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": details,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": details,
			"msg":     "Not Found",
		})
	}

}
