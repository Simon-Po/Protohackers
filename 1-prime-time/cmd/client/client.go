package client

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
}

func New() *Client {
	return &Client{}
}

func makePayload() []byte {

	return []byte("{\"method\":\"isPrime\",\"number\":23}\n{\"method\":\"isPrime\",\"number\":29393}\n")
}

// little test client
func (c *Client) Dial(port int) error {
	addr := fmt.Sprintf(":%d", port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.Write(makePayload())
	r := bufio.NewReader(conn)
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Println(s)
	}
}
