package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	localNics, _ := net.Interfaces()
	for _, ln := range localNics {
		addrs, e := ln.Addrs()
		if e != nil {
			fmt.Printf("get ip address error %s\n", e)
			os.Exit(1)
		}
		for _, add := range addrs {
			ipnet, ok := add.(*net.IPNet)
			if !ok {
				continue
			}
			ipstr := ipnet.IP.String()
			mask := net.IP(ipnet.Mask)
			maskStr := mask.String()
			vip := net.ParseIP("172.28.2.50")

			if ipnet.Contains(vip) {
				fmt.Printf("nic: %s IP: %s mask: %s, ipnet string %s ", ln.Name, ipstr, maskStr, ipnet.String())
				fmt.Printf("172.28.2.50 is in %s\n", ipnet.String())
			} else {
				fmt.Printf("nic: %s IP: %s mask: %s, ipnet string %s \n", ln.Name, ipstr, maskStr, ipnet.String())
			}
		}
	}
}
