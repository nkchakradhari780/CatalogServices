package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nkchakradhari780/catalogServices/internal/api"
	"github.com/nkchakradhari780/catalogServices/internal/cache"
	"github.com/nkchakradhari780/catalogServices/internal/config"
	"github.com/nkchakradhari780/catalogServices/internal/repository/storage/postgres"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()
	//Database Setup

	storage, err := postgres.New(cfg)
	if err != nil {
		log.Fatalf("Failed to Connect to database %s", err)
	}

	slog.Info("Connected to Database") 
	cache.InitRedis()

	//Router Setup
	router := http.NewServeMux() 

	router.HandleFunc("POST /admin/products", api.CreateNewProduct(storage))
	router.HandleFunc("PUT /admin/products/{id}", api.UpdateProductById(storage))
	router.HandleFunc("DELETE /admin/products/{id}", api.DeleteProductById(storage))
	router.HandleFunc("GET /products/{id}", api.GetProductById(storage))
	router.HandleFunc("GET /products/", api.GetProducts(storage))
	router.HandleFunc("GET /products/default", api.GetDefaultProducts(storage))
	router.HandleFunc("GET /products/filtered", api.GetFilteredProducts(storage))
	router.HandleFunc("GET /products/search", api.SearcProducts(storage))


	//Server Setup
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server %s\n", err.Error())
		}

	}()

	<-done
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server exited properly")

}
