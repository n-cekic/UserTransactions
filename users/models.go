package users

type createUserRequest struct {
	Email string `json:"email"`
}

type createUserResponse struct {
	Email string `json:"email"`
}
}
