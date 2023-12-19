package main

// Use 'main' package  for now. Do not need to create new package

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	Error string `json:"error"`
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

	router.HandleFunc("/user/{id}", WithJWTAuth(makeHTTPHandleFunc(s.handleUserByID), s.store))

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
	}
	return fmt.Errorf("unsupported method %s", r.Method)
}

func (s *APIServer) handleUserByID(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetUserById(w, r)
	case "DELETE":
		return s.handleDeleteUser(w, r)
	}
	return fmt.Errorf("unsupported method %s", r.Method)
}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetUsers()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}

func (s *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	log.Println("Get user with id: ", id)

	// Get user from database
	user, err := s.store.GetUserByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	// Note: new(x) and &x{} is the same. Returns a pointer to x
	// Note: Use new(x) if x is a basic type. Cannot use &int{0}.
	createUserRequest := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(createUserRequest); err != nil {
		return err
	}

	user := NewUser(createUserRequest.Username, createUserRequest.Email)
	if err := s.store.CreateUser(user); err != nil {
		return err
	}

	// Create token for user
	tokenString, err := CreateJWTTokenString(user)
	if err != nil {
		return err
	}

	fmt.Println("Token String: ", tokenString)

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	// TODO: implement the function
	return nil
}
func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	// Delete user from database
	id, err := getID(r)
	if err != nil {
		return err
	}

	log.Println("Delete user with id: ", id)

	if err := s.store.DeleteUser(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func getID(r *http.Request) (int, error) {
	// Note: ParseInt is faster, but requires additional parse from int64 -> int
	// Note: Atoi is slightly slower, but directly convert to int
	// id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 0) -> int64
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	// Input validation for 'id' field
	if err != nil {
		return 0, fmt.Errorf("Invalid id: %s", idStr)
	}

	return id, nil
}
