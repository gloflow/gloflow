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
	"strings"
	"github.com/fatih/color"
	
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/libp2p/go-libp2p/core/peer"

	discovery_routing "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	discovery_utils "github.com/libp2p/go-libp2p/p2p/discovery/util"
	multiaddr "github.com/multiformats/go-multiaddr"
	routed_host "github.com/libp2p/go-libp2p/p2p/host/routed"
	dht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GFp2pPeerPingFun func() ping.Result
type GFp2pPeerInitFun func(peer.ID) GFp2pPeerPingFun
type GFp2pAddrLst []multiaddr.Multiaddr

var logger = log.Default()

//-------------------------------------------------
func Init(pPortInt int,
	pRuntimeSys *gf_core.RuntimeSys) GFp2pStatusServerCh {

	// CONFIG
	config := getConfig()
	
	// INIT_LIBP2P
	node, dht := InitLibp2p(config, pPortInt, pRuntimeSys)

	// STATUS_SERVER
	statusServerCh := statusServer(node, dht, config, pRuntimeSys)

	return statusServerCh
}

//-------------------------------------------------
func InitLibp2p(pConfig GFp2pConfig,
	pPortInt    int,
	pRuntimeSys *gf_core.RuntimeSys) (host.Host, *dht.IpfsDHT) {

	blue := color.New(color.FgBlue).Add(color.BgWhite).SprintFunc()
	
	ctx := context.Background()

	//----------------
	// KEYS
	privateKey := loadKeys()

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
	identityOption := libp2p.Identity(privateKey)

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

		// libp2p.NATPortMap(),
		// libp2p.ChainOptions(
		// 	libp2p.Transport(ws.New),
		// 	libp2p.Transport(quic.NewTransport),
		// ),
		// libp2p.EnableRelay(),
		// libp2p.EnableAutoRelay(),
		// libp2p.EnableRelayService(),
	)

	if err != nil {
		panic(err)
	}
	// defer node.Close()

	fmt.Printf("node Listen addresses: %s\n", node.Addrs())
	fmt.Printf("node hosts ID is %s\n", blue(node.ID()))

	peerInfo := peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	peerAddrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	fmt.Println("libp2p node address:", peerAddrs[0])
	
	//----------------
	// DHT

	// start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	dht, err := dhtInit(node, ctx)
	if err != nil {
		panic(err)
	}

	//----------------
	// configure our own ping protocol
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)
	
	//-------------------------------------------------
	/*pingInitPeerFun := func(pPeerID peer.ID) GFp2pPeerPingFun {

		pingCh := pingService.Ping(context.Background(), pPeerID)
		pingPeerFun := func() ping.Result {
			res := <-pingCh
			return res
		}

		return GFp2pPeerPingFun(pingPeerFun)
	}*/

	//-------------------------------------------------

	//----------------

	InitStreamHandler(node, pConfig, pRuntimeSys)
	initPeerDiscovery(node, dht, pConfig, ctx, pRuntimeSys)
	

	return node, dht
}

