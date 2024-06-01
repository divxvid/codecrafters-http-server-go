package main

import (
	"bufio"
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

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	request, err := myhttp.FromReader(r)
	if err != nil {
		fmt.Println("Encountered an error", err)
		return
	}

	fmt.Println("TESTING")
	fmt.Println(request)

	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}
