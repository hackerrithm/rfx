package main

import (
	"fmt"

	"github.com/hackerrithm/longterm/rfx/internal/pkg/server"
)

// Hello returns a greeting
func Hello() string {
	return "rfx"
}

func main() {
	fmt.Println(Hello())
	server.StartServer()
}
