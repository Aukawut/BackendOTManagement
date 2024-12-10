package handler

import (
	"database/sql"
	"fmt"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"

	_ "github.com/denisenkom/go-mssqldb"
)

func GetUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"err": false,
	})
}

func GetPermissionByUsername(username string) model.UserEncepyt {
	var usersGenToken model.UserEncepyt

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

	stmtQuery := `SELECT  m.EMPLOYEE_CODE,m.NAME_ROLE,m.FACTORY_NAME,m.NAME_GROUP,hr.UHR_FullName_en as FULLNAME,hr.AD_Mail as MAIL,hr.UHR_Department as REAL_DEPT,m.ID_FACTORY,m.ID_GROUP_DEPT FROM (
		SELECT a.EMPLOYEE_CODE,r.NAME_ROLE,f.FACTORY_NAME,g.NAME_GROUP,f.ID_FACTORY,g.ID_GROUP_DEPT FROM (
		SELECT * FROM [dbo].[TBL_EMPLOYEE] WHERE [EMPLOYEE_CODE] COLLATE Thai_CI_AS =
		(SELECT [UHR_EmpCode] COLLATE Thai_CI_AS  FROM [dbo].[V_AllUserPSTH] WHERE  [AD_UserLogon]= @username))  a
		LEFT JOIN [dbo].[TBL_PERMISSION] p ON  a.EMPLOYEE_CODE = p.EMPLOYEE_CODE
		LEFT JOIN [dbo].[TBL_FACTORY] f ON  p.ID_FACTORY = f.ID_FACTORY
		LEFT JOIN [dbo].[TBL_ROLE] r ON  p.ID_ROLE = r.ID_ROLE
		LEFT JOIN [dbo].[TBL_GROUP_DEPT] g ON  p.ID_GROUP_DEPT = g.ID_GROUP_DEPT

		) m
		LEFT JOIN [dbo].[V_AllUserPSTH] hr ON m.EMPLOYEE_CODE COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS`

	// Execute SELECT query
	rows, errQuery := db.Query(stmtQuery, sql.Named("username", username))
	if errQuery != nil {
		fmt.Println("Query failed: " + errQuery.Error())
	}
	defer rows.Close()

	var roles []model.RoleEncepyt
	// Iterate over the result set
	for rows.Next() {

		var role model.RoleEncepyt

		errScan := rows.Scan(
			&role.EMPLOYEE_CODE,
			&role.NAME_ROLE,
			&role.FACTORY_NAME,
			&role.NAME_GROUP,
			&role.FULLNAME,
			&role.MAIL,
			&role.REAL_DEPT,
			&role.ID_FACTORY,
			&role.ID_GROUP_DEPT,
		)

		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			roles = append(roles, role)
		}

	}

	defer rows.Close()

	usersGenToken.EmployeeCode = roles[0].EMPLOYEE_CODE
	usersGenToken.Department = roles[0].REAL_DEPT
	usersGenToken.Fullname = roles[0].FULLNAME
	usersGenToken.Email = roles[0].MAIL

	usersGenToken.Role = roles

	return usersGenToken
}

func GetUserInfoByGroupId(c *fiber.Ctx) error {
	info := []model.UsersInfo{}
	var idGroup = c.Params("id")
	var factory = c.Params("factory")

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

	queryUser := fmt.Sprintf(`SELECT e.EMPLOYEE_CODE,e.PREFIX + e.FNAME_TH+ ' ' + e.LNAME_TH as FULLNAME,f.FACTORY_NAME,r.NAME_ROLE,g.NAME_GROUP,hr.UHR_Position as POSITION FROM TBL_EMPLOYEE e
		 LEFT JOIN TBL_FACTORY f ON  e.ID_FACTORY = f.ID_FACTORY
		 LEFT JOIN TBL_ROLE  r ON e.ID_ROLE = r.ID_ROLE
		 LEFT JOIN TBL_GROUP_DEPT g ON e.GROUP_ID = g.ID_GROUP_DEPT 
		 LEFT JOIN V_AllUserPSTH hr ON e.EMPLOYEE_CODE COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS
		 WHERE g.ID_GROUP_DEPT = %s AND f.[ID_FACTORY] = %s
		 ORDER BY g.NAME_GROUP,e.EMPLOYEE_CODE`, idGroup, factory)

	results, errorQueryser := db.Query(queryUser)

	if errorQueryser != nil {
		fmt.Println("Query failed: " + errorQueryser.Error())
	}

	for results.Next() {
		var user model.UsersInfo

		errScan := results.Scan(&user.EMPLOYEE_CODE, &user.FULLNAME, &user.FACTORY_NAME, &user.NAME_ROLE, &user.NAME_GROUP, &user.POSITION)
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, user)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}

}

