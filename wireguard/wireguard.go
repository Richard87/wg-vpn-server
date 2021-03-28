package wireguard

import (
	"fmt"
	"github.com/Richard87/wg-vpn-server/config"
	"github.com/Richard87/wg-vpn-server/database"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"net"
	"os"
)

func Init() {
	initPrivateKey()
	initPublicKey()
	// TODO: Initialise interface

}

func Close() {
	log.Println("Closing WG...")
	// <-wgDevice.Wait() // Should we wait for this?! Why?
	log.Println("Closing UAPI...")
	// clean up
	_ = uapi.Close()

	log.Println("Closing device...")
	wgDevice.Close()
	log.Println("WG Closed.")
}

var (
	uapi     net.Listener
	wgDevice *device.Device
)

func Run() {
	client, err := wgctrl.New()
	if err != nil {
		log.Fatalf("Could not connect to WireGuard Controller: %v", err)
	}
	allClients := []database.Client{}
	database.Connection.Find(&allClients)

	wgDevice, err := client.Device(config.Config.WgDeviceName)
	if err != nil {
		runEmbeddedWireGuard()
	}

	wgDevice, err = client.Device(config.Config.WgDeviceName)
	if err != nil {
		log.Fatalf("WG: Could not connect to Controller: %v", err)
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
		log.Printf("WG: Could not configure device %s: %s", config.Config.WgDeviceName, err)
	}

	for _, client := range allClients {
		key, err := wgtypes.ParseKey(client.PublicKey)
		if err != nil {
			log.Printf("WG: Could not add missing client %s: %s", client.Name, err)
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
		log.Printf("WG: Could not configure device %s: %s", config.Config.WgDeviceName, err)
	}

	config.Config.WgClient = client
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

func runEmbeddedWireGuard() {
	tunDevice, err := tun.CreateTUN(config.Config.WgDeviceName, config.MTU)

	if err != nil {
		log.Fatalf("WG: Could not create interface %s: %s", config.Config.WgDeviceName, err)
	}
	logger := device.NewLogger(1, fmt.Sprintf("WG: (%s) ", config.Config.WgDeviceName))
	wgDevice = device.NewDevice(tunDevice, logger)

	errs := make(chan error)

	fileUapi, err := ipc.UAPIOpen(config.Config.WgDeviceName)
	if err != nil {
		log.Printf("WG: Failed to openuapi socket: %v", err)
		os.Exit(1)
	}
	uapi, err = ipc.UAPIListen(config.Config.WgDeviceName, fileUapi)
	if err != nil {
		log.Printf("WG: Failed to listen on uapi socket: %v", err)
		os.Exit(1)
	}

	go func() {
		for {
			conn, err := uapi.Accept()
			if err != nil {
				errs <- err
				return
			}
			go wgDevice.IpcHandle(conn)
		}
	}()

	log.Println("WG: Running embedded server")
}

func RemoveClient(client *database.Client) {
	key, _ := wgtypes.ParseKey(client.PublicKey)

	cfg := wgtypes.Config{
		PrivateKey:   &config.Config.WgPrivateKey,
		ListenPort:   &config.Config.WgListenPort,
		ReplacePeers: false,
		Peers: []wgtypes.PeerConfig{{
			PublicKey:         key,
			Remove:            true,
			UpdateOnly:        false,
			ReplaceAllowedIPs: false,
			AllowedIPs:        getAllowedIpNets(client),
		}},
	}

	err := config.Config.WgClient.ConfigureDevice(config.Config.WgDeviceName, cfg)
	if err != nil {
		log.Printf("Could not configure device %s: %s", config.Config.WgDeviceName, err)
	}
}

func AddClient(client *database.Client) {

	key, err := wgtypes.ParseKey(client.PublicKey)
	if err != nil {
		log.Printf("Could not parse client key: %s", err)
	}
	newPeer := wgtypes.PeerConfig{
		PublicKey:         key,
		Remove:            false,
		UpdateOnly:        false,
		ReplaceAllowedIPs: false,
		AllowedIPs:        getAllowedIpNets(client),
	}

	cfg := wgtypes.Config{
		PrivateKey:   &config.Config.WgPrivateKey,
		ListenPort:   &config.Config.WgListenPort,
		ReplacePeers: false,
		Peers:        []wgtypes.PeerConfig{newPeer},
	}

	err = config.Config.WgClient.ConfigureDevice(config.Config.WgDeviceName, cfg)
	if err != nil {
		log.Printf("Could not configure device %s: %s", config.Config.WgDeviceName, err)
	}
}

func getAllowedIpNets(client *database.Client) []net.IPNet {
	var allowedIps []net.IPNet
	_, ipNet, err := net.ParseCIDR(client.AllowedIp4)
	if err != nil {
		log.Printf("Could not parse client ip (%s): %s", client.AllowedIp4, err)
	}

	allowedIps = append(allowedIps, *ipNet)
	return allowedIps
}
