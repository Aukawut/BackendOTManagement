package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

func GetWorkCellByGroup(c *fiber.Ctx) error {
	var group = c.Params("group")
	var workCellAll []model.WorkCell

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

	results, errorQuery := db.Query(`SELECT ID_WORK_CELL,NAME_WORKCELL FROM TBL_WORKCELL wc WHERE wc.ID_WORKGRP = @id`, sql.Named("id", group))

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var work model.WorkCell

		err := results.Scan(&work.ID_WORKGRP, &work.NAME_WORKGRP)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			workCellAll = append(workCellAll, work)
		}

	}

	if len(workCellAll) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": workCellAll,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": workCellAll,
		})
	}

}

func GetWorkCellByAll(c *fiber.Ctx) error {
	var workCellAll []model.WorkCellJoinFactory

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

	results, errorQuery := db.Query(`SELECT wc.ID_WORK_CELL,wc.NAME_WORKCELL,f.ID_FACTORY,f.FACTORY_NAME FROM TBL_WORKCELL wc LEFT JOIN 
TBL_FACTORY f ON wc.ID_FACTORY = f.ID_FACTORY ORDER BY [ID_WORK_CELL] DESC`)

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var work model.WorkCellJoinFactory

		err := results.Scan(&work.ID_WORK_CELL, &work.NAME_WORKCELL, &work.ID_FACTORY, &work.FACTORY_NAME)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			workCellAll = append(workCellAll, work)
		}

	}

	if len(workCellAll) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": workCellAll,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": workCellAll,
		})
	}

}
