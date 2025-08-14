package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"apiprofile/internal/domain"
	"apiprofile/internal/dto"
	"apiprofile/internal/domain/models"
	"apiprofile/internal/helper"
)

type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, *dto.UserResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
}

type authUsecase struct {
	repo      domain.UserRepository
	secret    string
	accessTTL int64
	refreshTTL int64
}

func NewAuthUsecase(repo domain.UserRepository, secret string, accessTTL, refreshTTL int64) AuthUsecase {
	return &authUsecase{repo: repo, secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (u *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	existingUser, err := u.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err // DB error
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}
	hashed, _ := helper.HashPassword(req.Password)
	user := &models.User{
		ID:       uuid.NewString(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Role:     "user",
	}
	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return &dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}, nil
}

func (u *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, *dto.UserResponse, error) {
	user, err := u.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}
	if !helper.CheckPasswordHash(req.Password, user.Password) {
		return nil, nil, errors.New("invalid credentials")
	}
	access, err := helper.GenerateJWTWithTTL(user.ID, user.Email, user.Role, u.secret, u.accessTTL)
	if err != nil {
		return nil, nil, err
	}
	refresh := uuid.NewString()
	expiry := time.Now().Add(time.Duration(u.refreshTTL) * time.Second)
	if err := u.repo.CreateRefreshToken(ctx, refresh, user.ID, expiry); err != nil {
		return nil, nil, err
	}
	return &dto.AuthResponse{AccessToken: access, RefreshToken: refresh}, &dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}, nil
}

func (u *authUsecase) Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	rt, err := u.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	if rt.ExpiresAt.Before(time.Now()) {
		// delete expired token
		_ = u.repo.DeleteRefreshToken(ctx, refreshToken)
		return nil, errors.New("refresh token expired")
	}
	user, err := u.repo.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}
	access, err := helper.GenerateJWTWithTTL(user.ID, user.Email, user.Role, u.secret, u.accessTTL)
	if err != nil {
		return nil, err
	}
	// rotate refresh token
	newRefresh := uuid.NewString()
	newExpiry := time.Now().Add(time.Duration(u.refreshTTL) * time.Second)
	if err := u.repo.UpdateRefreshToken(ctx, refreshToken, newRefresh, newExpiry); err != nil {
		return nil, err
	}
	return &dto.AuthResponse{AccessToken: access, RefreshToken: newRefresh}, nil
}
