package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"business-mgmt-backend/internal/models"

	"github.com/google/uuid"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) CreateUser(
	ctx context.Context,
	u *models.User,
) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO users (
			id, email, password_hash, full_name,
			phone, status, created_at, updated_at
		)
		VALUES ($1, LOWER($2), $3, $4, $5, $6, NOW(), NOW())
	`, u.ID, u.Email, u.PasswordHash, u.FullName, u.Phone, u.Status)
	return err
}

func (r *PostgresRepository) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.User, error) {
	var u models.User
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, email, password_hash, full_name, phone,
		       status, last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FullName,
		&u.Phone, &u.Status, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	u.Roles = []string{}
	rows, err := r.DB.QueryContext(ctx, `
		SELECT code FROM roles
		WHERE id IN (
			SELECT role_id FROM user_roles WHERE user_id = $1
		)
		ORDER BY code
	`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			rows.Close()
			return nil, err
		}
		u.Roles = append(u.Roles, code)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	u.Stores = []models.Store{}
	rows, err = r.DB.QueryContext(ctx, `
		SELECT s.id, s.name, s.code
		FROM stores s
		JOIN user_stores us ON us.store_id = s.id
		WHERE us.user_id = $1
		ORDER BY s.name
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var st models.Store
		if err := rows.Scan(&st.ID, &st.Name, &st.Code); err != nil {
			return nil, err
		}
		u.Stores = append(u.Stores, st)
	}
	return &u, rows.Err()
}

func (r *PostgresRepository) UpdateUser(
	ctx context.Context,
	id uuid.UUID,
	fullName, phone *string,
) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE users
		SET full_name = COALESCE($2, full_name),
		    phone = COALESCE($3, phone),
		    updated_at = NOW()
		WHERE id = $1
	`, id, fullName, phone)
	return err
}

func (r *PostgresRepository) UpdateUserStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE users SET status = $2, updated_at = NOW()
		WHERE id = $1
	`, id, status)
	return err
}

func (r *PostgresRepository) AssignRoles(
	ctx context.Context,
	id uuid.UUID,
	codes []string,
) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx,
		`DELETE FROM user_roles WHERE user_id = $1`, id); err != nil {
		return err
	}

	for _, code := range codes {
		code = strings.ToUpper(strings.TrimSpace(code))
		var roleID uuid.UUID

		err = tx.QueryRowContext(ctx,
			`SELECT id FROM roles WHERE code = $1`, code,
		).Scan(&roleID)
		if err != nil {
			return fmt.Errorf("unknown role %q: %w", code, err)
		}

		if _, err = tx.ExecContext(ctx, `
			INSERT INTO user_roles (user_id, role_id)
			VALUES ($1, $2) ON CONFLICT DO NOTHING
		`, id, roleID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *PostgresRepository) AssignStores(
	ctx context.Context,
	id uuid.UUID,
	ids []uuid.UUID,
) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx,
		`DELETE FROM user_stores WHERE user_id = $1`, id); err != nil {
		return err
	}

	for _, storeID := range ids {
		var active bool
		err = tx.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM stores
				WHERE id = $1 AND status = 'ACTIVE'
			)
		`, storeID).Scan(&active)
		if err != nil {
			return err
		}
		if !active {
			return fmt.Errorf("store %s does not exist or is inactive", storeID)
		}

		if _, err = tx.ExecContext(ctx, `
			INSERT INTO user_stores (user_id, store_id)
			VALUES ($1, $2) ON CONFLICT DO NOTHING
		`, id, storeID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *PostgresRepository) ListUsers(
	ctx context.Context,
	page, pageSize int,
	filters map[string]string,
) ([]models.User, int, error) {
	where := []string{"1=1"}
	args := []any{}

	add := func(condition string, value any) {
		args = append(args, value)
		where = append(where, fmt.Sprintf(condition, len(args)))
	}

	if v := filters["status"]; v != "" {
		add("u.status = $%d", v)
	}
	if v := filters["role"]; v != "" {
		add(`EXISTS (
			SELECT 1 FROM user_roles ur
			JOIN roles ro ON ro.id = ur.role_id
			WHERE ur.user_id = u.id AND ro.code = $%d
		)`, strings.ToUpper(v))
	}
	if v := filters["store_id"]; v != "" {
		add(`EXISTS (
			SELECT 1 FROM user_stores us
			WHERE us.user_id = u.id AND us.store_id = $%d::uuid
		)`, v)
	}
	if v := filters["q"]; v != "" {
		args = append(args, "%"+v+"%")
		n := len(args)
		where = append(where,
			fmt.Sprintf("(u.email ILIKE $%d OR u.full_name ILIKE $%d)", n, n))
	}

	clause := strings.Join(where, " AND ")

	var total int
	err := r.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM users u WHERE "+clause,
		args...,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listArgs := append([]any{}, args...)
	listArgs = append(listArgs, pageSize, (page-1)*pageSize)

	query := fmt.Sprintf(`
		SELECT u.id
		FROM users u
		WHERE %s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d
	`, clause, len(args)+1, len(args)+2)

	rows, err := r.DB.QueryContext(ctx, query, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, 0, err
		}
		ids = append(ids, id)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, 0, err
	}

	users := make([]models.User, 0, len(ids))
	for _, id := range ids {
		u, err := r.GetUserByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		// Never return password hashes from a list endpoint.
		u.PasswordHash = ""
		users = append(users, *u)
	}

	return users, total, nil
}
