package config

import (
	crand "crypto/rand"
	"flag"
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"os"
	"runtime"
	"strings"
)

type UsersFlag []string

func (i *UsersFlag) String() string {
	return "API User, can be repeated to create more users. For example: \n" +
		"-user 'admin:$argon2i$v=19$m=16,t=2,p=1$S1p3Z0FTQTViZkh0MURTVA$jxPFAzQ3kSrbEPSibCQIrg'\n" +
		"(If no users specified, a default admin password will be generated and printed to console"
}

func (i *UsersFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var Config = &ConfigStruct{}

const MTU = 1420

type ConfigStruct struct {
	WgCreateMissing    bool
	WgKey              string
	WgEndpoint         string
	WgListenPort       int
	WgRecommendedDns   string
	WgDeviceName       string
	WgBoringtunPath    string
	WgPublicKey        wgtypes.Key
	WgPrivateKey       wgtypes.Key
	ClientsSubnet      string
	Database           string
	Users              UsersFlag
	HttpsPort          string
	HttpsKey           string
	HttpsCrt           string
	HttpsCors          string
	HttpsJwtSigningKey []byte
	Help               bool
	WgClient           *wgctrl.Client
}

func InitFlags() {
	initVarFolder()
	defaultWgDeviceName := "wg0"
	if runtime.GOOS == "darwin" {
		defaultWgDeviceName = "utun0"
	}

	flag.BoolVar(&Config.WgCreateMissing, "wg-create-private-key-if-missing", false, "Set to generate private key if missing. WARNING, This will break existing clients!")
	flag.StringVar(&Config.WgKey, "wg-private-key", "./var/wg.private", "Specify WireGuard key file location")
	flag.StringVar(&Config.WgEndpoint, "wg-endpoint", "", "Specify WireGuard public IP and Port. For example 2.2.2.2")
	flag.IntVar(&Config.WgListenPort, "wg-listen-port", 51820, "Specify WireGuard Listen port")
	flag.StringVar(&Config.WgRecommendedDns, "wg-dns", "1.1.1.1", "Specify recommended DNS for clients.")
	flag.StringVar(&Config.WgDeviceName, "wg-device", defaultWgDeviceName, "WireGuard device name (must be utunX on Mac=")
	flag.StringVar(&Config.WgBoringtunPath, "wg-boringtun", "", "Path to boringtun")
	flag.StringVar(&Config.ClientsSubnet, "client-subnet", "10.0.0.0/24", "Specify default client subnet")
	flag.StringVar(&Config.Database, "database", "./var/wg.db", "Path to store clients.")
	flag.StringVar(&Config.HttpsPort, "https-port", "8443", "API Webserver port")
	flag.StringVar(&Config.HttpsKey, "https-key", "./var/server_key.pem", "Path to store PKCS8 webserver key (If missing new will be generated).")
	flag.StringVar(&Config.HttpsCrt, "https-crt", "./var/server_crt.pem", "Path to store webserver certificate (If missing new will be generated).")
	flag.StringVar(&Config.HttpsCors, "https-cors", "https://localhost:3000", "Which clients are allowed to connect (can be repeated)")
	flag.BoolVar(&Config.Help, "help", false, "Show this help")
	flag.Var(&Config.Users, "user", "API User, can be repeated to create more users. For example: \n-user 'admin:$argon2i$v=19$m=16,t=2,p=1$S1p3Z0FTQTViZkh0MURTVA$jxPFAzQ3kSrbEPSibCQIrg'\n(If no users specified, a default admin password will be generated and printed to console")

	flag.Parse()

	if Config.WgEndpoint == "" {
		log.Println("You must supply -wg-endpoint. For example vpn.example.com:51820 or 10.10.10:51820")
		os.Exit(1)
	}

	if runtime.GOOS == "darwin" && strings.HasPrefix(Config.WgDeviceName, "utun") == false {
		log.Fatalf("WG: Device name must be utun[0-9]*: %s", Config.WgDeviceName)
	}

	signingKey := make([]byte, 12)
	_, _ = crand.Read(signingKey)
	Config.HttpsJwtSigningKey = signingKey

	printConfiguration()
}

func printConfiguration() {
	if Config.Help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if Config.WgCreateMissing == true {
		fmt.Print("\n")
		fmt.Print("#################################################\n")
		fmt.Printf("#                                               #\n")
		fmt.Printf("#                !!! WARNING !!!                #\n")
		fmt.Printf("#                                               #\n")
		fmt.Printf("# Will generate missing wireguard private key!  #\n")
		fmt.Printf("# (This will break all existing clients!)       #\n")
		fmt.Print("#################################################\n")
		fmt.Print("\n")
	}

	log.Printf("Starting WireGuard VPN Server!")
	log.Printf("Using wg private key:   %s", Config.WgKey)
	log.Printf("Using wg public key:    %s", Config.WgPublicKey)
	log.Printf("Using wg DNS:           %s", Config.WgRecommendedDns)
	log.Printf("using wg endpoint:      %s", Config.WgEndpoint)
	log.Printf("using wg boringtun:     %s", Config.WgBoringtunPath)
	log.Printf("Using client subnet:    %s", Config.ClientsSubnet)
	log.Printf("Running webserver on:   https://0.0.0.0:%s", Config.HttpsPort)
	log.Printf("Using certificate:      %s (key: %s)", Config.HttpsCrt, Config.HttpsKey)
	log.Printf("Using CORS       :      %v", Config.HttpsCors)
	fmt.Println("\n")
}

func initVarFolder() {
	if _, err := os.Stat("./var"); os.IsNotExist(err) {
		if errDir := os.MkdirAll("./var", 0755); errDir != nil {
			log.Fatal(err)
		}
	}
}
