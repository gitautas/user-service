package models

const (
	MessageUserCreated = "userCreated"
	MessageUserUpdated = "userUpdated"
	MessageUserDeleted = "userDeleted"
)

type PubSubMessage struct {
	Message string `json:"message"`
	User    *User  `json:"user"`
}
