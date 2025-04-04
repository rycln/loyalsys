package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/services/mocks"
	"github.com/stretchr/testify/assert"
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
}
