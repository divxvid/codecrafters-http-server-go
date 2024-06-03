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
	router.GET("/", func(ctx *myhttp.HttpContext) myhttp.HttpResponse {
		return *myhttp.NewHttpResponseBuilder().Build()
	})
	router.GET("/echo/:data", func(ctx *myhttp.HttpContext) myhttp.HttpResponse {
		return *myhttp.NewHttpResponseBuilder().
			WithBody([]byte(ctx.PathParam("data"))).
			Build()
	})

	server := myhttp.NewServer(router)

	err := server.Start("0.0.0.0", 4221)
	if err != nil {
		fmt.Println("There was some error while starting on port 4221")
		os.Exit(1)
	}
}
