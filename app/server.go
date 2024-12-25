package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

const (
	ErrUnsupportedVersion = 35
)

type Request struct {
	Version       int16
	CorrelationId int32
}

func ParseRequest(data []byte) Request {
	return Request{
		Version:       int16(binary.BigEndian.Uint16(data[6:8])),
		CorrelationId: int32(binary.BigEndian.Uint32(data[8:12])),
	}
}

func WriteResponse(conn net.Conn, resp []byte) {
	// the first 4 bytes of the response should be the length of the response
	binary.Write(conn, binary.BigEndian, int32(len(resp)))
	// then we write the response
	binary.Write(conn, binary.BigEndian, resp)
}

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

	request := ParseRequest(data)

	if request.Version < 0 || request.Version > 4 {
		// Invalid Version
		resp := make([]byte, 6)
		// 4 bytes for correlation_id
		binary.BigEndian.PutUint32(resp[0:4], uint32(request.CorrelationId))
		// 2 bytes for error code
		binary.BigEndian.PutUint16(resp[4:6], uint16(ErrUnsupportedVersion))

		WriteResponse(conn, resp)
		os.Exit(1)
	}

	// build the response
	resp := make([]byte, 19)

	// 4 bytes for correlation_id
	binary.BigEndian.PutUint32(resp[0:4], uint32(request.CorrelationId))
	// 2 bytes for error code
	binary.BigEndian.PutUint16(resp[4:6], uint16(0))

	// TODO: some guy reversed engineered this on the forum, no official instructions yet
	// https://forum.codecrafters.io/t/question-about-handle-apiversions-requests-stage/1743/4
	resp[6] = 2
	binary.BigEndian.PutUint16(resp[7:9], 18)
	binary.BigEndian.PutUint16(resp[9:11], 3)
	binary.BigEndian.PutUint16(resp[11:13], 4)
	resp[13] = 0
	binary.BigEndian.PutUint32(resp[14:18], 0)
	resp[18] = 0

	WriteResponse(conn, resp)
}
