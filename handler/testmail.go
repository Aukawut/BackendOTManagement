package handler

import (
	"fmt"
	"os"

	"crypto/tls"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func TestingMail(c *fiber.Ctx) error {
	mail := c.Params("mail")
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	emailSender := fmt.Sprintf("Request OT <%s>", "it-system@prospira.local")

	subject := "แจ้งเตือน Request OT - PSTH"
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
				p {
					font-size: 18px;
					line-height: 5px;
				}
				th, td {
					border: 1px solid black;
					border-collapse: collapse;
					padding: 3px;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<p>12312312312</p>	
		</body>
		</html>
	`

	// SMTP server configuration
	smtpHost := "10.145.0.250"
	smtpPort := 25
	password := "Psth@min135"

	fmt.Println(password)
	fmt.Println(smtpHost)
	fmt.Println(smtpPort)
	fmt.Println(os.Getenv("MAIL_ADDRESS"))
	// Create a new message
	message := gomail.NewMessage()
	message.SetHeader("From", emailSender)
	message.SetHeader("To", mail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	// Create a new SMTP dialer with TLSConfig to skip certificate validation
	dialer := gomail.NewDialer(smtpHost, smtpPort, emailSender, password)
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Send mail error:", err)
		return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
	}

	return c.JSON(fiber.Map{"err": false, "msg": "Mail sent to " + mail})
}
