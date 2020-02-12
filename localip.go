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
				if !ip.IP.IsGlobalUnicast() {
					continue
				}
				ips = append(ips, ip.IP)
			}
		}
	}
	return
}
