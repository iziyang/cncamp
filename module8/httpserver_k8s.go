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

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	r := mux.NewRouter()

	// 创建日志文件
	logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// 设置日志输出到文件
	logger.SetOutput(logFile)

	// 设置日志级别
	logLevel := os.Getenv("LOG_LEVEL")
	fmt.Printf("loglevel is:", logLevel)
	if logLevel == "DEBUG" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.Debug("Debug level logging started") // Debug log

	r.HandleFunc("/localhost/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		logger.Debug("Begin health check, check url is /localhost/healthz")
		logger.Info("Health check passed") // Info log
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for name, values := range r.Header {
			for _, value := range values {
				w.Header().Add(name, value)
				logger.Debug("Header name value is:", name, value)
			}

		}
		logger.Info("Request header has been writed in response.")

		version := os.Getenv("VERSION")
		w.Header().Set("Version", version)

		logger.WithFields(logrus.Fields{
			"ClientIP":       r.RemoteAddr,
			"HTTPStatusCode": http.StatusOK,
		}).Info("Request processed")

		fmt.Fprintf(w, "Hello, World!")
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// 优雅终止
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		// 优雅关闭
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.WithError(err).Error("Server shutdown failed")
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.WithError(err).Error("Server closed unexpectedly")
		os.Exit(1)
	}

	<-idleConnsClosed
}
