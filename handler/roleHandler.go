package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gofiber/fiber/v2"
)

func GetAllRole(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()

	roles := []model.Role{}

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
	rows, err := db.Query("SELECT ID_ROLE,NAME_ROLE,CREATED_AT,UPDATED_AT FROM [DB_OT_MANAGEMENT].[dbo].[TBL_ROLE] ORDER BY ID_ROLE ASC")
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var role model.Role

		err := rows.Scan(&role.ID_ROLE, &role.NAME_ROLE, &role.CREATED_AT, &role.UPDATED_AT)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			roles = append(roles, role)
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

	if len(roles) > 0 {

		return c.JSON(fiber.Map{
			"err":     false,
			"results": roles,
			"status":  "Ok"})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": nil,
		})
	}

}
