package main

import (
	"log"
	"movieapp/rating/internal/controller/rating"
	httphandler "movieapp/rating/internal/handler/http"
	"movieapp/rating/internal/repository/memory"
	"net/http"
)

func main() {
	log.Println("Starting Rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/ratings", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
