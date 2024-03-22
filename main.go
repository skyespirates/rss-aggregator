package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load env variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is undefined")
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