package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

func AddMainPlan(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()

	var req model.MainPlan
	var plan []string

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
	}

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
	rows, err := db.Query("SELECT ID_FACTORY FROM TBL_PLAN_OVERTIME WHERE ID_FACTORY = @factory AND [MONTH] = @m AND [YEAR] = @y",
		sql.Named("factory", req.FactoryID),
		sql.Named("m", req.Month),
		sql.Named("y", req.Year),
	)
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var oldPlan string

		err := rows.Scan(&oldPlan)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			plan = append(plan, oldPlan)
		}

	}
	if len(plan) > 0 {
		return c.JSON(fiber.Map{"err": true, "msg": "Plan is duplicated!"})
	}

	_, errInsert := db.Exec(`INSERT INTO TBL_PLAN_OVERTIME (ID_FACTORY,[MONTH],[YEAR],[HOURS],[CREATED_BY],[STATUS_ACTIVE]) VALUES (@factory,@m,@y,@amount,@action,'Y')`,
		sql.Named("factory", req.FactoryID),
		sql.Named("m", req.Month),
		sql.Named("y", req.Year),
		sql.Named("amount", req.Hours),
		sql.Named("action", req.ActionBy),
	)

	defer rows.Close()

	if errInsert != nil {
		fmt.Println("Insert error : ", errInsert.Error())
		return c.JSON(fiber.Map{"err": true, "msg": errInsert.Error()})
	} else {
		return c.JSON(fiber.Map{"err": false, "status": "Ok", "msg": "Plan added successfully!"})
	}

}

func GetAllMainPlan(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()

	var info []model.ResultMainPlan

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
	rows, err := db.Query(`SELECT mp.ID_PLAN,mp.ID_FACTORY,f.FACTORY_NAME,mp.CREATED_AT,mp.[MONTH],
	mp.[YEAR],mp.[HOURS],mp.UPDATED_AT,hr.UHR_FirstName_th as FNAME FROM TBL_PLAN_OVERTIME mp 
	LEFT JOIN TBL_FACTORY f ON mp.ID_FACTORY = f.ID_FACTORY 
	LEFT JOIN V_AllUserPSTH hr ON mp.CREATED_BY COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS
	WHERE mp.STATUS_ACTIVE = 'Y' 
	ORDER BY [YEAR],[MONTH],ID_FACTORY DESC`)

	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var result model.ResultMainPlan

		err := rows.Scan(&result.ID_PLAN, &result.ID_FACTORY, &result.FACTORY_NAME, &result.CREATED_AT, &result.MONTH, &result.YEAR, &result.HOURS, &result.UPDATED_AT, &result.FNAME)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			info = append(info, result)
		}

	}

	defer rows.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{"err": false, "results": info, "status": "Ok"})
	} else {
		return c.JSON(fiber.Map{"err": true, "results": info, "msg": "Not Found"})
	}

}

func UpdateMainPlan(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()
	id := c.Params("id")
	var req model.MainPlan
	var plan []string

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
	}

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
	rows, err := db.Query("SELECT ID_FACTORY FROM TBL_PLAN_OVERTIME WHERE ID_FACTORY = @factory AND [MONTH] = @m AND [YEAR] = @y AND [ID_PLAN] <> @id",
		sql.Named("factory", req.FactoryID),
		sql.Named("m", req.Month),
		sql.Named("y", req.Year),
		sql.Named("id", id),
	)
	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var oldPlan string

		err := rows.Scan(&oldPlan)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			plan = append(plan, oldPlan)
		}

	}
	if len(plan) > 0 {
		return c.JSON(fiber.Map{"err": true, "msg": "Plan is duplicated!"})
	}

	_, errInsert := db.Exec(`UPDATE TBL_PLAN_OVERTIME SET ID_FACTORY = @factory,[MONTH] = @m,[YEAR] = @y,[HOURS] =  @amount,[UPDATED_AT] = GETDATE(),[UPDATED_BY] = @action WHERE [ID_PLAN] = @id`,
		sql.Named("factory", req.FactoryID),
		sql.Named("m", req.Month),
		sql.Named("y", req.Year),
		sql.Named("amount", req.Hours),
		sql.Named("id", id),
		sql.Named("action", req.ActionBy),
	)

	defer rows.Close()

	if errInsert != nil {
		fmt.Println("Update error : ", errInsert.Error())
		return c.JSON(fiber.Map{"err": true, "msg": errInsert.Error()})
	} else {
		return c.JSON(fiber.Map{"err": false, "status": "Ok", "msg": "Plan updated successfully!"})
	}

}
