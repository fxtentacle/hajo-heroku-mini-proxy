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

	allowed_ips := make(map[string]bool)

	auth := func(writer http.ResponseWriter, request *http.Request) {
		ip := request.Header.Get("X-FORWARDED-FOR")
		if ip == "" { ip = strings.SplitN(request.RemoteAddr,":",2)[0] }

		if request.URL.Path == "/proxyauth/"+os.Getenv("PROXYAUTH") {
			fmt.Println("AUTHORIZE IP: ", ip)
			allowed_ips[ip] = true
			writer.WriteHeader(200)
			writer.Write([]byte("200 OK PROXYAUTH\n"))
			return
		}

		if !allowed_ips[ip] {
			fmt.Println("IP NOT AUTHORIZED: ", ip)
			writer.WriteHeader(500)
			writer.Write([]byte("500 PROXYAUTH\n"))
			return
		}
		fmt.Println("IP is authorized: ", ip)

		proxy.ServeHTTP(writer, request)
	}

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(auth)))
}
