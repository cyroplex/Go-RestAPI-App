package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"product-app/controller/request"
	"product-app/controller/response"
	"product-app/service"
	"strconv"
)

type ProductController struct {
	productService service.IProductService
}

func NewProductController(productService service.IProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (productController *ProductController) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/v1/products/:id", productController.GetProductById)
	e.GET("/api/v1/products", productController.GetAllProducts)
	e.POST("/api/v1/products", productController.AddProduct)
	e.PUT("/api/v1/products/:id", productController.UpdatePrice)
	e.DELETE("/api/v1/products/:id", productController.DeleteProductById)
}

func (productController *ProductController) GetProductById(c echo.Context) error {
	productId := c.Param("id")
	atoiId, _ := strconv.Atoi(productId)
	product, err := productController.productService.GetById(int64(atoiId))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, response.ToResponse(product))
}

func (productController *ProductController) GetAllProducts(c echo.Context) error {
	store := c.QueryParam("store")
	if len(store) == 0 {
		allProducts := productController.productService.GetAllProducts()
		return c.JSON(http.StatusOK, response.ToResponseList(allProducts))
	}
	productsWithGivenStore := productController.productService.GetAllProductsByStore(store)
	return c.JSON(http.StatusOK, response.ToResponseList(productsWithGivenStore))
}

func (productController *ProductController) AddProduct(c echo.Context) error {
	var addProductRequest request.AddProductRequest
	bindErr := c.Bind(&addProductRequest)
	if bindErr != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: bindErr.Error(),
		})
	}
	addErr := productController.productService.Add(addProductRequest.ToModel())
	if addErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorDescription: addErr.Error(),
		})
	}
	return c.NoContent(http.StatusCreated)
}

func (productController *ProductController) UpdatePrice(c echo.Context) error {
	productId := c.Param("id")
	newPrice := c.QueryParam("newPrice")
	atoiId, _ := strconv.Atoi(productId)
	parsePrice, _ := strconv.ParseFloat(newPrice, 32)
	if len(newPrice) == 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: "newPrice should be sent",
		})
	}
	err := productController.productService.UpdatePrice(int64(atoiId), float32(parsePrice))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})
	}
	return c.NoContent(http.StatusCreated)
}

func (productController *ProductController) DeleteProductById(c echo.Context) error {
	productId := c.Param("id")
	atoiId, _ := strconv.Atoi(productId)
	err := productController.productService.DeleteById(int64(atoiId))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})
	}
	return c.NoContent(http.StatusOK)
}
