package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"product-app/common/app"
	"product-app/common/postgresql"
	"product-app/controller"
	"product-app/persistence"
	"product-app/service"
)

func main() {
	ctx := context.Background()
	e := echo.New()

	configrationManager := app.NewConfigrationManager()

	dbPool := postgresql.GetConnectionPool(ctx, configrationManager.PostgreSqlConfig)

	productRepository := persistence.NewProductRepository(dbPool)

	productService := service.NewProductService(productRepository)

	productController := controller.NewProductController(productService)

	productController.RegisterRoutes(e)

	e.Start("localhost:8080")
}
