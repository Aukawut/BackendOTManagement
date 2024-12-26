package handler

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"os"

	"gitgub.com/Aukawut/ServerOTManagement/config"
	"gitgub.com/Aukawut/ServerOTManagement/model"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendEMailToApprover(requestNo string, rev int, mail string, name string) bool {
	if mail == "" {
		return false
	}

	fmt.Println(mail)
	// Load database configuration
	strConfig := config.LoadDatabaseConfig()

	var users []model.UserBodyMail

	var TitleDescReq []model.RequestDetailBody

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Connect to SQL Server
	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating database connection: " + err.Error())
		return false
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error creating database connection: " + err.Error())
		return false
	}

	// SQL query
	stmt := `
		SELECT rh.REQUEST_NO, f.FACTORY_NAME, rh.REV, 
		FORMAT(CAST(rh.START_DATE AS DATETIME2), 'yyyy-MM-dd HH:mm') AS [START],
		FORMAT(CAST(rh.END_DATE AS DATETIME2), 'yyyy-MM-dd HH:mm') AS [END],
		CAST(DATEDIFF(MINUTE, rh.START_DATE, rh.END_DATE) / 60.00 AS DECIMAL(18, 2)) AS MINUTE_DIFF,
		s.NAME_STATUS,CONCAT('OT',t.HOURS_AMOUNT) as [OT_TYPE]
		FROM TBL_REQUESTS_HISTORY rh
		LEFT JOIN TBL_REQ_STATUS s ON rh.REQ_STATUS = s.ID_STATUS
		LEFT JOIN TBL_FACTORY f ON rh.ID_FACTORY = f.ID_FACTORY
		LEFT JOIN TBL_OT_TYPE t ON rh.ID_TYPE_OT = t.ID_TYPE_OT
		WHERE rh.REQUEST_NO = @requestNo AND rh.REV = @rev`

	stmtUser := `
		SELECT u.EMPLOYEE_CODE,hr.UHR_FullName_th as FULLNAME,hr.UHR_Department as DEPARTMENT FROM TBL_USERS_REQ u
LEFT JOIN (SELECT * FROM V_AllUserPSTH ) hr ON 
u.EMPLOYEE_CODE COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS
WHERE REQUEST_NO = @requestNo AND REV = @rev`

	// Execute the query
	headerRows, err := db.Query(stmt, sql.Named("requestNo", requestNo), sql.Named("rev", rev))
	if err != nil {
		fmt.Println("Query failed: " + err.Error())

		return false
	}
	defer headerRows.Close()

	// Parse query results
	for headerRows.Next() {
		var header model.RequestDetailBody
		err := headerRows.Scan(&header.REQUEST_NO, &header.FACTORY_NAME, &header.REV, &header.START, &header.END, &header.MINUTE_DIFF, &header.NAME_STATUS, &header.OT_TYPE)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())

		} else {
			TitleDescReq = append(TitleDescReq, header)

		}

	}

	// Build HTML table
	table := `<table
      style="width: 300px; border: 1px solid black; border-collapse: collapse">
		
		<tbody>`

	if len(TitleDescReq) > 0 {
		for _, v := range TitleDescReq {
			table += fmt.Sprintf(`
		<tr >
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">Request No.</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">Factory</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">Revise</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">Start</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">End</td>
        <td>%v</td>
      </tr>
	  <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">ประเภทโอที</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">ชั่วโมง / คน</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">สถานะ</td>
        <td>%v</td>
      </tr>
			
			`,
				v.REQUEST_NO, v.FACTORY_NAME, v.REV, v.START, v.END, v.OT_TYPE, v.MINUTE_DIFF, v.NAME_STATUS)
		}
	}

	table += `
		</tbody>
	</table>`

	// Execute the query
	rowUsers, errUser := db.Query(stmtUser, sql.Named("requestNo", requestNo), sql.Named("rev", rev))
	if errUser != nil {
		println("Query failed: " + errUser.Error())
		return false
	}
	defer rowUsers.Close()

	tableUsers := `  <table style="width: 500px; border: 1px solid black; border-collapse: collapse">
			<thead> 
        <tr style="background: #00A6B9;">
        
            <th style="font-weight: bold;">No.</th>
            <th style="font-weight: bold;">รหัสพนักงาน</th>
            <th style="font-weight: bold;">ชื่อ - สกุล</th>
            <th style="font-weight: bold;">หน่วยงาน</th>
        

        </tr>
      </thead>
		<tbody>`

	// Parse query results
	for rowUsers.Next() {
		var user model.UserBodyMail
		err := rowUsers.Scan(&user.EMPLOYEE_CODE, &user.FULLNAME, &user.DEPARTMENT)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return false
		}
		users = append(users, user)
	}

	if len(users) > 0 {
		for index, u := range users {
			tableUsers += fmt.Sprintf(`
			  <tr>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
        	  </tr>
			`, index+1, u.EMPLOYEE_CODE, u.FULLNAME, u.DEPARTMENT)

		}
	}
	tableUsers += `
		</tbody>
	</table>`

	emailSender := fmt.Sprintf("Request OT <%s>", os.Getenv("MAIL_ADDRESS"))

	to := "akawut.kamesuwan@prospira.com"

	subject := "แจ้งเตือน Request OT - PSTH"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Email</title>
			<style>
				body {
					font-family: 'Cordia New', sans-serif;
				}
				.custom-style {
					color: #333;
					font-size: 18px;
				}
					  p {
				font-size: 18px;
				line-height: 5px;
      					}
					th,
					td {
						border: 1px solid black;
						border-collapse: collapse;
						padding: 3px;
						text-align: center;
					}
		
			</style>
		</head>
		<body>
		<div>
      <p>
        เรียนคุณ %s
      </p>
      <p>
        กรุณาตรวจสอบคำขอหมายเลข : %s
      </p>
    </div>
    <h4 style="margin-bottom: 10px">รายละเอียดคำขอ</h4>
		<div>
				%s
		</div>
		<div>
	  <p>
       รายชื่อพนักงานที่เข้างาน ช่วง OT :  
      </p>
		%s

		<div><p> รวมเวลา : %v ชั่วโมง<p></div>

		<div><a href="http://localhost:5173/login">Go to Application</a></div>
		</div>
		</body>
		</html>
	`, name, TitleDescReq[0].REQUEST_NO, table, tableUsers, TitleDescReq[0].MINUTE_DIFF*float64(len(users)))

	// SMTP server configuration
	smtpHost := "10.145.0.250"

	smtpPort := 25
	password := os.Getenv("MAIL_PASSWORD")

	// Create a new message
	message := gomail.NewMessage()
	message.SetHeader("From", emailSender)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Skip certificate verification
	}

	// Create a new SMTP dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, emailSender, password)
	dialer.TLSConfig = tlsConfig
	// Send the email
	if err := dialer.DialAndSend(message); err != nil {

		fmt.Println("Send mail error:", err)
		return false

	}

	fmt.Println("Send sended to:", mail)
	return true

}

func SendEMailToRequestor(requestNo string, rev int, mail string, name string, status string) bool {
	if mail == "" {
		return false
	}

	fmt.Println(mail)
	// Load database configuration
	strConfig := config.LoadDatabaseConfig()

	var users []model.UserBodyMail

	var TitleDescReq []model.RequestDetailBody

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Connect to SQL Server
	db, err := sql.Open("sqlserver", strConfig)
	if err != nil {
		fmt.Println("Error creating database connection: " + err.Error())
		return false
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Error creating database connection: " + err.Error())
		return false
	}

	// SQL query
	stmt := `
		SELECT rh.REQUEST_NO, f.FACTORY_NAME, rh.REV, 
		FORMAT(CAST(rh.START_DATE AS DATETIME2), 'yyyy-MM-dd HH:mm') AS [START],
		FORMAT(CAST(rh.END_DATE AS DATETIME2), 'yyyy-MM-dd HH:mm') AS [END],
		CAST(DATEDIFF(MINUTE, rh.START_DATE, rh.END_DATE) / 60.00 AS DECIMAL(18, 2)) AS MINUTE_DIFF,
		s.NAME_STATUS,CONCAT('OT',t.HOURS_AMOUNT) as [OT_TYPE]
		FROM TBL_REQUESTS_HISTORY rh
		LEFT JOIN TBL_REQ_STATUS s ON rh.REQ_STATUS = s.ID_STATUS
		LEFT JOIN TBL_FACTORY f ON rh.ID_FACTORY = f.ID_FACTORY
		LEFT JOIN TBL_OT_TYPE t ON rh.ID_TYPE_OT = t.ID_TYPE_OT
		WHERE rh.REQUEST_NO = @requestNo AND rh.REV = @rev`

	stmtUser := `
		SELECT u.EMPLOYEE_CODE,hr.UHR_FullName_th as FULLNAME,hr.UHR_Department as DEPARTMENT FROM TBL_USERS_REQ u
