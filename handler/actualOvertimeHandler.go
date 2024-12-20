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
  WHERE OT_DATE BETWEEN '%s' AND '%s'%s%s  ORDER BY OT_DATE DESC`, start, end, conditionFactory, conditionUGroup)

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
