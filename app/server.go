package main

import (
	"fmt"
	"os"

	myhttp "github.com/codecrafters-io/http-server-starter-go/app/my_http"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	router := myhttp.NewRouter()
	router.GET("/", func(req *myhttp.HttpRequest) myhttp.HttpResponse {
		return *myhttp.NewHttpResponseBuilder().Build()
	})
	router.GET("/echo/%s", func(req *myhttp.HttpRequest) myhttp.HttpResponse {
		return *myhttp.NewHttpResponseBuilder().Build()
	})

	server := myhttp.NewServer(router)

	err := server.Start("0.0.0.0", 4221)
	if err != nil {
		fmt.Println("There was some error while starting on port 4221")
		os.Exit(1)
	}

	// for {
	// 	conn, err := l.Accept()
	// 	if err != nil {
	// 		fmt.Println("Error accepting connection: ", err.Error())
	// 		os.Exit(1)
	// 	}
	// 	defer conn.Close()
	//
	// 	request, err := myhttp.FromReader(conn)
	// 	if err != nil {
	// 		fmt.Println("Encountered an error", err)
	// 		os.Exit(1)
	// 	}
	//
	// 	pathParams := request.PathParams()
	//
	// 	if len(pathParams) == 0 {
	// 		//root path
	// 		myhttp.NewHttpResponseBuilder().
	// 			Build().
	// 			WriteToConn(conn)
	// 	} else if pathParams[0] == "echo" {
	// 		var value string
	// 		if len(pathParams) < 2 {
	// 			value = ""
	// 		} else {
	// 			value = pathParams[1]
	// 		}
	//
	// 		myhttp.NewHttpResponseBuilder().
	// 			WithHeader("Content-Type", "text/plain").
	// 			WithBody([]byte(value)).
	// 			Build().
	// 			WriteToConn(conn)
	// 	} else if pathParams[0] == "user-agent" {
	// 		body := []byte(request.Headers["User-Agent"])
	// 		myhttp.NewHttpResponseBuilder().
	// 			WithHeader("Content-Type", "text/plain").
	// 			WithBody(body).
	// 			Build().
	// 			WriteToConn(conn)
	// 	} else {
	// 		myhttp.NewHttpResponseBuilder().
	// 			WithStatusCode(404).
	// 			WithStatusText("Not Found").
	// 			Build().
	// 			WriteToConn(conn)
	// 	}
	//
	// }
}
