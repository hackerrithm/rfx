package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hackerrithm/longterm/rfx/task/engine"
	"github.com/hackerrithm/longterm/rfx/user/pkg/authenticating"
	"github.com/julienschmidt/httprouter"
)

type (
	task struct {
		engine.Task
	}
)

func newTask(f engine.Factory) *task {
	return &task{f.NewTask()}
}

// AddTask accepts a form post and creates a new
// greoting in the system. It could be made a
// lot smarter and automatically check for the
// content type to handle forms, JSON etc...
func (g task) AddTask(w http.ResponseWriter, r *http.Request) {
	// ctx := getContext(c)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Login failed!", http.StatusUnauthorized)
	}

	var userData map[string]string
	json.Unmarshal(body, &userData)

	req := &engine.AddTaskRequest{
		Author:  userData["Author"],
		Content: userData["Content"],
	}
	g.Add(ctx, req)
	// TODO: set a flash cookie for "added"
	// if this was a web request, otherwise
	// send a nice JSON response ...
	c.Redirect(http.StatusFound, "/")
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
