package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/toqueteos/webbrowser"
)

var (
	servePort      = ":8080"
	localServeAddr = "http://localhost" + servePort
)

func handler(w http.ResponseWriter, r *http.Request) {
	ch := make(chan string)
	go func() {
		j, _ := json.MarshalIndent(flows.ToStat(), "", "\t")
		ch <- string(j)
	}()
	fmt.Fprint(w, <-ch)
}

func Serve() {
	go func() {
		ticker := time.Tick(time.Second)
		for _ = range ticker {
			ConsoleClear()
			fmt.Println("Local Address : " + localServeAddr)
			fmt.Println(runtime.NumGoroutine(), "Goroutines is Running...")
		}
	}()

	webbrowser.Open(localServeAddr)

	http.HandleFunc("/", handler)
	http.ListenAndServe(servePort, nil)
}
