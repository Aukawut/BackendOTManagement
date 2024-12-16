package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

func GetUserGroup(c *fiber.Ctx) error {
	info := []model.UserGroup{}

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

	query := `SELECT [ID_UGROUP], [NAME_UGROUP] FROM [dbo].[TBL_UGROUP]`

	results, errorQuery := db.Query(query) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.UserGroup

		errScan := results.Scan(
			&result.ID_UGROUP,
			&result.NAME_UGROUP,
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
