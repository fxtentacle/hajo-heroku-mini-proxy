package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	auth := func(writer http.ResponseWriter, request *http.Request) {
		proxy.ServeHTTP(writer, request)
	}

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(auth)))
}
