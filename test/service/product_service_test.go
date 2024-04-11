package service

import (
	"github.com/stretchr/testify/assert"
	"os"
	"product-app/domain"
	"product-app/service"
	"product-app/service/model"
	"testing"
)

var productService service.IProductService

func TestMain(m *testing.M) {
	initialProducts := []domain.Product{
		{
			Id:    1,
			Name:  "AirFryer",
			Price: 1000.0,
			Store: "ABC TECH",
		},
		{
			Id:    2,
			Name:  "Ütü",
			Price: 4000.0,
			Store: "ABC TECH",
		},
	}
	fakeProductRepository := NewFakeProductRepository(initialProducts)
	productService = service.NewProductService(fakeProductRepository)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_ShouldGetAllProducts(t *testing.T) {
	t.Run("ShouldGetAllProducts", func(t *testing.T) {
		acualProducts := productService.GetAllProducts()
		assert.Equal(t, 2, len(acualProducts))
	})
}
func Test_WhenNoValidationErrorOccurrded_ShouldAddProduct(t *testing.T) {
	t.Run("Test_WhenNoValidationErrorOccurrded_ShouldAddProduct", func(t *testing.T) {
		productService.Add(model.ProductCreate{
			Name:     "Ütü",
			Price:    1000.0,
			Discount: 50,
			Store:    "ABC TECH",
		})
		actualProducts := productService.GetAllProducts()
		assert.Equal(t, 3, len(actualProducts))
		assert.Equal(t, domain.Product{
			Id:       3,
			Name:     "Ütü",
			Price:    1000.0,
			Discount: 50,
			Store:    "ABC TECH",
		}, actualProducts[len(actualProducts)-1])
	})
}

func Test_WhenDiscountIsHigherThan70_ShouldNotAddProduct(t *testing.T) {
	t.Run("Test_WhenNoValidationErrorOccurrded_ShouldAddProduct", func(t *testing.T) {
		validationErr := productService.Add(model.ProductCreate{
			Name:     "Ütü",
			Price:    1000.0,
			Discount: 80,
			Store:    "ABC TECH",
		})
		actualProducts := productService.GetAllProducts()
		assert.Equal(t, 2, len(actualProducts))
		assert.Equal(t, "Product discount can not be higher than 75!", validationErr.Error())
	})
}
