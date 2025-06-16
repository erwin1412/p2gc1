package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"p2gc1/config"
	"p2gc1/handler"
	"p2gc1/repository"
	"p2gc1/service"

	"github.com/joho/godotenv"
)

func main() {
	// Connect DB
	_ = godotenv.Load() // Load .env file

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Dependency Injection
	repo := repository.NewEmployeeRepository(db)
	svc := service.NewEmployeeService(repo)
	h := handler.NewEmployeeHandler(svc)

	// Routing
	http.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetAllEmployees(w, r)
		case http.MethodPost:
			h.CreateEmployee(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/employees/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetEmployeeByID(w, r)
		case http.MethodPut:
			h.UpdateEmployee(w, r)
		case http.MethodDelete:
			h.DeleteEmployee(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	fmt.Printf("Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
