package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

func GetAllFactory(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()

	allFactory := []model.Factory{}

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
	rows, err := db.Query("SELECT ID_FACTORY,FACTORY_NAME FROM [dbo].[TBL_FACTORY] ORDER BY ID_FACTORY ASC")
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var factory model.Factory

		err := rows.Scan(&factory.ID_FACTORY, &factory.FACTORY_NAME)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			allFactory = append(allFactory, factory)
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

	if len(allFactory) > 0 {

		return c.JSON(fiber.Map{
			"err":     false,
			"results": allFactory,
			"status":  "Ok"})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": nil,
		})
	}

}

func GetAllFactoryByGroup(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()
	var id = c.Params("group")
	allFactory := []model.Factory{}

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
	rows, err := db.Query("SELECT ID_FACTORY,FACTORY_NAME FROM [dbo].[TBL_FACTORY] WHERE [ID_GROUP_DEPT] = @id ORDER BY ID_FACTORY ASC", sql.Named("id", id))
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var factory model.Factory

		err := rows.Scan(&factory.ID_FACTORY, &factory.FACTORY_NAME)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			allFactory = append(allFactory, factory)
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

	if len(allFactory) > 0 {

		return c.JSON(fiber.Map{
			"err":     false,
			"results": allFactory,
			"status":  "Ok"})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"results": nil,
		})
	}

}
