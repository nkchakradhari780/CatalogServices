package storage

import "github.com/nkchakradhari780/catalogServices/internal/modules"

type Storage interface {
	CreateProduct(name string, price int, stock int, categoryId string, Brand string, Images []string) (int, error)
	GetProductById(id int) (modules.Product, error)
}