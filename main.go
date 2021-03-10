package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"embed"
	"encoding/pem"
	"flag"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"time"
)



var (
	//go:embed ui/build
 	embededFiles       embed.FS
 	Db                 *bolt.DB
	wgCreateMissingPtr = flag.Bool("wg-create-private-key-if-missing", false, "Set to generate private key if missing. WARNING, This will break existing clients!")
	wgKeyPtr           = flag.String("wg-private-key", "./var/wg.private", "Specify WireGuard key file location")
	wgEndpointPtr      = flag.String("wg-endpoint", "", "Specify WireGuard public IP and Port. For example 2.2.2.2:51820")
	wgRecommendedDns   = flag.String("wg-dns", "1.1.1.1", "Specify recommended DNS for clients.")
	wgPublicKey        string

	clientsSubnetPtr = flag.String("client-subnet", "10.0.0.0/24", "Specify default client subnet")
	databasePtr      = flag.String("database", "./var/wg.db", "Path to store clients.")

	usersPtr     = flag.String("user", "", "API User, can be repeated to create more users. For example: \n-user 'admin:$argon2i$v=19$m=16,t=2,p=1$S1p3Z0FTQTViZkh0MURTVA$jxPFAzQ3kSrbEPSibCQIrg'\n(If no users specified, a default admin password will be generated and printed to console")
	httpsPortPtr = flag.Int("https-port", 8443, "API Webserver port")
	httpsKeyPtr  = flag.String("https-key", "./var/server_key.pem", "Path to store PKCS8 webserver key (If missing new will be generated).")
	httpsCrtPtr  = flag.String("https-crt", "./var/server_crt.pem", "Path to store webserver certificate (If missing new will be generated).")
	httpsCorsPtr = flag.String("https-cors", "http://localhost:3000", "Which clients are allowed to connect (can be repeated)")
	helpPtr      = flag.Bool("help", false, "Show this help")
)

func main() {

	flag.Parse()

	if *wgEndpointPtr == "" {
		log.Println("You must supply -wg-endpoint. For example vpn.example.com:51820 or 10.10.10:51820")
		os.Exit(1)
	}

	initVarFolder()
	initPrivateKey()
	initPublicKey()
	printConfiguration()

	Db = initDatabase(databasePtr)
	defer Db.Close()

	router := NewRouter(httpsCorsPtr, httpsPortPtr)

	srv := initWebserver(httpsPortPtr, router, httpsCrtPtr, httpsKeyPtr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // Block until we receive our signal.

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("WireGuard: shutting down")
	os.Exit(0)
}

func initPublicKey() {
	privateKey, err := os.ReadFile(*wgKeyPtr)
	if err != nil {
		log.Fatalf("Could not read private key: %v", err)
	}
	key, err := wgtypes.ParseKey(string(privateKey))
	if err != nil {
		log.Fatalf("Could not parse private key: %v", err)
	}
	wgPublicKey = key.PublicKey().String()
}

func initVarFolder() {
	if _, err := os.Stat("./var"); os.IsNotExist(err) {
		if errDir := os.MkdirAll("./var", 0755); errDir != nil {
			log.Fatal(err)
		}
	}
}

func initPrivateKey() {
	if _, err := os.Stat(*wgKeyPtr); os.IsNotExist(err) {
		if *wgCreateMissingPtr {
			log.Print("Creating new WireGuard Private and public key! All existing clients will lose connection!")
			key, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				log.Fatalf("Could not create private key: %v", err)
			}
			err = os.WriteFile(*wgKeyPtr, []byte(key.String()), 0600)
			if err != nil {
				log.Fatalf("Could not save private key: %v", err)
			}
		} else {
			log.Fatalf("Private key does not exist! Use -wg-create-private-key-if-missing to generate it.")
		}
	}
}

func initWebserver(httpsPortPtr *int, router http.Handler, httpsCrtPtr *string, httpsKeyPtr *string) *http.Server {
	if _, err := os.Stat(*httpsKeyPtr); os.IsNotExist(err) {
		createPrivateKey(httpsKeyPtr)
	}

	if _, err := os.Stat(*httpsCrtPtr); os.IsNotExist(err) {
		createCertificate(httpsKeyPtr, httpsCrtPtr)
	}

	cert, err := tls.LoadX509KeyPair(*httpsCrtPtr, *httpsKeyPtr)
	if err != nil {
		log.Fatalf("Failed to open cert.pem for writing: %v", err)
	}

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", *httpsPortPtr),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	go func() {
		if err := srv.ListenAndServeTLS(*httpsCrtPtr, *httpsKeyPtr); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println("API Server: started")
	}()
	return srv
}
func createCertificate(httpsKeyPtr *string, httpsCrtPtr *string) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	SerialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
	template := x509.Certificate{
		SerialNumber: SerialNumber,
		Subject: pkix.Name{
			Organization: []string{"richard87/wg-vpn-server"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),
		IsCA:      true,

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	privateFile, _ := os.ReadFile(*httpsKeyPtr)
	block, _ := pem.Decode(privateFile)
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}
	certOut, err := os.Create(*httpsCrtPtr)
	if err != nil {
		log.Fatalf("Failed to open cert.pem for writing: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert.pem: %v", err)
	}
	_ = certOut.Close()
}
func createPrivateKey(httpsKeyPtr *string) {
	// path/to/whatever exists
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Errorf("could not generate private key! %v", err))
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(fmt.Errorf("could not generate private key! %v", err))
	}

	keyOut, err := os.OpenFile(*httpsKeyPtr, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to key.pem: %v", err)
	}

	if err != nil {
		panic(fmt.Errorf("could not write private key! %v", err))
	}
	_ = keyOut.Close()
}
func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func initDatabase(clientsPtr *string) *bolt.DB {
	thisDb, err := bolt.Open(*clientsPtr, 0666, nil)
	if err != nil {
		panic(err)
	}
	err = thisDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("clients"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return thisDb
}

func printConfiguration() {
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
	log.Printf("Using wg private key:   %s", *wgKeyPtr)
	log.Printf("Using wg public key:    %s", wgPublicKey)
	log.Printf("Using wg DNS:           %s", *wgRecommendedDns)
	log.Printf("Wireguard endpoint:     %s", *wgEndpointPtr)
	log.Printf("Using clients database: %s", *databasePtr)
	log.Printf("Using client subnet:    %s", *clientsSubnetPtr)
	log.Printf("Running webserver on:   https://0.0.0.0:%d", *httpsPortPtr)
	log.Printf("Using certificate:      %s (key: %s)", *httpsCrtPtr, *httpsKeyPtr)
	log.Printf("Using CORS       :      %v", *httpsCorsPtr)
	log.Printf("Using api-users:        %v", *usersPtr)
	fmt.Print("\n")
}
