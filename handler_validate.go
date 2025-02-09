package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type notValid struct {
		Error string `json:"error"`
	}

	type valid struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	ch := chirp{}
	if err := decoder.Decode(&ch); err != nil {
		log.Printf("error decoding request body: %s", err)
		w.WriteHeader(500)
		return
	}

	var data []byte
	var err error
	w.Header().Set("Content-Type", "application/json")

	if len(ch.Body) > 0 && len(ch.Body) <= 140 {
		res := valid{
			Valid: true,
		}
		data, err = json.Marshal(res)
		if err != nil {
			log.Printf("could not encode data: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	} else {
		res := notValid{
			Error: "Chirp is too long",
		}
		data, err = json.Marshal(res)
		if err != nil {
			log.Printf("could not encode data: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
	}

	w.Write(data)
}
