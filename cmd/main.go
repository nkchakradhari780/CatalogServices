package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nkchakradhari780/catalogServices/internal/config"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()
	//Database Setup
	//Router Setup
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to Catalog Services API"))
	})
	//Server Setup
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	fmt.Println("Server Started")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server %s\n", err.Error())
	}

}
