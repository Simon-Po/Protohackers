package main

import "mean2end/src/server"

func main() {
	s := server.New(8080)
	s.Run()
}
