package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{Success: true, Message: message, Data: data})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{Success: true, Message: message, Data: data})
}

func Error(c *gin.Context, status int, message string, errors interface{}) {
	c.JSON(status, APIResponse{Success: false, Message: message, Errors: errors})
}

func BadRequest(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusBadRequest, message, errors)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

func NotFound(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusNotFound, message, errors)
}

func ServerError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "Something went wrong. Please try again.", nil)
}

func InternalError(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusInternalServerError, message, errors)
}

func DBError(c *gin.Context, err error, notFoundMessage string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		NotFound(c, notFoundMessage, nil)
		return
	}
	ServerError(c)
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Paginated(c *gin.Context, message string, items interface{}, page int, limit int, total int64) {
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	OK(c, message, gin.H{
		"items": items,
		"meta": PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
