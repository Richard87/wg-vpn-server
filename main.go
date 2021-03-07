package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Client struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Ip        net.IPAddr `json:"ip"`
	PublicKey string     `json:"publicKey"`
}

//go:embed ui/build
var embededFiles embed.FS

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// You can use the serve file helper to respond to 404 with
	// your request file.
	w.Header().Add("Content-type", "text/html")

	f, _ := embededFiles.Open("ui/build/index.html")

	content, _ := ioutil.ReadAll(f)
	_, _ = w.Write(content)
	_ = f.Close()
}
func main() {
	wgCreateMissingPtr := flag.Bool("wg-create-private-key-if-missing", false, "Set to generate private key if missing. WARNING, This will break existing clients!")
	wgKeyPtr := flag.String("wg-private-key", "./var/wg.private", "Specify WireGuard key file location")

	clientsSubnetPtr := flag.String("client-subnet", "10.0.0.0/24", "Specify default client subnet")
	clientsPtr := flag.String("clients", "./var/clients.db", "Path to store clients.")

	usersPtr := flag.String("user", "", "API User, can be repeated to create more users. For example: \n-user 'admin:$argon2i$v=19$m=16,t=2,p=1$S1p3Z0FTQTViZkh0MURTVA$jxPFAzQ3kSrbEPSibCQIrg'\n(If no users specified, a default admin password will be generated and printed to console")
	httpsPortPtr := flag.Int("https-port", 8443, "API Webserver port")
	httpsKeyPtr := flag.String("https-key", "./var/server.key", "Path to store webserver key (If missing new will be generated).")
	httpsCrtPtr := flag.String("https-crt", "./var/server.crt", "Path to store webserver certificate (If missing new will be generated).")
	httpsCorsPtr := flag.String("https-cors", "http://localhost:3000", "Which clients are allowed to connect (can be repeated)")
	helpPtr := flag.Bool("help", false, "Show this help")
	flag.Parse()

	printConfiguration(helpPtr, wgCreateMissingPtr, wgKeyPtr, clientsPtr, clientsSubnetPtr, httpsPortPtr, httpsCrtPtr, httpsKeyPtr, usersPtr, httpsCorsPtr)

	db, err := bolt.Open(*clientsPtr, 0666, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	srv := runApiServer(httpsCorsPtr, httpsPortPtr, db)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // Block until we receive our signal.

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("WireGuard: shutting down")
	os.Exit(0)
}

func runApiServer(httpsCorsPtr *string, httpsPortPtr *int, db *bolt.DB) *http.Server {
	router := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-type"})
	originsOk := handlers.AllowedOrigins([]string{*httpsCorsPtr})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	handler := handlers.CORS(originsOk, headersOk, methodsOk)(router)
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", *httpsPortPtr),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler, // Pass our instance of gorilla/mux in.
	}
	router.StrictSlash(true)
	// router.PathPrefix("/").HandlerFunc(uiApp)

	subFs, _ := fs.Sub(embededFiles, "ui/build")
	router.Path("/api/clients").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
	router.PathPrefix("/").Handler(http.FileServer(http.FS(subFs)))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println("API Server: started")
	}()
	return srv
}

func printConfiguration(helpPtr *bool, wgCreateMissingPtr *bool, wgKeyPtr *string, clientsPtr *string, clientsSubnetPtr *string, httpsPortPtr *int, httpsCrtPtr *string, httpsKeyPtr *string, usersPtr *string, httpsCorsPtr *string) {
	if *helpPtr {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *wgCreateMissingPtr == true {
		fmt.Print("\n")
		fmt.Print("#################################################\n")
		fmt.Printf("#                                               #\n")
		fmt.Printf("#                !!! WARNIGN !!!                #\n")
		fmt.Printf("#                                               #\n")
		fmt.Printf("# Will generate missing wireguard private key!  #\n")
		fmt.Printf("# (This will break all existing clients!)       #\n")
		fmt.Print("#################################################\n")
		fmt.Print("\n")
	}

	log.Printf("Starting WireGuard VPN Server!")
	log.Printf("Using private key:      %s", *wgKeyPtr)
	log.Printf("Using clients database: %s", *clientsPtr)
	log.Printf("Using client subnet:    %s", *clientsSubnetPtr)
	log.Printf("Running webserver on:   https://0.0.0.0:%d", *httpsPortPtr)
	log.Printf("Using certificate:      %s (key: %s)", *httpsCrtPtr, *httpsKeyPtr)
	log.Printf("Using CORS       :      %v", *httpsCorsPtr)
	log.Printf("Using api-users:        %v", *usersPtr)
	fmt.Print("\n")
}
