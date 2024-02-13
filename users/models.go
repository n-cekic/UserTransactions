package users

type createUserRequest struct {
	Usernme string `json:"username"`
	Email   string `json:"email"`
}

type createUserResponse struct {
	Usernme string `json:"username"`
	Email   string `json:"email"`
}
