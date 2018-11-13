package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hackerrithm/longterm/rfx/task/engine"
)

// NewWebAdapter ...
func NewWebAdapter(f engine.Factory) http.Handler {
	r := mux.NewRouter()

	task := newTask(f)

	r.Handle("/v1/task", task.AddTask).Methods("POST")
	return r
}
