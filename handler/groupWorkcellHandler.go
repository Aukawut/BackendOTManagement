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

	results, errorQuery := db.Query(`SELECT  ID_WORKGRP,NAME_WORKGRP,[DESC] FROM [dbo].[TBL_WORK_GROUP] ORDER BY ID_WORKGRP ASC`)

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var group model.GroupWorkCell

		errScan := results.Scan(&group.ID_WORKGRP, &group.NAME_WORKGRP, &group.DESC)
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

func GetWorkcellByFactory(c *fiber.Ctx) error {
	workcellList := []model.Workcell{}
	id := c.Params("id")

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

	results, errorQuery := db.Query(`SELECT [ID_WORK_CELL],[NAME_WORKCELL] FROM TBL_WORKCELL WHERE ID_FACTORY = @id ORDER BY ID_WORK_CELL DESC`, sql.Named("id", id))

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var workcell model.Workcell

		errScan := results.Scan(&workcell.ID_WORK_CELL, &workcell.NAME_WORKCELL)
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			workcellList = append(workcellList, workcell)
		}
	}

	defer results.Close()

	if len(workcellList) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": workcellList,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": workcellList,
		})
	}

}
