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
			WHERE REQUEST_NO = @reqNo AND [STEP] = 1`, sql.Named("reqNo", running))
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

			SendEMailToApprover(running, 1, mail.Email)
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

	results, err := db.Query(`SELECT REQUEST_NO,SUM( CASE WHEN CODE_APPROVER IS NOT NULL THEN  1 ELSE 0 END) as APPROVED_COUNT,REQ_STATUS,STEP,CREATED_BY FROM (
	SELECT a.REQUEST_NO,r.CREATED_BY,REQ_STATUS,a.CODE_APPROVER,t.STEP FROM [dbo].[TBL_REQUESTS] r
	LEFT JOIN TBL_APPROVAL a ON r.REQUEST_NO = a.REQUEST_NO
	LEFT JOIN TBL_OT_TYPE t ON r.ID_TYPE_OT = t.ID_TYPE_OT
	WHERE  r.REQUEST_NO = @requestNo )m GROUP BY m.REQUEST_NO,CREATED_BY,REQ_STATUS,CODE_APPROVER,STEP`, sql.Named("requestNo", requestNo))

	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer results.Close()

	for results.Next() {
		var approved model.ResultCheckApproved

		errScan := results.Scan(&approved.REQUEST_NO, &approved.APPROVED_COUNT, &approved.REQ_STATUS, &approved.STEP, &approved.CREATED_BY)

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
		requestBy := resultsCheckApproved[0].CREATED_BY

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
		} else if user["employee_code"] != requestBy {

			return c.JSON(fiber.Map{
				"err": true,
				"msg": "Permission is't corrected.",
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

		} else if status != 1 && status == 5 {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": "Can't cancel Request No: Rejected.",
			})

		} else {

			// Update Status to Cencel
			stmtCancel := `UPDATE TBL_REQUESTS 
                SET REQ_STATUS = 4, UPDATED_AT = GETDATE(), UPDATED_BY = @actionBy 
                WHERE REQUEST_NO = @requestNo`

			update, errUpdate := db.Exec(stmtCancel,
				sql.Named("requestNo", requestNo),
				sql.Named("actionBy", user["employee_code"]))

			if errUpdate != nil {
				return c.JSON(fiber.Map{
					"err": true,
					"msg": errUpdate.Error(), // Use errUpdate instead of err
				})
			}

			// Retrieve RowsAffected and handle error
			rowsAffected, errRows := update.RowsAffected()
			if errRows != nil {
				return c.JSON(fiber.Map{
					"err": true,
					"msg": errRows.Error(),
				})
			}

			// Check if any rows were updated
			if rowsAffected == 0 {
				return c.JSON(fiber.Map{
					"err": true,
					"msg": "No rows were updated. Request might not exist or already processed.",
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
			"err": true,
			"msg": "Request isn't found!",
		})

	}

}

func RewiteRequestOvertime(c *fiber.Ctx) error {
	var req model.RewriteRequestOvertimeBody
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

		_, errorInsertHistory := db.Exec(`INSERT INTO [dbo].[TBL_REQUESTS_HISTORY] ([REQUEST_NO],[ID_GROUP_DEPT],
		[DATE_OT],[REMARK],[REQ_STATUS],[ID_TYPE_OT],[ID_FACTORY],[CREATED_AT],[CREATED_BY],[REV],[START_DATE],[END_DATE],[ID_WORKGRP],[ID_WORK_CELL])
		VALUES (@reqNo,@group,@date,@remark,@status,@type,@factory,GETDATE(),@actionBy,@rev,@start,@end,@groupWorkcell,@workcell)`,
			sql.Named("reqNo", req.RequestNo),
			sql.Named("group", req.GroupDept),
			sql.Named("date", req.OvertimeDate),
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
			WHERE REQUEST_NO = @reqNo AND [STEP] = 1 AND [REV] = @rev`, sql.Named("reqNo", req.RequestNo), sql.Named("rev", lastRevNo+1))
			if errorUpdate != nil {
				fmt.Println("Query failed: " + errResult.Error())
			}

		}

		for j := 0; j < len(users); j++ {
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
		[DATE_OT] = @date,[REMARK] = @remark,[REQ_STATUS] = @status,
		[ID_TYPE_OT] = @type,[ID_FACTORY] = @factory,[UPDATED_AT] = GETDATE(),[UPDATED_BY] = @actionBy,
		[START_DATE] = @start,[END_DATE] = @end,[ID_WORKGRP] = @groupWorkcell,[ID_WORK_CELL] = @workcell WHERE [REQUEST_NO] = @reqNo`,
			sql.Named("group", req.GroupDept),
			sql.Named("date", req.OvertimeDate),
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

		return c.JSON(fiber.Map{
			"err":    false,
			"msg":    "Rewrited successfully!",
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
			LEFT JOIN TBL_REQUESTS r ON a.REQUEST_NO = r.REQUEST_NO ) m
			WHERE YEAR(m.DATE_OT) = @year
			GROUP BY YEAR(m.DATE_OT),MONTH(m.DATE_OT) ORDER BY MONTH(m.DATE_OT) DESC`

	results, errorQueryser := db.Query(queryUser, sql.Named("year", year)) //Query

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

	query := `SELECT a.REQUEST_NO,a.REQ_STATUS,b.REV,s.NAME_STATUS FROM (
			SELECT REQUEST_NO,REQ_STATUS FROM TBL_REQUESTS WHERE CREATED_BY = @code 
			GROUP BY REQUEST_NO,REQ_STATUS ) a
			LEFT JOIN (SELECT REQUEST_NO,MAX(REV) as REV FROM TBL_REQUESTS_HISTORY h 
			GROUP BY REQUEST_NO) b ON a.REQUEST_NO = b.REQUEST_NO 
			LEFT JOIN TBL_REQ_STATUS s ON a.REQ_STATUS = s.ID_STATUS
			ORDER BY a.REQUEST_NO DESC`

	results, errorQuery := db.Query(query, sql.Named("code", empCode)) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.ResultRequestByUser

		errScan := results.Scan(&result.REQUEST_NO, &result.REQ_STATUS, &result.REV, &result.NAME_STATUS) // Scan เก็บข้อมูลใน Struct
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
WHERE a.CODE_APPROVER = @code AND h.REQUEST_NO = @requestNo AND REV = @rev`,
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

		_, errUpdate := db.Exec(`UPDATE [dbo].[TBL_APPROVAL] SET 
		[CODE_APPROVER] = @code,[ID_STATUS_APV] = @status,[REMARK] = @remark,[UPDATED_AT] = GETDATE() 
		WHERE [REQUEST_NO] = @requestNo AND [REV] = @rev AND [STEP] = @step`,
			sql.Named("code", req.ActionBy),
			sql.Named("status", req.Status),
			sql.Named("remark", req.Remark),
			sql.Named("requestNo", requestNo),
			sql.Named("rev", rev),
			sql.Named("step", currentStepApprover[0].STEP),
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

			// Update Main Request Table
			_, errUpdateMainReq := db.Exec(`UPDATE TBL_REQUESTS SET REMARK = @remark ,REQ_STATUS = @status,
			UPDATED_AT = GETDATE(),UPDATED_BY = @code WHERE REQUEST_NO = @requestNo`,
				sql.Named("code", req.ActionBy),
				sql.Named("status", jobStatus),
				sql.Named("remark", req.Remark),
				sql.Named("requestNo", requestNo),
			)

			if errUpdateMainReq != nil {
				fmt.Println("Query failed: " + errUpdateMainReq.Error())
			}
		}
		fmt.Println("completeStep", completeStep)
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

			// Update Main Request Table
			_, errUpdateMainReq := db.Exec(`UPDATE TBL_REQUESTS SET REMARK = @remark ,REQ_STATUS = @status,
			UPDATED_AT = GETDATE(),UPDATED_BY = @code WHERE REQUEST_NO = @requestNo`,
				sql.Named("code", req.ActionBy),
				sql.Named("status", jobStatus),
				sql.Named("remark", req.Remark),
				sql.Named("requestNo", requestNo),
			)

			if errUpdateMainReq != nil {
				fmt.Println("Query failed: " + errUpdateMainReq.Error())
			}

			//SendMail to Requestor

		}

		// 3 is Done (Status Process Approve)
		if (completeStep > approved) && req.Status == 3 {
			fmt.Println("Update Approval Status and Send mail to Next")

			CheckSendEmail(iRev, requestNo)

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

	query := `WITH CTE_MASTER AS (
				SELECT ms.ID_STATUS_APV,ms.NAME_STATUS,COUNT(*) as AMOUNT FROM (
				SELECT DISTINCT m.REQUEST_NO,s.NAME_STATUS,g.NAME_GROUP,f.FACTORY_NAME,apv.CODE_APPROVER,m.ID_STATUS_APV FROM (
				SELECT a.* FROM TBL_APPROVAL a RIGHT JOIN (SELECT MAX(REV) as REV,REQUEST_NO  FROM TBL_APPROVAL GROUP BY  REQUEST_NO) b
				ON a.REQUEST_NO = b.REQUEST_NO AND a.REV = b.REV ) m
				LEFT JOIN TBL_REQUESTS r ON m.REQUEST_NO = r.REQUEST_NO
				LEFT JOIN TBL_FACTORY f ON r.ID_FACTORY = f.ID_FACTORY
				LEFT JOIN TBL_GROUP_DEPT g ON r.ID_GROUP_DEPT = g.ID_GROUP_DEPT
				LEFT JOIN TBL_APPROVERS apv ON r.ID_GROUP_DEPT = apv.ID_GROUP_DEPT 
				AND r.ID_FACTORY = apv.ID_FACTORY
				LEFT JOIN TBL_APPROVE_STATUS s ON m.ID_STATUS_APV = s.ID_STATUS_APV
				WHERE  apv.CODE_APPROVER = @code ) ms GROUP BY ms.ID_STATUS_APV,ms.NAME_STATUS )
				SELECT ast.NAME_STATUS,ast.ID_STATUS_APV,ISNULL(c.AMOUNT,0) as [AMOUNT] FROM CTE_MASTER  c 
				RIGHT JOIN TBL_APPROVE_STATUS ast
				ON  c.ID_STATUS_APV = ast.ID_STATUS_APV`

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
				ORDER BY REQUEST_NO`

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
	wc.NAME_WORKCELL,us.USERS as USERS_AMOUNT,mi.SUM_MINUTE,r.START_DATE,r.END_DATE,f.ID_FACTORY,ISNULL(pln.SUM_PLAN,0) as [SUM_PLAN],ISNULL(pln.SUM_PLAN_OB,0) as [SUM_PLAN_OB],wc.ID_WORK_CELL FROM  [dbo].[TBL_REQUESTS] r LEFT JOIN 
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

   WHERE r.REQUEST_NO = @reqNo
   ORDER BY r.REQUEST_NO DESC`

	results, errorQuery := db.Query(query, sql.Named("reqNo", reqNo)) //Query

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
