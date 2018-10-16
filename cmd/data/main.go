package main

import (
	"fmt"
	"time"

	"github.com/hackerrithm/longterm/rfx/internal/pkg/authenticating"
	"github.com/hackerrithm/longterm/rfx/pkg/storage/asjson"
)

// Message ...
type Message interface{}

func main() {

	// error handling omitted for simplicity
	s, _ := asjson.NewStorage()

	// create the available services
	adder := authenticating.NewService(s) // adding "actor"
	//reviewer := reviewing.NewService(s) // reviewing "actor"

	resultsBeer := adder.AddSampleUsers(authenticating.DefaultUsers)
	//resultsReview := reviewer.AddSampleReviews(reviewing.DefaultReviews)

	go func() {
		for result := range resultsBeer {
			fmt.Printf("Added sample user with result %s.\n", result.GetMeaning()) // human-friendly
		}
	}()

	// go func(){
	// 	for result := range resultsReview {
	// 		fmt.Printf("Added sample review with result %d.\n", result) // machine-friendly
	// 	}
	// }()

	// main could have its own "mailbox" exposed, for example an HTTP endpoint,
	// so we could be waiting here for more sample data to be added
	// (but we'll just exit for simplicity)

	time.Sleep(2 * time.Second) // this is here just to get the output from goroutines printed

	fmt.Println("No more data to add!")
}
