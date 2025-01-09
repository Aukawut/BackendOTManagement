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

	if len(roles) > 0 {
		var approver []model.ApproverCheck
		fmt.Println(roles[0].EMPLOYEE_CODE)
		resultsApprover, errorResultsApprover := db.Query(`SELECT CODE_APPROVER,a.ID_FACTORY,g.ID_GROUP_DEPT,f.FACTORY_NAME,g.NAME_GROUP FROM TBL_APPROVERS  a
		LEFT JOIN TBL_FACTORY f ON a.ID_FACTORY = f.ID_FACTORY
		LEFT JOIN TBL_GROUP_DEPT g ON a.ID_GROUP_DEPT = g.ID_GROUP_DEPT
		WHERE CODE_APPROVER = @code`, sql.Named("code", roles[0].EMPLOYEE_CODE))

		if errorResultsApprover != nil {
			fmt.Println(errorResultsApprover.Error())
		} else {

			for resultsApprover.Next() {

				var approverInfo model.ApproverCheck

				errScanApproverInfo := resultsApprover.Scan(
					&approverInfo.CODE_APPROVER,
					&approverInfo.ID_FACTORY,
					&approverInfo.ID_GROUP_DEPT,
					&approverInfo.FACTORY_NAME,
					&approverInfo.NAME_GROUP,
				)

				if errScanApproverInfo != nil {
					fmt.Println("Row scan failed: " + errScanApproverInfo.Error())

				} else {
					approver = append(approver, approverInfo)
				}

			}

		}

		var roleUser model.RoleEncepyt
		if len(approver) > 0 {
			roleUser.EMPLOYEE_CODE = approver[0].CODE_APPROVER
			roleUser.ID_FACTORY = approver[0].ID_FACTORY
			roleUser.ID_GROUP_DEPT = approver[0].ID_GROUP_DEPT
			roleUser.FACTORY_NAME = approver[0].FACTORY_NAME
			roleUser.NAME_GROUP = approver[0].NAME_GROUP
			roleUser.MAIL = roles[0].MAIL
			roleUser.REAL_DEPT = roles[0].REAL_DEPT
			roleUser.NAME_ROLE = "APPROVER"
			roleUser.FULLNAME = roles[0].FULLNAME

			roles = append(roles, roleUser)
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

	queryUser := `SELECT a.CODE_APPROVER,a.NAME_APPROVER,a.ID_GROUP_DEPT,g.NAME_GROUP,a.ROLE,hr.UHR_Position,a.STEP,a.[ID_APPROVER],
f.FACTORY_NAME,hr.AD_Mail as MAIL,f.ID_FACTORY FROM TBL_APPROVERS  a  
				LEFT JOIN  V_AllUserPSTH hr ON a.CODE_APPROVER 
				COLLATE Thai_CI_AS = HR.UHR_EmpCode COLLATE Thai_CI_AS 
				LEFT JOIN TBL_GROUP_DEPT g ON a.ID_GROUP_DEPT = g.ID_GROUP_DEPT
				LEFT JOIN TBL_FACTORY f ON a.ID_FACTORY = f.ID_FACTORY
				WHERE a.ID_GROUP_DEPT = @groupId AND a.STATUS_ACTIVE = 'Y'  AND a.ID_FACTORY = @factory
				ORDER BY a.STEP DESC`

	results, errorQueryser := db.Query(queryUser, sql.Named("groupId", idGroup), sql.Named("factory", factory))

	if errorQueryser != nil {
		fmt.Println("Query failed: " + errorQueryser.Error())
	}

	for results.Next() {
		var user model.Approver

		errScan := results.Scan(&user.CODE_APPROVER, &user.NAME_APPROVER, &user.ID_GROUP_DEPT, &user.NAME_GROUP, &user.ROLE, &user.UHR_Position, &user.STEP, &user.ID_APPROVER, &user.FACTORY_NAME, &user.MAIL, &user.ID_FACTORY)
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
				ORDER BY a.STEP DESC`

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
	var approver []model.ApproverCheckDuplicated
	countDuplicated := 0

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
	// Check User is users of PSTH
	stmtCheckUser := `SELECT [UHR_FullName_en] FROM V_AllUserPSTH WHERE UHR_EmpCode = @code`

	rowsUser, errorSelectUser := db.Query(stmtCheckUser,
		sql.Named("code", req.EmpCode),
	)

	if errorSelectUser != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errorSelectUser.Error(),
		})
	}

	for rowsUser.Next() {

		var user model.ApproverCheckDuplicated

		if err := rowsUser.Scan(&user.FullName); err != nil {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": "Error scanning result: " + err.Error(),
			})
		} else {
			approver = append(approver, user)
		}
	}

	if len(approver) > 0 {

		stmtCheck := `SELECT COUNT(*) as [COUNT] FROM [DB_OT_MANAGEMENT].[dbo].[TBL_APPROVERS] 
	WHERE [ID_GROUP_DEPT] = @groupId AND [ID_FACTORY] = @factory AND STEP = @step`

		rows, errorSelect := db.Query(stmtCheck,
			sql.Named("groupId", req.GroupID),
			sql.Named("factory", req.Factory),
			sql.Named("step", req.Step),
		)

		if errorSelect != nil {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorSelect.Error(),
			})
		}

		for rows.Next() {

			if err := rows.Scan(&countDuplicated); err != nil {
				return c.JSON(fiber.Map{
					"err": true,
					"msg": "Error scanning result: " + err.Error(),
				})
			}
		}

		if countDuplicated > 0 {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": "Duplicate record found",
			})
		}

		// คำสั่ง SQL
		stmt := `INSERT INTO [dbo].[TBL_APPROVERS] 
		([CODE_APPROVER], [NAME_APPROVER], [ID_GROUP_DEPT], [ROLE], [STATUS_ACTIVE], [STEP],[ID_FACTORY], [CREATED_AT],[CREATED_BY])
		VALUES (@empCode, @name, @groupId, @roleId, @active, @step,@factory, GETDATE(),@createBy)`
		// Execute SQL statement]
		_, err = db.Exec(stmt,
			sql.Named("empCode", req.EmpCode),
			sql.Named("name", approver[0].FullName),
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

	} else {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "User isn't found.",
		})
	}

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

func GetUserType(c *fiber.Ctx) error {
	info := []model.UserType{}

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

	query := `SELECT [ID_UTYPE],[NAME_UTYPE] ,[CREATED_AT],[UPDATED_AT] FROM [DB_OT_MANAGEMENT].[dbo].[TBL_UTYPE]`

	results, errorQuery := db.Query(query) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.UserType

		errScan := results.Scan(
			&result.ID_UTYPE,
			&result.NAME_UTYPE,
			&result.CREATED_AT,
			&result.UPDATED_AT,
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

func GetEmployeeAll(c *fiber.Ctx) error {
	info := []model.Employee{}

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

	query := `SELECT e.EMPLOYEE_CODE,e.ID_FACTORY,e.GROUP_ID,e.PREFIX,e.FNAME_TH,e.LNAME_TH,e.FNAME_EN,e.LNAME_EN,e.ID_ROLE,e.TYPE_ID,e.ID_UGROUP,f.FACTORY_NAME,
  r.NAME_ROLE,g.NAME_GROUP,ug.NAME_UGROUP,ut.NAME_UTYPE,e.CREATED_AT,e.CREATED_BY,e.UPDATED_AT,e.UPDATED_BY FROM TBL_EMPLOYEE e
  LEFT JOIN TBL_FACTORY f ON e.ID_FACTORY = f.ID_FACTORY 
  LEFT JOIN TBL_ROLE r ON e.ID_ROLE = r.ID_ROLE 
  LEFT JOIN TBL_GROUP_DEPT g ON e.GROUP_ID = g.ID_GROUP_DEPT
  LEFT JOIN TBL_UGROUP ug ON e.ID_UGROUP = ug.ID_UGROUP
  LEFT JOIN TBL_UTYPE ut ON e.TYPE_ID = ut.ID_UTYPE ORDER BY f.FACTORY_NAME`

	results, errorQuery := db.Query(query) //Query

	if errorQuery != nil {
		fmt.Println("Query failed: " + errorQuery.Error())
	}

	for results.Next() {
		var result model.Employee

		errScan := results.Scan(
			&result.EMPLOYEE_CODE,
			&result.ID_FACTORY,
			&result.GROUP_ID,
			&result.PREFIX,
			&result.FNAME_TH,
			&result.LNAME_TH,
			&result.FNAME_EN,
			&result.LNAME_EN,
			&result.ID_ROLE,
			&result.TYPE_ID,
			&result.ID_UGROUP,
			&result.FACTORY_NAME,
			&result.NAME_ROLE,
			&result.NAME_GROUP,
			&result.NAME_UGROUP,
			&result.NAME_UTYPE,
			&result.CREATED_AT,
			&result.CREATED_BY,
			&result.UPDATED_AT,
			&result.UPDATED_BY,
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

func InsertEmployee(c *fiber.Ctx) error {
	var req model.BodyEmployee
	var CheckUser []model.EmployeeCode
	var employeeId string

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
	}

	if req.EmployeeCode == "" || req.Prefix == "" || req.FnameTH == "" || req.LnameTH == "" || req.FnameEN == "" || req.LnameEN == "" || req.ActionBy == "" {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Data is required!",
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

	user, errUser := db.Query(`SELECT UHR_EmpCode FROM [dbo].[V_AllUserPSTH] WHERE UHR_EmpCode = @code`, sql.Named("code", req.EmployeeCode))

	for user.Next() {
		errScanUser := user.Scan(&employeeId)
		if errScanUser != nil {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errScanUser.Error(),
			})
		}
	}

	defer user.Close()

	if errUser != nil {
		fmt.Println("Error executing query: " + err.Error())
		return c.JSON(fiber.Map{
			"err": true,
			"msg": err.Error(),
		})
	}

	if employeeId != "" {

		// คำสั่ง SQL
		stmt := `SELECT EMPLOYEE_CODE FROM [dbo].[TBL_EMPLOYEE] WHERE [EMPLOYEE_CODE] = @code`
		// Execute SQL statement
		rows, err := db.Query(stmt,
			sql.Named("code", req.EmployeeCode),
		)

		if err != nil {
			fmt.Println("Error executing query: " + err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err": true,
				"msg": err.Error(),
			})
		}

		for rows.Next() {
			var emp model.EmployeeCode
			errScan := rows.Scan(&emp.EMPLOYEE_CODE)
			if errScan != nil {
				fmt.Println("Error Scan", errScan.Error())
			} else {
				CheckUser = append(CheckUser, emp)
			}
		}

		if len(CheckUser) > 0 {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": "User is duplicated!",
			})
		}

		_, errorInsert := db.Exec(`INSERT INTO [dbo].[TBL_EMPLOYEE] ([EMPLOYEE_CODE],[ID_FACTORY],[GROUP_ID],[PREFIX],[FNAME_TH]
           ,[LNAME_TH],[FNAME_EN],[LNAME_EN],[ID_ROLE],[TYPE_ID],[ID_UGROUP],[CREATED_AT],[CREATED_BY]) 
		   VALUES (@code,@factory,@group,@prefix,@fnameTH,@lnameTH,@fnameEN,@lnameEN
		   ,@role,@type,@ugroup,GETDATE(),@action)`,
			sql.Named("code", req.EmployeeCode),
			sql.Named("factory", req.Factory),
			sql.Named("group", req.Group),
			sql.Named("prefix", req.Prefix),
			sql.Named("fnameTH", req.FnameTH),
			sql.Named("lnameTH", req.LnameTH),
			sql.Named("fnameEN", req.FnameEN),
			sql.Named("lnameEN", req.LnameEN),
			sql.Named("role", req.Role),
			sql.Named("type", req.Type),
			sql.Named("ugroup", req.Ugroup),
			sql.Named("action", req.ActionBy))

		if errorInsert != nil {
			return c.JSON(fiber.Map{
				"err": true,
				"msg": errorInsert.Error(),
			})

		}

		// Response
		return c.JSON(fiber.Map{
			"err":    false,
			"msg":    "Employee inserted successfully!",
			"status": "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Employee isn't found!",
		})
	}

}

func UpdateEmployee(c *fiber.Ctx) error {
	var req model.BodyUpdateEmployee
	code := c.Params("code")
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid request body",
		})
	}
	fmt.Println(req)
	fmt.Println(code)

	if req.ActionBy == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Data is required!",
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

	_, errorInsert := db.Exec(`UPDATE [dbo].[TBL_EMPLOYEE] SET [ID_FACTORY] = @factory,[GROUP_ID] = @group,
		[TYPE_ID] = @type,[ID_UGROUP] = @ugroup,[UPDATED_AT] = GETDATE(),[UPDATED_BY] = @action,[ID_ROLE] = @role WHERE [EMPLOYEE_CODE] = @code`,
		sql.Named("code", code),
		sql.Named("factory", req.Factory),
		sql.Named("group", req.Group),
		sql.Named("role", req.Role),
		sql.Named("type", req.Type),
		sql.Named("ugroup", req.Ugroup),
		sql.Named("action", req.ActionBy))

	if errorInsert != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errorInsert.Error(),
		})

	}

	// Response
	return c.JSON(fiber.Map{
		"err":    false,
		"msg":    "Employee updated successfully!",
		"status": "Ok",
	})

}

func DeleteEmployee(c *fiber.Ctx) error {

	code := c.Params("code")

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

	_, errorDelete := db.Exec(`DELETE FROM [dbo].[TBL_EMPLOYEE] WHERE [EMPLOYEE_CODE] = @code`,
		sql.Named("code", code))

	if errorDelete != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errorDelete.Error(),
		})

	}

	// Response
	return c.JSON(fiber.Map{
		"err":    false,
		"msg":    "Employee deleted successfully!",
		"status": "Ok",
	})

}

func GetEmployeeByCode(c *fiber.Ctx) error {

	code := c.Params("code")
	var detailUser []model.ResultEmployeeByCode

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

	userResult, errUser := db.Query(`SELECT UHR_Prefix_th as Prefix,UHR_EmpCode as EmployeeCode,UHR_FirstName_th as FnameTH,UHR_LastName_th as LnameTH,UHR_FirstName_en as FnameEN,UHR_LastName_en as LnameEN
  FROM [dbo].[V_AllUserPSTH] WHERE UHR_EmpCode = @code`, sql.Named("code", code))

	if errUser != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errUser.Error(),
		})
	}

	for userResult.Next() {
		var user model.ResultEmployeeByCode

		errScan := userResult.Scan(&user.Prefix, &user.EmployeeCode, &user.FnameTH, &user.LnameTH, &user.FnameEN, &user.LnameEN)
		if errScan != nil {
			fmt.Println(errScan.Error())

		} else {
			detailUser = append(detailUser, user)
		}
	}

	if len(detailUser) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": detailUser,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Employee isn't found!",
		})
	}

}

func GetAllEmployee(c *fiber.Ctx) error {

	code := c.Params("code")
	var detailUser []model.ResultEmployeeByCode

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

	userResult, errUser := db.Query(`SELECT UHR_Prefix_th as Prefix,UHR_EmpCode as EmployeeCode,UHR_FirstName_th as FnameTH,UHR_LastName_th as LnameTH,UHR_FirstName_en as FnameEN,UHR_LastName_en as LnameEN
  FROM [dbo].[V_AllUserPSTH] WHERE UHR_EmpCode = @code`, sql.Named("code", code))

	if errUser != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errUser.Error(),
		})
	}

	for userResult.Next() {
		var user model.ResultEmployeeByCode

		errScan := userResult.Scan(&user.Prefix, &user.EmployeeCode, &user.FnameTH, &user.LnameTH, &user.FnameEN, &user.LnameEN)
		if errScan != nil {
			fmt.Println(errScan.Error())

		} else {
			detailUser = append(detailUser, user)
		}
	}

	if len(detailUser) > 0 {
		return c.JSON(fiber.Map{
			"err":     false,
			"results": detailUser,
			"status":  "Ok",
		})
	} else {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Employee isn't found!",
		})
	}

}
