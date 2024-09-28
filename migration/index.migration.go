package migration

import (
	"first-app/database"
	"first-app/models"
	"log"
)

func RunMigration() {

	database.DB.AutoMigrate(&models.User{},
		&models.Role{},
		&models.Course{},
		&models.CourseUser{},
		&models.Lesson{},
		&models.Post{},
		&models.Assignment{},
		&models.StudentAssignment{},
		&models.FileAssignment{},
		&models.FilePost{})

	log.Println("Database Migrated")
}
