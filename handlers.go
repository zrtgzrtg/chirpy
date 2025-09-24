package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func handlerReady(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}
func handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", int(apiCfg.fileserverHits.Load()))))
}
func handlerReset(w http.ResponseWriter, req *http.Request) {
	apiCfg.reset(w, req)
}
func handlerValidate(w http.ResponseWriter, req *http.Request) {
	type jsonError struct {
		Error string `json:"error"`
	}
	w.Header().Set("Content-type", "application/json")
	body := req.Body
	bBody, err := io.ReadAll(body)
	if err != nil {
		log.Printf("Error reading req.Body: %s", err)
		errStruct := jsonError{"Something went wrong"}
		resp, _ := json.Marshal(&errStruct)
		w.Write(resp)
		w.WriteHeader(http.StatusBadRequest)
	}
	if len(bBody) > 140 {
		errStruct := jsonError{"Chirp is too long"}
		resp, _ := json.Marshal(&errStruct)
		//WriteHeader has to be called before Write
		w.WriteHeader(400)
		w.Write(resp)
	} else {
		type jsonValid struct {
			Valid bool `json:"valid"`
		}
		w.WriteHeader(200)
		validStruct := jsonValid{true}
		resp, _ := json.Marshal(&validStruct)
		w.Write(resp)
	}
}
