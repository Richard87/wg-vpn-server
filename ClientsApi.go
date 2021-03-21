package main

import (
	"github.com/gorilla/mux"
	"log"
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
	var newClient Client
	if err := parseJsonRequest(w, r, &newClient); err != nil {
		w.WriteHeader(http.StatusInternalServerError) // unprocessable entity
		log.Printf("API: Could not parse request: %s", err)
		return
	}

	client := *ClientCreate(&newClient)

	writeJsonResponse(w, client)
}

func apiClientsIndex(w http.ResponseWriter, r *http.Request) {
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
	writeJsonResponse(w, &allClients)
}
func apiClientShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientId, _ := strconv.Atoi(vars["clientId"])

	client := ClientFind(clientId)
	if client == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writeJsonResponse(w, &ClientInfo{
		Client:          *client,
		LatestHandshake: "2021-01-09T21:15:37.189Z",
		Endpoint:        "77.18.62.145:15427",
		SentBytes:       32165498,
		ReceivedBytes:   132546,
	})
}
func apiClientRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientId, _ := strconv.Atoi(vars["clientId"])

	ClientRemove(clientId)
	w.WriteHeader(http.StatusNoContent)
}
