// +build !appengine

package main

import (
	"net/http"

	"github.com/hackerrithm/longterm/rfx/adapters/web"
	post "github.com/hackerrithm/longterm/rfx/post/engine"
	postProvider "github.com/hackerrithm/longterm/rfx/post/providers/mongodb"
	user "github.com/hackerrithm/longterm/rfx/user/engine"

	userProvider "github.com/hackerrithm/longterm/rfx/user/providers/mongodb"

	"github.com/hackerrithm/longterm/rfx/user/providers/security"
)

// when running in traditional or 'standalone' mode
// we're going to use MongoDB as the storage provider
// and start the webserver running ourselves.
func main() {
	// s1 := greeterProvider.NewStorage("mongodb://localhost/test1")
	up := userProvider.NewStorage("mongodb://localhost/test1")
	pp := postProvider.NewStorage("mongodb://localhost/test1")
	sec := security.NewJWT()

	eUser := user.NewEngine(up, sec)
	ePost := post.NewEngine(pp)
	http.ListenAndServe(":7003", web.NewWebAdapter(eUser, ePost, true))
}