func GetUserByFactory(c *fiber.Ctx) error {
	info := []model.UsersInfo{}

	var factory = c.Params("factory")

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

	queryUser := fmt.Sprintf(`SELECT e.EMPLOYEE_CODE,e.PREFIX + e.FNAME_TH+ ' ' + e.LNAME_TH as FULLNAME,f.FACTORY_NAME,r.NAME_ROLE,g.NAME_GROUP,hr.UHR_Position as POSITION FROM TBL_EMPLOYEE e
		 LEFT JOIN TBL_FACTORY f ON  e.ID_FACTORY = f.ID_FACTORY
		 LEFT JOIN TBL_ROLE  r ON e.ID_ROLE = r.ID_ROLE
		 LEFT JOIN TBL_GROUP_DEPT g ON e.GROUP_ID = g.ID_GROUP_DEPT 
		 LEFT JOIN V_AllUserPSTH hr ON e.EMPLOYEE_CODE COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS
		 WHERE f.[ID_FACTORY] = %s
		 ORDER BY g.NAME_GROUP,e.EMPLOYEE_CODE`, factory)

	results, errorQueryser := db.Query(queryUser)

	if errorQueryser != nil {
		fmt.Println("Query failed: " + errorQueryser.Error())
	}

	for results.Next() {
		var user model.UsersInfo

		errScan := results.Scan(&user.EMPLOYEE_CODE, &user.FULLNAME, &user.FACTORY_NAME, &user.NAME_ROLE, &user.NAME_GROUP, &user.POSITION)
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, user)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}

}

func GetApproverByGroupId(c *fiber.Ctx) error {
	info := []model.Approver{}
	var idGroup = c.Params("id")
	var factory = c.Params("factory")

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

	queryUser := `SELECT a.CODE_APPROVER,a.NAME_APPROVER,a.ID_GROUP_DEPT,g.NAME_GROUP,a.ROLE,hr.UHR_Position,a.STEP,a.[ID_APPROVER] FROM TBL_APPROVERS  a  
				LEFT JOIN  V_AllUserPSTH hr ON a.CODE_APPROVER 
				COLLATE Thai_CI_AS = HR.UHR_EmpCode COLLATE Thai_CI_AS 
				LEFT JOIN TBL_GROUP_DEPT g ON a.ID_GROUP_DEPT = g.ID_GROUP_DEPT
				WHERE a.ID_GROUP_DEPT = @groupId AND a.STATUS_ACTIVE = 'Y'  AND a.ID_FACTORY = @factory
				ORDER BY a.STEP ASC`

	results, errorQueryser := db.Query(queryUser, sql.Named("groupId", idGroup), sql.Named("factory", factory))

	if errorQueryser != nil {
		fmt.Println("Query failed: " + errorQueryser.Error())
	}

	for results.Next() {
		var user model.Approver

		errScan := results.Scan(&user.CODE_APPROVER, &user.NAME_APPROVER, &user.ID_GROUP_DEPT, &user.NAME_GROUP, &user.ROLE, &user.UHR_Position, &user.STEP, &user.ID_APPROVER)
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, user)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}

}

func GetAllApprover(c *fiber.Ctx) error {
	info := []model.Approver{}

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

	queryUser := `SELECT a.CODE_APPROVER,a.NAME_APPROVER,a.ID_GROUP_DEPT,g.NAME_GROUP,a.ROLE,hr.UHR_Position,a.STEP,a.[ID_APPROVER] FROM TBL_APPROVERS  a  
				LEFT JOIN  V_AllUserPSTH hr ON a.CODE_APPROVER 
				COLLATE Thai_CI_AS = HR.UHR_EmpCode COLLATE Thai_CI_AS 
				LEFT JOIN TBL_GROUP_DEPT g ON a.ID_GROUP_DEPT = g.ID_GROUP_DEPT
				WHERE a.STATUS_ACTIVE = 'Y'
				ORDER BY a.STEP ASC`

	results, errorQueryser := db.Query(queryUser)

	if errorQueryser != nil {
		fmt.Println("Query failed: " + errorQueryser.Error())
	}

	for results.Next() {
		var user model.Approver

		errScan := results.Scan(&user.CODE_APPROVER, &user.NAME_APPROVER, &user.ID_GROUP_DEPT, &user.NAME_GROUP, &user.ROLE, &user.UHR_Position, &user.STEP, &user.ID_APPROVER)
		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, user)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}

}

