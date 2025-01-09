package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

func SaveActualOvertime(c *fiber.Ctx) error {
	var req model.BodyRequestSaveActual

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
			"res": err.Error(),
		})
	}

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

	stmtInsert := `INSERT INTO [dbo].[TBL_ACTUAL_OVERTIME] ([SCAN_IN],[SCAN_OUT],[OT_DATE],[SHIFT]
           ,[OT1_HOURS],[OT15_HOURS],[OT2_HOURS],[OT3_HOURS] ,[CREATED_BY],[CREATED_AT],[EMPLOYEE_CODE],[TOTAL_HOURS]) 
		   VALUES (@start,@end,@date,@shift,@ot1,@ot15,@ot2,@ot3,@action,GETDATE(),@code,@total)`

	for _, num := range req.Overtime {
		fmt.Println(req)
		_, errInsert := db.Exec(stmtInsert,
			sql.Named("start", num.Start),
			sql.Named("end", num.End),
			sql.Named("date", num.Date),
			sql.Named("shift", num.Shift),
			sql.Named("ot1", num.Overtime1),
			sql.Named("ot15", num.Overtime15),
			sql.Named("ot2", num.Overtime2),
			sql.Named("ot3", num.Overtime3),
			sql.Named("action", req.ActionBy),
			sql.Named("code", num.EmployeeCode),
			sql.Named("total", num.Total),
		)

		if errInsert != nil {
			fmt.Println(err.Error())
		}

	}
	// ลบข้อมูลซ้ำ
	_, errClear := db.Exec(`EXEC [dbo].[sProcDeleteDuplicateOvertimeRecords]`)

	if errClear != nil {
		fmt.Println(errClear.Error())
	}

	return c.JSON(fiber.Map{
		"err":    false,
		"msg":    "Actual inserted!",
		"status": "Ok",
	})
}

func GetActualOvertime(c *fiber.Ctx) error {
	var actualAll []model.AllActualOvertime

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

	stmt := `SELECT [Id]
      ,[EMPLOYEE_CODE]
      ,[SCAN_IN]
      ,[SCAN_OUT]
      ,[OT_DATE]
      ,[SHIFT]
      ,[OT1_HOURS]
      ,[OT15_HOURS]
      ,[OT2_HOURS]
      ,[OT3_HOURS]
      ,[TOTAL_HOURS]
      ,[UPDATED_AT]
      ,[CREATED_BY]
      ,[UPDATED_BY]
      ,[FACTORY_NAME]
      ,[NAME_UGROUP]
      ,[NAME_UTYPE]
  FROM [DB_OT_MANAGEMENT].[dbo].[V_ActualOvertime]`

	rows, errSelect := db.Query(stmt)

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.AllActualOvertime

		errorScan := rows.Scan(
			&actual.Id,
			&actual.EMPLOYEE_CODE,
			&actual.SCAN_IN,
			&actual.SCAN_OUT,
			&actual.OT_DATE,
			&actual.SHIFT,
			&actual.OT1_HOURS,
			&actual.OT15_HOURS,
			&actual.OT2_HOURS,
			&actual.OT3_HOURS,
			&actual.TOTAL_HOURS,
			&actual.UPDATED_AT,
			&actual.CREATED_BY,
			&actual.UPDATED_BY,
			&actual.FACTORY_NAME,
			&actual.NAME_UGROUP,
			&actual.NAME_UTYPE,
		)

		if errorScan != nil {
			fmt.Println(errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			actualAll = append(actualAll, actual)
		}
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"results": actualAll,
		"status":  "Ok",
	})
}

