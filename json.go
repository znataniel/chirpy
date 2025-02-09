package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondJson(w http.ResponseWriter, code int, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("error marshalling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func respondJsonError(w http.ResponseWriter, code int, err error, msg string) {
	type errJson struct {
		Error string `json:"error"`
	}

	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("responding with 5XX code: %v", code)
	}

	respondJson(w, code, errJson{
		Error: msg,
	})
}
