package main

import (
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"net"
)

func removeClient(client *Client) {
	key, _ := wgtypes.ParseKey(client.PublicKey)

	cfg := wgtypes.Config{
		PrivateKey:   &wgPrivateKey,
		ListenPort:   wgListenPortPtr,
		ReplacePeers: false,
		Peers: []wgtypes.PeerConfig{{
			PublicKey:         key,
			Remove:            true,
			UpdateOnly:        false,
			ReplaceAllowedIPs: false,
			AllowedIPs:        getAllowedIpNets(client),
		}},
	}

	err := wgClient.ConfigureDevice(*wgDeviceName, cfg)
	if err != nil {
		log.Printf("Could not configure device %s: %s", *wgDeviceName, err)
	}
}

func addClient(client *Client) {

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
		PrivateKey:   &wgPrivateKey,
		ListenPort:   wgListenPortPtr,
		ReplacePeers: false,
		Peers:        []wgtypes.PeerConfig{newPeer},
	}

	err = wgClient.ConfigureDevice(*wgDeviceName, cfg)
	if err != nil {
		log.Printf("Could not configure device %s: %s", *wgDeviceName, err)
	}
}

func getAllowedIpNets(client *Client) []net.IPNet {
	allowedIps := []net.IPNet{}
	for _, ip := range client.AllowedIps {
		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			log.Printf("Could not parse client ip (%s): %s", ip, err)
		}

		allowedIps = append(allowedIps, *ipNet)
	}
	return allowedIps
}
