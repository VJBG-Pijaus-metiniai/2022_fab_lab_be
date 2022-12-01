package controllers

import (
	"fablab-project/database"
	"fablab-project/models"
	"fablab-project/utils"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type RegisterKey struct {
	Key string `json:"register_key"`
}

type registerRequest struct {
	Secret string `json:"secret"`
	Password string `json:"password"`
	Email string `json:"email"`
	Name string `json:"name"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func createJWTTOken(user models.User) (string, time.Time, error) {
	exp := time.Now().Add(time.Hour * 30)
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = exp.Unix()
	t, err := token.SignedString([]byte(os.Getenv("SECRET_JWT")))

	if err != nil {
	  return "", time.Now(), err
	}

	return t, exp, nil
}

func getSecret() string {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

	return os.Getenv("SECRET_BETA")
}

func Register(c *fiber.Ctx) error {
	var req registerRequest
	secret := getSecret()
	var user models.User
	validUser := models.User{}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Couldn't parse request, please try again later");
	}

	if (req.Secret != secret) {
		return fiber.NewError(fiber.StatusBadRequest, "Wrong secret open beta key, contact your supervisor to get help")
	}

	if (req.Name == "" || req.Password == "" || req.Secret == "") {
		return fiber.NewError(fiber.StatusBadRequest, "Every field is required")
	}

	db := database.Database.DB
	db.Where("email = ?", req.Email).First(&validUser)
	if (validUser.ID != 0) {
		return fiber.NewError(fiber.StatusBadRequest, "Email already taken")
	}
	db.Where("name = ?", req.Name).First(&validUser)
	if (validUser.ID != 0) {
		return fiber.NewError(fiber.StatusBadRequest, "Name already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"message": "Couldn't generate hash",
		})

		return fiber.NewError(fiber.StatusBadGateway, "Couldn't generate hash")
	}

	user.Email = req.Email
	user.Name = req.Name
	user.Password = string(hash)

	db.Create(&user)

	token, exp, err := createJWTTOken(user)
	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"message": "Couldn't create token",
		})

		return fiber.NewError(fiber.StatusBadGateway, "Couldn't create token")
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    string(token),
		Expires:  exp,
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	c.Status(200).JSON(user)

	return nil
}

func Login(c *fiber.Ctx) error {
	db := database.Database.DB
	req := new(LoginRequest)
	var user models.User

	if err := c.BodyParser(&req); err != nil {
		c.Status(400).JSON(fiber.Map{"error": err.Error()})
		return fiber.NewError(fiber.StatusBadRequest, "Couldn't parse request body")
	}

	if req.Password == "" || req.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid sign in credentials")
	}

	db.Where("email = ?", req.Email).First(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid sign in credentials")
	}
	
	token, exp, err := createJWTTOken(user)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, "Couldn't create token")
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    string(token),
		Expires:  exp,
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	c.Status(200).JSON(user)

	return nil
}

func GetCurrentUser(c *fiber.Ctx) error {
	tokenStr := c.Cookies("jwt")
	db := database.Database.DB
	claims, err := utils.ExtractClaims(tokenStr)
	var user models.User

	log.Println(tokenStr)

	if !err {
		c.JSON(nil)
		return fiber.NewError(fiber.StatusBadRequest, "Unauthorized")
	}

	db.Where("ID = ?", claims["user_id"] ).First(&user)
	
	if user.ID == 0 {
		c.JSON(nil)
		return fiber.NewError(fiber.StatusBadRequest, "user not found!")
	}

	c.JSON(user)

	return nil
}

func LogOut(c *fiber.Ctx) error {
	c.ClearCookie("jwt");
	return c.JSON(fiber.Map{
		"success": true,
	})
}

func CheckRegisterKey(c *fiber.Ctx) error {
	var req RegisterKey
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(fiber.Map{
			"error": "Request key not specified",
			"success": false,
		})
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Failed to connect to db! \n", err.Error())
		os.Exit(2)
	}

	registerKey := os.Getenv("SECRET_BETA")

	if (req.Key != registerKey) {
		return c.JSON(fiber.Map{
			"error": "Request key not valid",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}