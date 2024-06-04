package main

import (
	"ebiznes-backend-app-for-sonar/controllers"
	"ebiznes-backend-app-for-sonar/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitializeDatabase Function to initialize the database with some products
func InitializeDatabase(db *gorm.DB) {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count == 0 {
		products := []models.Product{
			{Name: "Produkt 1", Price: 100.00},
			{Name: "Produkt 2", Price: 200.00},
			{Name: "Produkt 3", Price: 300.00},
			{Name: "Produkt 4", Price: 400.00},
			{Name: "Produkt 5", Price: 500.00},
		}
		for _, product := range products {
			db.Create(&product)
		}
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	db, err := gorm.Open(sqlite.Open("products.db"), &gorm.Config{})
	if err != nil {

		panic("failed to connect database" + err.Error())
	}

	err = db.AutoMigrate(&models.Product{}, &models.Cart{}, &models.CartItem{}, &models.Category{}, &models.Payment{})
	if err != nil {
		panic("failed to migrate database" + err.Error())
	}

	InitializeDatabase(db)

	const productsIdEndpoint = "/products/:id"

	e.GET("/", controllers.Home)
	e.GET("/products", controllers.GetProducts(db))
	e.GET(productsIdEndpoint, controllers.GetProduct(db))
	e.PUT(productsIdEndpoint, controllers.UpdateProduct(db))
	e.DELETE(productsIdEndpoint, controllers.DeleteProduct(db))
	e.POST("/products", controllers.AddProduct(db))
	e.GET("/products/category/:id", controllers.GetProductsByCategory(db))

	e.GET("/cart", controllers.GetCartItems(db))
	e.POST("/cart", controllers.AddCartItem(db))
	e.PUT("/cart/:id", controllers.UpdateCartItem(db))
	e.DELETE("/cart/:id", controllers.DeleteCartItem(db))

	e.GET("/categories", controllers.GetCategories(db))
	e.GET("/categories/:id", controllers.GetCategory(db))
	e.POST("/categories", controllers.AddCategory(db))
	e.DELETE("/categories/:id", controllers.DeleteCategory(db))

	e.GET("/payments", controllers.ListAllPayments(db))
	e.GET("/payments/:id", controllers.GetPayment(db))
	e.POST("/payments", controllers.AddPayment(db))
	e.DELETE("/payments/:id", controllers.DeletePayment(db))

	e.Logger.Fatal(e.Start(":8080"))
}
