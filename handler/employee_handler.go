package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"p2gc1/model"
	"p2gc1/service"
)

type EmployeeHandler struct {
	Service *service.EmployeeService
}

func NewEmployeeHandler(service *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{Service: service}
}

// POST /employees
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateEmployee(emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	emp.ID = int(id)
	response := map[string]interface{}{
		"message": "Employee created successfully",
		"data":    emp,
	}
	writeJSON(w, http.StatusCreated, response)
}

// json data post
// {
// 	"Name":  "John Doe",
// 	"Email": "John@gmail.com",
// 	"Phone": "1234567890"
// }

// GET /employees/:id
func (h *EmployeeHandler) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	emp, err := h.Service.GetEmployeeByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, emp)
}

// GET /employees
func (h *EmployeeHandler) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	emps, err := h.Service.GetAllEmployees()
	if err != nil {
		http.Error(w, "Failed to get employees", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, emps)
}

// PUT /employees/:id
func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedEmp, err := h.Service.UpdateEmployee(id, emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message": "Employee updated successfully",
		"data":    updatedEmp,
	}
	writeJSON(w, http.StatusOK, response)
}

// DELETE /employees/:id
func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	deletedEmp, err := h.Service.DeleteEmployee(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"message": "Employee deleted successfully",
		"data":    deletedEmp,
	}
	writeJSON(w, http.StatusOK, response)
}

// Helper: parse /employees/:id
func extractIDFromPath(path string) (int, error) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) < 2 {
		return 0, http.ErrMissingFile
	}
	return strconv.Atoi(segments[1])
}

// json data put
// {
// 	"Name":  "John Doe Updated",
// 	"Email": "John@gmail.com",
// 	"Phone": "0987654321"
// }

// Helper: write JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
