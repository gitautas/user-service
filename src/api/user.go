package api

import (
	"errors"
	"log"
	"user-service/src/models"
	"user-service/src/storage"

	"github.com/google/uuid"
)

// This is where our main business logic lives!

func CreateUser(user *models.User, db storage.Database) (*models.User, error) {
	user.Password = user.HashPassword(user.Password)
	user.Id = uuid.New().String() // Generate a new UUID.
	timestamp := user.GenerateTimestamp()

	user.CreatedAt = timestamp
	user.UpdatedAt = timestamp

	err := db.CreateUser(user)
	if err != nil {
		log.Println(err)
		return nil, err
	}


	user, err = db.GetUser(user.Id) // Get updated timestamps.
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func UpdateUser(user *models.User, db storage.Database) (*models.User, error) {
	oldUser, err := GetUser(user.Id, db)
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

	err = db.UpdateUser(user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	user, err = GetUser(user.Id, db) // Update timestamps
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func RemoveUser(userID string, db storage.Database) error {
	err := db.DeleteUser(userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetUser(userID string, db storage.Database) (user *models.User, err error){
	user, err = db.GetUser(userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func GetUsers(limit int, offset int, filter map[string]string, db storage.Database) (users []*models.User, err error) {
	users, err = db.GetUserList(limit, offset, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}
