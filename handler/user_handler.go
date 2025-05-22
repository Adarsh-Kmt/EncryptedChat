package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Adarsh-Kmt/EncryptedChat/helper"
	"github.com/Adarsh-Kmt/EncryptedChat/request"
	"github.com/Adarsh-Kmt/EncryptedChat/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type UserHandler struct {
	mapMutex     *sync.RWMutex
	userConnMap  map[string]*websocket.Conn
	connMutexMap map[*websocket.Conn]*sync.Mutex
	userService  *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {

	return &UserHandler{
		userService:  service,
		mapMutex:     &sync.RWMutex{},
		userConnMap:  make(map[string]*websocket.Conn),
		connMutexMap: make(map[*websocket.Conn]*sync.Mutex),
	}
}
func (handler *UserHandler) MuxSetup(router *mux.Router) *mux.Router {

	router.HandleFunc("/user/register", helper.MakeHTTPHandlerFunc(handler.RegisterUser)).Methods("POST")
	router.HandleFunc("/user/{username}/public-key", helper.MakeHTTPHandlerFunc(handler.GetPublicKey)).Methods("GET")
	router.HandleFunc("/user/{username}/message", helper.MakeHTTPHandlerFunc(handler.SendMessage))
	return router
}

func (handler *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	return nil
}

func (handler *UserHandler) GetPublicKey(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	return nil
}

func (handler *UserHandler) ConnectUser(username string, conn *websocket.Conn) {

	handler.mapMutex.Lock()
	defer handler.mapMutex.Unlock()
	handler.userConnMap[username] = conn
	handler.connMutexMap[conn] = &sync.Mutex{}

}

func (handler *UserHandler) DisconnectUser(username string, conn *websocket.Conn) {

	handler.mapMutex.Lock()
	defer handler.mapMutex.Unlock()

	delete(handler.userConnMap, username)
	delete(handler.connMutexMap, conn)

	if err := conn.Close(); err != nil {
		slog.Error(err.Error())
	}
}

func (handler *UserHandler) WriteMessage(conn *websocket.Conn, message string) error {

	handler.mapMutex.RLock()

	m := handler.connMutexMap[conn]

	handler.mapMutex.RUnlock()

	m.Lock()
	defer m.Unlock()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		return err
	}

	return nil
}

func (handler *UserHandler) GetConn(username string) *websocket.Conn {

	handler.mapMutex.RLock()
	defer handler.mapMutex.RUnlock()

	conn, exists := handler.userConnMap[username]

	if !exists {
		return nil
	}

	return conn
}

func (handler *UserHandler) SendMessage(w http.ResponseWriter, r *http.Request) *helper.HTTPError {

	pathVariables := mux.Vars(r)
	username := pathVariables["username"]

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return &helper.HTTPError{Error: "error while switching protocols", Status: 500}
	}

	handler.ConnectUser(username, conn)

	for {

		_, message, err := conn.ReadMessage()

		if err != nil {
			handler.DisconnectUser(username, conn)
			return &helper.HTTPError{Error: err.Error(), Status: 500}
		}

		var mr request.MessageRequest

		if err := json.Unmarshal(message, &mr); err != nil {

			handler.DisconnectUser(username, conn)
			return &helper.HTTPError{Error: err.Error(), Status: 500}
		}

		toConn := handler.GetConn(mr.Username)

		if toConn == nil {

			if err := handler.WriteMessage(toConn, mr.Message); err != nil {
				slog.Error(err.Error())
			}

		} else {

			handler.WriteMessage(conn, "user is offline")
		}

	}

}
