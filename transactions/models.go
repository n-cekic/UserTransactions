package transactions

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
