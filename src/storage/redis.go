package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
	"user-service/src/models"

	"github.com/go-redis/redis/v8"
)

type PubSub interface {
	Connect()
	PublishMessage(event string, user *models.User) *models.Status
}

type Redis struct {
	endpoint    string
	client      *redis.Client
	pubSubMutex *sync.Mutex
	pubSubChan  chan *models.PubSubMessage
	timeout     time.Duration
}

func NewRedis(endpoint string, timeout time.Duration) *Redis {
	return &Redis{
		endpoint:    endpoint,
		pubSubMutex: &sync.Mutex{},
		pubSubChan:  make(chan *models.PubSubMessage),
		timeout:     timeout,
	}
}

func (r *Redis) Connect() {
	client := redis.NewClient(&redis.Options{
		Addr: r.endpoint,
	})

	r.client = client
}

func (r *Redis) PublishMessage(event string, user *models.User) *models.Status {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	message, _ := json.Marshal(&models.PubSubMessage{
		Message: event,
		User:    user,
	})

	err := r.client.Publish(ctx, "user", message).Err()
	if err != nil {
		return models.NewStatus(http.StatusInternalServerError, "could not publish message")
	}
	return nil
}
