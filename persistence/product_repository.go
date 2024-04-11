package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
	"product-app/domain"
	"product-app/persistence/common"
)

type IProductRepository interface {
	GetAllProducts() []domain.Product
	GetAllProductsByStore(storeName string) []domain.Product
	AddProduct(product domain.Product) error
	GetById(productId int64) (domain.Product, error)
	DeleteById(productId int64) error
	UpdatePrice(productId int64, newPrice float32) error
}

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) IProductRepository {
	return &ProductRepository{
		dbPool: dbPool,
	}
}

func (productreposityory *ProductRepository) GetAllProducts() []domain.Product {
	ctx := context.Background()
	productRows, err := productreposityory.dbPool.Query(ctx, "select * from products")
	if err != nil {
		log.Error("Error while getting all products %v", err)
		return []domain.Product{}
	}
	return extractProductsFromRows(productRows)
}

func (productreposityory *ProductRepository) GetAllProductsByStore(storeName string) []domain.Product {
	ctx := context.Background()

	getProductsByStoreNameSql := `select * from products where store = $1`
	productRows, err := productreposityory.dbPool.Query(ctx, getProductsByStoreNameSql, storeName)
	if err != nil {
		log.Error("Error while getting all products %v", err)
		return []domain.Product{}
	}
	return extractProductsFromRows(productRows)
}

func (productreposityory *ProductRepository) AddProduct(product domain.Product) error {
	ctx := context.Background()
	insert_sql := `Insert into products (name, price, discount, store) VALUES ($1, $2, $3, $4)`
	addNewProduct, err := productreposityory.dbPool.Exec(ctx, insert_sql, product.Name, product.Price, product.Discount, product.Store)
	if err != nil {
		log.Error("Error while adding product %v", err)
		return err
	}
	log.Info(fmt.Printf("Product added with %v", addNewProduct))
	return nil
}

func extractProductsFromRows(productRows pgx.Rows) []domain.Product {
	var products = []domain.Product{}
	var id int64
	var name string
	var price float32
	var discount float32
	var store string

	for productRows.Next() {
		productRows.Scan(&id, &name, &price, &discount, &store)
		products = append(products, domain.Product{
			Id:       id,
			Name:     name,
			Price:    price,
			Discount: discount,
			Store:    store,
		})
	}
	return products
}

func (productreposityory *ProductRepository) GetById(productId int64) (domain.Product, error) {
	ctx := context.Background()
	getByIdSql := `select * from products where id = $1`
	queryRow := productreposityory.dbPool.QueryRow(ctx, getByIdSql, productId)

	var id int64
	var name string
	var price float32
	var discount float32
	var store string

	scanErr := queryRow.Scan(&id, &name, &price, &discount, &store)

	if scanErr != nil && scanErr.Error() == common.NOT_FOUND {
		return domain.Product{}, errors.New(fmt.Sprintf("Product not found with id %d", productId))
	}
	if scanErr != nil {
		return domain.Product{}, errors.New(fmt.Sprintf("Error while getting product by id %d", productId))
	}
	return domain.Product{
		Id:       id,
		Name:     name,
		Price:    price,
		Discount: discount,
		Store:    store,
	}, nil
}

func (productreposityory *ProductRepository) DeleteById(productId int64) error {
	ctx := context.Background()
	_, getErr := productreposityory.GetById(productId)
	if getErr != nil {
		return errors.New(fmt.Sprintf("Product not found with id %d", productId))
	}
	deleteSql := `DELETE FROM products WHERE id = $1`
	_, err := productreposityory.dbPool.Exec(ctx, deleteSql, productId)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while deleting product by id %d", productId))
	}
	log.Info("Product deleted with id %d", productId)
	return nil
}

func (productreposityory *ProductRepository) UpdatePrice(productId int64, newPrice float32) error {
	ctx := context.Background()
	_, getError := productreposityory.GetById(productId)
	if getError != nil {
		return errors.New(fmt.Sprintf("Product not found with id %d", productId))
	}
	updatePriceSql := `UPDATE products SET price = $1 WHERE id = $2`
	_, err := productreposityory.dbPool.Exec(ctx, updatePriceSql, newPrice, productId)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while updating product price by id %d", productId))
	}
	log.Info("Product %d price updated new price %v", productId, newPrice)
	return nil
}
