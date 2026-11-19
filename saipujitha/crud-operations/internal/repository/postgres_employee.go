package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saipujitha/crud-operations/internal/models"
)

var (
	ErrNotFound = errors.New("employee not found")
	ErrConflict = errors.New("employee email already exists")
)

type EmployeeRepository interface {
	List(ctx context.Context) ([]models.Employee, error)
	Get(ctx context.Context, id int64) (models.Employee, error)
	Create(ctx context.Context, input models.EmployeeInput) (models.Employee, error)
	Update(ctx context.Context, id int64, input models.EmployeeInput) (models.Employee, error)
	Delete(ctx context.Context, id int64) error
}

type PostgresEmployeeRepository struct {
	db *pgxpool.Pool
}

func NewPostgresEmployeeRepository(db *pgxpool.Pool) *PostgresEmployeeRepository {
	return &PostgresEmployeeRepository{db: db}
}

const employeeColumns = `id, name, email, department,
	to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),
	to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"')`

func (repo *PostgresEmployeeRepository) List(ctx context.Context) ([]models.Employee, error) {
	rows, err := repo.db.Query(ctx, `SELECT `+employeeColumns+` FROM employees ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query employees: %w", err)
	}
	defer rows.Close()

	employees := make([]models.Employee, 0)
	for rows.Next() {
		employee, err := scanEmployee(rows)
		if err != nil {
			return nil, fmt.Errorf("scan employee: %w", err)
		}
		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate employees: %w", err)
	}
	return employees, nil
}

func (repo *PostgresEmployeeRepository) Get(ctx context.Context, id int64) (models.Employee, error) {
	employee, err := scanEmployee(repo.db.QueryRow(ctx,
		`SELECT `+employeeColumns+` FROM employees WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Employee{}, ErrNotFound
	}
	if err != nil {
		return models.Employee{}, fmt.Errorf("get employee: %w", err)
	}
	return employee, nil
}

func (repo *PostgresEmployeeRepository) Create(ctx context.Context, input models.EmployeeInput) (models.Employee, error) {
	employee, err := scanEmployee(repo.db.QueryRow(ctx, `
		INSERT INTO employees (name, email, department)
		VALUES ($1, $2, $3)
		RETURNING `+employeeColumns, input.Name, input.Email, input.Department))
	if isUniqueViolation(err) {
		return models.Employee{}, ErrConflict
	}
	if err != nil {
		return models.Employee{}, fmt.Errorf("create employee: %w", err)
	}
	return employee, nil
}

func (repo *PostgresEmployeeRepository) Update(ctx context.Context, id int64, input models.EmployeeInput) (models.Employee, error) {
	employee, err := scanEmployee(repo.db.QueryRow(ctx, `
		UPDATE employees
		SET name = $1, email = $2, department = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING `+employeeColumns, input.Name, input.Email, input.Department, id))
	if isUniqueViolation(err) {
		return models.Employee{}, ErrConflict
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Employee{}, ErrNotFound
	}
	if err != nil {
		return models.Employee{}, fmt.Errorf("update employee: %w", err)
	}
	return employee, nil
}

func (repo *PostgresEmployeeRepository) Delete(ctx context.Context, id int64) error {
	tag, err := repo.db.Exec(ctx, `DELETE FROM employees WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete employee: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanEmployee(row rowScanner) (models.Employee, error) {
	var employee models.Employee
	err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.Email,
		&employee.Department,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	return employee, err
}

func isUniqueViolation(err error) bool {
	var postgresErr *pgconn.PgError
	return errors.As(err, &postgresErr) && postgresErr.Code == "23505"
}
