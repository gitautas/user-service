package storage

import (
	"user-service/src/models"

	"github.com/stretchr/testify/mock"
)

type RedisMock struct {
	mock.Mock
}

func NewRedisMock() *RedisMock {
	return &RedisMock{}
}

func (r *RedisMock) Connect() {
	return
}

func (r *RedisMock) PublishMessage(event string, user *models.User) error {
	args := r.Called()

	result := args.Get(0)
	if result == nil {
		return nil
	}

	return args.Error(0)
}
