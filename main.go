package main

import (
	"net/http"

	db "github.com/Adarsh-Kmt/EncryptedChat/db/config"
	"github.com/Adarsh-Kmt/EncryptedChat/handler"
	"github.com/Adarsh-Kmt/EncryptedChat/service"
	"github.com/gorilla/mux"
)

func main() {

	if err := db.PostgresDBClientInit(); err != nil {
		panic(err)
	}

	userService := service.NewUserService()

	userHandler := handler.NewUserHandler(userService)

	router := mux.NewRouter()

	router = userHandler.MuxSetup(router)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
