package controllers

import (
	"fablab-project/database"
	"fablab-project/models"
	"fablab-project/utils"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type createProjectRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Supervisor string `json:"supervisor"`
 	Images pq.StringArray `gorm:"type:text[]" json:"images"`
}

func CreateProject(c *fiber.Ctx) error {
	tokenStr := c.Cookies("jwt")
	claims, err := utils.ExtractClaims(tokenStr)
	db := database.Database.DB
	var req createProjectRequest

	if err := c.BodyParser(&req); err != nil {
		log.Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "failed parsing request body")
	}

	if !err {
		return fiber.NewError(fiber.StatusBadRequest, "Unauthorized")
	}

	var user models.User
	db.Where("ID = ?", claims["user_id"] ).First(&user)
	
	if user.ID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "user not found!")
	}

	project := new(models.Project)
	project.Title = req.Title
	project.Description = req.Description
	project.Supervisor = req.Supervisor
	project.Images = strings.Join(req.Images, `{`)
	project.Author = user.Name

	db.Create(&project)
	c.JSON(project)

	return nil
}

func GetProjects(c *fiber.Ctx) error {
	db := database.Database.DB
	var projects []models.Project
	db.Find(&projects);
	c.JSON(projects)

	return nil
}

func GetProject(c *fiber.Ctx) error {
	id := c.Params("id");
	db := database.Database.DB
	var project models.Project

	db.Where("ID = ?", id).First(&project);

	c.JSON(project)

	return nil
}

func DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	tokenStr := c.Cookies("jwt")
	claims, err := utils.ExtractClaims(tokenStr)
	db := database.Database.DB
	var user models.User
	var project models.Project

	if !err  {
		return fiber.NewError(fiber.StatusBadRequest, "Unauthorized")
	}
	
	db.Where("ID = ?", claims["user_id"]).First(&user)
	db.Where("ID = ?", id).First(&project)

	log.Println(user)

	if project.Author != user.Name {
		return fiber.NewError(fiber.StatusBadRequest, "Unauthorized")
	}

	db.Delete(&project)
	c.JSON(fiber.Map{
		"success": true,
	})
	return nil
}

func EditProject(c *fiber.Ctx) error {
	id := c.Params("id")
	tokenStr := c.Cookies("jwt")
	claims, err := utils.ExtractClaims(tokenStr)
	var req createProjectRequest
	db := database.Database.DB
	var user models.User
	var project models.Project

	if err := c.BodyParser(&req); err != nil {
		log.Println(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "failed parsing request body")
	}

	if !err  {
		return fiber.NewError(fiber.StatusBadRequest, "Unauthorized")
	}
	
	db.Where("ID = ?", claims["user_id"]).First(&user)
	db.Where("ID = ?", id).First(&project)

	if project.Author != user.Name {
		return fiber.NewError(fiber.StatusBadRequest, "Unauthorized")
	}

	if (req.Description != "") {
		project.Description = req.Description
	} else if (req.Supervisor != "") {
		project.Supervisor = req.Supervisor
	} else if (req.Title != "") {
		project.Title = req.Title
	}

	db.Save(&project)
	c.JSON(project)
	return nil
}