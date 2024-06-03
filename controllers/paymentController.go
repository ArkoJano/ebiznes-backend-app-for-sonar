package controllers

import (
	"ebiznes-backend-app-for-sonar/models"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func ListAllPayments(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payments []models.Payment
		limit := c.QueryParam("limit")
		offset := c.QueryParam("offset")
		if limit != "" {
			limitVal, _ := strconv.Atoi(limit)
			db = db.Limit(limitVal)
		}
		if offset != "" {
			offsetVal, _ := strconv.Atoi(offset)
			db = db.Offset(offsetVal)
		}
		orderID := c.QueryParam("order_id")
		if orderID != "" {
			db = db.Where("order_id = ?", orderID)
		}
		sortOrder := c.QueryParam("sort_order")
		if sortOrder == "" {
			sortOrder = "asc"
		}
		sortField := c.QueryParam("sort_field")
		if sortField == "" {
			sortField = "id"
		}
		db.Order(fmt.Sprintf("%s %s", sortField, sortOrder)).Find(&payments)
		return c.JSON(http.StatusOK, payments)
	}
}

func GetPayment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid ID")
		}

		var payment models.Payment
		result := db.First(&payment, id)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Payment not found")
		} else if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, "Error retrieving payment")
		}

		return c.JSON(http.StatusOK, payment)
	}
}

func AddPayment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payment models.Payment
		if err := c.Bind(&payment); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid data")
		}
		if payment.Amount <= 0 {
			return c.JSON(http.StatusBadRequest, "Invalid amount: "+strconv.FormatFloat(payment.Amount, 'f', -1, 64))
		}
		if payment.OrderID == 0 {
			return c.JSON(http.StatusBadRequest, "Order ID is required")
		}
		result := db.Create(&payment)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, "Error creating payment")
		}
		return c.JSON(http.StatusCreated, payment)
	}
}

func UpdatePayment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid ID")
		}

		var payment models.Payment
		result := db.First(&payment, id)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Payment not found")
		} else if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, "Error retrieving payment")
		}

		if err := c.Bind(&payment); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid data")
		}

		if payment.Amount <= 0 {
			return c.JSON(http.StatusBadRequest, "Invalid amount: "+strconv.FormatFloat(payment.Amount, 'f', -1, 64))
		}

		result = db.Save(&payment)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, "Error updating payment")
		}

		return c.JSON(http.StatusOK, payment)
	}
}

func DeletePayment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid ID")
		}

		result := db.Delete(&models.Payment{}, id)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, "Error deleting payment")
		}

		if result.RowsAffected == 0 {
			return c.JSON(http.StatusNotFound, "Payment not found")
		}

		return c.JSON(http.StatusOK, "Payment deleted")
	}
}
