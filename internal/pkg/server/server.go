package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	cfg "github.com/hackerrithm/longterm/rfx/configs"
	"github.com/hackerrithm/longterm/rfx/user/pkg/authenticating"
	"github.com/hackerrithm/longterm/rfx/user/pkg/http/rest"
	"github.com/hackerrithm/longterm/rfx/user/pkg/storage/mongo"
)

// Type defines available storage types
type Type int

const (
	// JSON will store data in JSON files saved on disk
	JSON Type = iota
	// Memory will store data in memory
	Memory
	// MONGO...
	MONGO Type = iota
	// REST will use rest
	REST Type = iota
	// GRAPHQL will make use of graphql
	GRAPHQL
)

// StartServer starts the main application server
func StartServer() {
	var wait time.Duration
	var r http.Handler

	// set up storage
	storageType := MONGO
	httpType := REST

	var authenticater authenticating.Service
	// var taskService engine.Task

	switch storageType {
	// case Memory:
	// 	s := new(inmemory.Storage)

	// 	authenticater = authenticating.NewService(s)

	// case JSON:
	// 	s, _ := asjson.NewStorage()

	// authenticater = authenticating.NewService(s)

	case MONGO:
		s, _ := mongo.NewStorage()

		authenticater = authenticating.NewService(s)
		// taskService = engine
	}

	switch httpType {
	case REST:
		r = rest.Handler(authenticater)
		// case GRAPHQL:
		// 	r = graphql.SetupMux()

	}

	router := r

	file, _ := os.Open("../configs/config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := cfg.Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(configuration.Port)

	srv := &http.Server{
		Handler:      router,
		Addr:         configuration.Address + ":" + configuration.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  120 * time.Minute,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	go srv.Shutdown(ctx)

	<-ctx.Done()

	log.Println("shutting down")
	os.Exit(0)
}
