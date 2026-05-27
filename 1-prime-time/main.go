package main

import (
	"1-prime-time/cmd/server"
)

func main() {
	server.New().Listen(8080)
}