func GetActualByDate(c *fiber.Ctx) error {
	var actualAll []model.AllActualOvertime
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

	stmt := fmt.Sprintf(`SELECT [Id]
      ,[EMPLOYEE_CODE]
      ,[SCAN_IN]
      ,[SCAN_OUT]
      ,[OT_DATE]
      ,[SHIFT]
      ,[OT1_HOURS]
      ,[OT15_HOURS]
      ,[OT2_HOURS]
      ,[OT3_HOURS]
      ,[TOTAL_HOURS]
      ,[UPDATED_AT]
      ,[CREATED_BY]
      ,[UPDATED_BY]
      ,[FACTORY_NAME]
      ,[NAME_UGROUP]
      ,[NAME_UTYPE]
  FROM [DB_OT_MANAGEMENT].[dbo].[V_ActualOvertime] WHERE CONVERT(DATE,SCAN_IN) BETWEEN '%s' AND '%s' ORDER BY CONVERT(DATE,SCAN_IN) DESC`, start, end)

	rows, errSelect := db.Query(stmt)

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.AllActualOvertime

		errorScan := rows.Scan(
			&actual.Id,
			&actual.EMPLOYEE_CODE,
			&actual.SCAN_IN,
			&actual.SCAN_OUT,
			&actual.OT_DATE,
			&actual.SHIFT,
			&actual.OT1_HOURS,
			&actual.OT15_HOURS,
			&actual.OT2_HOURS,
			&actual.OT3_HOURS,
			&actual.TOTAL_HOURS,
			&actual.UPDATED_AT,
			&actual.CREATED_BY,
			&actual.UPDATED_BY,
			&actual.FACTORY_NAME,
			&actual.NAME_UGROUP,
			&actual.NAME_UTYPE,
		)

		if errorScan != nil {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			actualAll = append(actualAll, actual)
		}
	}

	return c.JSON(fiber.Map{
		"err":     false,
		"results": actualAll,
		"status":  "Ok",
	})
}

func SummaryActualComparePlan(c *fiber.Ctx) error {
	var actualAll []model.SummaryActualComparePlan
	year := c.Params("year")

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

	rows, errSelect := db.Query(`EXEC sProcSummaryActualByYear 1,@year`, sql.Named("year", year))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.SummaryActualComparePlan

		errorScan := rows.Scan(&actual.MONTH_NO, &actual.MONTH_NAME, &actual.MONTH, &actual.SUM_OT_ACTUAL, &actual.SUM_OT_PLANWC, &actual.SUM_OT_PLANOB)

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

func SummaryActualByDurationAndFac(c *fiber.Ctx) error {
	var actualAll []model.SummaryActualByFactory
	start := c.Params("start")
	end := c.Params("end")
	ugroup := c.Params("ugroup")

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

	rows, errSelect := db.Query(`DECLARE @table TABLE (
    ID_FACTORY INT,
    FACTORY_NAME NVARCHAR(255),
    SUM_PLAN_OB DECIMAL(18, 2),
    SUM_PLAN DECIMAL(18, 2),
    SUM_ACTUAL DECIMAL(18, 2),
    [YEAR] INT,
    [MONTH] INT
);
INSERT INTO @table
EXEC sProcGetMasterPlanActualReport  @start, @end, @ugroup;
SELECT t.ID_FACTORY,t.FACTORY_NAME,SUM(t.SUM_ACTUAL) as [SUM_ACTUAL],SUM(t.SUM_PLAN) as [SUM_PLAN],
SUM(t.SUM_PLAN_OB) as [SUM_PLAN_OB] FROM @table t GROUP BY 
t.ID_FACTORY,t.FACTORY_NAME`, sql.Named("start", start), sql.Named("end", end), sql.Named("ugroup", ugroup))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.SummaryActualByFactory

		errorScan := rows.Scan(&actual.ID_FACTORY, &actual.FACTORY_NAME, &actual.SUM_ACTUAL, &actual.SUM_PLAN, &actual.SUM_PLAN_OB)

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

func GetCountActualOvertime(c *fiber.Ctx) error {
	var count []model.CountActualOvertime
	start := c.Params("start")
	end := c.Params("end")
	ugroup := c.Params("ugroup")

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

	rows, errSelect := db.Query(`SELECT COUNT(*) as COUNT_OT FROM TBL_ACTUAL_OVERTIME a
LEFT JOIN (SELECT * FROM Func_GetTempTableAllEmployee()) hr ON a.EMPLOYEE_CODE = hr.EMPLOYEE_CODE_INT
WHERE hr.ID_UGROUP = @ugroup AND
CONVERT(DATE,SCAN_IN) BETWEEN @start AND @end
`, sql.Named("start", start), sql.Named("end", end), sql.Named("ugroup", ugroup))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.CountActualOvertime

		errorScan := rows.Scan(&actual.COUNT_OT)

		if errorScan != nil {
			fmt.Println("Error Scan : ", errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			count = append(count, actual)
		}
	}

	if len(count) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": count,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": count,
			"msg":     "Not Found",
		})
	}

}

