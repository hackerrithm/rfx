package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/hackerrithm/longterm/rfx/internal/pkg/authenticating"
	"github.com/hackerrithm/longterm/rfx/internal/pkg/listing"
	"github.com/julienschmidt/httprouter"
)

// Handler acts as router
func Handler(a authenticating.Service, l listing.Service) http.Handler {
	router := httprouter.New()

	router.POST("/users", addUser(a))
	router.GET("/users", getUsers(l))
	router.GET("/users/:id", getUser(l))
	router.POST("/auth/login", loginUser(a))

	q := handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router)
	return q
}

// addUser returns a handler for POST /users requests
func addUser(s authenticating.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		decoder := json.NewDecoder(r.Body)

		var newUser authenticating.User
		err := decoder.Decode(&newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.AddUser(newUser)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("New user added.")
	}
}

// getUsers returns a handler for GET /users requests
func getUsers(s listing.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		list := s.GetUsers()
		json.NewEncoder(w).Encode(list)
	}
}

// getUser returns a handler for GET /users/:id requests
func getUser(s listing.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ID, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s is not a valid user ID, it must be a number.", p.ByName("id")), http.StatusBadRequest)
			return
		}

		user, err := s.GetUser(ID)
		if err == listing.ErrNotFound {
			http.Error(w, "The user you requested does not exist.", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// loginUser returns a handler for POST /auth/login
func loginUser(s authenticating.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		var userData map[string]string
		json.Unmarshal(body, &userData)

		user, err := s.Login(userData["username"], userData["password"])
		if err == authenticating.ErrDuplicate {
			http.Error(w, "The user you requested does not exist.", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
