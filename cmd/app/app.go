// +build !appengine

package main

import (
	"net/http"

	"github.com/hackerrithm/longterm/rfx/user/adapters/web"
	user "github.com/hackerrithm/longterm/rfx/user/engine"
	userProvider "github.com/hackerrithm/longterm/rfx/user/providers/mongodb"
	"github.com/hackerrithm/longterm/rfx/user/providers/security"
)

// when running in traditional or 'standalone' mode
// we're going to use MongoDB as the storage provider
// and start the webserver running ourselves.
func main() {
	// s1 := greeterProvider.NewStorage("mongodb://localhost/test1")
	s1 := userProvider.NewStorage("mongodb://localhost/test1")
	s2 := security.NewJWT()

	// eGreeter := greeter.NewEngine(s1)
	eUser := user.NewEngine(s1, s2)
	http.ListenAndServe(":7003", web.NewWebAdapter(eUser, true))
}
