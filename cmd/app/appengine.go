// +build appengine

package main

import (
	"net/http"

	"github.com/hackerrithm/longterm/rfx/greeter/adapters/web"
	"github.com/hackerrithm/longterm/rfx/greeter/engine"
	"github.com/hackerrithm/longterm/rfx/greeter/providers/appengine"
)

// for appengine we don't use main to start the server
// because that is done for us by the platform. Instead
// we attach to the standard mux router. Note that we're
// using the appengine provider for storage and wiring
// it up to the engine and then the engine to the web.
func init() {
	s := appengine.NewStorage()
	e := engine.NewEngine(s)
	http.Handle("/", web.NewWebAdapter(e, false))
}
