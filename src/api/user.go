package api

import (
	"errors"
	"log"
	"user-service/src/models"
	"user-service/src/storage"

	"github.com/google/uuid"
)

// This is where our main business logic lives!

type IUserService interface {
    CreateUser(user *models.User, db storage.Database, pubSub storage.PubSub) (*models.User, error)
	UpdateUser(user *models.User, db storage.Database, pubSub storage.PubSub) (*models.User, error)
	RemoveUser(userID string, db storage.Database, pubSub storage.PubSub) error
	GetUser(userID string, db storage.Database) (user *models.User, err error)
	GetUserList(limit int, offset int, filter map[string]string, db storage.Database) (users []*models.User, err error)
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


func (us *UserService) CreateUser(user *models.User) (*models.User, error) {
	user.Password = user.HashPassword(user.Password)
	user.Id = uuid.New().String() // Generate a new UUID.
	timestamp := user.GenerateTimestamp()

	user.CreatedAt = timestamp
	user.UpdatedAt = timestamp

	err := us.db.CreateUser(user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	user, err = us.db.GetUser(user.Id) // Get updated timestamps.
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = us.pubSub.PublishMessage(models.MessageUserCreated, user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (us *UserService) UpdateUser(user *models.User) (*models.User, error) {
	oldUser, err := us.GetUser(user.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	user.CreatedAt = oldUser.CreatedAt
	user.CreatedAt = user.GenerateTimestamp()

	if oldUser == nil {
		return nil, errors.New("user not found")
	}

	if user.Password != oldUser.Password {
		user.Password = user.HashPassword(user.Password) // I might move password updates to a separate endpoint to avoid this extra logic.
	}

	err = us.db.UpdateUser(user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	user, err = us.GetUser(user.Id) // Update timestamps
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = us.pubSub.PublishMessage(models.MessageUserUpdated, user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (us *UserService) RemoveUser(userID string) error {
	err := us.db.DeleteUser(userID)
	if err != nil {
		log.Println(err)
		return err
	}

	err = us.pubSub.PublishMessage(models.MessageUserDeleted, &models.User{Id: userID})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (us *UserService) GetUser(userID string) (user *models.User, err error) {
	user, err = us.db.GetUser(userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetUserList(limit int, offset int, filter map[string]string) (users []*models.User, err error) {
	users, err = us.db.GetUserList(limit, offset, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}
