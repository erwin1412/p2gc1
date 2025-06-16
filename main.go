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

// CORS Middleware Global
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS Headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Jika preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Lanjut ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load .env file
	_ = godotenv.Load()

	// Connect ke database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Inisialisasi dependency
	repo := repository.NewEmployeeRepository(db)
	svc := service.NewEmployeeService(repo)
	h := handler.NewEmployeeHandler(svc)

	// Routing menggunakan DefaultServeMux
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

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	fmt.Printf("Server started at :%s\n", port)

	// Bungkus semua handler dengan CORS dan jalankan server
	log.Fatal(http.ListenAndServe(":"+port, withCORS(http.DefaultServeMux)))
}
