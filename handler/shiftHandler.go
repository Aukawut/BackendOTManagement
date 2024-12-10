package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

func GetAllShift(c *fiber.Ctx) error {

	strConfig := config.LoadDatabaseConfig()

	shiftAll := []model.Shift{}

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
	rows, err := db.Query("SELECT [SHIFT_CODE],[START],[END],[STATUS_ACTIVE] FROM [DB_OT_MANAGEMENT].[dbo].[TBL_SHIFT] ORDER BY SHIFT_CODE ASC")
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var shift model.Shift

		err := rows.Scan(&shift.SHIFT_CODE, &shift.START, &shift.END, &shift.STATUS_ACTIVE)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			shiftAll = append(shiftAll, shift)
		}

	}

	if len(shiftAll) > 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": shiftAll,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"err": true,
			"msg": "Not Found",
		})
	}

}

func GetAllShiftActive(c *fiber.Ctx) error {
	status := c.Params("status")

	strConfig := config.LoadDatabaseConfig()

	shiftAll := []model.Shift{}

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
	rows, err := db.Query("SELECT [SHIFT_CODE],[START],[END],[STATUS_ACTIVE] FROM [DB_OT_MANAGEMENT].[dbo].[TBL_SHIFT] WHERE [STATUS_ACTIVE] = @status ORDER BY SHIFT_CODE ASC", sql.Named("status", status))
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var shift model.Shift

		err := rows.Scan(&shift.SHIFT_CODE, &shift.START, &shift.END, &shift.STATUS_ACTIVE)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			shiftAll = append(shiftAll, shift)
		}

	}

	if len(shiftAll) > 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": shiftAll,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"err": true,
			"msg": "Not Found",
		})
	}

}
