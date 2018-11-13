package web

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	usr "github.com/hackerrithm/longterm/rfx/usr/engine"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// NewWebAdapter creates a new web adaptor which will
// handle the web interface and pass calls on to the
// engine to do the real work (that's why the engine
// factory is passed in - so anything that *it* needs
// is unknown to this).
// Because the web adapter ends up quite lightweight
// it easier to replace. We could use any one of the
// Go web routers / frameworks (Gin, Echo, Goji etc...)
// or just stick with the standard framework. Changing
// should be far less costly.
func NewWebAdapter(f1 usr.EngineFactory, log bool) http.Handler {
	var e *gin.Engine
	if log {
		e = gin.Default()

		e.Use(cors.Default())
	} else {
		e = gin.New()
	}

	e.LoadHTMLGlob("templates/*")

	initUsers(e, f1, "/")

	return e
}
