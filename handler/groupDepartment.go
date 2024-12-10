package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

func GetGroupDepartment(c *fiber.Ctx) error {

	strConfig := config.LoadDatabaseConfig()

	allGroup := []model.GroupDepartment{}

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
	rows, err := db.Query("SELECT ID_GROUP_DEPT,NAME_GROUP FROM [dbo].[TBL_GROUP_DEPT] ORDER BY ID_GROUP_DEPT ASC")
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var group model.GroupDepartment

		err := rows.Scan(&group.ID_GROUP_DEPT, &group.NAME_GROUP)

		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			allGroup = append(allGroup, group)
		}
	}
	return c.JSON(fiber.Map{
		"err":     false,
		"msg":     "Not Found",
		"results": allGroup,
	})
}

func GetGroupDepartmentByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	strConfig := config.LoadDatabaseConfig()

	allGroup := []model.GroupDepartment{}

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
	rows, err := db.Query("SELECT ID_GROUP_DEPT,NAME_GROUP FROM [dbo].[TBL_GROUP_DEPT] WHERE [STATUS_ACTIVE] = @status ORDER BY ID_GROUP_DEPT ASC",
		sql.Named("status", status))

	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}

	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var group model.GroupDepartment

		err := rows.Scan(&group.ID_GROUP_DEPT, &group.NAME_GROUP)

		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			allGroup = append(allGroup, group)
		}
	}
	if len(allGroup) > 0 {

		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": allGroup,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": allGroup,
		})
	}

}
