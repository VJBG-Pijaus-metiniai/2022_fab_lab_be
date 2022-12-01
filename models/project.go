package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Title string `json:"title"`
	Description string `json:"description"`
	Supervisor string `json:"supervisor"`
	Author string `json:"author"`
	Images string `json:"images"`
}