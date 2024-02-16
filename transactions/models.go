package transactions

import "time"

type depositRequest struct {
	UserId *int    `json:"userId"`
	Amount float64 `json:"amount"`
}

type depositResponse struct {
	Balance float64 `json:"balance"`
}

type transferRequest struct {
	FromUserId *int    `json:"from"`
	ToUserId   *int    `json:"to"`
	Amount     float64 `json:"amount"`
}

type BalanceNATSResponse struct {
	Balance float64 `jsaon:"balance"`
}

type kafkaMessage struct {
	Id        int       `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
