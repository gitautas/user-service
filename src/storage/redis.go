package storage

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"user-service/src/models"

	"github.com/go-redis/redis/v8"
)

type PubSub interface {
    Connect()
	PublishMessage(event string, user *models.User) error
}

type Redis struct {
    endpoint string
	client *redis.Client
	pubSubMutex *sync.Mutex
	pubSubChan chan *models.PubSubMessage
	timeout time.Duration
}

func NewRedis(endpoint string, timeout time.Duration) *Redis {
	return &Redis{
		endpoint:   endpoint,
		pubSubMutex: &sync.Mutex{},
		pubSubChan:  make(chan *models.PubSubMessage),
		timeout: timeout,
	}
}

func (r *Redis) Connect() {
	client := redis.NewClient(&redis.Options{
		Addr:    r.endpoint,
	})

	r.client = client
}


func (r *Redis) PublishMessage(event string, user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	message, err := json.Marshal(&models.PubSubMessage{
		Message: event,
		User: user,
	})
	if err != nil {
		return err
	}
	r.client.Publish(ctx, "user", message)
	return nil
}
