package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Adarsh-Kmt/EncryptedChat/request"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type HTTPError struct {
	Status int
	Error  any
}

type HTTPFunc func(w http.ResponseWriter, r *http.Request) *HTTPError

func MakeHTTPHandlerFunc(f HTTPFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if httpError := f(w, r); httpError != nil {

			WriteJSON(w, map[string]any{"error": httpError.Error}, httpError.Status)
		}
	}
}

func WriteJSON(w http.ResponseWriter, body any, status int) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type UserHandler struct {
	mapMutex     *sync.RWMutex
	userConnMap  map[string]*websocket.Conn
	connMutexMap map[*websocket.Conn]*sync.Mutex
}

func (handler *UserHandler) muxSetup(router *mux.Router) *mux.Router {

	router.HandleFunc("/user/register", MakeHTTPHandlerFunc(handler.RegisterUser)).Methods("POST")
	router.HandleFunc("/user/{username}/public-key", MakeHTTPHandlerFunc(handler.GetPublicKey)).Methods("GET")
	router.HandleFunc("/user/{username}/message", MakeHTTPHandlerFunc(handler.SendMessage))
	return router
}

func (handler *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) *HTTPError {

	return nil
}

func (handler *UserHandler) GetPublicKey(w http.ResponseWriter, r *http.Request) *HTTPError {

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

func (handler *UserHandler) SendMessage(w http.ResponseWriter, r *http.Request) *HTTPError {

	pathVariables := mux.Vars(r)
	username := pathVariables["username"]

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return &HTTPError{Error: "error while switching protocols", Status: 500}
	}

	handler.ConnectUser(username, conn)

	for {

		_, message, err := conn.ReadMessage()

		if err != nil {
			handler.DisconnectUser(username, conn)
			return &HTTPError{Error: err.Error(), Status: 500}
		}

		var mr request.MessageRequest

		if err := json.Unmarshal(message, &mr); err != nil {

			handler.DisconnectUser(username, conn)
			return &HTTPError{Error: err.Error(), Status: 500}
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
