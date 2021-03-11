package main

import (
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"net"
)

func removeClient(client *Client) {
	device, err := wgClient.Device(*wgDeviceName)
	if err != nil {
		log.Fatalf("WireGuard: Could not configure device %s: %s", *wgDeviceName, err)
	}
	peers := []wgtypes.PeerConfig{}

	for _, peer := range device.Peers {
		shouldRemove := peer.PublicKey.String() == client.PublicKey
		peer := wgtypes.PeerConfig{
			PublicKey:         peer.PublicKey,
			Remove:            shouldRemove,
			UpdateOnly:        true,
			ReplaceAllowedIPs: true,
			AllowedIPs:        peer.AllowedIPs,
		}
		peers = append(peers, peer)
	}

	cfg := wgtypes.Config{
		PrivateKey:   &wgPrivateKey,
		ListenPort:   wgListenPortPtr,
		FirewallMark: &device.FirewallMark,
		ReplacePeers: false,
		Peers:        peers,
	}

	err = wgClient.ConfigureDevice(*wgDeviceName, cfg)
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

	device, err := wgClient.Device(*wgDeviceName)
	if err != nil {
		log.Fatalf("WireGuard: Could not configure device %s: %s", *wgDeviceName, err)
	}
	peers := []wgtypes.PeerConfig{}

	for _, peer := range device.Peers {
		peer := wgtypes.PeerConfig{
			PublicKey:         peer.PublicKey,
			Remove:            false,
			UpdateOnly:        false,
			ReplaceAllowedIPs: false,
			AllowedIPs:        peer.AllowedIPs,
		}
		peers = append(peers, peer)
	}
	peers = append(peers, newPeer)

	cfg := wgtypes.Config{
		PrivateKey:   &wgPrivateKey,
		ListenPort:   wgListenPortPtr,
		FirewallMark: &device.FirewallMark,
		ReplacePeers: false,
		Peers:        peers,
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