func SummaryActualByDate(c *fiber.Ctx) error {
	var count []model.SummaryActualByDuration
	start := c.Params("start")
	end := c.Params("end")
	ugroup := c.Params("ugroup")

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

	rows, errSelect := db.Query(`SELECT CONVERT(DATE,SCAN_IN) as DATE_OT, SUM(TOTAL_HOURS) as SUM_TOTAL,DAY(CONVERT(DATE,SCAN_IN)) as DAY_OT,COUNT(*) as COUNT_OT FROM TBL_ACTUAL_OVERTIME a
LEFT JOIN (SELECT * FROM Func_GetTempTableAllEmployee()) hr ON a.EMPLOYEE_CODE = hr.EMPLOYEE_CODE_INT
WHERE hr.ID_UGROUP = @ugroup AND
CONVERT(DATE,SCAN_IN) BETWEEN @start AND @end
GROUP BY CONVERT(DATE,SCAN_IN)`, sql.Named("start", start), sql.Named("end", end), sql.Named("ugroup", ugroup))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.SummaryActualByDuration

		errorScan := rows.Scan(&actual.DATE_OT, &actual.SUM_TOTAL, &actual.DAY_OT, &actual.COUNT_OT)

		if errorScan != nil {
			fmt.Println("Error Scan : ", errorScan.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorScan.Error(),
			})
		} else {
			count = append(count, actual)
		}
	}

	if len(count) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": count,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": count,
			"msg":     "Not Found",
		})
	}

}

