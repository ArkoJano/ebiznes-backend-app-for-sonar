package controllers

import (
	"ebiznes/models"
	"github.com/labstack/echo/v4"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

const (
	invalidIDMessage       = "Invalid ID"
	paymentNotFoundMessage = "Payment not found"
	invalidAmountMessage   = "Invalid amount: "
	orderIDRequiredMessage = "Order ID is required"
	invalidDataMessage     = "Invalid data"
	errorCreatingMessage   = "Error creating payment"
	errorUpdatingMessage   = "Error updating payment"
	errorDeletingMessage   = "Error deleting payment"
	paymentDeletedMessage  = "Payment deleted"
)

func ListAllPayments(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payments []models.Payment
		db = applyQueryParams(c, db)

		result := db.Find(&payments)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, result.Error)
		}
		return c.JSON(http.StatusOK, payments)
	}
}

func applyQueryParams(c echo.Context, db *gorm.DB) *gorm.DB {
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
	return db.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
}

func GetPayment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := parseID(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, invalidIDMessage)
		}

		var payment models.Payment
		result := db.First(&payment, id)
		if err := handleDBError(c, result, "Payment"); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, payment)
	}
}

func AddPayment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payment models.Payment
		if err := c.Bind(&payment); err != nil {
			return c.JSON(http.StatusBadRequest, invalidDataMessage)
		}
		if err := validatePayment(payment); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		result := db.Create(&payment)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, errorCreatingMessage)
		}
		return c.JSON(http.StatusCreated, payment)
	}
}

func UpdatePayment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := parseID(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, invalidIDMessage)
		}

		var payment models.Payment
		result := db.First(&payment, id)
		if err := handleDBError(c, result, "Payment"); err != nil {
			return err
		}

		if err := c.Bind(&payment); err != nil {
			return c.JSON(http.StatusBadRequest, invalidDataMessage)
		}
		if err := validatePayment(payment); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		result = db.Save(&payment)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, errorUpdatingMessage)
		}

		return c.JSON(http.StatusOK, payment)
	}
}

func validatePayment(payment models.Payment) error {
	if payment.Amount <= 0 {
		return fmt.Errorf(invalidAmountMessage + strconv.FormatFloat(payment.Amount, 'f', -1, 64))
	}
	if payment.OrderID == 0 {
		return errors.New(orderIDRequiredMessage)
	}
	return nil
}

func DeletePayment(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := parseID(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, invalidIDMessage)
		}

		result := db.Delete(&models.Payment{}, id)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, errorDeletingMessage)
		}

		if result.RowsAffected == 0 {
			return c.JSON(http.StatusNotFound, paymentNotFoundMessage)
		}

		return c.JSON(http.StatusOK, paymentDeletedMessage)
	}
}
