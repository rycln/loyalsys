package services

import (
	"context"

	"github.com/rycln/loyalsys/internal/auth"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/storage"
)

type userStorager interface {
	AddUser(context.Context, storage.UserDB) (models.UserID, error)
	GetUserByLogin(context.Context, string) (storage.UserDB, error)
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
	userDB := &storage.UserDB{
		Login:        user.Login,
		PasswordHash: hash,
	}
	id, err := s.strg.AddUser(ctx, *userDB)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// добавить отдельные ошибки для отсутсвтия юзера и на неправильный пароль
func (us *UserService) UserAuth(ctx context.Context, user *models.User) error {
	userDB, err := us.strg.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return err
	}
	err = auth.CompareHashAndPassword(userDB.PasswordHash, user.Password)
	if err != nil {
		return err
	}
	return nil
}
