package main

import (
	// internals
	"fmt"
	"os"

	// externals
	"github.com/GlenDC/go-external-ip"
	//"github.com/libp2p/go-libp2p"
	//"github.com/libp2p/go-libp2p-crypto"
	//"github.com/libp2p/go-libp2p-host"
	//"github.com/libp2p/go-libp2p-net"
	//"github.com/libp2p/go-libp2p-peer"
	//"github.com/libp2p/go-libp2p-peerstore"
	//"github.com/multiformats/go-multiaddr"
)

func initHybridModeServer() {

	//var addr string

	if question("are you sure?", []string{"yes", "no"}) == "no" {

		return

	}

	fmt.Printf("retrieving external ip address... ")
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {

		fmt.Printf("failed\n")
		fmt.Printf("[err]: unable to retrieve external ip address...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	fmt.Printf("done. ip is %s\n", ip.String())
	return

}
