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
	"context"
	"github.com/libp2p/go-libp2p/core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/gloflow/gloflow/go/gf_core"

	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GFp2pStatus struct {
	RendezvousSymbolStr string
	ProtocolIDstr       string

	// PEERS
	BootstrapPeers []GFp2pPeerInfo
	PeersNumberInt int
	PeersIDsLst    []string

	// DHT
	DHTmodeInt int
	DHTmodeStr string
	DHTtestKeyStr string
	DHTtestValStr string

	// SELF
	SelfPeerIDstr          string
	SelfAddrsSerializedLst []string
}

type GFp2pStatusServerCh chan GFp2pGetStatusMsg
type GFp2pGetStatusMsg struct {
	responseCh chan GFp2pStatus
}

//-------------------------------------------------
func statusServer(pNode host.Host,
	pDHT        *dht.IpfsDHT,
	pConfig     GFp2pConfig,
	pRuntimeSys *gf_core.RuntimeSys) GFp2pStatusServerCh {
	
	statusMngrCh := make(chan GFp2pGetStatusMsg, 10)
	go func() {

		ctx := context.Background()

		for {
			select {
			case getStatusMsg := <-statusMngrCh:

				status := getStatus(pNode,
					pDHT,
					pConfig,
					ctx)

				getStatusMsg.responseCh <- *status
			}
		}
	}()
	return statusMngrCh
}

//-------------------------------------------------
func GetStatusFromServer(pStatusServerCh GFp2pStatusServerCh) GFp2pStatus {

	responseCh := make(chan GFp2pStatus)
	pStatusServerCh <- GFp2pGetStatusMsg{
		responseCh: responseCh,
	}

	status := <-responseCh
	return status
}

//-------------------------------------------------
func getStatus(pNode host.Host,
	pDHT    *dht.IpfsDHT,
	pConfig GFp2pConfig,
	pCtx    context.Context) *GFp2pStatus {

	//----------------
	// PEERS
	bootstrapPeers           := pConfig.BootstrapPeers
	bootstrapPeersSerialized := serializePeersInfo(bootstrapPeers)

	// get the peerstore of the current p2p host
	peers := pNode.Peerstore().Peers()
	peersNumberInt := len(peers)
	
	peerstorePeerIDsLst := []string{}
	for _, peerID := range peers {
		peerstorePeerIDsLst = append(peerstorePeerIDsLst, string(peerID))
	}

	//----------------
	// DHT

	// dht mode
	dhtModeInt := int(pDHT.Mode())
	dhtModeStr := "not_set"
	switch dhtModeInt {
	case int(dht.ModeClient):
		dhtModeStr = "client"
	case int(dht.ModeServer):
		dhtModeStr = "server"
	}

	// dh routing-table diversity stats
	// :[]peerdiversity.CplDiversityStats
	stats := pDHT.GetRoutingTableDiversityStats()
	spew.Dump(stats)

	//----------------
	// DHT READ/WRITE

	dhtTestKeyStr := "/gf/0.0.1/key_test_1"

	err := dhtPut(dhtTestKeyStr, "test_val_1", pDHT, pCtx)
	if err != nil {

	}

	dhtTestValStr, err := dhtGet("/gf/0.0.1/key_test_1", pDHT, pCtx)
	if err != nil {
		dhtTestValStr = fmt.Sprint(err)
	}

	//----------------
	// DHT FIND_SELF

	selfPeerIDstr, selfAddrsSerializedLst := dhtFindSelf(pNode, pDHT)

	//----------------

	status := &GFp2pStatus{
		RendezvousSymbolStr: pConfig.RendezvousSymbolStr,
		ProtocolIDstr:       pConfig.ProtocolIDstr,
		BootstrapPeers:      bootstrapPeersSerialized,

		// PEERS
		PeersNumberInt: peersNumberInt,
		PeersIDsLst:    peerstorePeerIDsLst,

		// DHT
		DHTmodeInt: dhtModeInt,
		DHTmodeStr: dhtModeStr,
		DHTtestKeyStr: dhtTestKeyStr,
		DHTtestValStr: dhtTestValStr,

		// SELF
		SelfPeerIDstr:          selfPeerIDstr,
		SelfAddrsSerializedLst: selfAddrsSerializedLst,
	}
	return status
}