package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/lib/pq"
	"github.com/nkchakradhari780/catalogServices/internal/cache"
	"github.com/nkchakradhari780/catalogServices/internal/config"
	"github.com/nkchakradhari780/catalogServices/internal/modules"
)

type Postgres struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
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

	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id     SERIAL PRIMARY KEY,            
			name        TEXT NOT NULL,
			email       TEXT UNIQUE NOT NULL, 
			password    TEXT NOT NULL,
			phone       TEXT,
			role        TEXT CHECK (role IN ('admin', 'user')) DEFAULT 'user',
			address     TEXT,	
			created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS products (
			product_id   SERIAL PRIMARY KEY,           
			name         VARCHAR(255) NOT NULL,        
			price        INT NOT NULL,                 
			stock        INT NOT NULL,                 
			category_id  INT NOT NULL,                 
			quantity     INT NOT NULL, 				
			brand        VARCHAR(100) NOT NULL,        
			images       TEXT[]
		)`,

		`CREATE TABLE IF NOT EXISTS cartTable (
			cart_id     SERIAL PRIMARY KEY,			
			user_id     INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
			created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,			
			updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,			
			status      TEXT CHECK (status IN ('active', 'ordered', 'abandoned')) DEFAULT 'active'
		)`,

		`CREATE TABLE IF NOT EXISTS cartItems (
			cart_item_id  SERIAL PRIMARY KEY, 		
			cart_id       INT NOT NULL REFERENCES cartTable(cart_id) ON DELETE CASCADE,
			product_id    INT NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
			quantity      INT NOT NULL,
			price_at_time FLOAT NOT NULL,
			discount      FLOAT DEFAULT NULL,
			subtotal      FLOAT NOT NULL,
			added_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS wishList (
			wish_list_id SERIAL PRIMARY KEY, 		
			product_id   INT NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
			user_id      INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
			added_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_user_product UNIQUE (user_id, product_id)
		)`,
	}

	for _, query := range tables {
		if _, err := db.Exec(query); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	slog.Info("✅ Tables created successfully")

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{Db: db}, nil
}

func (p *Postgres) CreateProduct(name string, price int, stock int, categoryId string, quantity int, Brand string, Images []string) (int, error) {
	stmt, err := p.Db.Prepare("INSERT INTO products (name, price, stock, category_id, quantity, brand, images) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING product_id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var lastId int
	err = stmt.QueryRow(name, price, stock, categoryId, quantity, Brand, pq.Array(Images)).Scan(&lastId)
	if err != nil {
		return 0, err
	}

	InvalidateProductCache()

	return int(lastId), nil
}

func (p *Postgres) GetProductById(id int) (modules.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)

	if cached, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result(); err == nil {
		var product modules.Product
		if unmarshalErr := json.Unmarshal([]byte(cached), &product); unmarshalErr == nil {
			fmt.Println("Cache Memory Hit")
			return product, nil
		}
	}

	stmt, err := p.Db.Prepare("SELECT * FROM products WHERE id = $1")
	if err != nil {
		return modules.Product{}, err
	}
	defer stmt.Close()

	var product modules.Product

	err = stmt.QueryRow(id).Scan(&product.ProductId, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.Brand, pq.Array(&product.Images))
	if err != nil {
		if err == sql.ErrNoRows {
			return modules.Product{}, fmt.Errorf("product with id %d not found", id)
		}
		return modules.Product{}, fmt.Errorf("error fetching product: %v", err)
	}

	data, _ := json.Marshal(product)
	cache.Rdb.Set(cache.Ctx, cacheKey, data, 7*24*time.Hour)

	fmt.Println("Cache missed Fetching data from DB.")

	return product, nil
}

func (p *Postgres) GetProducts() ([]modules.Product, error) {
	cacheKey := "products_cache_key"

	if val, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result(); err == nil {
		var products []modules.Product
		if jsonErr := json.Unmarshal([]byte(val), &products); jsonErr == nil {
			fmt.Println("Redis cache hit")
			return products, nil
		}
	}

	stmt, err := p.Db.Prepare("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []modules.Product

	for rows.Next() {
		var product modules.Product
		err := rows.Scan(&product.ProductId, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.Quantity, &product.Brand, pq.Array(&product.Images))
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	data, _ := json.Marshal(products)
	cache.Rdb.Set(cache.Ctx, cacheKey, data, 7*24*time.Hour)

	fmt.Println(" Cache missed data fetched from DB")

	return products, nil
}

func (p *Postgres) GetDefaultProducts() ([]modules.Product, error) {
	cacheKey := "default_products"

	cached, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err == nil {
		var products []modules.Product
		if unmarshalErr := json.Unmarshal([]byte(cached), &products); unmarshalErr == nil {
			fmt.Println("Redis cache hit")
			return products, nil
		}
	}

	stms, err := p.Db.Prepare(`SELECT product_id, name, price, stock, category_id, quantity, brand, images
					FROM products
					ORDER BY RANDOM()
					LIMIT 50;
				`)
	if err != nil {
		return nil, err
	}
	defer stms.Close()

	rows, err := stms.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []modules.Product

	for rows.Next() {
		var product modules.Product
		err := rows.Scan(&product.ProductId, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.Quantity, &product.Brand, pq.Array(&product.Images))
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	data, _ := json.Marshal(products)
	cache.Rdb.Set(cache.Ctx, cacheKey, data, 7*24*time.Hour)

	fmt.Println("Cache missed fetching data from DB.")
	return products, nil
}

func (p *Postgres) GetFilteredProducts(filters map[string][]string) ([]modules.Product, error) {

	cacheKey := "products:filtered:"
	for k, v := range filters {
		cacheKey += fmt.Sprintf("%s=%s;", k, v)
	}

	if cached, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result(); err == nil {
		var products []modules.Product
		if unmarshalErr := json.Unmarshal([]byte(cached), &products); unmarshalErr == nil {
			fmt.Println("Cache hit for filters:", cacheKey)
			return products, nil
		}
	}

	query := `SELECT product_id, name, price, stock, category_id, brand, images FROM products WHERE 1=1 `
	args := []any{}
	argID := 1

	// Building Dynamic query
	if name, ok := filters["name"]; ok {
		query += fmt.Sprintf("AND name ILIKE $%d ", argID)
		args = append(args, "%"+name[0]+"%")
		argID++
	}

	if brand, ok := filters["brand"]; ok {
		query += fmt.Sprintf("AND brand = $%d ", argID)
		args = append(args, brand[0])
		argID++
	}

	if category, ok := filters["category_id"]; ok {
		query += fmt.Sprintf("AND category_id = $%d ", argID)
		args = append(args, category[0])
		argID++
	}

	if minPrice, ok := filters["min_price"]; ok {
		query += fmt.Sprintf("AND price >= $%d ", argID)
		args = append(args, minPrice[0])
		argID++
	}

	if maxPrice, ok := filters["max_price"]; ok {
		query += fmt.Sprintf("AND price <= $%d ", argID)
		args = append(args, maxPrice[0])
		argID++
	}

	if stockGT, ok := filters["stock_gt"]; ok {
		query += fmt.Sprintf("AND stock > $%d ", argID)
		args = append(args, stockGT[0])
		argID++
	}

	query += "ORDER BY product_id DESC LIMIT 50"

	rows, err := p.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []modules.Product
	for rows.Next() {
		var product modules.Product
		err := rows.Scan(&product.ProductId, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.Brand, pq.Array(&product.Images))
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if len(products) > 0 {
		data, _ := json.Marshal(products)
		cache.Rdb.Set(cache.Ctx, cacheKey, data, 7*24*time.Hour)
		fmt.Println("Cache missed data stored for filters:", cacheKey)
	}

	return products, nil
}

func (p *Postgres) SearchProducts(qureyStr string) ([]modules.Product, error) {
	sqlQurey := `SELECT product_id, name, price, stock, category_id, brand, images FROM products WHERE name ILIKE $1 OR brand ILIKE $1 ORDER BY product_id DESC LIMIT 50;`

	rows, err := p.Db.Query(sqlQurey, "%"+qureyStr+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []modules.Product
	for rows.Next() {
		var product modules.Product
		if err := rows.Scan(&product.ProductId, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.Brand, pq.Array(&product.Images)); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil

}

func (p *Postgres) UpdateProductById(id int, name string, price int, stock int, categoryId string, quantity int, Brand string, Images []string) (modules.Product, error) {
	stmt, err := p.Db.Prepare("UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4, quantity=$5, brand = $6, images = $7 WHERE id = $8 RETURNING *")
	if err != nil {
		return modules.Product{}, err
	}
	defer stmt.Close()

	var product modules.Product
	err = stmt.QueryRow(name, price, stock, categoryId, quantity, Brand, pq.Array(Images), id).Scan(&product.ProductId, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.Quantity, &product.Brand, pq.Array(&product.Images))
	if err != nil {
		if err == sql.ErrNoRows {
			return modules.Product{}, fmt.Errorf("product with id %d not found", id)
		}
		return modules.Product{}, fmt.Errorf("error updating product: %v", err)
	}

	InvalidateProductCache()
	return product, nil
}

func (p *Postgres) DeleteProductById(id int) error {
	stmt, err := p.Db.Prepare("DELETE FROM products WHERE id = $1")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("error deleting product %w", err)
	}

	InvalidateProductCache()

	return nil

}

func InvalidateProductCache() {
	err := cache.Rdb.FlushDB(cache.Ctx).Err()
	if err != nil {
		fmt.Println("Error Clearing Cache: ", err)
		return
	}

	fmt.Println("Data Erased from cache memory")
}

func (p *Postgres) CreateUser(name string, email string, password string, phone string, role string, address string) (int, error) {
	stmt, err := p.Db.Prepare("INSERT INTO users (name, email, password,phone, role, address) VALUES ($1, $2, $3, $4,$5, $6) RETURNING user_id")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var userId int
	err = stmt.QueryRow(name, email, password, phone, role, address).Scan(&userId)

	if err != nil {
		return 0, err
	}

	return int(userId), nil
}

func (p *Postgres) AddToWishList(user_id int, product_id int) (int, error) {
	stmt, err := p.Db.Prepare("INSERT INTO wishList (product_id, user_id) VALUES ($1, $2) ON CONFLICT (user_id, product_id) DO NOTHING RETURNING wish_list_id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var wishListId int

	err = stmt.QueryRow(product_id, user_id).Scan(&wishListId)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("product already added to wish list")
	}

	if err != nil {
		return 0, err
	}

	return int(wishListId), nil
}
func (p *Postgres) AddToCart(user_id int, product_id int, quantity int, discount int) (int, error) {
	var cartID int
	err := p.Db.QueryRow(`
		SELECT cart_id FROM cartTable WHERE user_id = $1 AND status = 'active' LIMIT 1
	`, user_id).Scan(&cartID)

	if err == sql.ErrNoRows {
		err = p.Db.QueryRow(`
			INSERT INTO cartTable (user_id, status) VALUES ($1, 'active') RETURNING cart_id
		`, user_id).Scan(&cartID)
		if err != nil {
			return 0, fmt.Errorf("failed to create cart: %w", err)
		}
	} else if err != nil {
		return 0, fmt.Errorf("failed to fetch cart: %w", err)
	}

	var price float64
	var stock int
	err = p.Db.QueryRow(`
		SELECT price, stock FROM products WHERE product_id = $1
	`, product_id).Scan(&price, &stock)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch product details: %w", err)
	}

	if quantity > stock {
		return 0, fmt.Errorf("cannot add %d items, only %d available in stock", quantity, stock)
	}

	var existingQty int
	var cartItemID int
	err = p.Db.QueryRow(`
		SELECT cart_item_id, quantity FROM cartItems WHERE cart_id = $1 AND product_id = $2
	`, cartID, product_id).Scan(&cartItemID, &existingQty)

	if err == nil && existingQty+quantity > stock {
		return 0, fmt.Errorf("cannot add %d more, only %d left in stock", quantity, stock-existingQty)
	}

	subtotal := (price * float64(quantity)) - float64(discount)
	if subtotal < 0 {
		subtotal = 0
	}

	switch err {
	case sql.ErrNoRows:
		err = p.Db.QueryRow(`
			INSERT INTO cartItems (cart_id, product_id, quantity, price_at_time, discount, subtotal)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING cart_item_id
		`, cartID, product_id, quantity, price, discount, subtotal).Scan(&cartItemID)
		if err != nil {
			return 0, fmt.Errorf("failed to add item: %w", err)
		}
	case nil:
		_, err = p.Db.Exec(`
			UPDATE cartItems
			SET quantity = quantity + $1,
			    discount = $2,
			    subtotal = subtotal + $3
			WHERE cart_item_id = $4
		`, quantity, discount, subtotal, cartItemID)
		if err != nil {
			return 0, fmt.Errorf("failed to update cart item: %w", err)
		}
	default:
		return 0, fmt.Errorf("failed to check cart item: %w", err)
	}

	return cartItemID, nil
}
