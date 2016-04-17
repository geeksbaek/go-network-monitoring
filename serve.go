package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/toqueteos/webbrowser"
)

var (
	servePort      = ":8080"
	localServeAddr = "http://localhost" + servePort
)

func handler(w http.ResponseWriter, r *http.Request) {
	ch := make(chan string)
	go func() {
		j, _ := json.Marshal(flow.GetStats())
		ch <- string(j)
	}()
	fmt.Fprint(w, <-ch)
}

func Serve() {
	fmt.Print(CLEAR)
	fmt.Printf("Running... %s\n", localServeAddr)

	webbrowser.Open(localServeAddr)

	http.HandleFunc("/", handler)
	http.ListenAndServe(servePort, nil)
}
