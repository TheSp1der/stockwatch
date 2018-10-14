package main

import (
	"strconv"
	"time"

	"net/http"

	"github.com/TheSp1der/goerror"
)

func webRoot(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	resp.Header().Add("Expires", "0")	
	resp.Header().Add("Content-Type", "text/html")
	if req.Method == "GET" {
		resp.Write([]byte(displayWeb(sData)))
	}
}

func webListener(port int) {
	ws := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: ws,
		// time to read request headers
		ReadTimeout: time.Duration(15 * time.Second),
		// time from accept to end of response
		WriteTimeout: time.Duration(10 * time.Second),
		// time a Keep-Alive connection will be kept idle
		IdleTimeout: time.Duration(120 * time.Second),
	}

	ws.HandleFunc("/", webRoot)

	if err := srv.ListenAndServe(); err != nil {
		goerror.Fatal(err)
	}
}
