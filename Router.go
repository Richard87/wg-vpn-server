package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
	"io/fs"
	"log"
	"net/http"
	"time"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(httpsCorsPtr *string, db *bolt.DB) http.Handler {

	router := mux.NewRouter().StrictSlash(true)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-type"})
	originsOk := handlers.AllowedOrigins([]string{*httpsCorsPtr})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	subFs, _ := fs.Sub(embededFiles, "ui/build")
	router.PathPrefix("/").Handler(http.FileServer(http.FS(subFs)))
	return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}

var routes = Routes{
	Route{
		"Client List",
		"GET",
		"/api/clients",
		apiClientsIndex,
	},
	Route{
		"Client Create",
		"POST",
		"/api/clients",
		apiClientCreate,
	},
	Route{
		"TodoShow",
		"GET",
		"/api/clients/{clientId}",
		apiClientShow,
	},
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
