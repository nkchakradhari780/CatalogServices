package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/nkchakradhari780/catalogServices/internal/config"
	"github.com/nkchakradhari780/catalogServices/internal/modules"
)

type Postgres struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error)  {

	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products (
	id          SERIAL PRIMARY KEY,               -- unique product id
	name        VARCHAR(255) NOT NULL,          -- product name
	price       INT NOT NULL,                   -- product price
	stock       INT NOT NULL,                   -- product stock count
	category_id INT NOT NULL,                  -- category reference
	brand       VARCHAR(100) NOT NULL,          -- product brand
	images      TEXT[]                          -- array of image URLs
);
`)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err 
	}

	return &Postgres{
		Db: db,
	}, nil
}

func (p *Postgres) CreateProduct(name string, price int, stock int, categoryId string, Brand string, Images []string) (int, error) {
	stmt, err := p.Db.Prepare("INSERT INTO products (name, price, stock, category_id, brand, images) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()


	var lastId int
	err = stmt.QueryRow(name, price, stock, categoryId, Brand, pq.Array(Images)).Scan(&lastId)
	if err != nil {
		return 0, err
	}

	return int(lastId), nil
}

func (p *Postgres) GetProductById(id int) (modules.Product, error) {
	stmt, err := p.Db.Prepare("SELECT * FROM products WHERE id = $1")
	if err != nil {
		return modules.Product{}, err
	}
	defer stmt.Close()

	var product modules.Product

	err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.Brand, pq.Array(&product.Images))
	if err != nil {
		if err == sql.ErrNoRows {
			return modules.Product{}, fmt.Errorf("product with id %d not found", id)
		}
		return modules.Product{}, fmt.Errorf("error fetching product: %v", err)
	}

	return product, nil
}
