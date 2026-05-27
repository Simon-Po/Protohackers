package main

import (
	"fmt"
	"io"
	"log"
	"net"
)



func handleConnection(c net.Conn) {
	var msg []byte
	defer c.Close()
	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf)
		if n > 0 {
			msg = append(msg,buf[:n]...)
		}
		if err != nil {
			if err == io.EOF {
				c.Write(msg)
				fmt.Println("client closed connection")
			} else {
				fmt.Println("read error:", err)
			}
			return
		}
	}

}

func acceptLoop(l net.Listener) {
	for {
		c,err := l.Accept()
		if err != nil {
			fmt.Println("Error could not accept conenction")
			continue
		}
		go handleConnection(c)
	}

}

func main() {
	listener,err := net.Listen("tcp",":8080")
	if err != nil {
		log.Fatal(err)
	}
	acceptLoop(listener)

}
