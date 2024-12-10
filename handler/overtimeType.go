package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

func GetOvertimeType(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()

	allOvertime := []model.Overtime{}

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
	rows, err := db.Query("SELECT [ID_TYPE_OT],[HOURS_AMOUNT],[CREATED_AT],[CREATED_BY] FROM [DB_OT_MANAGEMENT].[dbo].[TBL_OT_TYPE] ORDER BY [HOURS_AMOUNT] ASC")
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var overtime model.Overtime

		err := rows.Scan(&overtime.ID_TYPE_OT, &overtime.HOURS_AMOUNT, &overtime.CREATED_AT, &overtime.CREATED_AT)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			allOvertime = append(allOvertime, overtime)
		}

	}

	// Check for any error during iteration
	if err = rows.Err(); err != nil {
		fmt.Println("Row iteration error: " + err.Error())
		return c.JSON(fiber.Map{
			"err": true,
			"msg": err.Error(),
		})

	}

	if len(allOvertime) > 0 {

		return c.JSON(fiber.Map{
			"err":     false,
			"results": allOvertime,
			"status":  "Ok"})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": nil,
		})
	}

}
