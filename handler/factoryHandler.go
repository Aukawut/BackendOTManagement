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
	rows, err := db.Query(`SELECT ID_FACTORY,f.ID_GROUP_DEPT,FACTORY_NAME,f.CREATED_AT,g.NAME_GROUP FROM TBL_FACTORY f 
LEFT JOIN TBL_GROUP_DEPT g 
ON f.ID_GROUP_DEPT = g.ID_GROUP_DEPT
ORDER BY f.ID_FACTORY DESC`)
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var factory model.Factory

		err := rows.Scan(&factory.ID_FACTORY, &factory.ID_GROUP_DEPT, &factory.FACTORY_NAME, &factory.CREATED_AT, &factory.NAME_GROUP)
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

type ResponseFactory struct {
	ID_FACTORY int
}

func InsertFactory(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()
	var factoryCheck []ResponseFactory

	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating connection: " + err.Error())
	}
	defer db.Close()

	var req model.BodyInsertFactory

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
	}

	// Execute SELECT query
	rows, err := db.Query("SELECT [ID_FACTORY] FROM [dbo].[TBL_FACTORY] WHERE [ID_GROUP_DEPT] = @groupDept AND [FACTORY_NAME] = @facName",
		req.FactoryName, req.GroupDept,
	)
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var factory ResponseFactory

		err := rows.Scan(&factory.ID_FACTORY)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			factoryCheck = append(factoryCheck, factory)
		}

	}

	if len(factoryCheck) > 0 {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Factory Duplicated",
		})
	} else {
		_, errInsert := db.Exec(`INSERT INTO [dbo].[TBL_FACTORY] ([ID_GROUP_DEPT],[FACTORY_NAME],[CREATED_AT]) 
		VALUES (@group,@facName,GETDATE())`, sql.Named("group", req.GroupDept), sql.Named("group", req.GroupDept))

		if errInsert != nil {
			return c.JSON(fiber.Map{"err": true, "msg": errInsert.Error()})
		} else {
			return c.JSON(fiber.Map{"err": false, "msg": "Factory inserted", "status": "Ok"})
		}
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
