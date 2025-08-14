package usecase

import (
	"context"

	"apiprofile/internal/domain"
	"apiprofile/internal/domain/models"
	"apiprofile/internal/dto"
	"apiprofile/internal/helper"
)

type UserUsecase interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetAll(ctx context.Context) ([]dto.UserResponse, error)
	GetByID(ctx context.Context, id string) (*dto.UserResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id string) error
}

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) UserUsecase { return &userUsecase{repo: r} }

func (u *userUsecase) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	hashed, _ := helper.HashPassword(req.Password)
	user := &models.User{ID: "", Name: req.Name, Email: req.Email, Password: hashed, Role: req.Role}
	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return &dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}, nil
}

func (u *userUsecase) GetAll(ctx context.Context) ([]dto.UserResponse, error) {
	list, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []dto.UserResponse
	for _, it := range list {
		out = append(out, dto.UserResponse{ID: it.ID, Name: it.Name, Email: it.Email, Role: it.Role, CreatedAt: it.CreatedAt, UpdatedAt: it.UpdatedAt})
	}
	return out, nil
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}, nil
}

func (u *userUsecase) Update(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Password != "" {
		h, _ := helper.HashPassword(req.Password)
		user.Password = h
	}
	if err := u.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return &dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}, nil
}

func (u *userUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
