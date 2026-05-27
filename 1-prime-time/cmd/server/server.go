package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

type PrimeResponse struct {
	Method string `json:"method"`
	// was the number prime or not
	Prime bool `json:"prime"`
}

type ExpectedRequest struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

func sendMalformed(c net.Conn) {
	c.Write([]byte("{\"error\":true\n"))
}
func isNumberPrime(maybeF float64) bool {
	// if it was a float with a decimal it was not prime
	if math.Trunc(maybeF) != maybeF {
		return false
	}
	n := int(maybeF)
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func handleConnection(c net.Conn) {
	defer c.Close()
	r := bufio.NewReader(c)
	for {
		msg, err := r.ReadBytes('\n')
		fmt.Println("got msg: ",string(msg))
		if err != nil {
			if errors.Is(io.EOF, err) {
				fmt.Println("EOF: done reading closing connection")
				return
			} else {
				fmt.Printf("Error while reading from conenction: %v\n", err)
				return
			}
		}
		var exr ExpectedRequest
		err = json.Unmarshal(msg, &exr)
		if err != nil {
			fmt.Println("got a bad message: ", string(msg))
			sendMalformed(c)
			return
		}
		if exr.Method == nil || exr.Number == nil {
			sendMalformed(c)
			continue
		}

		if *exr.Method != "isPrime" {
			sendMalformed(c)
			continue
		}
		respBody := PrimeResponse{
			Method: "isPrime",
			Prime:  isNumberPrime(*exr.Number),
		}
		m, err := json.Marshal(respBody)
		if err != nil {
			sendMalformed(c)
			continue
		}
		fmt.Println("sending answer: ",string(m))
		c.Write(append(m, '\n'))
	}
}

func acceptLoop(l net.Listener) error {
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		go handleConnection(c)
	}
}

func (s *Server) Listen(port int) error {
	addr := fmt.Sprintf(":%v", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	acceptLoop(l)

	return nil
}
