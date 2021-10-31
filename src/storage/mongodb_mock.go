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

func (m *MongoMock) Connect() error {
	args := m.Called()

	result := args.Get(0)
	if result == nil {
		return nil
	}

	return args.Error(0)
}

func (m *MongoMock) Disconnect() error {
	args := m.Called()

	result := args.Get(0)
	if result == nil {
		return nil
	}

	return args.Error(0)
}

func (m *MongoMock) CreateUser(user *models.User) error {
	args := m.Called()

	result := args.Get(0)
	if result == nil {
		return nil
	}

	return args.Error(0)
}

func (m *MongoMock) UpdateUser(user *models.User) error {
	args := m.Called()

	result := args.Get(0)
	if result == nil {
		return nil
	}

	return args.Error(0)
}

func (m *MongoMock) DeleteUser(userID string) error {
	args := m.Called()

	result := args.Get(0)
	if result == nil {
		return nil
	}

	return args.Error(0)
}

func (m *MongoMock) GetUser(userID string) (file *models.User, err error) {
	args := m.Called()

	result := args.Get(0)
	user, ok := result.(*models.User)
	if ok {
		return user, nil
	}

	return nil, args.Error(1)
}

func (m *MongoMock) GetUserList(limit int, skip int, filter map[string]string) (users []*models.User, err error) {
	args := m.Called()

	result := args.Get(0)
	users, ok := result.([]*models.User)
	if ok {
		return users, nil
	}

	return nil, args.Error(1)
}
