package controllers

import (
	"ebiznes-backend-app-for-sonar/models"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func GetProduct(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		productID, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid product ID")
		}

		var product models.Product
		result := db.First(&product, productID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Product not found")
		} else if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		return c.JSON(http.StatusOK, product)
	}

}

func GetProductsByCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		categoryIDParam := c.Param("id")
		categoryID, err := strconv.Atoi(categoryIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid or missing category ID")
		}

		var products []models.Product
		result := db.Scopes(models.ProductsByCategory(uint(categoryID))).Find(&products)

		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, "Error fetching products")
		}

		return c.JSON(http.StatusOK, products)
	}
}

func AddProduct(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var product models.Product
		if err := c.Bind(&product); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if result := db.Create(&product); result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		return c.JSON(http.StatusCreated, product)
	}
}

func UpdateProduct(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		productID, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid product ID")
		}

		var product models.Product
		result := db.First(&product, productID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Product not found")
		} else if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		if err := c.Bind(&product); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if result := db.Save(&product); result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		return c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		productID, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid product ID")
		}

		var product models.Product
		result := db.First(&product, productID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Product not found")
		} else if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		if result := db.Delete(&product); result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		return c.JSON(http.StatusNoContent, nil)
	}
}

func GetProducts(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var products []models.Product
		if result := db.Find(&products); result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}
		return c.JSON(http.StatusOK, products)
	}
}
