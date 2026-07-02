package main

import (
	"context"
	"go-samurai/internal/blogs"
	"go-samurai/internal/db"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	blogHandler := blogs.NewHandler(blogs.NewService(blogs.NewPostgresRepository(pool)))
	router := chi.NewRouter()
	router.Route("/blogs", blogHandler.Routes)
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
