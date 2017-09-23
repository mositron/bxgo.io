package main

import (
	"log"
	"net/http"
)

type Http struct {
	server http.Server
}

func (w *Http) Listen(port int64) {
	mux := http.NewServeMux()
	handler := &MyHandler{}
	mux.Handle("/", handler)
	server := http.Server{Handler: mux, Addr: ":" + _is(port)}
	log.Fatal(server.ListenAndServe())
}
