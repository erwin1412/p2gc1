package service

import (
	"errors"
	"p2gc1/model"
	"p2gc1/repository"
	"strings"
)

type EmployeeService struct {
	Repo *repository.EmployeeRepository
}

func NewEmployeeService(repo *repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{Repo: repo}
}

func (s *EmployeeService) CreateEmployee(emp model.Employee) (int64, error) {
	// Validasi input
	if strings.TrimSpace(emp.Name) == "" || strings.TrimSpace(emp.Email) == "" || strings.TrimSpace(emp.Phone) == "" {
		return 0, errors.New("name, email, and phone cannot be empty")
	}

	// Cek apakah email sudah ada
	employees, err := s.Repo.GetAll()
	if err != nil {
		return 0, err
	}
	for _, e := range employees {
		if strings.EqualFold(e.Email, emp.Email) {
			return 0, errors.New("email already exists")
		}
	}

	return s.Repo.Create(emp)
}

func (s *EmployeeService) GetEmployeeByID(id int) (*model.Employee, error) {
	emp, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}
	return emp, nil
}

func (s *EmployeeService) GetAllEmployees() ([]model.Employee, error) {
	return s.Repo.GetAll()
}

func (s *EmployeeService) UpdateEmployee(id int, emp model.Employee) (*model.Employee, error) {
	// Validasi input
	if strings.TrimSpace(emp.Name) == "" || strings.TrimSpace(emp.Email) == "" || strings.TrimSpace(emp.Phone) == "" {
		return nil, errors.New("name, email, and phone cannot be empty")
	}

	// Cek apakah karyawan ada
	if _, err := s.Repo.GetByID(id); err != nil {
		return nil, errors.New("employee not found")
	}

	// Cek email jika diubah dan sudah dipakai orang lain
	employees, err := s.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	for _, e := range employees {
		if e.ID != id && strings.EqualFold(e.Email, emp.Email) {
			return nil, errors.New("email already exists")
		}
	}

	err = s.Repo.Update(id, emp)
	if err != nil {
		return nil, err
	}

	return s.Repo.GetByID(id) // return updated employee
}

func (s *EmployeeService) DeleteEmployee(id int) (*model.Employee, error) {
	emp, err := s.Repo.Delete(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}
	return emp, nil
}
