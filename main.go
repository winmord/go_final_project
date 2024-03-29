package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func getPort() int {
	port := 7540
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(eport)
		}
	}

	return port
}

func main() {
	webDir := "web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(fmt.Sprintf(":%d", getPort()), nil)
	if err != nil {
		panic(err)
	}
}
