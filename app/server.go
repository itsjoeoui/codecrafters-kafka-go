package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

const (
	UnsupportedVersion = 35
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

	resp := make([]byte, 10)

	// message_size
	copy(resp[0:4], []byte{0, 0, 0, 0})
	// corelation_id
	copy(resp[4:8], data[8:13])

	version := int16(binary.BigEndian.Uint16(data[6:8]))
	if version > 4 || version < 0 {
		// error_code
		copy(resp[8:10], []byte{0, 35})
	}

	conn.Write(resp)
}
