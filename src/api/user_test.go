package api

import (
	"net/http"
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
	mongoMock   *storage.MongoMock
	redisMock   *storage.RedisMock
	userService *UserService
}

func (s *UserServiceTestSuite) SetupTest() {
	s.mongoMock = storage.NewMongoMock()
	s.redisMock = storage.NewRedisMock()
	s.userService = NewUserService(s.mongoMock, s.redisMock)
}

func (s *UserServiceTestSuite) TestCreateUser() {
	expectedUser := &models.User{
		FirstName: "Alice",
		LastName:  "Bob",
		Nickname:  "abobby27",
		Password:  "hunter2",
		Email:     "abobby27@mail.com",
		Country:   "LT",
	}

	s.mongoMock.On("CreateUser", expectedUser).Return(nil)
	s.redisMock.On("PublishMessage", models.MessageUserCreated, expectedUser).Return(nil)

	user, err := s.userService.CreateUser(expectedUser)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedUser.FirstName, user.FirstName)
	assert.NotNil(s.T(), user.Id)

	s.mongoMock.AssertExpectations(s.T())
	s.redisMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestCreateUserErr() {
	expectedUser := &models.User{
		FirstName: "Alice",
		LastName:  "Bob",
		Nickname:  "abobby27",
		Password:  "hunter2",
		Email:     "abobby27@mail.com",
		Country:   "LT",
	}

	s.mongoMock.On("CreateUser", expectedUser).Return(models.NewStatus(http.StatusInternalServerError, "could not create user"))

	user, err := s.userService.CreateUser(expectedUser)

	assert.Nil(s.T(), user)
	assert.NotNil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestCreateUserPublishErr() {
	expectedUser := &models.User{
		FirstName: "Alice",
		LastName:  "Bob",
		Nickname:  "abobby27",
		Password:  "hunter2",
		Email:     "abobby27@mail.com",
		Country:   "LT",
	}

	s.mongoMock.On("CreateUser", expectedUser).Return(nil)
	s.redisMock.On("PublishMessage", models.MessageUserCreated, expectedUser).Return(models.NewStatus(http.StatusInternalServerError, "could not publish message"))

	user, err := s.userService.CreateUser(expectedUser)

	assert.Nil(s.T(), user)
	assert.NotNil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
	s.redisMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestUpdateUser() {
	old := &models.User{
		Id:        "e2ee0654-26c5-41b7-8de2-5d937b462fa5",
		FirstName: "Alice",
		LastName:  "Bob",
		Nickname:  "abobby27",
		Password:  "hunter2",
		Email:     "abobby27@mail.com",
		Country:   "LT",
		CreatedAt: "2019-10-12T07:20:50.52Z",
	}

	new := *old
	new.FirstName = "Alissa"
	new.Password = "archbtw"
	new.UpdatedAt = new.GenerateTimestamp()

	s.mongoMock.On("GetUser", old.Id).Return(old)
	s.mongoMock.On("UpdateUser", &new).Return(nil)
	s.redisMock.On("PublishMessage", models.MessageUserUpdated, &new).Return(nil)

	newUser, err := s.userService.UpdateUser(&new)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), newUser)
	assert.Equal(s.T(), old.Id, newUser.Id)
	assert.Equal(s.T(), old.Country, newUser.Country)
	assert.NotEqual(s.T(), old.FirstName, newUser.FirstName)
	assert.NotEqual(s.T(), old.UpdatedAt, newUser.UpdatedAt)
	assert.NotEqual(s.T(), old.Password, newUser.Password)

	s.mongoMock.AssertExpectations(s.T())
	s.redisMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestUpdateUserErr() {
	old := &models.User{
		Id:        "e2ee0654-26c5-41b7-8de2-5d937b462fa5",
		FirstName: "Alice",
		LastName:  "Bob",
		Nickname:  "abobby27",
		Password:  "fspf892fffheeiudjsn",
		Email:     "abobby27@mail.com",
		Country:   "LT",
		CreatedAt: "2019-10-12T07:20:50.52Z",
	}

	new := *old
	new.FirstName = "Alissa"
	new.UpdatedAt = new.GenerateTimestamp()

	s.mongoMock.On("GetUser", old.Id).Return(old)
	s.mongoMock.On("UpdateUser", &new).Return(models.NewStatus(http.StatusInternalServerError, "database error"))

	newUser, err := s.userService.UpdateUser(&new)
	assert.Nil(s.T(), newUser)
	assert.NotNil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestUpdateUserGetUserErr() {
	old := &models.User{
		Id:        "e2ee0654-26c5-41b7-8de2-5d937b462fa5",
		FirstName: "Alice",
		LastName:  "Bob",
		Nickname:  "abobby27",
		Password:  "fspf892fffheeiudjsn",
		Email:     "abobby27@mail.com",
		Country:   "LT",
		CreatedAt: "2019-10-12T07:20:50.52Z",
	}

	new := *old
	new.FirstName = "Alissa"
	new.UpdatedAt = new.GenerateTimestamp()

	s.mongoMock.On("GetUser", old.Id).Return(nil, models.NewStatus(http.StatusNotFound, "user not found"))

	newUser, err := s.userService.UpdateUser(&new)

	assert.Nil(s.T(), newUser)
	assert.NotNil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestUpdateUserPublishErr() {
	old := &models.User{
		Id:        "e2ee0654-26c5-41b7-8de2-5d937b462fa5",
		FirstName: "Alice",
		LastName:  "Bob",
		Nickname:  "abobby27",
		Password:  "hunter2",
		Email:     "abobby27@mail.com",
		Country:   "LT",
		CreatedAt: "2019-10-12T07:20:50.52Z",
	}

	new := *old
	new.FirstName = "Alissa"
	new.Password = "archbtw"
	new.UpdatedAt = new.GenerateTimestamp()

	s.mongoMock.On("GetUser", old.Id).Return(old)
	s.mongoMock.On("UpdateUser", &new).Return(nil)
	s.redisMock.On("PublishMessage", models.MessageUserUpdated, &new).Return(models.NewStatus(http.StatusInternalServerError, "failed to publish message"))

	newUser, err := s.userService.UpdateUser(&new)

	assert.Nil(s.T(), newUser)
	assert.NotNil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
	s.redisMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestRemoveUser() {
	old := &models.User{
		Id: "e2ee0654-26c5-41b7-8de2-5d937b462fa5",
	}

	s.mongoMock.On("DeleteUser", old.Id).Return(nil)
	s.redisMock.On("PublishMessage", models.MessageUserDeleted, old).Return(nil)

	err := s.userService.RemoveUser(old.Id)

	assert.Nil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
	s.redisMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestRemoveUserErr() {
	old := &models.User{
		Id: "e2ee0654-26c5-41b7-8de2-5d937b462fa5",
	}

	s.mongoMock.On("DeleteUser", old.Id).Return(models.NewStatus(http.StatusInternalServerError, "failed to remove user"))

	err := s.userService.RemoveUser(old.Id)

	assert.NotNil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestRemoveUserPublishErr() {
	old := &models.User{
		Id: "e2ee0654-26c5-41b7-8de2-5d937b462fa5",
	}

	s.mongoMock.On("DeleteUser", old.Id).Return(nil)
	s.redisMock.On("PublishMessage", models.MessageUserDeleted, old).Return(models.NewStatus(http.StatusInternalServerError, "failed"))

	err := s.userService.RemoveUser(old.Id)

	assert.NotNil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
	s.redisMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestGetUserList() {
	expectedUsers := []*models.User{
		{
			FirstName: "Alice",
			LastName:  "Bob",
			Nickname:  "abobby27",
			Password:  "hunter2",
			Email:     "abobby27@mail.com",
			Country:   "LT",
		},
		{
			FirstName: "Arina",
			LastName:  "Boubette",
			Nickname:  "abobby12",
			Password:  "hunter111",
			Email:     "abobby7@mail.com",
			Country:   "GB",
		},
	}

	limit := 10
	offset := 10
	filter := make(map[string]string)

	s.mongoMock.On("GetUserList", limit, offset, filter).Return(expectedUsers)

	users, err := s.userService.GetUserList(limit, offset, filter)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), users)
	assert.Equal(s.T(), len(users), 2)

	s.mongoMock.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestGetUserListErr() {
	limit := 10
	offset := 10
	filter := make(map[string]string)

	s.mongoMock.On("GetUserList", limit, offset, filter).Return(nil, models.NewStatus(http.StatusInternalServerError, "failed to get user list"))

	users, err := s.userService.GetUserList(limit, offset, filter)

	assert.Nil(s.T(), users)
	assert.NotNil(s.T(), err)

	s.mongoMock.AssertExpectations(s.T())
}
