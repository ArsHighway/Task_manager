package main

import (
	"log"
	"os"

	"github.com/ArsHighway/Tasks-PSQL/internal/config"
	taskHandler "github.com/ArsHighway/Tasks-PSQL/internal/handlers/taskHandler"
	userHandler "github.com/ArsHighway/Tasks-PSQL/internal/handlers/userHandler"
	taskRepository "github.com/ArsHighway/Tasks-PSQL/internal/repository/taskRepository"
	userRepository "github.com/ArsHighway/Tasks-PSQL/internal/repository/userRepository"
	"github.com/ArsHighway/Tasks-PSQL/internal/routers"
	taskService "github.com/ArsHighway/Tasks-PSQL/internal/service/taskService"
	userService "github.com/ArsHighway/Tasks-PSQL/internal/service/userService"
	"github.com/gin-gonic/gin"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/app?sslmode=disable"
	}

	pool := config.NewPostgressPool(databaseURL)
	defer pool.Close()

	userRepo := userRepository.NewUserRepository(pool)
	taskRepo := taskRepository.NewTaskRepository(pool)
	taskServ := taskService.NewTaskService(taskRepo)
	userServ := userService.NewUserService(userRepo, taskRepo)

	uh := userHandler.NewUserHandler(userServ)
	th := taskHandler.NewTaskHandler(taskServ)

	gin.SetMode(gin.ReleaseMode)
	r := routers.RegisterRoutes(uh, th)

	log.Println("Server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
