// +build !linux

package wireguard

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func configureInterface(networkCidr string, ifaceName string) error {
	availableIps := findAvailableIps(networkCidr)
	if len(availableIps) == 0 {
		return fmt.Errorf("no availble ips in client subnet (%s)", networkCidr)
	}

	iface, _ := netlink.LinkByName(ifaceName)
	addr, _ := netlink.ParseAddr(ip)
	netlink.AddrAdd(lo, addr)
	return nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func findAvailableIps(networkCidr string) []string {
	ip, ipnet, err := net.ParseCIDR(networkCidr)
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
