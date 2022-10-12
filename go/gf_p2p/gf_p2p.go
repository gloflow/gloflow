/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_p2p

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"
	"bufio"
	"sync"
	"log"
	"flag"
	"strings"
	
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/tls"

	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/libp2p/go-libp2p/core/peer"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	multiaddr "github.com/multiformats/go-multiaddr"
)

//-------------------------------------------------
type GFp2pPeerPingFun func() ping.Result
type GFp2pPeerInitFun func(peer.ID) GFp2pPeerPingFun
type GFp2pAddrLst []multiaddr.Multiaddr

type GFp2pConfig struct {
	RendezvousString string
	BootstrapPeers   GFp2pAddrLst
	ListenAddresses  GFp2pAddrLst
	ProtocolID       string
}

var logger = log.Default()

//-------------------------------------------------
func ParseFlags() (GFp2pConfig, error) {
	config := GFp2pConfig{}

	flag.StringVar(&config.RendezvousString, "rendezvous", "meet me here",
		"unique string to identify group of nodes. share this with node operators to connect to GF network")

	flag.Var(&config.BootstrapPeers, "peer", "Adds a peer multiaddress to the bootstrap list")
	flag.Var(&config.ListenAddresses, "listen", "Adds a multiaddress to the listen list")
	flag.StringVar(&config.ProtocolID, "pid", "/chat/1.1.0", "Sets a protocol id for stream headers")
	flag.Parse()

	if len(config.BootstrapPeers) == 0 {
		config.BootstrapPeers = dht.DefaultBootstrapPeers
	}

	return config, nil
}

func (al *GFp2pAddrLst) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *GFp2pAddrLst) Set(value string) error {
	addr, err := multiaddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}

func StringsToAddrs(addrStrings []string) (maddrs []multiaddr.Multiaddr, err error) {
	for _, addrString := range addrStrings {
		addr, err := multiaddr.NewMultiaddr(addrString)
		if err != nil {
			return maddrs, err
		}
		maddrs = append(maddrs, addr)
	}
	return
}

//-------------------------------------------------
func Init(pPortInt int) (host.Host, GFp2pPeerInitFun) {

	//----------------
	// KEYPAIR
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	//----------------
	// CONNECTION_MANAGER
	// prevent peer from having too many connections
	connmgr, err := connmgr.NewConnManager(
		100, // low water
		400, // high water
		connmgr.WithGracePeriod(time.Minute), // grace period
	)
	if err != nil {
		panic(err)
	}

	//----------------
	// OPTIONS
	// https://github.com/libp2p/go-libp2p/blob/master/options.go

	// configures libp2p to use the given private key to identify itself
	identityOption := libp2p.Identity(priv)

	// listen addresses
	addressOption := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", pPortInt), // regular tcp connections
	)

	// configures libp2p to use the given security transport
	// support TLS connections
	securityOption := libp2p.Security(libp2ptls.ID, libp2ptls.New)

	connectionManagerOption := libp2p.ConnectionManager(connmgr)

	//----------------

	node, err := libp2p.New(
		identityOption,
		addressOption,
		securityOption,
		connectionManagerOption,
		
		// disable built-in ping protocol
		libp2p.Ping(false),
	)

	if err != nil {
		panic(err)
	}
	// defer node.Close()

	fmt.Printf("node Listen addresses: %s\n", node.Addrs())
	fmt.Printf("node hosts ID is %s\n", node.ID())

	peerInfo := peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	peerAddrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	fmt.Println("libp2p node address:", peerAddrs[0])
	



	//----------------
	// configure our own ping protocol
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)
	

	//-------------------------------------------------
	pingInitPeerFun := func(pPeerID peer.ID) GFp2pPeerPingFun {

		pingCh := pingService.Ping(context.Background(), pPeerID)
		pingPeerFun := func() ping.Result {
			res := <-pingCh
			return res
		}

		return GFp2pPeerPingFun(pingPeerFun)
	}

	//-------------------------------------------------

	//----------------


	


	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	initPeerDiscovery(node, config)
	InitStreamHandler(node)

	return node, GFp2pPeerInitFun(pingInitPeerFun)
}

//-------------------------------------------------
func initPeerDiscovery(pNode host.Host,
	pConfig GFp2pConfig) {

	// CONFIG
	bootstrapPeers   := pConfig.BootstrapPeers
	randezvousString := pConfig.RendezvousString
	protocolID       := pConfig.ProtocolID



	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	ctx := context.Background()
	kademliaDHT, err := dht.New(ctx, pNode)
	if err != nil {
		panic(err)
	}



	// bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	logger.Print("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}



	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup

	

	for _, peerAddr := range bootstrapPeers {

		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := pNode.Connect(ctx, *peerinfo); err != nil {
				logger.Print(err)
			} else {
				logger.Print("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()



	//----------------
	// ANNOUNCING RANDEZVOUS

	logger.Print("Announcing ourselves...")
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)

	// makes this node announce that it can provide a value for the given key,
	// key being the "randezvous string"
	
	dutil.Advertise(ctx, routingDiscovery, randezvousString)

	logger.Print("Successfully announced!")

	//----------------
	// look for others peers who have announced
	logger.Print("Searching for other peers...")

	peerChan, err := routingDiscovery.FindPeers(ctx, randezvousString)
	if err != nil {
		panic(err)
	}

	for peer := range peerChan {
		
		if peer.ID == pNode.ID() {
			continue
		}
		logger.Print("Found peer:", peer)

		logger.Print("Connecting to:", peer)
		stream, err := pNode.NewStream(ctx, peer.ID, protocol.ID(protocolID))

		if err != nil {
			logger.Print("Connection failed:", err)
			continue
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeData(rw)
			go readData(rw)
		}

		logger.Print("Connected to:", peer)
	}

	select {}
}

//-------------------------------------------------
func InitStreamHandler(pNode host.Host) {

	//-------------------------------------------------
	streamHandlerFun := func(pStream network.Stream) {

		// create a buffer stream for non blocking read and write
		rw := bufio.NewReadWriter(bufio.NewReader(pStream), bufio.NewWriter(pStream))
	
		go readData(rw)
		go writeData(rw)
	
		// 'stream' will stay open until you close it (or the other side closes it).
	}

	//-------------------------------------------------
	pNode.SetStreamHandler("/gf/0.0.1", streamHandlerFun)
}

//-------------------------------------------------
func readData(pReadWriter *bufio.ReadWriter) {

	for {
		str, err := pReadWriter.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {

			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

//-------------------------------------------------
func writeData(pReadWriter *bufio.ReadWriter) {
	
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = pReadWriter.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = pReadWriter.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

//-------------------------------------------------
func InitShutdownOnSignal(pNode host.Host) {

	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-ch
	fmt.Println("Received signal, shutting down...")

	// shut the node down
	if err := pNode.Close(); err != nil {
		panic(err)
	}
}