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
	rows, err := db.Query("SELECT ID_WORK_CELL FROM TBL_PLAN_OVERTIME WHERE ID_WORK_CELL = @work AND [MONTH] = @m AND [YEAR] = @y AND [STATUS_ACTIVE] = 'Y'",
		sql.Named("work", req.WorkcellID),
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

	_, errInsert := db.Exec(`INSERT INTO TBL_PLAN_OVERTIME ([ID_WORK_CELL],[MONTH],[YEAR],[HOURS],[CREATED_BY],[STATUS_ACTIVE],[TYPE]) VALUES (@work,@m,@y,@amount,@action,'Y',@type)`,

		sql.Named("work", req.WorkcellID),
		sql.Named("m", req.Month),
		sql.Named("y", req.Year),
		sql.Named("amount", req.Hours),
		sql.Named("action", req.ActionBy),
		sql.Named("type", req.UserGroup),
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
	rows, err := db.Query(`SELECT mp.ID_PLAN,mp.ID_WORK_CELL,wc.NAME_WORKCELL,f.FACTORY_NAME,ug.NAME_UGROUP,ug.ID_UGROUP,mp.CREATED_AT,mp.[MONTH],
	mp.[YEAR],mp.[HOURS],mp.UPDATED_AT,hr.UHR_FirstName_th as FNAME FROM TBL_PLAN_OVERTIME mp 
	LEFT JOIN TBL_WORKCELL wc ON mp.ID_WORK_CELL = wc.ID_WORK_CELL
	LEFT JOIN TBL_FACTORY f ON wc.ID_FACTORY = f.ID_FACTORY 
	LEFT JOIN V_AllUserPSTH hr ON mp.CREATED_BY COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS
	LEFT JOIN TBL_UGROUP ug ON mp.TYPE = ug.ID_UGROUP
	WHERE mp.STATUS_ACTIVE = 'Y' 
	ORDER BY [YEAR],[MONTH],f.ID_FACTORY DESC`)

	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var result model.ResultMainPlan

		err := rows.Scan(&result.ID_PLAN, &result.ID_WORK_CELL, &result.NAME_WORKCELL, &result.FACTORY_NAME, &result.NAME_UGROUP, &result.ID_UGROUP, &result.CREATED_AT, &result.MONTH, &result.YEAR, &result.HOURS, &result.UPDATED_AT, &result.FNAME)
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
	rows, err := db.Query("SELECT ID_WORK_CELL FROM TBL_PLAN_OVERTIME WHERE ID_WORK_CELL = @work AND [MONTH] = @m AND [YEAR] = @y AND [ID_PLAN] <> @id",
		sql.Named("work", req.WorkcellID),
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

	_, errInsert := db.Exec(`UPDATE TBL_PLAN_OVERTIME SET [ID_WORK_CELL] = @work,[MONTH] = @m,[YEAR] = @y,[HOURS] =  @amount,[UPDATED_AT] = GETDATE(),[UPDATED_BY] = @action,[TYPE] =  @type WHERE [ID_PLAN] = @id`,
		sql.Named("work", req.WorkcellID),
		sql.Named("m", req.Month),
		sql.Named("y", req.Year),
		sql.Named("amount", req.Hours),
		sql.Named("id", id),
		sql.Named("action", req.ActionBy),
		sql.Named("type", req.UserGroup),
	)

	defer rows.Close()

	if errInsert != nil {
		fmt.Println("Update error : ", errInsert.Error())
		return c.JSON(fiber.Map{"err": true, "msg": errInsert.Error()})
	} else {
		return c.JSON(fiber.Map{"err": false, "status": "Ok", "msg": "Plan updated successfully!"})
	}

}

func GetPlanByFactory(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()
	year := c.Params("year")
	month := c.Params("month")
	id := c.Params("id")

	var info []model.ResultPlanByFactory

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
	rows, err := db.Query(`SELECT [YEAR],[MONTH],f.ID_FACTORY,SUM(p.HOURS) as SUM_HOURS FROM TBL_PLAN_OVERTIME p 
LEFT JOIN TBL_WORKCELL wc ON p.ID_WORK_CELL = wc.ID_WORK_CELL
LEFT JOIN TBL_FACTORY f ON wc.ID_FACTORY  = f.ID_FACTORY
WHERE  [YEAR] = @year AND [MONTH] = @month AND f.ID_FACTORY = @factory
GROUP BY [YEAR],[MONTH],f.ID_FACTORY`, sql.Named("year", year), sql.Named("month", month), sql.Named("factory", id))

	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var result model.ResultPlanByFactory

		err := rows.Scan(&result.YEAR, &result.MONTH, &result.ID_FACTORY, &result.SUM_HOURS)
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

func GetPlanByWorkcell(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()
	year := c.Params("year")
	month := c.Params("month")
	id := c.Params("id")

	var info []model.ResultPlanByFactory

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
	rows, err := db.Query(`SELECT [YEAR],[MONTH],f.ID_FACTORY,SUM(p.HOURS) as SUM_HOURS FROM TBL_PLAN_OVERTIME p 
LEFT JOIN TBL_WORKCELL wc ON p.ID_WORK_CELL = wc.ID_WORK_CELL
LEFT JOIN TBL_FACTORY f ON wc.ID_FACTORY  = f.ID_FACTORY
WHERE  [YEAR] = @year AND [MONTH] = @month AND p.ID_WORK_CELL = @work
GROUP BY [YEAR],[MONTH],f.ID_FACTORY`, sql.Named("year", year), sql.Named("month", month), sql.Named("work", id))

	if err != nil {
		fmt.Println("Query failed: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var result model.ResultPlanByFactory

		err := rows.Scan(&result.YEAR, &result.MONTH, &result.ID_FACTORY, &result.SUM_HOURS)
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

func DeletePlan(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()
	id := c.Params("id")

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
	_, errUpdate := db.Exec(`UPDATE [TBL_PLAN_OVERTIME] SET [STATUS_ACTIVE] = 'N' WHERE [ID_PLAN] = @id`, sql.Named("id", id))

	if errUpdate != nil {
		fmt.Println("Execute failed: " + errUpdate.Error())
	}

	return c.JSON(fiber.Map{
		"err":    false,
		"msg":    "Plan Deleted Succesfully!",
		"status": "Ok",
	})

}
