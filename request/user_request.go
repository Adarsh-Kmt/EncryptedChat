package request

import (
	"encoding/json"
	"net/http"

	"github.com/Adarsh-Kmt/EncryptedChat/helper"
)

type Validator interface {
	Validate() (errorMap map[string]any)
}

func DecodeAndValidate[T Validator](r *http.Request) (T, *helper.HTTPError) {

	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, &helper.HTTPError{Status: 400, Error: "bad request"}
	}
	if errorMap := v.Validate(); len(errorMap) > 0 {
		return v, &helper.HTTPError{Status: 400, Error: errorMap}
	}
	return v, nil
}

type RegisterUserRequest struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

func (request RegisterUserRequest) Validate() (errorMap map[string]any) {

	errorMap = make(map[string]any)

	if request.Username == "" {
		errorMap["username"] = "username cannot be empty"
	}
	if request.PublicKey == "" {
		errorMap["public_key"] = "public key cannot be empty"
	}

	if len(errorMap) != 0 {
		return errorMap
	}
	return nil
}

type MessageRequest struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func (request MessageRequest) Validate() (errorMap map[string]any) {

	errorMap = make(map[string]any)

	if request.Username == "" {
		errorMap["username"] = "username cannot be empty"
	}
	if request.Message == "" {
		errorMap["message"] = "public key cannot be empty"
	}

	if len(errorMap) != 0 {
		return errorMap
	}
	return nil
}
