package main

import (
	"crud/config"
	"crud/handlers"
	"crud/repository"
	"crud/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// ১. কনফিগারেশন লোড
	cfg := config.LoadConfig()
	// data base initialize
	db := config.InitDB(cfg.DBDSN)
	defer db.Close()

	//handler initialize
	// ২. ডিপেনডেন্সি ওয়্যারিং (Clean Architecture Chain)
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)

	// route setup
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", userHandler.GetAll)
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)
	mux.HandleFunc("PUT /users/{id}", userHandler.Update)

	fmt.Println("Server running clean on http://localhost:8080")
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
