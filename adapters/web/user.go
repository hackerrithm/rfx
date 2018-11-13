package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hackerrithm/longterm/rfx/user/engine"
)

type (
	user struct {
		engine.User
	}
)

// wire up the greetings routes
func initUsers(e *gin.Engine, f engine.EngineFactory, endpoint string) {
	user := &user{f.NewUser()}
	g := e.Group(endpoint)
	{
		g.GET("/auth/list", user.list)
		g.POST("/auth/user/signup", user.add)
		g.POST("/auth/user/login", user.login)
		g.GET("/auth/user/profile", user.profile)
	}
}

// list converts the parameters into an engine
// request and then marshalls the results based
// on the format requested, returning either an
// html rendered page or JSON (to simulate basic
// content negotiation). It's simpler if the UI
// is a SPA and the web interface is just an API.
func (g user) list(c *gin.Context) {
	ctx := getContext(c)
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil || count == 0 {
		count = 5
	}
	req := &engine.ListUsersRequest{
		Count: count,
	}
	res := g.List(ctx, req)
	if c.Query("format") == "json" {
		fmt.Println("here")
		c.JSON(http.StatusOK, res.Users)
	} else {
		c.HTML(http.StatusOK, "guestbook.html", res)
	}
}

// add accepts a form post and creates a new
// greoting in the system. It could be made a
// lot smarter and automatically check for the
// content type to handle forms, JSON etc...
func (g user) add(c *gin.Context) {
	ctx := getContext(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	var userData map[string]string
	json.Unmarshal(body, &userData)

	req := &engine.AddUserRequest{
		UserName:  userData["username"],
		Password:  userData["password"],
		FirstName: userData["firstname"],
		LastName:  userData["lastname"],
		Gender:    userData["gender"],
	}
	res := g.Add(ctx, req)
	c.JSON(http.StatusOK, res.Token)
	// TODO: set a flash cookie for "added"
	// if this was a web request, otherwise
	// send a nice JSON response ...
	// c.Redirect(http.StatusFound, "/add")
	c.Header("Content-Type", "application/json") //.Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(tokenValue)
}

func (g user) login(c *gin.Context) {
	ctx := getContext(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	var userData map[string]string
	json.Unmarshal(body, &userData)

	req := &engine.AuthenticateUserRequest{
		UserName: userData["username"],
		Password: userData["password"],
	}
	res := g.Read(ctx, req)
	c.JSON(http.StatusOK, res.Token)
	// TODO: set a flash cookie for "added"
	// if this was a web request, otherwise
	// send a nice JSON response ...
	// c.Redirect(http.StatusFound, "/add")
	c.Header("Content-Type", "application/json") //.Set("Content-Type", "application/json")

	// json.NewEncoder(w).Encode(tokenValue)
}

// profile ...
func (g user) profile(c *gin.Context) {
	ctx := getContext(c)
	authToken := c.GetHeader("Authorization")
	authArr := strings.Split(authToken, " ")
	UUID := queryValue("uuid", c.Request)
	if len(authArr) != 2 {
		log.Println("Authentication header is invalid: " + authToken)
		http.Error(c.Writer, "Request failed!", http.StatusUnauthorized)
	}

	jwtToken := authArr[1]

	req := &engine.ProfileRequest{
		Token: jwtToken,
		ID:    UUID,
	}
	res := g.Profile(ctx, req)
	c.Header("Content-Type", "application/json") //.Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, res)
}