LEFT JOIN (SELECT * FROM V_AllUserPSTH ) hr ON 
u.EMPLOYEE_CODE COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS
WHERE REQUEST_NO = @requestNo AND REV = @rev`

	// Execute the query
	headerRows, err := db.Query(stmt, sql.Named("requestNo", requestNo), sql.Named("rev", rev))
	if err != nil {
		fmt.Println("Query failed: " + err.Error())

		return false
	}
	defer headerRows.Close()

	// Parse query results
	for headerRows.Next() {
		var header model.RequestDetailBody
		err := headerRows.Scan(&header.REQUEST_NO, &header.FACTORY_NAME, &header.REV, &header.START, &header.END, &header.MINUTE_DIFF, &header.NAME_STATUS, &header.OT_TYPE)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())

		} else {
			TitleDescReq = append(TitleDescReq, header)

		}

	}

	// Build HTML table
	table := `<table
      style="width: 300px; border: 1px solid black; border-collapse: collapse">
		
		<tbody>`

	if len(TitleDescReq) > 0 {
		for _, v := range TitleDescReq {
			table += fmt.Sprintf(`
		<tr >
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">Request No.</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">Factory</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">Revise</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">Start</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">End</td>
        <td>%v</td>
      </tr>
	  <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">ประเภทโอที</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">ชั่วโมง / คน</td>
        <td>%v</td>
      </tr>
      <tr>
        <td style="font-weight: bold;background: #00A6B9;color: #011d21;width: 120px;">สถานะ</td>
        <td>%v</td>
      </tr>
			
			`,
				v.REQUEST_NO, v.FACTORY_NAME, v.REV, v.START, v.END, v.OT_TYPE, v.MINUTE_DIFF, status)
		}
	}

	table += `
		</tbody>
	</table>`

	// Execute the query
	rowUsers, errUser := db.Query(stmtUser, sql.Named("requestNo", requestNo), sql.Named("rev", rev))
	if errUser != nil {
		println("Query failed: " + errUser.Error())
		return false
	}
	defer rowUsers.Close()

	tableUsers := `  <table style="width: 500px; border: 1px solid black; border-collapse: collapse">
			<thead> 
        <tr style="background: #00A6B9;">
        
            <th style="font-weight: bold;">No.</th>
            <th style="font-weight: bold;">รหัสพนักงาน</th>
            <th style="font-weight: bold;">ชื่อ - สกุล</th>
            <th style="font-weight: bold;">หน่วยงาน</th>
        

        </tr>
      </thead>
		<tbody>`

	// Parse query results
	for rowUsers.Next() {
		var user model.UserBodyMail
		err := rowUsers.Scan(&user.EMPLOYEE_CODE, &user.FULLNAME, &user.DEPARTMENT)
		if err != nil {
			fmt.Println("Row scan failed: " + err.Error())
			return false
		}
		users = append(users, user)
	}

	if len(users) > 0 {
		for index, u := range users {
			tableUsers += fmt.Sprintf(`
			  <tr>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
				<td>%v</td>
        	  </tr>
			`, index+1, u.EMPLOYEE_CODE, u.FULLNAME, u.DEPARTMENT)

		}
	}
	tableUsers += `
		</tbody>
	</table>`

	emailSender := fmt.Sprintf("Request OT <%s>", os.Getenv("MAIL_ADDRESS"))

	to := "akawut.kamesuwan@prospira.com"

	subject := "แจ้งเตือน Request OT - PSTH"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Email</title>
			<style>
				body {
					font-family: 'Cordia New', sans-serif;
				}
				.custom-style {
					color: #333;
					font-size: 18px;
				}
					  p {
				font-size: 18px;
				line-height: 5px;
      					}
					th,
					td {
						border: 1px solid black;
						border-collapse: collapse;
						padding: 3px;
						text-align: center;
					}
		
			</style>
		</head>
		<body>
		<div>
      <p>
        เรียนคุณ %s
      </p>
      <p>
        คำขอของท่านสถานะ : %s
      </p>
    </div>
    <h4 style="margin-bottom: 10px">รายละเอียดคำขอ</h4>
		<div>
				%s
		</div>
		<div>
	  <p>
       รายชื่อพนักงานที่เข้างาน ช่วง OT :  
      </p>
		%s

		<div><p> รวมเวลา : %v ชั่วโมง<p></div>

		<div><a href="http://localhost:5173/login">Go to Application</a></div>
		</div>
		</body>
		</html>
	`, name, status, table, tableUsers, TitleDescReq[0].MINUTE_DIFF*float64(len(users)))

	// SMTP server configuration
	smtpHost := "10.145.0.250"

	smtpPort := 25
	password := os.Getenv("MAIL_PASSWORD")

	// Create a new message
	message := gomail.NewMessage()
	message.SetHeader("From", emailSender)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Skip certificate verification
	}

	// Create a new SMTP dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, emailSender, password)
	dialer.TLSConfig = tlsConfig
	// Send the email
	if err := dialer.DialAndSend(message); err != nil {

		fmt.Println("Send mail error:", err)
		return false

	}

	fmt.Println("Send sended to:", mail)
	return true

}
