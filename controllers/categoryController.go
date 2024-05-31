package controllers

import (
	"backend/models"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var categories []models.Category

func GetCategories(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		db.Find(&categories)
		return c.JSON(http.StatusOK, categories)
	}
}

func GetCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid ID")
		}

		category := models.Category{}
		result := db.First(&category, id)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Category not found")
		}

		return c.JSON(http.StatusOK, category)
	}

}

func AddCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		category := models.Category{}
		if err := c.Bind(&category); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid data")
		}
		db.Create(&category)
		return c.JSON(http.StatusCreated, category)
	}
}

func DeleteCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid ID")
		}

		category := models.Category{}
		db.First(&category, id)
		if category.ID == 0 {
			return c.JSON(http.StatusNotFound, "Category not found")
		}

		db.Delete(&category)
		return c.JSON(http.StatusOK, "Category deleted")
	}
}
