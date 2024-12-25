package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gitgub.com/Aukawut/ServerOTManagement/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateToken(user model.UserEncepyt) (string, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err) // Lowercase error message
	}

	// Retrieve the secret key from the environment variables
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		return "", fmt.Errorf("secret_key not set in .env file") // Lowercase error message
	}

	// Convert the JWT secret key to a byte slice
	secretKey := []byte(jwtSecret)

	// Create claims with user data
	claims := jwt.MapClaims{

		"employee_code": user.EmployeeCode,
		"role":          user.Role,
		"exp":           time.Now().Add(time.Hour * 3).Unix(), // Token expiration (3 Hours)
	}

	// Create the token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err) // Lowercase error message
	}

	return signedToken, nil
}

// VerifyToken verifies the provided JWT token and returns the claims
func VerifyToken(tokenString string) (jwt.MapClaims, error) {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// Retrieve the secret key from environment variables
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		return nil, fmt.Errorf("secret_key not set in .env file")
	}

	// Convert the JWT secret key to a byte slice
	secretKey := []byte(jwtSecret)

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	// Check if the token is valid
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if the token is valid and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func DecodeToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Authorization header must start with 'Bearer '",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	decoded, errToken := VerifyToken(token) // Ensure VerifyToken is implemented correctly

	if errToken != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": errToken.Error(),
		})
	}

	// Store decoded claims in context for downstream handlers
	c.Locals("user", decoded)
	return c.Next()
}

func DecodeTokenAdmin(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	isAdmin := 0

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": "Authorization header must start with 'Bearer '",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	decoded, errToken := VerifyToken(token) // Ensure VerifyToken is implemented correctly

	if errToken != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": true,
			"msg": errToken.Error(),
		})
	}

	roleData, ok := decoded["role"]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Role not found in user data",
		})
	}

	roles, ok := roleData.([]interface{})

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": true,
			"msg": "Invalid role data format",
		})
	}

	// วนลูปแสดงข้อมูลใน roles
	for _, r := range roles {
		role, ok := r.(map[string]interface{})

		if !ok {
			continue
		}
		if role["NAME_ROLE"] == "ADMIN" {
			isAdmin++
		}

	}
	if isAdmin > 0 {
		c.Locals("user", decoded)
		return c.Next()
	} else {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Token incorrect!",
		})
	}
	// Store decoded claims in context for downstream handlers

}

func CheckToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": "Authorization header must start with 'Bearer '",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	decoded, errToken := VerifyToken(token) // Ensure VerifyToken is implemented correctly

	if errToken != nil {
		return c.JSON(fiber.Map{
			"err": true,
			"msg": errToken.Error(),
		})
	}

	// Auth Success
	return c.JSON(fiber.Map{
		"err":     false,
		"msg":     "Auth success",
		"decoded": decoded,
	})

}
