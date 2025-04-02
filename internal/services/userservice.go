package services

import (
	"context"

	"github.com/rycln/loyalsys/internal/auth"
	"github.com/rycln/loyalsys/internal/models"
)

type userStorager interface {
	AddUser(context.Context, *models.UserDB) (models.UserID, error)
	GetUserByLogin(context.Context, string) (*models.UserDB, error)
}

type UserService struct {
	strg userStorager
}

func NewUserService(strg userStorager) *UserService {
	return &UserService{strg: strg}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (models.UserID, error) {
	hash, err := auth.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	userDB := &models.UserDB{
		Login:        user.Login,
		PasswordHash: hash,
	}
	uid, err := s.strg.AddUser(ctx, userDB)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func (s *UserService) UserAuth(ctx context.Context, user *models.User) (models.UserID, error) {
	userDB, err := s.strg.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return 0, err
	}
	err = auth.CompareHashAndPassword(userDB.PasswordHash, user.Password)
	if err != nil {
		return 0, err
	}
	return models.UserID(userDB.ID), nil
}
