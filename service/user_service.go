package service

import (
	"context"
	"log/slog"

	db "github.com/Adarsh-Kmt/EncryptedChat/db/config"
	"github.com/Adarsh-Kmt/EncryptedChat/db/sqlc"
	"github.com/Adarsh-Kmt/EncryptedChat/helper"

	"github.com/Adarsh-Kmt/EncryptedChat/request"
	"github.com/Adarsh-Kmt/EncryptedChat/response"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}
func (service *UserService) RegisterUser(r *request.RegisterUserRequest) *helper.HTTPError {

	params := sqlc.RegisterUserParams{
		Username:  &r.Username,
		PublicKey: r.PublicKey,
	}
	err := db.Client.RegisterUser(context.TODO(), params)

	if err != nil {
		slog.Error(err.Error(), "msg", "error while registering user")
		return &helper.HTTPError{Status: 500, Error: "internal server error"}
	}

	return nil
}

func (service *UserService) GetPublicKey(username string) (*response.PublicKeyResponse, *helper.HTTPError) {

	publicKey, err := db.Client.GetPublicKey(context.TODO(), &username)

	if err != nil {
		slog.Error(err.Error(), "msg", "error while retrieving public key")
		return nil, &helper.HTTPError{Status: 500, Error: "internal server error"}
	}
	return &response.PublicKeyResponse{
		PublicKey: publicKey,
	}, nil
}
