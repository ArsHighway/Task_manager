package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/internal/errs"
	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	user "github.com/ArsHighway/Tasks-PSQL/internal/service/userService"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type userHandler struct {
	serv user.UserService
}

func NewUserHandler(serv user.UserService) *userHandler {
	return &userHandler{serv: serv}
}

type UserHandler interface {
	CreateUser(c *gin.Context)
	GetUserWithID(c *gin.Context)
	PatchUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetTaskWithUserID(c *gin.Context)
}

func (h *userHandler) CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	log := slog.With("handler", "CreateUser", "method", c.Request.Method)

	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	log.Info("Creating user", "name", u.Name, "email", u.Email)
	user, err := h.serv.CreateUser(ctx, &u)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create user")
		log.Error("Create user failed", "error", err)
		return
	}

	c.JSON(http.StatusCreated, user)
	log.Info("User created successfully", "userID", user.ID)
}

func (h *userHandler) GetUserWithID(c *gin.Context) {
	log := slog.With("handler", "GetUserWithID", "request_method", c.Request.Method)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusNotFound, "Сonversion error")
		log.Warn("Сonversion error", "error", err)
		return
	}
	log.Info("Get user", "userID", id)
	u, err := h.serv.GetUserWithID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.String(http.StatusNotFound, "User not found")
		} else {
			c.String(http.StatusInternalServerError, "Failed to get user")
		}
		log.Warn("Failed to get user", "error", err)
		return
	}
	c.JSON(http.StatusOK, u)
	log.Info("user received", "user", u.Name)
}

func (h *userHandler) PatchUser(c *gin.Context) {
	log := slog.With("handler", "PatchUser", "request_method", c.Request.Method)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusNotFound, "Сonversion error")
		log.Warn("Сonversion error", "error", err)
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	log.Info("Patch user", "userID", id)
	u, err := h.serv.PatchUser(ctx, id, updates)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			c.String(http.StatusNotFound, "User not found")
		case errors.Is(err, errs.ErrNotValidFieldsUser):
			c.String(http.StatusBadRequest, "No valid fields to update")
		default:
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		log.Warn("Failed to patch user", "error", err)
		return
	}
	c.JSON(http.StatusOK, u)
	log.Info("user updated", "user", u.Name)
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	log := slog.With("handler", "DeleteUser", "request_method", c.Request.Method)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusNotFound, "Сonversion error")
		log.Warn("Сonversion error", "error", err)
		return
	}
	log.Info("Delete user", "userID", id)
	err = h.serv.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			c.String(http.StatusNotFound, "User not found")
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		log.Warn("Failed to delete user", "error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"userID":  id,
	})

	log.Info("User deleted", "userID", id)
}

func (h *userHandler) GetTaskWithUserID(c *gin.Context) {
	log := slog.With("handler", "GetTaskWithUserID", "request_method", c.Request.Method)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusNotFound, "Сonversion error")
		log.Warn("Сonversion error", "error", err)
		return
	}
	log.Info("Get tasks for user", "userID", id)
	tasks, err := h.serv.GetUserTasks(ctx, id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		log.Warn("Failed to get user tasks", "error", err)
		return
	}
	c.JSON(http.StatusOK, tasks)
	log.Info("user tasks received", "count", len(tasks))
}
