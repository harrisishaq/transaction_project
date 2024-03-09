package main

import (
	"fmt"
	"log"
	"test_project/config"
	"test_project/controller"
	"test_project/entity"
	"test_project/helpers"
	"test_project/model"
	"test_project/repository"
	"test_project/service"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	var e = echo.New()
	e.Use(middleware.CORS())
	e.Validator = &helpers.EchoValidator{Validator: validator.New()}

	config.LoadConfig()
	dbInit := config.InitDBConnection(config.GormConfig.DBHost, config.GormConfig.DBUser, config.GormConfig.DBPassword, config.GormConfig.DBName, config.GormConfig.DBPort)

	// setup repositories
	usersRepo := repository.NewUserRepository(dbInit)

	// setup services
	usersService := service.NewUserService(usersRepo)

	// setup controller
	usersController := controller.NewUserController(usersService)

	// migrate database
	migrate := gormigrate.New(dbInit, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&entity.User{},
					&entity.UserLog{},
					&entity.Category{},
					&entity.CategoryLog{},
					&entity.Product{},
					&entity.ProductLog{},
				); err != nil {
					return err
				}

				timeNow := time.Now()

				hashPassword, err := bcrypt.GenerateFromPassword([]byte(config.AppConfig.DefaultPassword), bcrypt.DefaultCost)
				if err != nil {
					fmt.Println("Error encrypting password:", err)
					return model.NewError("500", "Internal server error.")
				}

				var newData = &entity.User{
					ID:       uuid.New(),
					Name:     "ADMIN",
					Email:    config.AppConfig.DefaultEmail,
					Password: fmt.Sprintf("%x", hashPassword),
					Audit: &entity.Audit{
						CurrNo:    1,
						CreatedAt: &timeNow,
						CreatedBy: "SYSTEM",
					},
				}
				if err := tx.Create(newData).Error; err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					&entity.User{},
					&entity.UserLog{},
					&entity.Category{},
					&entity.CategoryLog{},
					&entity.Product{},
					&entity.ProductLog{})
			},
		},
	})

	if err := migrate.Migrate(); err != nil {
		log.Fatalf("Migrate failed: %+v\n", err)
	} else {
		log.Println("Success Migration")
	}

	usersController.UserRoutes(e)

	e.Logger.Fatal(e.Start(":3003"))
}
