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

	results, errorQuery := db.Query(`SELECT ID_WORKGRP,NAME_WORKGRP,[DESC] FROM [dbo].[TBL_WORK_GROUP] ORDER BY ID_WORKGRP ASC`)

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

func InsertWorkcell(c *fiber.Ctx) error {
	var req model.ReqWorkCellBody

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

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

	query := `SELECT COUNT(*) FROM TBL_WORKCELL WHERE ID_WORKGRP = @work AND ID_FACTORY = @factory AND NAME_WORKCELL = @nameWorkcell`
	var count int
	err = db.QueryRow(query,
		sql.Named("work", req.ID_WORKGRP),
		sql.Named("factory", req.ID_FACTORY),
		sql.Named("nameWorkcell", req.NAME_WORKCEL),
	).Scan(&count)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": true,
			"msg": "Error checking workcell: " + err.Error(),
		})
	}

	if count > 0 {
		// Workcell already exists
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Workcell already exists",
		})
	} else {

		// Insert
		_, errInsert := db.Exec(`
		INSERT INTO [dbo].[TBL_WORKCELL] ([NAME_WORKCELL],[ID_WORKGRP],[ID_FACTORY],[CREATED_AT]) VALUES (@name,@workgroup,@factory,GETDATE())
		`,
			sql.Named("name", req.NAME_WORKCEL),
			sql.Named("workgroup", req.ID_WORKGRP),
			sql.Named("factory", req.ID_FACTORY),
		)

		if errInsert != nil {
			return c.JSON(fiber.Map{"err": true, "msg": errInsert.Error()})
		}

		return c.JSON(fiber.Map{"err": false, "msg": "Workcell Inserted", "status": "Ok"})

	}

}

func UpdateWorkcell(c *fiber.Ctx) error {
	var req model.ReqWorkCellBody
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

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

	query := `SELECT COUNT(*) FROM TBL_WORKCELL WHERE ID_WORKGRP = @work AND ID_FACTORY = @factory AND NAME_WORKCELL = @nameWorkcell AND [ID_WORK_CELL] <> @id`
	var count int
	err = db.QueryRow(query,
		sql.Named("work", req.ID_WORKGRP),
		sql.Named("factory", req.ID_FACTORY),
		sql.Named("nameWorkcell", req.NAME_WORKCEL),
		sql.Named("id", id),
	).Scan(&count)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": true,
			"msg": "Error checking workcell: " + err.Error(),
		})
	}

	if count > 0 {
		// Workcell already exists
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Workcell already exists",
		})
	} else {

		// Update
		_, errUpdate := db.Exec(`
		UPDATE [dbo].[TBL_WORKCELL] SET [NAME_WORKCELL] = @name,[ID_WORKGRP] = @workgroup,[ID_FACTORY] = @factory,[UPDATED_AT] = GETDATE() WHERE [ID_WORK_CELL] = @id
		`,
			sql.Named("name", req.NAME_WORKCEL),
			sql.Named("workgroup", req.ID_WORKGRP),
			sql.Named("factory", req.ID_FACTORY),
			sql.Named("id", id),
		)

		if errUpdate != nil {
			fmt.Println(errUpdate.Error())
			return c.JSON(fiber.Map{"err": true, "msg": errUpdate.Error()})
		}
		defer db.Close()

		return c.JSON(fiber.Map{"err": false, "msg": "Workcell Updated", "status": "Ok"})

	}

}

func DeleteWorkcell(c *fiber.Ctx) error {
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

	query := `SELECT COUNT(*) FROM TBL_WORKCELL WHERE [ID_WORK_CELL] = @id`
	var count int
	err = db.QueryRow(query,
		sql.Named("id", id),
	).Scan(&count)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": true,
			"msg": "Error checking workcell: " + err.Error(),
		})
	}

	if count > 0 {

		// Delete
		_, errDelete := db.Exec(`
		DELETE FROM [dbo].[TBL_WORKCELL] WHERE [ID_WORK_CELL] = @id
		`,
			sql.Named("id", id),
		)

		if errDelete != nil {
			fmt.Println(errDelete.Error())
			return c.JSON(fiber.Map{"err": true, "msg": errDelete.Error()})
		}

		defer db.Close()

		return c.JSON(fiber.Map{"err": false, "msg": "Workcell Deleted", "status": "Ok"})
	} else {

		return c.JSON(fiber.Map{"err": true, "msg": "Workcell isn't found."})
	}

}
