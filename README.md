
---

# 📦 Catalog Services

Catalog Services is a backend service built with **Go**, **PostgreSQL**, and **Redis** that provides a product catalog system with caching support.
It follows a clean architecture pattern with clear separation between API handlers, repository (DB), cache, and domain models.

---

## 🚀 Features

* ✅ Product CRUD operations (Create, Read, Update, Delete)
* ✅ Filtered product queries (brand, price range, category, etc.)
* ✅ Search products by query string
* ✅ Default product listing (random 50 products)
* ✅ Redis caching with TTL (1 week) for frequently accessed data
* ✅ PostgreSQL as persistent storage
* ⚡ Graceful server shutdown with `context`

> 🛠️ **Work in Progress:** Cart and Wishlist modules are defined but not yet implemented.

---

## 📂 Project Structure

```
CATALOGSERVICES/
├── cmd/
│   ├── main.go                  # Application entry point
│   └── tmp/                     # Temporary build artifacts
│
├── config/
│   └── local.yaml               # Config file (DB, Redis, server)
│
├── internal/
│   ├── api/                     # HTTP handlers
│   │   └── handlers.go
│   │
│   ├── cache/                   # Redis client + cache logic
│   │   └── cache.go
│   │
│   ├── config/                  # Config loader
│   │   └── config.go
│   │
│   ├── modules/                 # Domain models (structs)
│   │   ├── cart.go              # (Not implemented yet)
│   │   ├── product.go
│   │   └── wishList.go          # (Not implemented yet)
│   │
│   ├── repository/              # Data access layer
│   │   └── storage/
│   │       ├── postgres.go
│   │       └── storage.go
│   │
│   └── utils/                   # Helpers (response, etc.)
│
├── tmp/                         # Binary output dir
│   └── main.exe
│
├── .gitignore
├── go.mod
├── go.sum
└── Taskfile.yml                 # Task runner config
```

---

## ⚙️ Tech Stack

* **Go (Golang)** – Backend framework
* **PostgreSQL** – Relational database
* **Redis** – In-memory cache
* **Taskfile** – Task runner for building and running

---

## 📑 API Endpoints

### Product APIs

| Method   | Endpoint                  | Description                                       |
| -------- | ------------------------- | ------------------------------------------------- |
| `POST`   | `/admin/products`         | Create a new product                              |
| `PUT`    | `/admin/products/{id}`    | Update product by ID                              |
| `DELETE` | `/admin/products/{id}`    | Delete product by ID                              |
| `GET`    | `/products/{id}`          | Get product by ID (with Redis cache)              |
| `GET`    | `/products/`              | Get all products                                  |
| `GET`    | `/products/default`       | Get 50 random products (with Redis cache)         |
| `GET`    | `/products/filtered`      | Get filtered products (brand, price, stock, etc.) |
| `GET`    | `/products/search?q=text` | Search products by name                           |

---

## 🛠️ Setup Instructions

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/nkchakradhari780/CatalogServices.git
cd CatalogServices
```

### 2️⃣ Setup PostgreSQL

```sql
CREATE DATABASE catalogdb;
```

### 3️⃣ Configure the Project

Edit `config/local.yaml`:

```yaml
env: "dev"

http_server:
  address: "localhost:8081"

database:
  host: "localhost"
  port: 5432
  name: "catalogdb"
  username: --Your User Name--
  password: --Your Password--
  sslmode: "disable"
```

### 4️⃣ Run Redis

Make sure Redis is running:

```bash
redis-server
```

### 5️⃣ Install Dependencies

```bash
go mod tidy
```

### 6️⃣ Run the Application

Using **Taskfile**:

```bash
task dev
```

Or directly:

```bash
go run ./cmd/main.go
```

---

## 🧪 Example Requests

### Create Product

```http
POST http://localhost:8081/admin/products
Content-Type: application/json

{
  "name": "iPhone 15",
  "price": 120000,
  "stock": 10,
  "category_id": "1",
  "brand": "Apple",
  "images": ["https://example.com/iphone15.jpg"]
}
```

### Search Products

```http
GET http://localhost:8081/products/search?q=iphone
```

Response:

```json
[
  {
    "id": 1,
    "name": "iPhone 15",
    "price": 120000,
    "stock": 10,
    "category_id": "1",
    "brand": "Apple",
    "images": ["https://example.com/iphone15.jpg"]
  }
]
```
---

## 📝 License

This project is licensed under the MIT License.

---
