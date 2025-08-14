package gorm

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"apiprofile/internal/domain"
	"apiprofile/internal/domain/models"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	// Use Exec with native SQL to insert
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO users (id, name, email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, now(), now())",
		u.ID, u.Name, u.Email, u.Password, u.Role).Error
}

// func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
// 	var u models.User
// 	// using Raw and Scan
// 	if err := r.db.WithContext(ctx).
// 		Raw("SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email = ?", email).
// 		Scan(&u).Error; err != nil {
// 		return nil, err
// 	}
// 	return &u, nil
// }

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).
		Raw("SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email = ?", email).
		Scan(&u).Error
	if err != nil {
		return nil, err
	}

	// kalau tidak ada data, return nil, nil
	if u.ID == "" {
		return nil, nil
	}

	return &u, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).
		Raw("SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetAll(ctx context.Context) ([]models.User, error) {
	var list []models.User
	if err := r.db.WithContext(ctx).
		Raw("SELECT id, name, email, password, role, created_at, updated_at FROM users").
		Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *userRepo) Update(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).
		Exec("UPDATE users SET name = ?, email = ?, password = ?, role = ?, updated_at = now() WHERE id = ?", u.Name, u.Email, u.Password, u.Role, u.ID).
		Error
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Exec("DELETE FROM users WHERE id = ?", id).Error
}

// Refresh tokens
func (r *userRepo) CreateRefreshToken(ctx context.Context, token string, userID string, expiresAt time.Time) error {
	id := uuid.NewString()
	return r.db.WithContext(ctx).
		Exec("INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at) VALUES (?, ?, ?, ?, now())", id, userID, token, expiresAt).
		Error
}

func (r *userRepo) GetRefreshToken(ctx context.Context, token string) (*domain.RefreshTokenRow, error) {
	var row domain.RefreshTokenRow
	if err := r.db.WithContext(ctx).
		Raw("SELECT token, user_id, expires_at FROM refresh_tokens WHERE token = ?", token).
		Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *userRepo) UpdateRefreshToken(ctx context.Context, oldToken string, newToken string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).
		Exec("UPDATE refresh_tokens SET token = ?, expires_at = ? WHERE token = ?", newToken, expiresAt, oldToken).
		Error
}

func (r *userRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Exec("DELETE FROM refresh_tokens WHERE token = ?", token).Error
}
