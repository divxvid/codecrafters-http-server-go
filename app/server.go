package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	myhttp "github.com/codecrafters-io/http-server-starter-go/app/my_http"
)

var directory *string

func main() {
	directory = flag.String("directory", "/tmp", "root directory of the folder")
	flag.Parse()

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	router := myhttp.NewRouter()
	router.GET("/", func(ctx *myhttp.HttpContext) myhttp.HttpResponse {
		return *myhttp.NewHttpResponseBuilder().Build()
	})
	router.GET("/echo/:data", handleEcho)
	router.GET("/user-agent", func(ctx *myhttp.HttpContext) myhttp.HttpResponse {
		return *myhttp.NewHttpResponseBuilder().
			WithHeader("Content-Type", "text/plain").
			WithBody([]byte(ctx.GetRequestHeader("User-Agent"))).
			Build()
	})
	router.GET("/files/:filename", handleGetFile)
	router.POST("/files/:filename", handlePostFile)

	server := myhttp.NewServer(router)

	err := server.Start("0.0.0.0", 4221)
	if err != nil {
		fmt.Println("There was some error while starting on port 4221")
		os.Exit(1)
	}
}

func handleEcho(ctx *myhttp.HttpContext) myhttp.HttpResponse {
	encoding := ctx.GetRequestHeader("Accept-Encoding")
	if encoding != "gzip" {
		return *myhttp.NewHttpResponseBuilder().
			WithHeader("Content-Type", "text/plain").
			WithBody([]byte(ctx.PathParam("data"))).
			Build()
	}

	return *myhttp.NewHttpResponseBuilder().
		WithHeader("Content-Type", "text/plain").
		WithHeader("Content-Encoding", "gzip").
		WithBody([]byte(ctx.PathParam("data"))).
		Build()
}

func handlePostFile(ctx *myhttp.HttpContext) myhttp.HttpResponse {
	fileName := ctx.PathParam("filename")
	if directory == nil {
		fmt.Printf("Error parsing the root directory, using /tmp instead")
		*directory = "/tmp"
	}
	fullPath := filepath.Join(*directory, fileName)

	contents := ctx.GetRequestBody()
	err := os.WriteFile(fullPath, contents, 0777)
	if err != nil {
		fmt.Printf("Error writing to the file. err: %v\n", err)
		return *myhttp.NewHttpResponseBuilder().
			WithStatusCode(500).
			WithStatusText("File Write failed").
			Build()
	}

	return *myhttp.NewHttpResponseBuilder().
		WithStatusCode(201).
		WithStatusText("Created").
		Build()
}

func handleGetFile(ctx *myhttp.HttpContext) myhttp.HttpResponse {
	notFoundResponse := *myhttp.NewHttpResponseBuilder().
		WithStatusCode(404).
		WithStatusText("Not Found").
		Build()

	fileName := ctx.PathParam("filename")

	if directory == nil {
		fmt.Printf("Error parsing the root directory")
		return notFoundResponse
	}

	fullPath := filepath.Join(*directory, fileName)

	if !checkFileExists(fullPath) {
		fmt.Printf("File does not exist: %s\n", fullPath)
		return notFoundResponse
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Printf("Error encountered while reading the file: %v\n", err)
		return notFoundResponse
	}

	return *myhttp.NewHttpResponseBuilder().
		WithHeader("Content-Type", "application/octet-stream").
		WithBody(content).
		Build()
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}
