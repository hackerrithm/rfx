package web

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	pstWeb "github.com/hackerrithm/longterm/rfx/post/adapters/web"
	pst "github.com/hackerrithm/longterm/rfx/post/engine"
	usrWeb "github.com/hackerrithm/longterm/rfx/user/adapters/web"
	usr "github.com/hackerrithm/longterm/rfx/user/engine"
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
func NewWebAdapter(a usr.EngineFactory, b pst.EngineFactory, log bool) http.Handler {
	var e *gin.Engine
	if log {
		e = gin.Default()

		e.Use(cors.Default())
	} else {
		e = gin.New()
	}

	usrWeb.InitUsers(e, a, "/")
	pstWeb.InitPosts(e, b, "/")

	return e
}
