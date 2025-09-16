package storage

import "github.com/nkchakradhari780/catalogServices/internal/modules"

type Storage interface {
	CreateProduct(name string, price int, stock int, categoryId string, quantity int, Brand string, Images []string) (int, error)
	GetProductById(id int) (modules.Product, error)
	GetProducts() ([]modules.Product, error)
	GetDefaultProducts() ([]modules.Product, error)
	GetFilteredProducts(filters map[string][]string) ([]modules.Product, error)
	UpdateProductById(id int, name string, price int, stock int, categoryId string, quantity int, brand string, images []string) (modules.Product, error)
	DeleteProductById(id int) error
	SearchProducts(qureyStr string) ([]modules.Product, error)

	CreateUser(name string, email string, password string, phone string, role string, address string) (int, error)
}
