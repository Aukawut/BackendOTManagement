package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

func GetAllGroupWorkcell(c *fiber.Ctx) error {
	groupWorkCell := []model.GroupWorkCell{}

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

	results, errorQuery := db.Query(`SELECT  ID_WORKGRP,NAME_WORKGRP FROM [dbo].[TBL_WORK_GROUP] ORDER BY ID_WORKGRP ASC`)

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var group model.GroupWorkCell

		errScan := results.Scan(&group.ID_WORKGRP, &group.NAME_WORKGRP)
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			groupWorkCell = append(groupWorkCell, group)
		}
	}

	defer results.Close()

	if len(groupWorkCell) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": groupWorkCell,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": groupWorkCell,
		})
	}

}
