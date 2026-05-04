package routers

import (
	task "github.com/ArsHighway/Tasks-PSQL/internal/handlers/taskHandler"
	user "github.com/ArsHighway/Tasks-PSQL/internal/handlers/userHandler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(userHandler user.UserHandler, taskHandler task.TaskHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	tasks := r.Group("/tasks")
	{
		tasks.POST("/", taskHandler.CreateTask)
		tasks.GET("/", taskHandler.GetTasks)
		tasks.GET("/:id", taskHandler.GetTaskWithID)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.PATCH("/:id", taskHandler.PatchTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}

	users := r.Group("/users")
	{
		users.POST("/", userHandler.CreateUser)
		users.GET("/:id/tasks", userHandler.GetTaskWithUserID)
		users.GET("/:id", userHandler.GetUserWithID)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.PATCH("/:id", userHandler.PatchUser)
	}

	return r
}
