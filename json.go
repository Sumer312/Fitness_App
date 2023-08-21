package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMessage struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Fatal("internal server error")
	}
	respondWithJson(w, 400, ErrorMessage{
		Error: msg,
	})

}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Error Marshalling to JSON")
		w.WriteHeader(500)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}
