package main

// Use 'main' package  for now. Do not need to create new package

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Note: NewEcoder takes in a io.Writter and http.ResponseWriter implements the same interface
	return json.NewEncoder(w).Encode(v)
}

// Create custom function signature to include errors
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

// Decorate all apiFuncs into a http handlerFunc
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// TODO: Handle the error appropriately
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}

}

type APIServer struct {
	listenAddress string
	store         Storage
}

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

func (s *APIServer) Shutdown() error {
	return nil
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/user", makeHTTPHandleFunc(s.handleUser))

	router.HandleFunc("/user/{id}", makeHTTPHandleFunc(s.handleGetUser))

	log.Println("JSON API server started on port: ", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)

}

// handleUser is a function that handles user requests
func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetUser(w, r)
	case "POST":
		return s.handleCreateUser(w, r)
	case "PUT":
		return s.handleUpdateUser(w, r)
	case "DELETE":
		return s.handleDeleteUser(w, r)
	}
	return fmt.Errorf("unsupported method %s", r.Method)
}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	// TODO: implement the function
	user := NewUser("testuser", "testemail@example.com")

	id := mux.Vars(r)["id"]
	log.Println("Get user with id: ", id)

	jsonUser := WriteJSON(w, http.StatusOK, user)

	return jsonUser
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	// TODO: implement the function
	return nil
}

func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	// TODO: implement the function
	return nil
}
func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	// TODO: implement the function
	return nil
}
