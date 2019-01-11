package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"strings"
)

type FastCGIServer struct{}

func (s *FastCGIServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if nil != err {
		fmt.Fprintf(os.Stderr, "[ERROR]: %v", err)
	}

	fmt.Println("FormData:")
	for k, v := range req.Form {
		fmt.Printf("%v: %v\n", k, strings.Join(v, ", "))
	}

	fmt.Println("Host: ", req.Host)
	fmt.Println("Header:")
	for k, v := range req.Header {
		fmt.Fprintf(os.Stdout, "%v: %v\n", k, strings.Join(v, ","))
	}

	var buffer [1024]byte
	fmt.Println("BODY:")
	for {
		len, err := req.Body.Read(buffer[:])
		if nil != err {
			if io.EOF != err {
				fmt.Fprintf(os.Stderr, "[ERROR]: %v\n", err)
			} else {
				break
			}
		}

		fmt.Print(string(buffer[:len]))
	}

	w.Write([]byte("This is a FastCGI example server.\n"))
}

func main() {
	fmt.Println("Starting server...")
	l, _ := net.Listen("tcp", "127.0.0.1:9000")
	h := new(FastCGIServer)
	fcgi.Serve(l, h)
}
