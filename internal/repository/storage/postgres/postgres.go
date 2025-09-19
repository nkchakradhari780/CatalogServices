package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/nkchakradhari780/catalogServices/internal/cache"
	"github.com/nkchakradhari780/catalogServices/internal/config"
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

func InvalidateProductCache() {
	err := cache.Rdb.FlushDB(cache.Ctx).Err()
	if err != nil {
		fmt.Println("Error Clearing Cache: ", err)
		return
	}

	fmt.Println("Data Erased from cache memory")
}



