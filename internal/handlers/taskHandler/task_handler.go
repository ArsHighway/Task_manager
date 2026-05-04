package task

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/internal/errs"
	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	task "github.com/ArsHighway/Tasks-PSQL/internal/service/taskService"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type taskHandler struct {
	serv task.TaskService
}

func NewTaskHandler(serv task.TaskService) *taskHandler {
	return &taskHandler{serv: serv}
}

type TaskHandler interface {
	CreateTask(c *gin.Context)
	GetTaskWithID(c *gin.Context)
	UpdateTask(c *gin.Context)
	PatchTask(c *gin.Context)
	DeleteTask(c *gin.Context)
	GetTasks(c *gin.Context)
	GetTaskByUserID(c *gin.Context)
}

func (h *taskHandler) CreateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	log := slog.With("method_func", "CreateTask", "method", c.Request.Method)

	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	log.Info("Creating task",
		"title", t.Title,
		"user_id", t.UserID,
	)
	task, err := h.serv.CreateTask(ctx, &t)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create task")
		log.Error("Create task failed", "error", err)
		return
	}

	c.JSON(http.StatusCreated, task)
	log.Info("Task created successfully", "taskID", task.ID)
}

func (h *taskHandler) GetTaskWithID(c *gin.Context) {
	log := slog.With("method_func", "GetTaskWithID",
		"request_method", c.Request.Method)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Сonversion error")
		log.Warn("Сonversion error", "error", err)
		return
	}
	log.Info("Get task", "taskID", id)
	t, err := h.serv.GetTaskWithID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.String(http.StatusNotFound, "Task not found")
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		log.Warn("Failed to get Task", "error", err)
		return
	}
	c.JSON(http.StatusOK, t)
	log.Info("task received", "task", t.Title)
}

func (h *taskHandler) UpdateTask(c *gin.Context) {
	log := slog.With("method_func", "UpdateTask",
		"request_method", c.Request.Method)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Сonversion error")
		log.Warn("Сonversion error", "error", err)
		return
	}
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	log.Info("Update task", "taskID", id)
	t, err := h.serv.UpdateTask(ctx, id, &task)
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			c.String(http.StatusNotFound, "Task not found")
		} else if errors.Is(err, errs.ErrInvalidTask) {
			c.String(http.StatusBadRequest, "Invalid task")
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		log.Warn("Failed to update Task", "error", err)
		return
	}
	c.JSON(http.StatusOK, t)
	log.Info("task updated", "task", t.Title)
}

func (h *taskHandler) PatchTask(c *gin.Context) {
	log := slog.With("method_func", "UpdateTask",
		"request_method", c.Request.Method)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Сonversion error")
		log.Warn("Сonversion error", "error", err)
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	log.Info("Patch task", "taskID", id)
	t, err := h.serv.PatchTask(ctx, id, updates)
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			c.String(http.StatusNotFound, "Task not found")
		} else if errors.Is(err, errs.ErrInvalidTask) || errors.Is(err, errs.ErrNotValidFields) {
			c.String(http.StatusBadRequest, "Invalid patch fields")
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		log.Warn("Failed to patch Task", "error", err)
		return
	}
	c.JSON(http.StatusOK, t)
	log.Info("task patch", "task", t.Title)
}

func (h *taskHandler) DeleteTask(c *gin.Context) {
	log := slog.With("method_func", "DeleteTask",
		"request_method", c.Request.Method)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Сonversion error")
		log.Warn("Сonversion error", "error", err)
		return
	}
	log.Info("Update task", "taskID", id)
	err = h.serv.DeleteTask(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			c.String(http.StatusNotFound, "Task not found")
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		log.Warn("Failed to delete Task", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
		"taskID":  id,
	})

	log.Info("Task delete", "task", id)
}

func (h *taskHandler) GetTasks(c *gin.Context) {
	log := slog.With("method_func", "GetTasks", "method", c.Request.Method)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	params := c.Request.URL.Query()
	log.Info("Get tasks")
	t, err := h.serv.GetTasks(ctx, params)
	if err != nil {
		if errors.Is(err, errs.ErrBadConvertation) || errors.Is(err, errs.ErrNotValidFields) || errors.Is(err, errs.ErrInvalidTask) {
			c.String(http.StatusBadRequest, "Bad request")
		} else if errors.Is(err, errs.ErrTaskNotFound) || errors.Is(err, pgx.ErrNoRows) {
			c.String(http.StatusNotFound, "Task not found")
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		log.Warn("Failed to take Task", "error", err)
		return
	}
	c.JSON(http.StatusOK, t)
	log.Info("tasks received")
}

func (h *taskHandler) GetTaskByUserID(c *gin.Context) {
	log := slog.With("method_func", "GetTaskByUserID", "method")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.String(http.StatusBadRequest, "Сonversion error")
		log.Warn("Сonversion error", "error", err)
		return
	}
	task, err := h.serv.GetTasksByUserID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) || errors.Is(err, pgx.ErrNoRows) {
			c.String(http.StatusNotFound, "Task not found")
		} else if errors.Is(err, errs.ErrInvalidTask) {
			c.String(http.StatusBadRequest, "Invalid task")
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		log.Warn("Failed to get task by user id", "error", err)
		return
	}
	c.JSON(http.StatusOK, task)
	log.Info("task got by user id")
}
