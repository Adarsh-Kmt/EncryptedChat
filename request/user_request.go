package request

type RegisterUserRequest struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

type MessageRequest struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
