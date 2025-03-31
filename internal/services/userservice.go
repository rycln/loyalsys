package services

import (
	"context"

	"github.com/rycln/loyalsys/internal/auth"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/storage"
)

type userStorager interface {
	AddUser(context.Context, storage.UserDB) error
	GetPasswordHashByLogin(context.Context, string) (string, error)
}

type UserService struct {
	strg userStorager
}

func NewUserService(strg userStorager) *UserService {
	return &UserService{strg: strg}
}

func (us *UserService) CreateUser(ctx context.Context, user *models.User) error {
	hash, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}
	userDB := &storage.UserDB{
		Login:        user.Login,
		PasswordHash: hash,
	}
	err = us.strg.AddUser(ctx, *userDB)
	if err != nil {
		return err
	}
	return nil
}

// добавить отдельные ошибки для отсутсвтия юзера и на неправильный пароль
func (us *UserService) UserAuth(ctx context.Context, user *models.User) error {
	hash, err := us.strg.GetPasswordHashByLogin(ctx, user.Login)
	if err != nil {
		return err
	}
	err = auth.CompareHashAndPassword(hash, user.Password)
	if err != nil {
		return err
	}
	return nil
}
