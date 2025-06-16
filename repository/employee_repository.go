package repository

import (
	"database/sql"
	"p2gc1/model"
)

type EmployeeRepository struct {
	DB *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r *EmployeeRepository) Create(emp model.Employee) (int64, error) {
	result, err := r.DB.Exec("INSERT INTO employees (name, email, phone) VALUES (?, ?, ?)",
		emp.Name, emp.Email, emp.Phone)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *EmployeeRepository) GetByID(id int) (*model.Employee, error) {
	var emp model.Employee
	row := r.DB.QueryRow("SELECT * FROM employees WHERE id = ?", id)
	err := row.Scan(&emp.ID, &emp.Name, &emp.Email, &emp.Phone, &emp.CreatedAt, &emp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepository) GetAll() ([]model.Employee, error) {
	rows, err := r.DB.Query("SELECT id, name, email FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []model.Employee
	for rows.Next() {
		var emp model.Employee
		err := rows.Scan(&emp.ID, &emp.Name, &emp.Email)
		if err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}
	return employees, nil
}

func (r *EmployeeRepository) Update(id int, emp model.Employee) error {
	_, err := r.DB.Exec("UPDATE employees SET name = ?, email = ?, phone = ? WHERE id = ?",
		emp.Name, emp.Email, emp.Phone, id)
	return err
}

func (r *EmployeeRepository) Delete(id int) (*model.Employee, error) {
	emp, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	_, err = r.DB.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return emp, nil
}
