package service

import (
	"context"
	"log/slog"

	db "github.com/Adarsh-Kmt/EncryptedChat/db/config"
	"github.com/Adarsh-Kmt/EncryptedChat/db/sqlc"
	"github.com/Adarsh-Kmt/EncryptedChat/handler"
	"github.com/Adarsh-Kmt/EncryptedChat/request"
	"github.com/Adarsh-Kmt/EncryptedChat/response"
)

type UserService struct {
}

func (service *UserService) RegisterUser(r *request.RegisterUserRequest) *handler.HTTPError {

	params := sqlc.RegisterUserParams{
		Username:  &r.Username,
		PublicKey: r.PublicKey,
	}
	err := db.Client.RegisterUser(context.TODO(), params)

	if err != nil {
		slog.Error(err.Error(), "msg", "error while registering user")
		return &handler.HTTPError{Status: 500, Error: "internal server error"}
	}

	return nil
}

func (service *UserService) GetPublicKey(username string) (*response.PublicKeyResponse, *handler.HTTPError) {

	publicKey, err := db.Client.GetPublicKey(context.TODO(), &username)

	if err != nil {
		slog.Error(err.Error(), "msg", "error while retrieving public key")
		return nil, &handler.HTTPError{Status: 500, Error: "internal server error"}
	}
	return &response.PublicKeyResponse{
		PublicKey: publicKey,
	}, nil
}
