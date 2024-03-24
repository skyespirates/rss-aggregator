package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/skyespirates/rss-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load env variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is undefined")
	}

	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("DB_URL is undefined")
	}

	conn, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Cant connect to database")
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
 
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)


	server := &http.Server{
		Handler: router,
		Addr: 		":" + port,
	}

	log.Printf("Server running on port %v", port) 
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}