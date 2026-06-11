package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParsePagination(c *gin.Context) (int, int, int) {
	page := parsePositiveInt(c.Query("page"), 1)
	limit := parsePositiveInt(c.Query("limit"), 20)
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	return page, limit, offset
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func QueryInt(c *gin.Context, key string, fallback int) int {
	val := c.Query(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return parsed
}
