package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	// 创建日志文件
	logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// 设置日志输出到文件
	log.SetOutput(logFile)

	mux.HandleFunc("/localhost/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for name, values := range r.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		version := os.Getenv("GOPATH")
		w.Header().Set("GOPATH", version)

		log.Printf("Client IP: %s, HTTP Status Code: %d\n", r.RemoteAddr, http.StatusOK)
		fmt.Fprintf(w, "Hello, World!")
	})

	http.ListenAndServe(":8080", mux)
}
