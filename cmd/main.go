package main

import (
	"github.com/hackerrithm/longterm/rfx/internal/pkg/server"
)

// Hello returns a greeting
func Hello() string {
	return "rfx"
}

func main() {
	server.StartServer()
}
