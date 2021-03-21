package main

import (
	"encoding/json"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	RequireAuth bool
}

type Routes []Route

func NewRouter(httpsCorsPtr *string, httpsPortPtr *int) http.Handler {

	router := mux.NewRouter().StrictSlash(true)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-type", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{*httpsCorsPtr, fmt.Sprintf("https://localhost:%d", *httpsPortPtr)})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	credentialsOk := handlers.AllowCredentials()

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		if route.RequireAuth {
			handler = Authenticator(handler)
		}
		handler = Logger(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	subFs, _ := fs.Sub(embededFiles, "ui/build")
	router.PathPrefix("/").Handler(http.FileServer(http.FS(subFs)))
	return handlers.CORS(originsOk, headersOk, methodsOk, credentialsOk)(router)
}

type Config struct {
	Endpoint         string `json:"endpoint"`
	NextAvailableIp4 string `json:"nextAvailableIp4"`
	PublicKey        string `json:"publicKey"`
	RecommendedDNS   string `json:"recommendedDNS"`
	Mtu              int    `json:"mtu"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
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

func apiAuthenticate(w http.ResponseWriter, r *http.Request) {
	var login = Login{}

	if err := parseJsonRequest(w, r, &login); err != nil {
		w.WriteHeader(http.StatusInternalServerError) // unprocessable entity
		log.Printf("API: Could not check password hash!: %s", err)
		return
	}

	if login.Password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var user = &User{}
			err := json.Unmarshal(v, user)
			if err != nil {
				return fmt.Errorf("error in users database! Failed user: %d (%v): %s", k, v, err)
			}

			if user.Username != login.Username {
				continue
			}

			match, err := argon2id.ComparePasswordAndHash(login.Password, user.Hash)
			if err != nil {
				return fmt.Errorf("could not check password hash!: %s", err)
			}

			if !match {
				w.WriteHeader(http.StatusUnauthorized)
				return nil
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": login.Username,
				"exp":      time.Now().Add(time.Minute * 115).UnixNano(),
			})

			tokenString, err := token.SignedString(httpsJwtSigningKey)
			if err != nil {
				return fmt.Errorf("could not sign jwt: %s", err)
			}

			parts := strings.Split(tokenString, ".")
			http.SetCookie(w, &http.Cookie{
				Name:     "auth",
				Value:    parts[2],
				Expires:  time.Now().Add(time.Hour * 2),
				MaxAge:   3600 * 2,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			tokenString = fmt.Sprintf("%s.%s", parts[0], parts[1])
			writeJsonResponse(w, &LoginResponse{Token: tokenString})
			return nil
		}

		w.WriteHeader(http.StatusUnauthorized)
		return nil
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // unprocessable entity
		log.Printf("API: %s", err)
		return
	}

}

func parseJsonRequest(w http.ResponseWriter, r *http.Request, newObject interface{}) error {
	content, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalf("Could not read from request stream: %s", err)
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusBadRequest) // unprocessable entity
		return fmt.Errorf("API: Could not parse json: %s", err)
	}
	if err := json.Unmarshal(content, newObject); err != nil {
		w.WriteHeader(http.StatusBadRequest) // unprocessable entity
		return fmt.Errorf("\"API: Could not parse body: %s\"", err)
	}

	return nil
}

func writeJsonResponse(w http.ResponseWriter, object interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(object); err != nil {
		w.WriteHeader(http.StatusInternalServerError) // unprocessable entity
		log.Printf("API: Could not check password hash!: %s", err)
		return
	}
}

func apiGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	ips := findAvailableIps()
	if len(ips) <= 1 {
		log.Print("No IP's in range!")
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

	var nextAvailableIp4 string
	if len(ips) > 0 {
		nextAvailableIp4 = ips[0]
	}
	_ = json.NewEncoder(w).Encode(Config{
		Endpoint:         *wgEndpointPtr + ":" + strconv.Itoa(*wgListenPortPtr),
		NextAvailableIp4: nextAvailableIp4,
		PublicKey:        wgPublicKey.String(),
		RecommendedDNS:   *wgRecommendedDns,
		Mtu:              MTU,
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
		"Auhenticate",
		"POST",
		"/api/authenticate",
		apiAuthenticate,
		false,
	},
	Route{
		"Config",
		"GET",
		"/api/config",
		apiGetConfig,
		true,
	},
	Route{
		"Client List",
		"GET",
		"/api/clients",
		apiClientsIndex,
		true,
	},
	Route{
		"Client Create",
		"POST",
		"/api/clients",
		apiClientCreate,
		true,
	},
	Route{
		"Client",
		"GET",
		"/api/clients/{clientId}",
		apiClientShow,
		true,
	},
	Route{
		"Client Delete",
		"DELETE",
		"/api/clients/{clientId}",
		apiClientRemove,
		true,
	},
}

func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"API: %s %s (duration: %s)",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func Authenticator(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		authParts := strings.Split(authorization, " ")
		if len(authParts) != 2 || authParts[0] != "bearer" || authParts[1] == "" {
			log.Printf("API: %s %s (AUTH DENIED HEADER)", r.Method, r.RequestURI)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		authCookie, err := r.Cookie("auth")
		if err != nil {
			log.Printf("API: %s %s (AUTH DENIED COOKIE)", r.Method, r.RequestURI)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		tokenString := fmt.Sprintf("%s.%s", authParts[1], authCookie.Value)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Header["alg"] != "HS256" {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return httpsJwtSigningKey, nil
		})
		if err != nil || !token.Valid {
			log.Printf("API: %s %s (AUTH DENIED TOKEN: $s)", r.Method, r.RequestURI, err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		inner.ServeHTTP(w, r)
	})
}
