package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/auth"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateUser(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mStrg := mocks.NewMockuserStorager(ctrl)

		testUser := &models.User{
			Login:    "test",
			Password: "secret",
		}
		testUserID := models.UserID(1)
		mStrg.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(testUserID, nil)

		s := NewUserService(mStrg)
		uid, err := s.CreateUser(context.Background(), testUser)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, uid)
	})

	t.Run("AddUser error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mStrg := mocks.NewMockuserStorager(ctrl)

		testUser := &models.User{
			Login:    "test",
			Password: "secret",
		}
		mStrg.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(models.UserID(0), errors.New("test error"))

		s := NewUserService(mStrg)
		_, err := s.CreateUser(context.Background(), testUser)
		assert.Error(t, err)
	})

	t.Run("password hash failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mStrg := mocks.NewMockuserStorager(ctrl)

		testUser := &models.User{
			Login:    "test",
			Password: string(make([]byte, 100)),
		}

		s := NewUserService(mStrg)
		_, err := s.CreateUser(context.Background(), testUser)
		assert.Error(t, err)
	})
}

func TestUserService_UserAuth(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mStrg := mocks.NewMockuserStorager(ctrl)

		testUser := &models.User{
			Login:    "test",
			Password: "secret",
		}
		testPasswordHash, err := auth.HashPassword(testUser.Password)
		require.NoError(t, err)
		testUserID := models.UserID(1)
		testUserDB := &models.UserDB{
			ID:           testUserID,
			PasswordHash: testPasswordHash,
		}
		mStrg.EXPECT().GetUserByLogin(context.Background(), testUser.Login).Return(testUserDB, nil)

		s := NewUserService(mStrg)
		uid, err := s.UserAuth(context.Background(), testUser)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, uid)
	})

	t.Run("GetUserByLogin error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mStrg := mocks.NewMockuserStorager(ctrl)

		testUser := &models.User{
			Login:    "test",
			Password: "secret",
		}
		mStrg.EXPECT().GetUserByLogin(context.Background(), testUser.Login).Return(nil, errors.New("test err"))

		s := NewUserService(mStrg)
		_, err := s.UserAuth(context.Background(), testUser)
		assert.Error(t, err)
	})

	t.Run("password hash is not the same", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mStrg := mocks.NewMockuserStorager(ctrl)

		testUser := &models.User{
			Login:    "test",
			Password: "secret",
		}
		testUserID := models.UserID(1)
		testUserDB := &models.UserDB{
			ID:           testUserID,
			PasswordHash: "wrong hash",
		}
		mStrg.EXPECT().GetUserByLogin(context.Background(), testUser.Login).Return(testUserDB, nil)

		s := NewUserService(mStrg)
		_, err := s.UserAuth(context.Background(), testUser)
		assert.Error(t, err)
	})
}
