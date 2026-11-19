package models

import "strings"

type Employee struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type EmployeeInput struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
}

func (input EmployeeInput) Normalize() EmployeeInput {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Department = strings.TrimSpace(input.Department)
	return input
}

func (input EmployeeInput) Validate() string {
	switch {
	case input.Name == "":
		return "name is required"
	case input.Email == "":
		return "email is required"
	case !strings.Contains(input.Email, "@"):
		return "email must be a valid email address"
	case input.Department == "":
		return "department is required"
	default:
		return ""
	}
}
