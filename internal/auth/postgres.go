package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"business-mgmt-backend/internal/models"

	"github.com/google/uuid"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	var u models.User

	err := r.DB.QueryRowContext(ctx, `
		SELECT id, email, password_hash, full_name, phone,
		       status, last_login_at, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FullName,
		&u.Phone, &u.Status, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := r.loadRelations(ctx, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) GetUserByID(
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

	if err := r.loadRelations(ctx, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) loadRelations(
	ctx context.Context,
	u *models.User,
) error {
	u.Roles = []string{}
	u.Permissions = []string{}
	u.Stores = []models.Store{}

	rows, err := r.DB.QueryContext(ctx, `
		SELECT DISTINCT ro.code
		FROM roles ro
		JOIN user_roles ur ON ur.role_id = ro.id
		WHERE ur.user_id = $1
		ORDER BY ro.code
	`, u.ID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			rows.Close()
			return err
		}
		u.Roles = append(u.Roles, code)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}

	rows, err = r.DB.QueryContext(ctx, `
		SELECT DISTINCT p.code
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = $1
		ORDER BY p.code
	`, u.ID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			rows.Close()
			return err
		}
		u.Permissions = append(u.Permissions, code)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}

	rows, err = r.DB.QueryContext(ctx, `
		SELECT s.id, s.name, s.code
		FROM stores s
		JOIN user_stores us ON us.store_id = s.id
		WHERE us.user_id = $1 AND s.status = 'ACTIVE'
		ORDER BY s.name
	`, u.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Store
		if err := rows.Scan(&s.ID, &s.Name, &s.Code); err != nil {
			return err
		}
		u.Stores = append(u.Stores, s)
	}
	return rows.Err()
}

func (r *PostgresUserRepository) UpdateLastLogin(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE users
		SET last_login_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

func (r *PostgresUserRepository) UpdatePassword(
	ctx context.Context,
	id uuid.UUID,
	passwordHash string,
) error {
	result, err := r.DB.ExecContext(ctx, `
		UPDATE users
		SET password_hash = $2, updated_at = NOW()
		WHERE id = $1
	`, id, passwordHash)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

type PostgresSessionRepository struct {
	DB *sql.DB
}

func NewPostgresSessionRepository(db *sql.DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{DB: db}
}

func (r *PostgresSessionRepository) CreateSession(
	ctx context.Context,
	s *models.UserSession,
) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO user_sessions (
			id, user_id, refresh_token_hash, device_info,
			ip_address, expires_at, created_at
		)
		VALUES (
			$1, $2, $3, $4, NULLIF($5, '')::inet, $6, $7
		)
	`, s.ID, s.UserID, s.RefreshTokenHash,
		s.DeviceInfo, ipOnly(s.IPAddress),
		s.ExpiresAt, s.CreatedAt)
	return err
}

func (r *PostgresSessionRepository) GetSessionByTokenHash(
	ctx context.Context,
	hash string,
) (*models.UserSession, error) {
	var s models.UserSession
	var ip sql.NullString
	var revoked sql.NullTime

	err := r.DB.QueryRowContext(ctx, `
		SELECT id, user_id, refresh_token_hash,
		       COALESCE(device_info, ''), host(ip_address),
		       expires_at, revoked_at, created_at
		FROM user_sessions
		WHERE refresh_token_hash = $1
	`, hash).Scan(
		&s.ID, &s.UserID, &s.RefreshTokenHash,
		&s.DeviceInfo, &ip, &s.ExpiresAt,
		&revoked, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if ip.Valid {
		s.IPAddress = ip.String
	}
	if revoked.Valid {
		s.RevokedAt = &revoked.Time
	}
	return &s, nil
}

func (r *PostgresSessionRepository) RevokeSession(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE user_sessions
		SET revoked_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL
	`, id)
	return err
}

func ipOnly(value string) string {
	host, _, err := net.SplitHostPort(value)
	if err == nil {
		return host
	}
	return value
}
