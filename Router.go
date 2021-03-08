package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

func NewRouter(httpsCorsPtr *string, httpsPortPtr *int) http.Handler {

	router := mux.NewRouter().StrictSlash(true)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-type", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{*httpsCorsPtr, fmt.Sprintf("https://localhost:%d", *httpsPortPtr)})
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
		"Client",
		"GET",
		"/api/clients/{clientId}",
		apiClientShow,
	},
	Route{
		"Client Delete",
		"DELETE",
		"/api/clients/{clientId}",
		apiClientRemove,
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
