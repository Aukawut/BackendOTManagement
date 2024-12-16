package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
)

func InsertUserPermission(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()
	var permission []model.Permission

	var req model.BodyPermission

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

	row, errorRow := db.Query(`SELECT ID_PERMISSION,[EMPLOYEE_CODE],ID_FACTORY,ID_ROLE,ID_GROUP_DEPT,CREATED_AT,UPDATED_AT FROM [dbo].[TBL_PERMISSION] WHERE [EMPLOYEE_CODE] = @code AND [ID_FACTORY] = @factory 
	AND [ID_ROLE] = @role AND [ID_GROUP_DEPT] = @group`,

		sql.Named("code", req.EmployeeCode),
		sql.Named("factory", req.Factory),
		sql.Named("role", req.Role),
		sql.Named("group", req.GroupDept),
	)

	if errorRow != nil {
		fmt.Println("Insert error : ", errorRow.Error())
		return c.JSON(fiber.Map{"err": true, "msg": errorRow.Error()})
	}

	if row.Next() {
		var result model.Permission
		err := row.Scan(&result.ID_PERMISSION, &result.EMPLOYEE_CODE, &result.ID_FACTORY, &result.ID_ROLE, &result.ID_GROUP_DEPT, &result.CREATED_AT, &result.UPDATED_AT)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			permission = append(permission, result)
		}
	}

	if len(permission) > 0 {
		return c.JSON(fiber.Map{"err": true, "msg": "Permission is duplicated!", "results": permission})
	}
	_, errInsert := db.Exec(`INSERT INTO [dbo].[TBL_PERMISSION] ([EMPLOYEE_CODE],[ID_FACTORY],[ID_ROLE],[ID_GROUP_DEPT],[CREATED_AT])
	VALUES (@code,@factory,@role,@group,GETDATE())`,

		sql.Named("code", req.EmployeeCode),
		sql.Named("factory", req.Factory),
		sql.Named("role", req.Role),
		sql.Named("group", req.GroupDept),
	)

	if errInsert != nil {
		fmt.Println("Insert error : ", errInsert.Error())
		return c.JSON(fiber.Map{"err": true, "msg": errInsert.Error()})
	} else {
		return c.JSON(fiber.Map{"err": false, "status": "Ok", "msg": "Inserted successfully!"})
	}

}

func UpdateUserPermission(c *fiber.Ctx) error {
	strConfig := config.LoadDatabaseConfig()
	var permission []model.Permission

	id := c.Params("id")

	var req model.BodyPermission

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

	row, errorRow := db.Query(`SELECT ID_PERMISSION,[EMPLOYEE_CODE],ID_FACTORY,ID_ROLE,ID_GROUP_DEPT,CREATED_AT,UPDATED_AT FROM [dbo].[TBL_PERMISSION] WHERE [EMPLOYEE_CODE] = @code AND [ID_FACTORY] = @factory 
	AND [ID_ROLE] = @role AND [ID_GROUP_DEPT] = @group AND [ID_PERMISSION] <> @id`,

		sql.Named("code", req.EmployeeCode),
		sql.Named("factory", req.Factory),
		sql.Named("role", req.Role),
		sql.Named("group", req.GroupDept),
		sql.Named("id", id),
	)

	if errorRow != nil {
		fmt.Println("Select error : ", errorRow.Error())
		return c.JSON(fiber.Map{"err": true, "msg": errorRow.Error()})
	}

	if row.Next() {
		var result model.Permission
		err := row.Scan(&result.ID_PERMISSION, &result.EMPLOYEE_CODE, &result.ID_FACTORY, &result.ID_ROLE, &result.ID_GROUP_DEPT, &result.CREATED_AT, &result.UPDATED_AT)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return c.JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		} else {
			permission = append(permission, result)
		}
	}

	if len(permission) > 0 {
		return c.JSON(fiber.Map{"err": true, "msg": "Permission is duplicated!", "results": permission})
	}
	_, errInsert := db.Exec(`UPDATE [dbo].[TBL_PERMISSION] SET [EMPLOYEE_CODE] = @code,[ID_FACTORY] = @factory
	,[ID_ROLE] = @role,[ID_GROUP_DEPT] = @group,[UPDATED_AT] = GETDATE() WHERE [ID_PERMISSION] = @id`,

		sql.Named("code", req.EmployeeCode),
		sql.Named("factory", req.Factory),
		sql.Named("role", req.Role),
		sql.Named("group", req.GroupDept),
		sql.Named("id", id),
	)

	if errInsert != nil {
		fmt.Println("Update error : ", errInsert.Error())
		return c.JSON(fiber.Map{"err": true, "msg": errInsert.Error()})
	} else {
		return c.JSON(fiber.Map{"err": false, "status": "Ok", "msg": "Updated successfully!"})
	}

}

func DeleteUserPermission(c *fiber.Ctx) error {
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

	_, errDelete := db.Exec(`DELETE FROM [dbo].[TBL_PERMISSION] WHERE ID_PERMISSION = @id`,
		sql.Named("id", id),
	)

	if errDelete != nil {
		fmt.Println("Update error : ", errDelete.Error())
		return c.JSON(fiber.Map{"err": true, "msg": errDelete.Error()})
	} else {
		return c.JSON(fiber.Map{"err": false, "status": "Ok", "msg": "Deleted successfully!"})
	}

}
