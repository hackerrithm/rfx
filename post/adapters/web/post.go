package web

import (
	"fmt"
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
		p.PUT("/auth/post/edit", post.edit)
		p.GET("/auth/post/list", post.list)
		p.GET("/auth/post/read", post.read)
		p.DELETE("/auth/post/delete", post.delete)
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

	c.JSON(http.StatusOK, res.Posts)
	c.Header("Content-Type", "application/json")
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

	id := queryValue("id", c.Request)

	req := &engine.ReadPostRequest{
		ID: id,
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

func (p post) delete(c *gin.Context) {
	ctx := getContext(c)

	id := queryValue("id", c.Request)

	req := &engine.RemovePostRequest{
		ID: id,
	}
	repo := p.Remove(ctx, req)
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
	ctx := getContext(c)

	id := queryValue("id", c.Request)

	fileName, err := FileUpload(c.Writer, c.Request)
	if err != nil {
		log.Println("error bya")
	}

	req := &engine.UpdatePostRequest{
		Author:       c.Request.FormValue("author"),
		Topic:        c.Request.FormValue("topic"),
		Category:     c.Request.FormValue("category"),
		ContentText:  c.Request.FormValue("contentText"),
		ContentPhoto: string(fileName),
	}

	repo := p.Update(ctx, req, id)

	// TODO: token stuff

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, repo)
}
