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

// JWTData is a struct with the structure of the jwt data
// type JWTData struct {
// 	// Standard claims are the standard jwt claims from the IETF standard
// 	// https://tools.ietf.org/html/rfc7519
// 	jwt.StandardClaims
// 	CustomClaims map[string]string `json:"custom,omitempty"`
// }

// wire up the greetings routes
func initUsers(e *gin.Engine, f engine.EngineFactory, endpoint string) {
	user := &user{f.NewUser()}
	u := e.Group(endpoint)
	{
		u.POST("/auth/user/signup", user.add)
		u.POST("/auth/user/login", user.login)
		u.GET("/auth/list", user.list)
		u.GET("/auth/user/profile", user.profile)
	}
}

// list converts the parameters into an engine
// request and then marshalls the results based
// on the format requested, returning either an
// html rendered page or JSON (to simulate basic
// content negotiation). It's simpler if the UI
// is a SPA and the web interface is just an API.
func (u user) list(c *gin.Context) {
	ctx := getContext(c)
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil || count == 0 {
		count = 5
	}
	req := &engine.ListUsersRequest{
		Count: count,
	}
	res := u.List(ctx, req)
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
func (u user) add(c *gin.Context) {
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
	repo := u.Add(ctx, req)
	fmt.Println(repo)
	// TODO: token stuff

	res, err := u.GenerateToken()
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, res["token"])
	c.Header("Content-Type", "application/json")
}

func (u user) login(c *gin.Context) {
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
	repo := u.Read(ctx, req)
	fmt.Println(repo)

	// TODO: token stuff

	res, err := u.GenerateToken()
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, res["token"])
	c.Header("Content-Type", "application/json")
}

// profile ...
func (u user) profile(c *gin.Context) {
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

	repo := u.Profile(ctx, req)

	// TODO: token stuff

	res, err := u.ParseToken(jwtToken)
	if err != nil {
		log.Println("err", err)
	}

	res["user"] = repo.Payload

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, res)
}
