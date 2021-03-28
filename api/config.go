package api

import (
	"github.com/Richard87/wg-vpn-server/config"
	"github.com/Richard87/wg-vpn-server/database"
	"github.com/gofiber/fiber/v2"
	"log"
	"net"
	"net/http"
	"strconv"
)

const MTU = 1420

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

func findAvailableIps() []string {
	ip, ipnet, err := net.ParseCIDR("10.0.0.0/24")
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

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
func GetConfig(c *fiber.Ctx) error {
	ips := findAvailableIps()
	if len(ips) <= 1 {
		log.Print("No IP's in range!")
	}
	// Keep first ip for gateway
	ips = ips[1:]

	clients := []database.Client{}
	database.Connection.Find(&clients)
	for _, c := range clients {
		remove(ips, c.AllowedIp4)
	}

	var nextAvailableIp4 string
	if len(ips) > 0 {
		nextAvailableIp4 = ips[0]
	}

	return c.Status(http.StatusOK).JSON(Config{
		Endpoint:         config.Config.WgEndpoint + ":" + strconv.Itoa(config.Config.WgListenPort),
		NextAvailableIp4: nextAvailableIp4,
		PublicKey:        config.Config.WgPublicKey.String(),
		RecommendedDNS:   config.Config.WgRecommendedDns,
		Mtu:              MTU,
	})

}
