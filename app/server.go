package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
			WithHeader("Content-Type", "text/plain").
			WithBody([]byte(ctx.PathParam("data"))).
			Build()
	})
	router.GET("/user-agent", func(ctx *myhttp.HttpContext) myhttp.HttpResponse {
		return *myhttp.NewHttpResponseBuilder().
			WithHeader("Content-Type", "text/plain").
			WithBody([]byte(ctx.GetRequestHeader("User-Agent"))).
			Build()
	})
	router.GET("/files/:filename", handleFiles)

	server := myhttp.NewServer(router)

	err := server.Start("0.0.0.0", 4221)
	if err != nil {
		fmt.Println("There was some error while starting on port 4221")
		os.Exit(1)
	}
}

func handleFiles(ctx *myhttp.HttpContext) myhttp.HttpResponse {
	fileName := ctx.PathParam("filename")
	fullPath := filepath.Join("tmp", fileName)

	notFoundResponse := *myhttp.NewHttpResponseBuilder().
		WithStatusCode(404).
		WithStatusText("Not Found").
		Build()

	if !checkFileExists(fullPath) {
		return notFoundResponse
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Printf("Error encountered while reading the file: %v\n", err)
		return notFoundResponse
	}

	return *myhttp.NewHttpResponseBuilder().
		WithHeader("Content-Type", "text/octet-stream").
		WithBody(content).
		Build()
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}