func InsertApprover(c *fiber.Ctx) error {
	var req model.ApproverRequest

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

	// คำสั่ง SQL
	stmt := `INSERT INTO [dbo].[TBL_APPROVERS] 
		([CODE_APPROVER], [NAME_APPROVER], [ID_GROUP_DEPT], [ROLE], [STATUS_ACTIVE], [STEP],[ID_FACTORY], [CREATED_AT],[CREATED_BY])
		VALUES (@empCode, @name, @groupId, @roleId, @active, @step,@factory, GETDATE(),@createBy)`
	// Execute SQL statement]
	_, err = db.Exec(stmt,
		sql.Named("empCode", req.EmpCode),
		sql.Named("name", req.Name),
		sql.Named("groupId", req.GroupID),
		sql.Named("roleId", req.RoleID),
		sql.Named("active", "Y"),
		sql.Named("step", req.Step),
		sql.Named("factory", req.Factory),
		sql.Named("createBy", req.CreatedBy),
	)

	if err != nil {
		fmt.Println("Error executing query: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": true,
			"msg": err.Error(),
		})
	}

	// Response
	return c.JSON(fiber.Map{
		"err":    false,
		"msg":    "Approver added successfully",
		"status": "Ok",
	})

}

func UpdateApprover(c *fiber.Ctx) error {
	var req model.ApproverRequestUpdate
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

	// คำสั่ง SQL
	stmt := `UPDATE [dbo].[TBL_APPROVERS] SET [CODE_APPROVER] = @empCode, 
	[NAME_APPROVER] = @name, [ID_GROUP_DEPT] = @groupId, [ROLE] = @roleId, [STATUS_ACTIVE] = @active, [STEP] = @step
	,[UPDATED_AT] = GETDATE(),[UPDATED_BY] = @updateBy WHERE [ID_APPROVER] = @id`
	// Execute SQL statement
	_, err = db.Exec(stmt,
		sql.Named("empCode", req.EmpCode),
		sql.Named("name", req.Name),
		sql.Named("groupId", req.GroupID),
		sql.Named("roleId", req.RoleID),
		sql.Named("active", "Y"),
		sql.Named("step", req.Step),
		sql.Named("updateBy", req.UpdateBy),
		sql.Named("id", id),
	)

	if err != nil {
		fmt.Println("Error executing query: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": true,
			"msg": err.Error(),
		})
	}

	// Response
	return c.JSON(fiber.Map{
		"err":    false,
		"msg":    "Approver updated successfully",
		"status": "Ok",
	})

}

func GetUserByRequestNoAndRev(c *fiber.Ctx) error {
	info := []model.ListUserByRequest{}
	reqNo := c.Params("requestNo")
	rev := c.Params("rev")

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

	query := `SELECT EMPLOYEE_CODE,hr.UHR_FullName_th  as FULLNAME,hr.UHR_Position as POSITION FROM TBL_USERS_REQ u
LEFT JOIN V_AllUserPSTH hr ON 
u.EMPLOYEE_CODE COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS 
WHERE u.REQUEST_NO = @reqNo AND u.REV = @rev ORDER BY EMPLOYEE_CODE ASC`

	results, errorQuery := db.Query(query, sql.Named("reqNo", reqNo), sql.Named("rev", rev)) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.ListUserByRequest

		errScan := results.Scan(
			&result.EMPLOYEE_CODE,
			&result.FULLNAME,
			&result.POSITION,
		) // Scan เก็บข้อมูลใน Struct

		if errScan != nil {
			fmt.Println("Row scan failed: " + errScan.Error())

		} else {
			info = append(info, result)
		}
	}

	defer results.Close()

	if len(info) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"status":  "Ok",
			"results": info,
		})
	} else {
		return c.JSON(fiber.Map{
			"err":     true,
			"msg":     "Not Found",
			"results": info,
		})
	}
}
