package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	r := mux.NewRouter()

	// 创建日志文件
	//logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer logFile.Close()

	// 设置日志输出到文件
	logger.SetOutput(os.Stdout)

	// 设置日志级别
	logLevel := os.Getenv("LOG_LEVEL")
	logger.Debug("loglevel is:", logLevel)
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

	r.HandleFunc("/", rootHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":80",
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

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("entering v2 root handler")

	delay := randInt(10, 20)
	time.Sleep(time.Millisecond * time.Duration(delay))
	io.WriteString(w, "===================Details of the http request header:============\n")
	req, err := http.NewRequest("GET", "http://service1", nil)
	if err != nil {
		fmt.Printf("%s", err)
	}
	lowerCaseHeader := make(http.Header)
	for key, value := range r.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	logger.Info("headers:", lowerCaseHeader)
	req.Header = lowerCaseHeader
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info("HTTP get failed with error: ", "error", err)
	} else {
		logger.Info("HTTP get succeeded")
	}
	if resp != nil {
		resp.Write(w)
	}
	logger.Infof("Respond in %d ms", delay)
}
