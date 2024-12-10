package auth

import (
	"fmt"
	"os"

	jwtAuth "gitgub.com/Aukawut/ServerOTManagement/auth/jwt"
	"gitgub.com/Aukawut/ServerOTManagement/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nmcclain/ldap"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthenticateUserDomain(username string, password string) (bool, string) {

	// Load ENV from .env file
	if err := godotenv.Load(); err != nil {
		return false, fmt.Sprintf("error loading .env file: %v", err)
	}
	// Get LDAP server from environment variable
	ldapServer := os.Getenv("LDAP_SERVER")

	// Make sure LDAP_SERVER is set
	if ldapServer == "" {
		return false, fmt.Sprintf("LDAP_SERVER not defined in .env file,%v", "")
	}

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, 389))
	if err != nil {

		return false, fmt.Sprintf("failed to connect to LDAP server: %v", err)
	}
	defer l.Close() // Close the connection once done

	err = l.Bind(username, password)
	if err != nil {

		return false, fmt.Sprintf("failed to authenticate user: %v", username)
	}

	// user is authenticated
	return true, ""
}

func LoginDomain(c *fiber.Ctx) error {
	var body LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Unable to parse JSON body",
		})
	}

	if body.Username == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Username and Password is required!",
		})
	}

	usernameAd := body.Username + os.Getenv("LDAP_DNS") // awk@psth.com

	// Domain Login
	verifiesDomain, errorMsg := AuthenticateUserDomain(usernameAd, body.Password)

	// Login success
	if verifiesDomain && errorMsg == "" {

		// Get Detail User and Permission
		userDetail := handler.GetPermissionByUsername(body.Username)

		token, _ := jwtAuth.GenerateToken(userDetail)

		return c.JSON(fiber.Map{
			"err":     false,
			"msg":     "Success",
			"status":  "Ok",
			"results": userDetail,
			"token":   token,
		})

	}

	// Login Failed
	return c.JSON(fiber.Map{
		"err": true,
		"msg": "Login failed!",
	})

}
