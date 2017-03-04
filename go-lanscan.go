package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

const timeoutSeconds = 2
const outboundPort = "80"

// go run *.go -remote "4.4.4.4" -subnet "/25"

func main() {
	remoteIPParameter := flag.String("remote", "8.8.8.8", "The remote IP Address to use to detect outbound connectivity")
	subnetParameter := flag.String("subnet", "", "The CIDR subnet to search, the default is to use the host outbound ip subnet")
	flag.Parse()

	myOutIP, _ := GetOutboundIPAddress(*remoteIPParameter + ":" + outboundPort)
	fmt.Println("Current outbound IP Address:", myOutIP)
	var subnet *net.IPNet
	var err error
	if *subnetParameter == "" {
		subnet, err = GetHostSubnet(myOutIP)
		if err != nil {
			log.Fatal("ERROR unable to get subnet from", myOutIP)
		}
	} else {
		_, subnet, err = net.ParseCIDR(myOutIP + *subnetParameter)
	}

	fmt.Println("Searching subnet:", subnet)
	// TODO check if invalid subnet override
	firstIP := GetFirstIPAddress(*subnet)
	fmt.Printf("%v subnet.Contains(%v): %v \n", subnet, firstIP, subnet.Contains(firstIP))
	lastIP := GetLastIP(*subnet)
	fmt.Printf("%v subnet.Contains(%v): %v \n", subnet, lastIP, subnet.Contains(lastIP))

	addresses, err := GetAllSubnetAddresses(*subnet)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(subnet, "has", len(addresses), "addresses from", addresses[0], "to", addresses[len(addresses)-1])

	fmt.Println("done")
}

// GetHostSubnet gets the subnet of the ip address from the current host https://golang.org/pkg/net/#ParseCIDR
// https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing#Subnet_masks
func GetHostSubnet(ip string) (*net.IPNet, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var myOutAddress net.Addr
	for _, a := range addresses {
		if strings.HasPrefix(a.String(), ip) {
			myOutAddress = a
		}
	}
	_, network, err := net.ParseCIDR(myOutAddress.String())
	if err != nil {
		return nil, err
	}
	return network, err
}

// GetOutboundIPAddress uses the outbound connection interface to return the current ip address
func GetOutboundIPAddress(remoteIP string) (string, error) {
	var ip string
	// https://golang.org/pkg/net/#pkg-overview
	conn, err := net.Dial("udp", remoteIP)
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().String()
		idx := strings.LastIndex(localAddr, ":")
		ip = localAddr[0:idx]
	}
	return ip, err
}

// GetFirstIPAddress gets the first IP Address from the subnet, assumes IPV4 https://golang.org/pkg/net/#IP.To4
func GetFirstIPAddress(network net.IPNet) net.IP {
	networkIPbytes := network.IP.To4()
	firstIP := networkIPbytes.Mask(network.Mask)
	return firstIP
}

// GetLastIP gets the last IP Address from the subnet, assumes IPv4
func GetLastIP(network net.IPNet) net.IP {
	networkIPbytes := network.IP.To4()
	lastIP := net.IPv4(0, 0, 0, 0).To4()
	for i := 0; i < len(lastIP); i++ {
		lastIP[i] = networkIPbytes[i] | ^network.Mask[i]
	}
	return lastIP
}

// GetAllSubnetAddresses returns a list of strings of every address in the subnet
func GetAllSubnetAddresses(subnet net.IPNet) ([]string, error) {
	var err error
	var addresses []string
	firstIP := GetFirstIPAddress(subnet)
	lastIP := GetLastIP(subnet)
	// https://golang.org/pkg/net/#IPNet.Contains
	for ip := firstIP; subnet.Contains(ip); nextAddress(ip) {
		// fmt.Println(ip)
		addresses = append(addresses, firstIP.String())
	}
	if lastIP.String() != addresses[len(addresses)-1] {
		err = fmt.Errorf("Last Subnet IP Address: %s does not equal the incremented last IP Address: %s", lastIP.String(), addresses[len(addresses)-1])
	}
	return addresses, err
}

// nextAddress assumes IPV4 https://en.wikipedia.org/wiki/IP_address#IPv4_addresses
// since an IP address is a slice of length 4 each containing a byte https://golang.org/pkg/net/#IP
func nextAddress(ip net.IP) net.IP {
	// increment the smallest octet first, 2^8 = 256
	if ip[3] < 255 {
		ip[3]++
	} else if ip[2] < 255 {
		ip[3] = 0
		ip[2]++
	} else if ip[1] < 255 {
		ip[3] = 0
		ip[2] = 0
		ip[1]++
	} else {
		ip[3] = 0
		ip[2] = 0
		ip[1] = 0
		ip[0]++
	}
	return ip
}