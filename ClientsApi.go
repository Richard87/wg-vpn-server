package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func apiClientCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	var newClient Client
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &newClient); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	client := *ClientCreate(&newClient)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(client); err != nil {
		panic(err)
	}
}
func apiClientsIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	clients := ClientList()

	_ = json.NewEncoder(w).Encode(clients)
}
func apiClientShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientId, _ := strconv.Atoi(vars["clientId"])

	client := ClientFind(clientId)
	if client != nil {
		w.Header().Add("Content-type", "application/json")
		json.NewEncoder(w).Encode(*client)
	} else {
		w.WriteHeader(404)
	}
}
func apiClientRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientId, _ := strconv.Atoi(vars["clientId"])

	ClientRemove(clientId)
	w.WriteHeader(204)
}
