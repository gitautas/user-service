package storage

import (
	"user-service/src/models"

	"github.com/stretchr/testify/mock"
)

type MongoMock struct {
	mock.Mock
}

func NewMongoMock() *MongoMock {
	return &MongoMock{}
}

func (m *MongoMock) Connect() *models.Status {
	args := m.Called()

	result := args.Get(0)
	status, ok := result.(*models.Status)
	if ok {
		return status
	}

	return nil
}

func (m *MongoMock) Disconnect() *models.Status {
	args := m.Called()
	result := args.Get(0)
	status, ok := result.(*models.Status)
	if ok {
		return status
	}

	return nil
}

func (m *MongoMock) CreateUser(user *models.User) *models.Status {
	args := m.Called(user)
	result := args.Get(0)
	status, ok := result.(*models.Status)
	if ok {
		return status
	}

	return nil
}

func (m *MongoMock) UpdateUser(user *models.User) *models.Status {
	args := m.Called(user)
	result := args.Get(0)
	status, ok := result.(*models.Status)
	if ok {
		return status
	}

	return nil
}

func (m *MongoMock) DeleteUser(userID string) *models.Status {
	args := m.Called(userID)
	result := args.Get(0)
	status, ok := result.(*models.Status)
	if ok {
		return status
	}

	return nil
}

func (m *MongoMock) GetUser(userID string) (user *models.User, err *models.Status) {
	args := m.Called(userID)

	result := args.Get(0)
	user, ok := result.(*models.User)
	if ok {
		return user, nil
	}

	return nil, args.Get(1).(*models.Status)
}

func (m *MongoMock) GetUserList(limit int, skip int, filter map[string]string) (users []*models.User, err *models.Status) {
	args := m.Called(limit, skip, filter)

	result := args.Get(0)
	users, ok := result.([]*models.User)
	if ok {
		return users, nil
	}

	return nil, args.Get(1).(*models.Status)
}
