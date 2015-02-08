package main

import (
	"log"
	"net"
)

func localIP() (ips []net.IP) {
	ifaces, err := net.Interfaces()
	check(err)
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.Println("error getting addrs for", iface)
			continue
		}
		for _, a := range addrs {
			if ip, ok := a.(*net.IPNet); ok {
				if ipv4 := ip.IP.To4(); ipv4 != nil {
					ips = append(ips, ipv4)
				}
			}
		}
	}
	return
}
