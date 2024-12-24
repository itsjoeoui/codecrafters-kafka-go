package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	data := make([]byte, 1024)
	conn.Read(data)

	resp := make([]byte, 8)
	copy(resp, []byte{0, 0, 0, 0})
	copy(resp, data[8:13])

	conn.Write(resp)
}
