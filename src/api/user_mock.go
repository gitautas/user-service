package api

import (
	"user-service/src/models"

	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{}
}

func (um *UserServiceMock) CreateUser(user *models.User) (*models.User, *models.Status) {
	args := um.Called(user)

	result := args.Get(0)

	user, ok := result.(*models.User)
	if ok {
		return user, nil
	}

	return nil, args.Get(1).(*models.Status)
}

func (um *UserServiceMock) UpdateUser(user *models.User) (*models.User, *models.Status) {
	args := um.Called(user)

	result := args.Get(0)

	user, ok := result.(*models.User)
	if ok {
		return user, nil
	}

	return nil, args.Get(1).(*models.Status)
}

func (um *UserServiceMock) RemoveUser(userID string) *models.Status {
	args := um.Called(userID)

	result := args.Get(0)

	err, ok := result.(*models.Status)
	if ok {
		return err
	}

	return nil
}

func (um *UserServiceMock) GetUser(userID string) (user *models.User, status *models.Status) {
	args := um.Called(userID)

	result := args.Get(0)

	user, ok := result.(*models.User)
	if ok {
		return user, nil
	}

	return nil, args.Get(1).(*models.Status)
}

func (um *UserServiceMock) GetUserList(limit int, offset int, filter map[string]string) (users []*models.User, status *models.Status) {
	args := um.Called(limit, offset, filter)

	result := args.Get(0)

	users, ok := result.([]*models.User)
	if ok {
		return users, nil
	}

	return nil, args.Get(1).(*models.Status)
}
