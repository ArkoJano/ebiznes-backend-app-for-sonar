package controllers

import (
	"ebiznes/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

const invalidProductIDMessage = "Invalid product ID"

func GetProduct(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		productID, err := parseID(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, invalidProductIDMessage)
		}

		var product models.Product
		result := db.First(&product, productID)
		if err := handleDBError(c, result, "Product"); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, product)
	}
}

func GetProductsByCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		categoryID, err := parseID(c)
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

		result := db.Create(&product)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		return c.JSON(http.StatusCreated, product)
	}
}

func UpdateProduct(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		productID, err := parseID(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, invalidProductIDMessage)
		}

		var product models.Product
		result := db.First(&product, productID)
		if err := handleDBError(c, result, "Product"); err != nil {
			return err
		}

		if err := c.Bind(&product); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		result = db.Save(&product)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		return c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		productID, err := parseID(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, invalidProductIDMessage)
		}

		var product models.Product
		result := db.First(&product, productID)
		if err := handleDBError(c, result, "Product"); err != nil {
			return err
		}

		result = db.Delete(&product)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		return c.JSON(http.StatusNoContent, nil)
	}
}

func GetProducts(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var products []models.Product
		result := db.Find(&products)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}
		return c.JSON(http.StatusOK, products)
	}
}
