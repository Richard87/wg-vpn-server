package main

import (
	"embed"
	"fmt"
	"github.com/Richard87/wg-vpn-server/api"
	"github.com/Richard87/wg-vpn-server/config"
	"github.com/Richard87/wg-vpn-server/database"

	_ "github.com/mdlayher/genetlink"
	_ "github.com/mdlayher/netlink"
	_ "github.com/mdlayher/netlink/nlenc"
	"golang.zx2c4.com/wireguard/device"
	_ "golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//go:embed ui/build
var embededFiles embed.FS

func main() {
	config.InitFlags()
	database.InitDatabase()
	database.InitUsers()

	initWireguard()
	api.InitRouter(embededFiles)

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM)
	signal.Notify(termSignal, os.Interrupt)
	<-termSignal // Block until we receive our signal.
	log.Println("API: shutting down")

	_ = api.Router.Shutdown()
	os.Exit(0)
}

func initPublicKey() {
	privateKey, err := os.ReadFile(config.Config.WgKey)
	if err != nil {
		log.Fatalf("Could not read private key: %v", err)
	}
	key, err := wgtypes.ParseKey(string(privateKey))
	if err != nil {
		log.Fatalf("Could not parse private key: %v", err)
	}
	config.Config.WgPublicKey = key.PublicKey()
	config.Config.WgPrivateKey = key
}

func initPrivateKey() {
	if _, err := os.Stat(config.Config.WgKey); os.IsNotExist(err) {
		if config.Config.WgCreateMissing {
			log.Print("Creating new WireGuard Private and public key! All existing clients will lose connection!")
			key, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				log.Fatalf("Could not create private key: %v", err)
			}
			err = os.WriteFile(config.Config.WgKey, []byte(key.String()), 0600)
			if err != nil {
				log.Fatalf("Could not save private key: %v", err)
			}
		} else {
			log.Fatalf("Private key does not exist! Use -wg-create-private-key-if-missing to generate it.")
		}
	}
}

func initWireguard() {
	initPrivateKey()
	initPublicKey()

	// TODO:
	// [X] Currently using embedded Wireguard if interface doesn't exist
	// [ ] Check if interface exist
	// [ ] If it doesn't exist, create it
	// [ ] Use native wireguard if it exist
	//   [ ] If native doesnt exist, start boringtun if exist in path / boringtunPath variable
	//   [ ] if boringtun doesn't exist, use embedded

	client, err := wgctrl.New()
	if err != nil {
		log.Fatalf("Could not connect to WireGuard Controller: %v", err)
	}
	allClients := []database.Client{}
	database.Connection.Find(&allClients)

	wgDevice, err := client.Device(config.Config.WgDeviceName)
	if err != nil {
		ready := make(chan bool, 1)
		go runEmbeddedWireguard(ready)
		<-ready
	}

	wgDevice, err = client.Device(config.Config.WgDeviceName)
	if err != nil {
		log.Fatalf("Could not connect to WireGuard Controller: %v", err)
	}

	cfg := wgtypes.Config{
		PrivateKey:   &config.Config.WgPrivateKey,
		ListenPort:   &config.Config.WgListenPort,
		ReplacePeers: false,
		Peers:        []wgtypes.PeerConfig{},
	}
	for _, peer := range wgDevice.Peers {
		cfg.Peers = append(cfg.Peers, wgtypes.PeerConfig{
			PublicKey:         peer.PublicKey,
			Remove:            true,
			UpdateOnly:        false,
			ReplaceAllowedIPs: false,
			AllowedIPs:        peer.AllowedIPs,
		})
	}
	err = client.ConfigureDevice(config.Config.WgDeviceName, cfg)
	if err != nil {
		log.Printf("Could not configure device %s: %s", config.Config.WgDeviceName, err)
	}

	for _, client := range allClients {
		key, err := wgtypes.ParseKey(client.PublicKey)
		if err != nil {
			log.Printf("Could not add missing WireGuard Client %s: %s", client.Name, err)
			continue
		}

		newPeer := wgtypes.PeerConfig{
			PublicKey:         key,
			Remove:            false,
			UpdateOnly:        false,
			ReplaceAllowedIPs: false,
			AllowedIPs:        getAllowedIpNets(&client),
		}

		cfg.Peers = append(cfg.Peers, newPeer)
	}

	err = client.ConfigureDevice(config.Config.WgDeviceName, cfg)
	if err != nil {
		log.Printf("Could not configure device %s: %s", config.Config.WgDeviceName, err)
	}

	config.Config.WgClient = client
}

func runEmbeddedWireguard(ready chan bool) {
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt)
	tunDevice, err := tun.CreateTUN(config.Config.WgDeviceName, config.MTU)

	if err != nil {
		log.Fatalf("WireGuard: Could not create interface %s: %s", config.Config.WgDeviceName, err)
	}
	logger := device.NewLogger(
		3,
		fmt.Sprintf("(%s) ", config.Config.WgDeviceName),
	)
	wgInternalDevice := device.NewDevice(tunDevice, logger)

	errs := make(chan error)

	fileUapi, err := ipc.UAPIOpen(config.Config.WgDeviceName)
	if err != nil {
		log.Printf("Failed to openuapi socket: %v", err)
		os.Exit(1)
	}
	uapi, err := ipc.UAPIListen(config.Config.WgDeviceName, fileUapi)
	if err != nil {
		log.Printf("Failed to listen on uapi socket: %v", err)
		os.Exit(1)
	}

	go func() {
		for {
			conn, err := uapi.Accept()
			if err != nil {
				errs <- err
				return
			}
			go wgInternalDevice.IpcHandle(conn)
		}
	}()

	ready <- true
	// wait for program to terminate

	select {
	case <-term:
	case <-errs:
	case <-wgInternalDevice.Wait():
	}

	// clean up
	_ = uapi.Close()
	wgInternalDevice.Close()
}
