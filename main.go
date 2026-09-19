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

	//public
	mux.HandleFunc("POST /auth/register", userHandler.Register)
	mux.HandleFunc("POST /auth/login", userHandler.Login)

	//protect with middleware
	mux.HandleFunc("GET /users", handlers.Protected(userHandler.GetAll))
	mux.HandleFunc("POST /users", handlers.Protected(userHandler.Create))
	mux.HandleFunc("GET /users/{id}", handlers.Protected(userHandler.GetByID))
	mux.HandleFunc("DELETE /users/{id}", handlers.Protected(userHandler.Delete))
	mux.HandleFunc("PUT /users/{id}", handlers.Protected(userHandler.Update))

	fmt.Println("Server running clean on http://localhost:8080")
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
