package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/saipujitha/crud-operations/internal/models"
	"github.com/saipujitha/crud-operations/internal/repository"
)

type EmployeeHandler struct {
	repository repository.EmployeeRepository
}

func NewEmployeeHandler(repository repository.EmployeeRepository) *EmployeeHandler {
	return &EmployeeHandler{repository: repository}
}

func (handler *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	employees, err := handler.repository.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list employees")
		return
	}
	writeJSON(w, http.StatusOK, employees)
}

func (handler *EmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := employeeID(w, r)
	if !ok {
		return
	}
	employee, err := handler.repository.Get(r.Context(), id)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, employee)
}

func (handler *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	input, ok := decodeInput(w, r)
	if !ok {
		return
	}
	employee, err := handler.repository.Create(r.Context(), input)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/api/employees/%d", employee.ID))
	writeJSON(w, http.StatusCreated, employee)
}

func (handler *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := employeeID(w, r)
	if !ok {
		return
	}
	input, ok := decodeInput(w, r)
	if !ok {
		return
	}
	employee, err := handler.repository.Update(r.Context(), id, input)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, employee)
}

func (handler *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := employeeID(w, r)
	if !ok {
		return
	}
	if err := handler.repository.Delete(r.Context(), id); err != nil {
		writeRepositoryError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func employeeID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "employee ID must be a positive integer")
		return 0, false
	}
	return id, true
}

func decodeInput(w http.ResponseWriter, r *http.Request) (models.EmployeeInput, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input models.EmployeeInput
	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "request body must be valid employee JSON")
		return models.EmployeeInput{}, false
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "request body must contain a single JSON object")
		return models.EmployeeInput{}, false
	}

	input = input.Normalize()
	if message := input.Validate(); message != "" {
		writeError(w, http.StatusBadRequest, message)
		return models.EmployeeInput{}, false
	}
	return input, true
}

func writeRepositoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		writeError(w, http.StatusNotFound, "employee not found")
	case errors.Is(err, repository.ErrConflict):
		writeError(w, http.StatusConflict, "an employee with that email already exists")
	default:
		writeError(w, http.StatusInternalServerError, "could not complete employee operation")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		return
	}
}
