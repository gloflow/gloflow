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

package main

import (
	// "os"
	// "fmt"
	// "context"
	"net/http"
	// multiaddr "github.com/multiformats/go-multiaddr"
	// "github.com/libp2p/go-libp2p/core/peer"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_p2p"
	
)

//-------------------------------------------------
func main() {


	portP2Pint := 0 // start on a random port

	// RUNTIME
	logFun, _ := gf_core.InitLogs()
	runtimeSys := &gf_core.RuntimeSys{
		Service_name_str: "gf_p2p_tester",
		LogFun:          logFun,
	}


	
	

	
	go func() {
		_, _ = gf_p2p.Init(portP2Pint, runtimeSys)
	}()

	/*if len(os.Args) > 1 {

		// if a remote peer has been passed on the command line, connect to it
		// and send it 5 ping messages
		peerAddr, err := multiaddr.NewMultiaddr(os.Args[1])
		if err != nil {
			panic(err)
		}

		peer, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			panic(err)
		}


		// connect to peer
		if err := node.Connect(context.Background(), *peer); err != nil {
			panic(err)
		}



		pingPeerFun := pingInitPeerFun(peer.ID)

		fmt.Println("sending 5 ping messages to", peerAddr)
		for i := 0; i < 5; i++ {
			res := pingPeerFun()
			fmt.Println("pinged", peerAddr, "in", res.RTT)
		}

	} else {
		defer gf_p2p.InitShutdownOnSignal(node)
	}*/


	//---------------------
	// HTTP_SERVICE
	gfHTTPmux := http.NewServeMux()

	portStatusServiceInt := 3000
	initService(portStatusServiceInt, gfHTTPmux, runtimeSys)

	//---------------------
	

}