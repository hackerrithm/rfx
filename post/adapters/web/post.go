package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hackerrithm/longterm/rfx/post/engine"
)

type (
	post struct {
		engine.Post
	}
)

// InitPosts wire up the posts routes
func InitPosts(e *gin.Engine, f engine.EngineFactory, endpoint string) {
	post := &post{f.NewPost()}
	p := e.Group(endpoint)
	{
		p.POST("/auth/post/add", post.add)
		p.POST("/auth/post/edit", post.edit)
		p.GET("/auth/post/list", post.list)
		p.GET("/auth/post/read", post.read)
	}
}

// list converts the parameters into an engine
// request and then marshalls the results based
// on the format requested, returning either an
// html rendered page or JSON (to simulate basic
// content negotiation). It's simpler if the UI
// is a SPA and the web interface is just an API.
func (p post) list(c *gin.Context) {
	ctx := getContext(c)
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil || count == 0 {
		count = 5
	}
	req := &engine.ListPostsRequest{
		Count: count,
	}
	res := p.List(ctx, req)
	if c.Query("format") == "json" {
		fmt.Println("here")
		c.JSON(http.StatusOK, res.Posts)
	} else {
		c.HTML(http.StatusOK, "guestbook.html", res)
	}
}

// add accepts a form post and creates a new
// greoting in the system. It could be made a
// lot smarter and automatically check for the
// content type to handle forms, JSON etc...
func (p post) add(c *gin.Context) {
	ctx := getContext(c)

	fileName, err := FileUpload(c.Writer, c.Request)
	if err != nil {
		log.Println("error bya")
	}

	req := &engine.AddPostRequest{
		Author:       c.Request.FormValue("author"),
		Topic:        c.Request.FormValue("topic"),
		Category:     c.Request.FormValue("category"),
		ContentText:  c.Request.FormValue("contentText"),
		ContentPhoto: string(fileName),
	}

	repo := p.Add(ctx, req)

	// TODO: token stuff

	// res, err := u.GenerateToken()
	// if err != nil {
	// 	log.Println(err)
	// }

	c.JSON(http.StatusOK, repo)
	c.Header("Content-Type", "application/json")
}

func (p post) read(c *gin.Context) {
	ctx := getContext(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	var postData map[string]string
	json.Unmarshal(body, &postData)

	req := &engine.ReadPostRequest{
		ID: postData["id"],
	}
	repo := p.Read(ctx, req)
	fmt.Println(repo)

	// TODO: token stuff

	// res, err := u.GenerateToken()
	// if err != nil {
	// 	log.Println(err)
	// }

	c.JSON(http.StatusOK, repo)
	c.Header("Content-Type", "application/json")
}

// edit ...
func (p post) edit(c *gin.Context) {
	// ctx := getContext(c)
	// authToken := c.GetHeader("Authorization")
	// authArr := strings.Split(authToken, " ")
	// UUID := queryValue("uuid", c.Request)
	// if len(authArr) != 2 {
	// 	log.Println("Authentication header is invalid: " + authToken)
	// 	http.Error(c.Writer, "Request failed!", http.StatusUnauthorized)
	// }

	// jwtToken := authArr[1]

	// req := &engine.ProfileRequest{
	// 	Token: jwtToken,
	// 	ID:    UUID,
	// }

	// repo := u.Profile(ctx, req)

	// // TODO: token stuff

	// res, err := u.ParseToken(jwtToken)
	// if err != nil {
	// 	log.Println("err", err)
	// }

	// res["Post"] = repo.Payload

	// c.Header("Content-Type", "application/json")
	// c.JSON(http.StatusOK, res)
}
