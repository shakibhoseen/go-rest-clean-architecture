package main

import (
	"crud/config"
	"crud/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// ১. কনফিগারেশন লোড
	cfg := config.LoadConfig()
	// data base initialize
	db := config.InitDB(cfg.DBPath)
	defer db.Close()

	//handler initialize
	userHandler := handlers.NewUserHandler(db)

	// route setup
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", userHandler.GetAllUsers)
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)
	mux.HandleFunc("PUT /users/{id}", userHandler.Update)

	fmt.Println("Server running clean on http://localhost:8080")
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
