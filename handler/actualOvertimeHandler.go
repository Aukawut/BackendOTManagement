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