//-------------------------------------------------
func initPeerDiscovery(pNode host.Host,
	pDHT        *dht.IpfsDHT,
	pConfig     GFp2pConfig,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) {

	yellow := color.New(color.FgYellow).SprintFunc()
	green  := color.New(color.FgGreen).SprintFunc()
	greenAndWhiteBg := color.New(color.FgGreen).Add(color.BgWhite).SprintFunc()

	// CONFIG
	bootstrapPeers      := pConfig.BootstrapPeers
	randezvousSymbolStr := pConfig.RendezvousSymbolStr
	protocolIDstr       := pConfig.ProtocolIDstr

	peersNamespaceStr := randezvousSymbolStr
	routedHost := routed_host.Wrap(pNode, pDHT)

	// connect to the bootstrap nodes first, to receive info about other nodes in the network
	var wg sync.WaitGroup

	for _, peerAddr := range bootstrapPeers {

		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := pNode.Connect(pCtx, *peerinfo); err != nil {
				logger.Print(err)
			} else {
				logger.Print("connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	//----------------
	// ANNOUNCING RANDEZVOUS

	logger.Print("announcing peer...")
	routingDiscovery := discovery_routing.NewRoutingDiscovery(pDHT)

	// makes this node announce that it can provide a value for the given key,
	// key being the "randezvous string"
	discovery_utils.Advertise(pCtx, routingDiscovery, peersNamespaceStr)

	logger.Print("peer announced...")

	fmt.Printf("DHT mode: %s\n", pDHT.Mode())

	//----------------
	// look for others peers who have announced (continuously)
	
	go func() {

		logger.Print(fmt.Sprintf("%s:", yellow("searching for other peers")))

		ctx := context.Background()
		
		type GFp2pPeersFailedToDial struct {
			addrLst     []multiaddr.Multiaddr
			attemptsInt int64
		}
		type GFp2pPeersConnected struct {

		}
		peersFailToDialMap := map[peer.ID]*GFp2pPeersFailedToDial{}
		peersConnectedMap  := map[peer.ID]*GFp2pPeersConnected{}
		
		for {

			// FindPeers() - provides us with all the peers that have 
			//               been discovered at the rendezvous-point.
			peersCh, err := routingDiscovery.FindPeers(ctx, peersNamespaceStr)
			if err != nil {
				panic(err)
			}

			for peerAddrInfo := range peersCh {
				
				// skip peer if its this node
				if peerAddrInfo.ID == pNode.ID() {
					continue
				}

				// if discovered peer has no addresses then skip it
				if len(peerAddrInfo.Addrs) == 0 {
					continue
				}

				// peer is already connected to, so skip it
				if _, ok := peersConnectedMap[peerAddrInfo.ID]; ok {
					continue
				}

				logger.Print(fmt.Sprintf("%s:", green("new peer discovered")), peerAddrInfo)

				stream, err := routedHost.NewStream(ctx, peerAddrInfo.ID, protocol.ID(protocolIDstr))
				if err != nil {

					// register all peers that this peer failed to connect with.
					// on the next try it will still be attempted to connect to that peer
					if strings.HasPrefix(fmt.Sprint(err), "failed to dial") {

						if attempt, ok := peersFailToDialMap[peerAddrInfo.ID]; ok {
							attempt.attemptsInt++
						} else {
							logger.Print("peer connection (first attempt) failed:", err)

							peersFailToDialMap[peerAddrInfo.ID] = &GFp2pPeersFailedToDial{
								addrLst:     peerAddrInfo.Addrs,
								attemptsInt: 1,
							}
						}
					}

					fmt.Println(err, strings.HasPrefix(fmt.Sprint(err), "failed to dial"))
					
					continue
				}

				// success
				rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
				go writeDataToStream(rw, pRuntimeSys)
				go readDataFromStream(rw, pRuntimeSys)

				peersConnectedMap[peerAddrInfo.ID] = &GFp2pPeersConnected{}
				logger.Print(fmt.Sprintf("%s:", greenAndWhiteBg("connected to peer")), peerAddrInfo)
			}

			// sleep and then try to discover peers again
			time.Sleep(10 * time.Second)
		}
	}()
}

//-------------------------------------------------
func InitStreamHandler(pNode host.Host,
	pConfig     GFp2pConfig,
	pRuntimeSys *gf_core.RuntimeSys) {

	//-------------------------------------------------
	streamHandlerFun := func(pStream network.Stream) {

		// create a buffer stream for non blocking read and write
		// stream will stay open until you close it (or the other side closes it)
		rw := bufio.NewReadWriter(bufio.NewReader(pStream), bufio.NewWriter(pStream))
	
		go readDataFromStream(rw, pRuntimeSys)
		go writeDataToStream(rw, pRuntimeSys)
	}

	//-------------------------------------------------
	pNode.SetStreamHandler(protocol.ID(pConfig.ProtocolIDstr), streamHandlerFun)
}

//-------------------------------------------------
func readDataFromStream(pReadWriter *bufio.ReadWriter,
	pRuntimeSys *gf_core.RuntimeSys) {

	for {
		lineStr, err := pReadWriter.ReadString('\n')

		if err != nil {
			fmt.Println("error reading line from buffer")
			panic(err)
		}

		if lineStr == "" {
			return
		}
		if lineStr != "\n" {
			fmt.Printf("\x1b[32m%s\x1b[0m> ", lineStr)
		}

		// PARSE
		msg, gfErr := gf_core.ParseJSONfromString(lineStr, pRuntimeSys)
		if gfErr != nil {

		}

		//----------------
		// PROTOCOL
		
		msgMap := msg.(map[string]interface{})
		gfErr = handleIncomingMsg(msgMap, pRuntimeSys)
		if gfErr != nil {

		}

		//----------------

	}
}

//-------------------------------------------------
func writeDataToStream(pReadWriter *bufio.ReadWriter,
	pRuntimeSys *gf_core.RuntimeSys) {
	
	msgMap := map[string]interface{}{}

	for {

		msgBytesLst := gf_core.EncodeJSONfromData(msgMap)
		
		_, err := pReadWriter.Write(msgBytesLst)
		if err != nil {
			fmt.Println("error writing to buffer")
			panic(err)
		}
		err = pReadWriter.Flush()
		if err != nil {
			fmt.Println("error flushing buffer")
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