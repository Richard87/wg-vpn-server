package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strconv"
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
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

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

type Config struct {
	Endpoint         string `json:"endpoint"`
	NextAvailableIp4 string `json:"nextAvailableIp4"`
	PublicKey        string `json:"publicKey"`
	RecommendedDNS   string `json:"recommendedDNS"`
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func apiGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	ips := findAvailableIps()
	if len(ips) <= 1 {
		log.Print("No IP's in range!")
		w.WriteHeader(400)
		return
	}

	// Keep first ip for gateway
	ips = ips[1:]

	_ = Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("clients"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var client = &Client{}
			err := json.Unmarshal(v, client)
			if err != nil {
				log.Print(err)
				continue
			}

			for _, ip := range client.AllowedIps {
				remove(ips, ip)
			}
		}

		return nil
	})

	_ = json.NewEncoder(w).Encode(Config{
		Endpoint:         *wgEndpointPtr + ":" + strconv.Itoa(*wgListenPortPtr),
		NextAvailableIp4: ips[0],
		PublicKey:        wgPublicKey.String(),
		RecommendedDNS:   *wgRecommendedDns,
	})
}

func findAvailableIps() []string {
	ip, ipnet, err := net.ParseCIDR(*clientsSubnetPtr)
	if err != nil {
		return []string{}
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String()+"/32")
	}

	if len(ips) <= 2 || ips == nil {
		return []string{}
	}

	ips = ips[1 : len(ips)-1]
	return ips
}

var routes = Routes{
	Route{
		"Config",
		"GET",
		"/api/config",
		apiGetConfig,
	},
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
