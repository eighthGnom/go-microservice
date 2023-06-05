package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	bindAddr = "80"
)

type Config struct {
}

func main() {

	app := &Config{}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", bindAddr),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
