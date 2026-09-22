package main

// @title           Todo REST API
// @version         1.0
// @description     Production-ready Go REST API with JWT Auth and Clean Architecture.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
import (
	"crud/config"
	_ "crud/docs" // swag init চালানোর পর এই ফোল্ডারটি তৈরি হবে
	"crud/handlers"
	"crud/repository"
	"crud/service"
	"fmt"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
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

	// 1. Todo Wiring
	todoRepo := repository.NewTodoRepository(db)
	todoSvc := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoSvc)

	// 2. Protected Routes
	mux.HandleFunc("POST /todos", handlers.Protected(todoHandler.Create))
	mux.HandleFunc("GET /todos", handlers.Protected(todoHandler.GetAll))
	mux.HandleFunc("DELETE /todos/{id}", handlers.Protected(todoHandler.Delete))
	mux.HandleFunc("PUT /todos/{id}", handlers.Protected(todoHandler.Update))

	// Swagger Documentation Route
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	crosHandler := handlers.EnableCORS(mux)
	globalHandler := handlers.RecoverMiddleware(crosHandler)

	fmt.Println("Server running clean on http://localhost:8080")
	log.Fatal(http.ListenAndServe(cfg.Port, globalHandler))

}
