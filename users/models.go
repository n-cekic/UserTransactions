package users

import "time"

type createUserRequest struct {
	Email string `json:"email"`
}

type createUserResponse struct {
	Email string `json:"email"`
}

type newUserKafkaMessage struct {
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
