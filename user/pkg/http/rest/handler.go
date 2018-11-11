package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/hackerrithm/longterm/rfx/user/pkg/authenticating"
	"github.com/julienschmidt/httprouter"
)

// Handler acts as router
func Handler(a authenticating.Service) http.Handler {
	router := httprouter.New()

	router.POST("/auth/login", loginUser(a))
	router.POST("/auth/signup", signUpUser(a))
	router.GET("/auth/profile", profile(a))

	q := handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router)
	return q
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

		tokenValue, err := s.Login(userData["username"], userData["password"])
		if err == authenticating.ErrDuplicate {
			http.Error(w, "The user you requested does not exist.", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenValue)
	}
}

// signUpUser returns a handler for POST /auth/signup
func signUpUser(s authenticating.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		var userData map[string]string
		json.Unmarshal(body, &userData)

		user, err := s.SignUp(userData["username"], userData["password"], userData["firstname"], userData["lastname"])
		if err == authenticating.ErrDuplicate {
			http.Error(w, "The user you requested does not exist.", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// profile returns a user profile for GET /auth/profile
func profile(s authenticating.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		authToken := r.Header.Get("Authorization")
		authArr := strings.Split(authToken, " ")
		UUID := r.URL.Query()["uuid"] //r.URL.Query().Get("uuid")
		if UUID[0] != "" {
			if len(authArr) != 2 {
				log.Println("Authentication header is invalid: " + authToken)
				http.Error(w, "Request failed!", http.StatusUnauthorized)
			}

			jwtToken := authArr[1]
			jsonData, err := s.Profile(jwtToken, UUID[0])
			if err != nil {
				log.Println(err)
				http.Error(w, "Request failed!", http.StatusUnauthorized)
			}

			w.Write(jsonData)
		} else {
			log.Println("request error")
		}
	}
}
