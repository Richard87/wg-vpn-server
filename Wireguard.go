package main

import (
	"github.com/Richard87/wg-vpn-server/config"
	"github.com/Richard87/wg-vpn-server/database"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"net"
)

func removeClient(client *database.Client) {
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

func addClient(client *database.Client) {

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
