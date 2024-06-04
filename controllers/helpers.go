package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func parseID(c echo.Context) (int, error) {
	id := c.Param("id")
	return strconv.Atoi(id)
}

func handleDBError(c echo.Context, result *gorm.DB, entityName string) error {
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, entityName+" not found")
	} else if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	return nil
}
