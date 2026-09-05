package service

import (
	"context"
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		Repo: repo,
	}
}

func (s *UserServiceImpl) GetProfile(ctx context.Context, userID uint) (*response.UserResponse, error) {
	user, err := s.Repo.FindUserByID(ctx, userID)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &response.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	}, nil
}

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID uint, req request.UpdateUserRequest) (*response.UserResponse, error) {
	user, err := s.Repo.FindUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	user.Name = req.Name
	user.PhoneNumber = req.PhoneNumber

	if err := s.Repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	}, nil
}

func (s *UserServiceImpl) DeleteProfile(ctx context.Context, userID uint) error {
	if _, err := s.Repo.FindUserByID(ctx, userID); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}

		return err
	}

	return s.Repo.DeleteUser(ctx, userID)
}