func SummaryActualOvertime(c *fiber.Ctx) error {
	var actualAll []model.OvertimeActual
	start := c.Params("start")
	end := c.Params("end")
	ugroup := c.Params("ugroup")
	fac := c.Params("fac")

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

	conditionFactory := ` AND ID_FACTORY = @factory`
	conditionUGroup := ` AND ID_UGROUP = @ugroup`

	// User isn't filter
	if fac == "0" {
		conditionFactory = ""
	}
	if ugroup == "0" {
		conditionUGroup = ""
	}

	stmt := fmt.Sprintf(`SELECT [EMPLOYEE_CODE]
      ,[OT_DATE]
      ,[SCAN_IN]
      ,[SCAN_OUT]
      ,[HOURS]
      ,[FACTORY_NAME]
      ,[NAME_UGROUP]
      ,[UHR_Department]
      ,[HOURS_AMOUNT]
      ,[NAME_UTYPE]
      ,[ID_FACTORY]
      ,[ID_UTYPE]
      ,[ID_UGROUP]
      ,[ID_TYPE_OT]
  FROM [DB_OT_MANAGEMENT].[dbo].[V_Actual_RowsFormat] 
  WHERE CONVERT(DATE,OT_DATE) BETWEEN '%s' AND '%s'%s%s  ORDER BY OT_DATE DESC`, start, end, conditionFactory, conditionUGroup)

	rows, errSelect := db.Query(stmt, sql.Named("start", start), sql.Named("end", end), sql.Named("ugroup", ugroup), sql.Named("factory", fac))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.OvertimeActual

		errorScan := rows.Scan(
			&actual.EMPLOYEE_CODE,
			&actual.OT_DATE,
			&actual.SCAN_IN,
			&actual.SCAN_OUT,
			&actual.HOURS,
			&actual.FACTORY_NAME,
			&actual.NAME_UGROUP,
			&actual.UHR_Department,
			&actual.HOURS_AMOUNT,
			&actual.NAME_UTYPE,
			&actual.ID_FACTORY,
			&actual.ID_UTYPE,
			&actual.ID_UGROUP,
			&actual.ID_TYPE_OT,
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

func SummaryActualOvertimeGroupFac(c *fiber.Ctx) error {

	var actualAll []model.OvertimeActualByFac
	start := c.Params("start")
	end := c.Params("end")
	ugroup := c.Params("ugroup")
	fac := c.Params("fac")

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

	conditionFactory := ` AND ID_FACTORY = @factory`
	conditionUGroup := ` AND ID_UGROUP = @ugroup`

	// User isn't filter
	if fac == "0" {
		conditionFactory = ""
	}
	if ugroup == "0" {
		conditionUGroup = ""
	}

	stmt := fmt.Sprintf(`SELECT f.FACTORY_NAME,f.ID_FACTORY,ISNULL(a.INLINE_HOURS,0) as [INLINE_HOURS],
ISNULL(a.OFFLINE_HOURS,0) as [OFFLINE_HOURS]   FROM TBL_FACTORY f LEFT JOIN (SELECT 
    FACTORY_NAME,
    ID_FACTORY,
    SUM(CASE WHEN NAME_UGROUP = 'INLINE' THEN HOURS ELSE 0 END) AS INLINE_HOURS,
    SUM(CASE WHEN NAME_UGROUP = 'OFFLINE' THEN HOURS ELSE 0 END) AS OFFLINE_HOURS
FROM [DB_OT_MANAGEMENT].[dbo].[V_Actual_RowsFormat]
WHERE CONVERT(DATE, OT_DATE) BETWEEN '%s' AND '%s'%s%s
GROUP BY FACTORY_NAME, ID_FACTORY) a ON f.ID_FACTORY = a.ID_FACTORY`, start, end, conditionFactory, conditionUGroup)

	rows, errSelect := db.Query(stmt, sql.Named("start", start), sql.Named("end", end), sql.Named("ugroup", ugroup), sql.Named("factory", fac))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.OvertimeActualByFac

		errorScan := rows.Scan(
			&actual.FACTORY_NAME,
			&actual.ID_FACTORY,
			&actual.INLINE_HOURS,
			&actual.OFFLINE_HOURS,
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

func SummaryActualOvertimeByType(c *fiber.Ctx) error {

	var actualAll []model.OvertimeActualByType
	start := c.Params("start")
	end := c.Params("end")
	ugroup := c.Params("ugroup")
	fac := c.Params("fac")

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

	conditionFactory := ` AND ID_FACTORY = @factory`
	conditionUGroup := ` AND ID_UGROUP = @ugroup`

	// User isn't filter
	if fac == "0" {
		conditionFactory = ""
	}
	if ugroup == "0" {
		conditionUGroup = ""
	}

	stmt := fmt.Sprintf(`SELECT t.ID_TYPE_OT,CONCAT('OT',t.HOURS_AMOUNT) as [HOURS_AMOUNT],
ISNULL(a.SUM_HOURS,0) as [SUM_HOURS] FROM TBL_OT_TYPE t
LEFT JOIN (
SELECT ID_TYPE_OT,SUM(HOURS) as SUM_HOURS,HOURS_AMOUNT FROM [DB_OT_MANAGEMENT].[dbo].[V_Actual_RowsFormat] 
WHERE CONVERT(DATE, OT_DATE) BETWEEN '%s' AND '%s'%s%s
GROUP BY ID_TYPE_OT,HOURS_AMOUNT) a ON a.ID_TYPE_OT = t.ID_TYPE_OT`, start, end, conditionFactory, conditionUGroup)

	rows, errSelect := db.Query(stmt, sql.Named("start", start), sql.Named("end", end), sql.Named("ugroup", ugroup), sql.Named("factory", fac))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.OvertimeActualByType

		errorScan := rows.Scan(
			&actual.ID_TYPE_OT,
			&actual.HOURS_AMOUNT,
			&actual.SUM_HOURS,
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

func SummaryActualOvertimeByDate(c *fiber.Ctx) error {

	var actualAll []model.OvertimeActualByDate
	start := c.Params("start")
	end := c.Params("end")
	ugroup := c.Params("ugroup")
	fac := c.Params("fac")

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

	conditionFactory := ` AND f.[ID_FACTORY] = @factory`
	conditionUGroup := ` AND ug.[ID_UGROUP] = @ugroup`

	// User isn't filter
	if fac == "0" {
		conditionFactory = ""
	}
	if ugroup == "0" {
		conditionUGroup = ""
	}

	stmt := fmt.Sprintf(`SELECT SUM(CASE WHEN  ug.NAME_UGROUP  = 'INLINE' THEN a.TOTAL_HOURS ELSE 0 END) AS INLINE_HOURS,
    SUM(CASE WHEN  ug.NAME_UGROUP = 'OFFLINE' THEN a.TOTAL_HOURS ELSE 0 END) AS OFFLINE_HOURS,
	CONVERT(DATE,SCAN_IN)  as DATE_OT
FROM TBL_ACTUAL_OVERTIME  a
LEFT JOIN (SELECT * FROM Func_GetTempTableAllEmployee()) u ON a.EMPLOYEE_CODE = u.EMPLOYEE_CODE_INT
LEFT JOIN TBL_FACTORY f ON u.ID_FACTORY = f.ID_FACTORY
LEFT JOIN TBL_UGROUP ug ON u.ID_UGROUP = ug.ID_UGROUP
WHERE (CONVERT(DATE,SCAN_IN) BETWEEN '%s' AND '%s' )%s%s
GROUP BY CONVERT(DATE,SCAN_IN) ORDER BY DATE_OT ASC`, start, end, conditionFactory, conditionUGroup)

	rows, errSelect := db.Query(stmt, sql.Named("ugroup", ugroup), sql.Named("factory", fac))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.OvertimeActualByDate

		errorScan := rows.Scan(
			&actual.INLINE_HOURS,
			&actual.OFFLINE_HOURS,
			&actual.DATE_OT,
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

func CalActualByFactory(c *fiber.Ctx) error {
	var actualAll []model.CalActualByFac
	year := c.Params("year")
	month := c.Params("month")
	fac := c.Params("fac")

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

	stmt := `SELECT SUM(TOTAL_HOURS) as SUM_HOURS FROM TBL_ACTUAL_OVERTIME a
LEFT JOIN (SELECT * FROM Func_GetTempTableAllEmployee()) hr ON a.EMPLOYEE_CODE = hr.EMPLOYEE_CODE_INT
WHERE MONTH(CONVERT(DATE,a.SCAN_IN)) = @month  AND hr.ID_FACTORY = @factory
AND YEAR(CONVERT(DATE,a.SCAN_IN)) = @year GROUP BY ID_FACTORY`

	rows, errSelect := db.Query(stmt, sql.Named("month", month), sql.Named("factory", fac), sql.Named("year", year))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.CalActualByFac

		errorScan := rows.Scan(
			&actual.SUM_HOURS,
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

func GetActualCompareWorkgroup(c *fiber.Ctx) error {
	var actualAll []model.ActualCompareWorkgroup
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

	stmt := fmt.Sprintf(`SELECT [OT_DATE],EMPLOYEE_CODE,NAME_WORKGRP,NAME_UGROUP,NAME_WORKCELL,[HOURS],[HOURS_AMOUNT] FROM V_Actual_CompareWorkGroup
WHERE OT_DATE BETWEEN '%s' AND '%s' ORDER BY OT_DATE DESC`, start, end)

	rows, errSelect := db.Query(stmt)

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.ActualCompareWorkgroup

		errorScan := rows.Scan(
			&actual.OT_DATE,
			&actual.EMPLOYEE_CODE,
			&actual.NAME_WORKGRP,
			&actual.NAME_UGROUP,
			&actual.NAME_WORKCELL,
			&actual.HOURS,
			&actual.HOURS_AMOUNT,
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

func GetActualCompareGroupWorkCell(c *fiber.Ctx) error {
	var actualAll []model.ActualGroupByWorkcell
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

	stmt := fmt.Sprintf(`SELECT SUM(HOURS) as SUM_HOURS,NAME_WORKCELL FROM V_Actual_CompareWorkGroup
WHERE OT_DATE BETWEEN '%s' AND '%s' 
GROUP BY NAME_WORKCELL
ORDER BY SUM(HOURS) DESC`, start, end)

	rows, errSelect := db.Query(stmt)

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.ActualGroupByWorkcell

		errorScan := rows.Scan(
			&actual.SUM_HOURS,
			&actual.NAME_WORKCELL,
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

func GetActualCompareGroupWorkGroup(c *fiber.Ctx) error {
	var actualAll []model.ActualGroupByWorkgroup
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

	stmt := fmt.Sprintf(`SELECT SUM(HOURS) as SUM_HOURS,NAME_WORKGRP FROM V_Actual_CompareWorkGroup
WHERE OT_DATE BETWEEN '%s' AND '%s' 
GROUP BY NAME_WORKGRP
ORDER BY SUM(HOURS) DESC`, start, end)

	rows, errSelect := db.Query(stmt)

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.ActualGroupByWorkgroup

		errorScan := rows.Scan(
			&actual.SUM_HOURS,
			&actual.NAME_WORKGRP,
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

func CalActualByWorkcell(c *fiber.Ctx) error {
	var actualAll []model.CalActualByWorkcell
	requestNo := c.Params("requestNo")
	rev := c.Params("rev")
	year := c.Params("year")
	month := c.Params("month")

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

	stmt := `SELECT SUM(HOURS) as SUM_HOURS,WORKCELL_ID
  FROM [DB_OT_MANAGEMENT].[dbo].[V_Actual_CompareWorkGroup] ag WHERE [WORKCELL_ID] IN (
  SELECT wc.ID_WORK_CELL FROM TBL_REQUESTS_HISTORY h LEFT JOIN TBL_WORKCELL wc
  ON wc.ID_WORK_CELL = h.ID_WORK_CELL WHERE REQUEST_NO = @requestNo AND REV = @rev
  AND YEAR(ag.OT_DATE) = @year AND MONTH(ag.OT_DATE) = @month 
  ) GROUP BY WORKCELL_ID`

	rows, errSelect := db.Query(stmt, sql.Named("requestNo", requestNo), sql.Named("rev", rev), sql.Named("year", year), sql.Named("month", month))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.CalActualByWorkcell

		errorScan := rows.Scan(
			&actual.SUM_HOURS,
			&actual.WORKCELL_ID,
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

func DeleteActualById(c *fiber.Ctx) error {

	id := c.Params("id")

	//user["employee_code"]
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

	stmt := `DELETE FROM [dbo].[TBL_ACTUAL_OVERTIME] WHERE [Id] = @id`

	_, errSelect := db.Exec(stmt, sql.Named("id", id))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	return c.JSON(fiber.Map{"err": false, "msg": "Deleted!", "status": "Ok"})

}

func SummaryActualByWorkcell(c *fiber.Ctx) error {
	var actualAll []model.CalActualByWorkcell

	year := c.Params("year")
	month := c.Params("month")
	idWorkcell := c.Params("idWorkcell")

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

	stmt := `SELECT SUM(HOURS) as [SUM_HOURS],[WORKCELL_ID] FROM  [dbo].[V_Actual_CompareWorkGroup] 
WHERE  YEAR(OT_DATE) = @y AND MONTH(OT_DATE) = @m AND WORKCELL_ID = @id
GROUP BY WORKCELL_ID`

	rows, errSelect := db.Query(stmt, sql.Named("y", year), sql.Named("m", month), sql.Named("id", idWorkcell))

	if errSelect != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errSelect.Error(),
		})
	}

	for rows.Next() {
		var actual model.CalActualByWorkcell

		errorScan := rows.Scan(
			&actual.SUM_HOURS,
			&actual.WORKCELL_ID,
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
