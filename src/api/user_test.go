package api

import (
	"errors"
	"testing"
	"user-service/src/models"
	"user-service/src/storage"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, &UserServiceTestSuite{})
}

type UserServiceTestSuite struct {
	suite.Suite
	mongoMock *storage.MongoMock
	redisMock *storage.RedisMock
	userService *UserService
}

func (s *UserServiceTestSuite) SetupTest() {
	s.mongoMock = storage.NewMongoMock()
	s.redisMock = storage.NewRedisMock()
	s.userService = NewUserService(s.mongoMock, s.redisMock)
}

func (s *UserServiceTestSuite) TestCreateUser() {
	expectedUser := &models.User{
		FirstName:            "Alice",
		LastName:             "Bob",
		Nickname:             "abobby27",
		Password:             "hunter2",
		Email:                "abobby27@mail.com",
		Country:              "LT",
	}

	s.mongoMock.On("CreateUser", expectedUser).Return(nil)
	s.redisMock.On("PublishMessage", models.MessageUserCreated, expectedUser).Return(nil)

	user, err := s.userService.CreateUser(expectedUser)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedUser.FirstName, user.FirstName)
	assert.NotNil(s.T(), user.Id)

	s.mongoMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestCreateUserErr() {
	expectedUser := &models.User{
		FirstName:            "Alice",
		LastName:             "Bob",
		Nickname:             "abobby27",
		Password:             "hunter2",
		Email:                "abobby27@mail.com",
		Country:              "LT",
	}

	s.mongoMock.On("CreateUser", expectedUser).Return(errors.New("couldn't create user"))

	user, err := s.userService.CreateUser(expectedUser)

	assert.Nil(s.T(), user)
	assert.NotNil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestCreateUserPublishErr() {
	expectedUser := &models.User{
		FirstName:            "Alice",
		LastName:             "Bob",
		Nickname:             "abobby27",
		Password:             "hunter2",
		Email:                "abobby27@mail.com",
		Country:              "LT",
	}

	s.mongoMock.On("CreateUser", expectedUser).Return(nil)
	s.redisMock.On("PublishMessage", models.MessageUserCreated, expectedUser).Return(nil)


	user, err := s.userService.CreateUser(expectedUser)

	assert.Nil(s.T(), user)
	assert.NotNil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
}
