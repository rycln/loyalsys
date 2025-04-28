package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

const testPasswordHash = "abcdefg"

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockuserStorager(ctrl)
	mHasher := mocks.NewMockpasswordHasher(ctrl)

	t.Run("valid test", func(t *testing.T) {
		testUser := &models.User{
			Login:    "test",
			Password: "secret",
		}

		mStrg.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(testUserID, nil)
		mHasher.EXPECT().Hash(testUser.Password).Return(testPasswordHash, nil)

		s := NewUserService(mStrg, mHasher)
		uid, err := s.CreateUser(context.Background(), testUser)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, uid)
	})

	t.Run("AddUser error", func(t *testing.T) {
		testUser := &models.User{
			Login:    "test",
			Password: "secret",
		}

		mStrg.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(models.UserID(0), errTest)
		mHasher.EXPECT().Hash(testUser.Password).Return(testPasswordHash, nil)

		s := NewUserService(mStrg, mHasher)
		_, err := s.CreateUser(context.Background(), testUser)
		assert.Error(t, err)
	})

	t.Run("password hash failed", func(t *testing.T) {
		testUser := &models.User{
			Login:    "test",
			Password: "wrong_password",
		}

		mHasher.EXPECT().Hash(testUser.Password).Return("", errTest)

		s := NewUserService(mStrg, mHasher)
		_, err := s.CreateUser(context.Background(), testUser)
		assert.Error(t, err)
	})
}

func TestUserService_UserAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockuserStorager(ctrl)
	mHasher := mocks.NewMockpasswordHasher(ctrl)

	t.Run("valid test", func(t *testing.T) {
		testUser := &models.User{
			Login:    "test",
			Password: "secret",
		}

		testUserDB := &models.UserDB{
			ID:           testUserID,
			PasswordHash: testPasswordHash,
		}

		mStrg.EXPECT().GetUserByLogin(context.Background(), testUser.Login).Return(testUserDB, nil)
		mHasher.EXPECT().Compare(testUserDB.PasswordHash, testUser.Password).Return(nil)

		s := NewUserService(mStrg, mHasher)
		uid, err := s.UserAuth(context.Background(), testUser)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, uid)
	})

	t.Run("GetUserByLogin error", func(t *testing.T) {
		testUser := &models.User{
			Login:    "test",
			Password: "secret",
		}

		mStrg.EXPECT().GetUserByLogin(context.Background(), testUser.Login).Return(nil, errors.New("test err"))

		s := NewUserService(mStrg, mHasher)
		_, err := s.UserAuth(context.Background(), testUser)
		assert.Error(t, err)
	})

	t.Run("password hash is not the same", func(t *testing.T) {
		testUser := &models.User{
			Login:    "test",
			Password: "secret",
		}
		testUserDB := &models.UserDB{
			ID:           testUserID,
			PasswordHash: "wrong hash",
		}
		mStrg.EXPECT().GetUserByLogin(context.Background(), testUser.Login).Return(testUserDB, nil)
		mHasher.EXPECT().Compare(testUserDB.PasswordHash, testUser.Password).Return(errTest)

		s := NewUserService(mStrg, mHasher)
		_, err := s.UserAuth(context.Background(), testUser)
		assert.Error(t, err)
	})
}
