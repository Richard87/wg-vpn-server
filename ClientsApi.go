package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ClientInfo struct {
	Client
	LatestHandshake string `json:"latestHandshake"`
	Endpoint        string `json:"endpoint"`
	SentBytes       int64  `json:"sentBytes"`
	ReceivedBytes   int64  `json:"receivedBytes"`
}
type allClientInfo = []ClientInfo

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
	var allClients = make(allClientInfo, len(*clients))
	for i, client := range *clients {
		allClients[i] = ClientInfo{
			Client:          client,
			LatestHandshake: "2021-01-09T21:15:37.189Z",
			Endpoint:        "77.18.62.145:15427",
			SentBytes:       32165498,
			ReceivedBytes:   132546,
		}
	}

	_ = json.NewEncoder(w).Encode(allClients)
}
func apiClientShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientId, _ := strconv.Atoi(vars["clientId"])

	client := ClientFind(clientId)
	if client != nil {
		w.Header().Add("Content-type", "application/json")
		_ = json.NewEncoder(w).Encode(ClientInfo{
			Client:          *client,
			LatestHandshake: "2021-01-09T21:15:37.189Z",
			Endpoint:        "77.18.62.145:15427",
			SentBytes:       32165498,
			ReceivedBytes:   132546,
		})
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
