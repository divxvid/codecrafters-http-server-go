package main

import (
	"fmt"
	"net"
	"os"

	myhttp "github.com/codecrafters-io/http-server-starter-go/app/my_http"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		defer conn.Close()

		request, err := myhttp.FromReader(conn)
		if err != nil {
			fmt.Println("Encountered an error", err)
			os.Exit(1)
		}

		if request == nil {
			//for the first test
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			continue
		}

		if request.RequestLine.Target == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

	}
}
