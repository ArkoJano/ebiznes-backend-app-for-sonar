package controllers

import (
	"ebiznes/models"
	"github.com/labstack/echo/v4"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// checkIfCartExistsAndReturn Function to check if a cart exists and return it.
func checkIfCartExistsAndReturn(db *gorm.DB) (*models.Cart, error) {
	var cart models.Cart
	result := db.Where("id = ?", 1).First(&cart)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		cart = models.Cart{
			Items: []models.CartItem{},
		}

		if createResult := db.Create(&cart); createResult.Error != nil {
			return nil, createResult.Error
		}
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &cart, nil
}

func GetCartItems(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cart, err := checkIfCartExistsAndReturn(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		var items []models.CartItem
		err = db.Where("cart_id = ?", cart.ID).Find(&items).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, items)
	}
}

func AddCartItem(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		cart, err := checkIfCartExistsAndReturn(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		var newItem models.CartItem
		if err := c.Bind(&newItem); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		var existingItem models.CartItem
		result := db.Where("product_id = ? AND cart_id = ?", newItem.ProductID, cart.ID).First(&existingItem)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newItem.CartID = cart.ID
			if err := db.Create(&newItem).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusCreated, newItem)
		} else if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		} else {
			existingItem.Quantity += newItem.Quantity
			if err := db.Save(&existingItem).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, existingItem)
		}
	}
}

func UpdateCartItem(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		itemID, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(400, "Invalid ID")
		}

		var item models.CartItem
		result := db.First(&item, itemID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Item not found")
		} else if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		if err := c.Bind(&item); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if err := db.Save(&item).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, item)
	}
}

func DeleteCartItem(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		itemID, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid ID")
		}

		var item models.CartItem
		result := db.First(&item, itemID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Item not found")
		} else if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}

		if err := db.Delete(&item).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusNoContent, nil)
	}
}
