package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net"
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

		if request.URL.Path == "/client_ip" {
			writer.WriteHeader(200)
			writer.Write([]byte(ip))
			return
		}

		if strings.HasPrefix(request.URL.Path, "/resolve/") {
			resolveMe := request.URL.Path[9:]
			ips, err := net.LookupIP(resolveMe)
			ipStrs := []string{}
			for _, ip := range ips {
				ipStrs = append(ipStrs, ip.String())
			}
			if(err != nil) {
				fmt.Println("ERROR: Could not resolve ", resolveMe)
				writer.WriteHeader(500)
				writer.Write([]byte("500 RESOLV ERROR\n"))
			} else {
				fmt.Println("RESOLVED: ", resolveMe, " -> ", )
				writer.WriteHeader(200)
				writer.Write([]byte(strings.Join(ipStrs, ",")))
			}
			return
		}


		proxy.ServeHTTP(writer, request)
	}

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(auth)))
}
