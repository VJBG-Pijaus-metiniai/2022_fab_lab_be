package main

import (
	"fablab-project/controllers"
	"fablab-project/database"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		c.Send([]byte("Welcome to the vjbg_fablab projects API"));

		return nil
	})

	app.Post("/login", controllers.Login)
	app.Post("/register", controllers.Register)
	app.Post("/project", controllers.CreateProject)
	app.Get("/project", controllers.GetProjects)
	app.Get("/project/:id", controllers.GetProject)
	app.Delete("/project/:id", controllers.DeleteProject)
	app.Patch("/project/:id", controllers.EditProject)
	app.Get("/user", controllers.GetCurrentUser)
	app.Delete("/login", controllers.LogOut)
	app.Post("/reg", controllers.CheckRegisterKey)
}

func main() {
	app := fiber.New()
	database.ConnectDB()

	app.Use(cors.New(cors.Config{
    	AllowOrigins: "*",
    	AllowHeaders:  "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	setupRoutes(app)

	if os.Getenv("PORT") == "" {
		log.Fatal(app.Listen(":8000"))
	} else {
		log.Fatal(app.Listen(":"+os.Getenv("PORT")))
	}
}