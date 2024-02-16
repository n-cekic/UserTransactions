package users

import "time"

type createUserRequest struct {
	Email string `json:"email"`
}

type newUserKafkaMessage struct {
	Id        int       `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type getBalanceRequest struct {
	Email string `json:"email"`
}

type getBalanceResponse struct {
	Email   string  `json:"email"`
	Balance float64 `json:"balance"`
}
