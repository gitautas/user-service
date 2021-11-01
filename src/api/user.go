package api

import (
	"log"
	"net/http"
	"user-service/src/models"
	"user-service/src/storage"

	"github.com/google/uuid"
)

// This is where our main business logic lives!

type IUserService interface {
    CreateUser(user *models.User, db storage.Database, pubSub storage.PubSub) (*models.User, *models.Status)
	UpdateUser(user *models.User, db storage.Database, pubSub storage.PubSub) (*models.User, *models.Status)
	RemoveUser(userID string, db storage.Database, pubSub storage.PubSub) *models.Status
	GetUser(userID string, db storage.Database) (user *models.User, status *models.Status)
	GetUserList(limit int, offset int, filter map[string]string, db storage.Database) (users []*models.User, status *models.Status)
}

type UserService struct {
    db storage.Database
	pubSub storage.PubSub
}

func NewUserService(db storage.Database, pubSub storage.PubSub) *UserService {
	return &UserService{
		db:     db,
		pubSub: pubSub,
	}
}


func (us *UserService) CreateUser(user *models.User) (*models.User, *models.Status) {
	user.Password = user.HashPassword(user.Password)
	user.Id = uuid.New().String() // Generate a new UUID.
	timestamp := user.GenerateTimestamp()

	user.CreatedAt = timestamp
	user.UpdatedAt = timestamp

	status := us.db.CreateUser(user)
	if status != nil {
		log.Println(status.Message)
		return nil, status
	}

	status = us.pubSub.PublishMessage(models.MessageUserCreated, user)
	if status != nil {
		log.Println(status.Message)
		return nil, status
	}

	return user, nil
}

func (us *UserService) UpdateUser(user *models.User) (*models.User, *models.Status) {
	oldUser, status := us.GetUser(user.Id)
	if status != nil {
		log.Println(status.Message)
		return nil, status
	}

	user.CreatedAt = oldUser.CreatedAt
	user.UpdatedAt = user.GenerateTimestamp()

	if oldUser == nil {
		return nil, models.NewStatus(http.StatusNotFound, "user not found")
	}

	if user.Password != oldUser.Password {
		user.Password = user.HashPassword(user.Password)
	}

	status = us.db.UpdateUser(user)
	if status != nil {
		log.Println(status.Message)
		return nil, status
	}

	status = us.pubSub.PublishMessage(models.MessageUserUpdated, user)
	if status != nil {
		log.Println(status.Message)
		return nil, status
	}

	return user, nil
}

func (us *UserService) RemoveUser(userID string) *models.Status {
	status := us.db.DeleteUser(userID)
	if status != nil {
		log.Println(status.Message)
		return status
	}

	status = us.pubSub.PublishMessage(models.MessageUserDeleted, &models.User{Id: userID})
	if status != nil {
		log.Println(status.Message)
		return status
	}

	return nil
}

func (us *UserService) GetUser(userID string) (user *models.User, status *models.Status) {
	user, status = us.db.GetUser(userID)
	if status != nil {
		log.Println(status.Message)
		return nil, status
	}

	return user, nil
}

func (us *UserService) GetUserList(limit int, offset int, filter map[string]string) (users []*models.User, status *models.Status) {
	users, status = us.db.GetUserList(limit, offset, filter)
	if status != nil {
		log.Println(status.Message)
		return nil, status
	}

	return users, nil
}
