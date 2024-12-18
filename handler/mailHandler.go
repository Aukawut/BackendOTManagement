package handler

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendMailToApprover(c *fiber.Ctx) error {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")

	}
	emailSender := fmt.Sprintf("Request OT <%s>", os.Getenv("MAIL_ADDRESS"))

	to := "akawut.kamesuwan@prospira.com"

	subject := "Test Email with HTML and Cordia New Font"
	body := `
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
			</style>
		</head>
		<body>
			<p class="custom-style">This is a test email using the Cordia New font!</p>
		</body>
		</html>
	`

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

	// Create a new SMTP dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, emailSender, password)

	stmt := `SELECT DISTINCT REQUEST_NO,g.NAME_GROUP,f.FACTORY_NAME,s.NAME_STATUS,h.REV,wc.NAME_WORKCELL,wg.NAME_WORKGRP,h.REMARK,
hr.UHR_FullName_th as REQUESTOR,CONVERT(date,h.START_DATE) as DATE_OT,ISNULL(pwc.SUM_HOURS,0) as PLAN_WORKCELL,
CONVERT(TIME,h.START_DATE) as TIME_START,CONVERT(TIME,h.END_DATE) as TIME_END,
CAST(DATEDIFF(MINUTE,h.START_DATE,h.END_DATE) / 60 as decimal(18,2)) as HOURS_TOTAL,
ISNULL(pfc.SUM_HOURS_FAC,0) as PLAN_FACTORY
FROM TBL_REQUESTS_HISTORY h
LEFT JOIN TBL_GROUP_DEPT g ON h.ID_GROUP_DEPT = g.ID_GROUP_DEPT
LEFT JOIN TBL_FACTORY f ON h.ID_FACTORY = f.ID_FACTORY
LEFT JOIN TBL_REQ_STATUS s ON h.REQ_STATUS = s.ID_STATUS
LEFT JOIN TBL_WORKCELL wc ON h.ID_WORK_CELL = wc.ID_WORK_CELL
LEFT JOIN TBL_WORK_GROUP wg ON wc.ID_WORKGRP = wg.ID_WORKGRP
LEFT JOIN V_AllUserPSTH hr ON h.CREATED_BY COLLATE Thai_CI_AS = hr.UHR_EmpCode COLLATE Thai_CI_AS
LEFT JOIN (
SELECT p.ID_WORK_CELL,w.NAME_WORKCELL ,p.[YEAR],p.MONTH,SUM(HOURS) as SUM_HOURS
FROM TBL_PLAN_OVERTIME p LEFT JOIN TBL_WORKCELL w ON p.ID_WORK_CELL = w.ID_WORK_CELL
LEFT JOIN TBL_FACTORY f ON w.ID_FACTORY = f.ID_FACTORY GROUP BY 
p.ID_WORK_CELL,w.NAME_WORKCELL,p.YEAR,p.MONTH
) pwc ON h.ID_WORK_CELL = pwc.ID_WORK_CELL AND YEAR(h.START_DATE) = pwc.YEAR AND MONTH(h.START_DATE) = pwc.MONTH

LEFT JOIN (
SELECT f.ID_FACTORY,f.FACTORY_NAME,p.YEAR,p.MONTH,SUM(HOURS) as SUM_HOURS_FAC
FROM TBL_PLAN_OVERTIME p LEFT JOIN TBL_WORKCELL w ON p.ID_WORK_CELL = w.ID_WORK_CELL
LEFT JOIN TBL_FACTORY f ON w.ID_FACTORY = f.ID_FACTORY GROUP BY 
f.ID_FACTORY,f.FACTORY_NAME,p.YEAR,p.MONTH
) pfc ON h.ID_FACTORY = pfc.ID_FACTORY AND YEAR(h.START_DATE) = pfc.YEAR AND MONTH(h.START_DATE) = pfc.MONTH
WHERE REQUEST_NO = 'RQ202412100001' AND REV = 1`

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println(err)
		return c.JSON(fiber.Map{
			"err":  true,
			"msg":  err,
			"stmt": stmt,
		})

	}

	return c.JSON(fiber.Map{
		"err": false,
		"msg": "Email sent successfully!",
	})

}
