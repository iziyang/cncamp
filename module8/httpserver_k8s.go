package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

		version := os.Getenv("VERSION")
		w.Header().Set("Version", version)

		log.Printf("Client IP: %s, HTTP Status Code: %d\n", r.RemoteAddr, http.StatusOK)
		fmt.Fprintf(w, "Hello, World!")
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 创建一个信号通道，用于接收终止信号
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		// 启动HTTP服务器
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server listen: %s\n", err)
		}
	}()

	// 等待终止信号
	<-stop

	// 创建一个5秒的超时上下文，用于优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown: %s\n", err)
	}

	log.Println("HTTP server gracefully stopped")
}
