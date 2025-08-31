package storage

type Storage interface {
	CreateProduct(name string, price int, stock int, categoryId string, Brand string, Images []string) (int, error)
}