package api

import (
	"errors"
	"fmt"
	"user-service/src/models"
	"user-service/src/storage"

	"github.com/google/uuid"
)

// This is where our main business logic lives!

func CreateUser(user *models.User, mysql storage.Database) (*models.User, error) {
	user.Password = user.HashPassword(user.Password)
	user.ID = uuid.New().String() // Generate a new UUID.
	err := mysql.CreateUser(user)
	if err != nil {
		return nil, err
	}

	user, err = mysql.GetUser(user.ID) // Get updated timestamps.
	if err != nil {
		return nil, err
	}

	return user, nil
	// I could return the function here, but because this is
	// where any additional logic would be placed in
	// any more complex service, therefore I choose to
	// keep the additional lines here in these service functions.
	//
	// return mysql.CreateUser(user)
}

func UpdateUser(user *models.User, mysql storage.Database) (*models.User, error) {
	oldUser, err := GetUser(user.ID, mysql)
	if err != nil {
		fmt.Println(err) // FIXME
		return nil, err
	}

	if oldUser == nil {
		fmt.Println("notfound") // FIXME
		return nil, errors.New("user not found")
	}

	if user.Password != oldUser.Password {
		user.Password = user.HashPassword(user.Password) // I might move password updates to a separate endpoint to avoid this extra logic.
	}

	err = mysql.UpdateUser(user)
	if err != nil {
		fmt.Println(err) // FIXME
		return nil, err
	}

	user, err = GetUser(user.ID, mysql) // Update timestamps
	if err != nil {
		fmt.Println(err) // FIXME
		return nil, err
	}

	return user, nil
}

func RemoveUser(userID string, mysql storage.Database) error {
	err := mysql.RemoveUser(userID)
	if err != nil {
		fmt.Println(err) // FIXME
		return err
	}
	return nil
}

func GetUser(userID string, mysql storage.Database) (user *models.User, err error){
	user, err = mysql.GetUser(userID)
	if err != nil {
		fmt.Println(err) // FIXME
		return nil, err
	}
	return user, nil
}

func GetUsers(limit int, offset int, mysql storage.Database) (users []*models.User, err error) {
	users, err = mysql.GetUsers(limit, offset)
	if err != nil {
		fmt.Println(err) // FIXME
		return nil, err
	}
	return users, nil
}
