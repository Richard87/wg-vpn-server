package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

//go:embed ui/build
var embededFiles embed.FS

func uiApp(w http.ResponseWriter, r *http.Request) {
	filename := r.RequestURI
	if filename == "/" || filename == "" {
		filename = "/index.html"
	}

	filename = "ui/build" + filename
	extension := filepath.Ext(filename)
	contentType := "text/html"
	switch extension {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "text/javascript"
	case ".svg":
		contentType = "image/svg+xml"
	case ".map":
		contentType = "text/plain"
	}
	w.Header().Add("Content-type", contentType)

	log.Printf("%s (%s: %s)", filename, extension, contentType)
	file, err := embededFiles.Open(filename)
	if err != nil {
		log.Printf("API: Warning: %v", err)
		return
	}
	data := make([]byte, 512)
	for {
		data = data[:cap(data)]
		n, err := file.Read(data)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		data = data[:n]

		_, err = w.Write(data)
		if err != nil {
			log.Printf("API: Warning: %v", err)
			break
		}
	}

	err = file.Close()
	if err != nil {
		log.Printf("API: Warning: %v", err)
	}
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
	helpPtr := flag.Bool("help", false, "Show this help")
	flag.Parse()

	printConfiguration(helpPtr, wgCreateMissingPtr, wgKeyPtr, clientsPtr, clientsSubnetPtr, httpsPortPtr, httpsCrtPtr, httpsKeyPtr, usersPtr)
	router := mux.NewRouter()
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", *httpsPortPtr),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}
	router.StrictSlash(true)
	router.PathPrefix("/").HandlerFunc(uiApp)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println("API Server: started")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // Block until we receive our signal.

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("WireGuard: shutting down")
	os.Exit(0)
}

func printConfiguration(helpPtr *bool, wgCreateMissingPtr *bool, wgKeyPtr *string, clientsPtr *string, clientsSubnetPtr *string, httpsPortPtr *int, httpsCrtPtr *string, httpsKeyPtr *string, usersPtr *string) {
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
	log.Printf("Using api-users:        %v", *usersPtr)
	fmt.Print("\n")
}
