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

func (r *RedisMock) PublishMessage(event string, user *models.User) *models.Status {
	args := r.Called(event, user)

	result := args.Get(0)
	status, ok := result.(*models.Status)
	if ok {
		return status
	}

	return nil
}
