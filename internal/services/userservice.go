package services

import (
	"context"

	"github.com/rycln/loyalsys/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type userStorager interface {
	AddUser(context.Context, *models.UserDB) (models.UserID, error)
	GetUserByLogin(context.Context, string) (*models.UserDB, error)
}

type passwordHasher interface {
	Hash(string) (string, error)
	Compare(string, string) error
}

type UserService struct {
	strg   userStorager
	hasher passwordHasher
}

func NewUserService(strg userStorager, hasher passwordHasher) *UserService {
	return &UserService{
		strg:   strg,
		hasher: hasher,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (models.UserID, error) {
	hash, err := s.hasher.Hash(user.Password)
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
	err = s.hasher.Compare(userDB.PasswordHash, user.Password)
	if err != nil {
		return 0, err
	}
	return models.UserID(userDB.ID), nil
}
