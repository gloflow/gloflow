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
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GFp2pPeerInfo struct {
	IDstr         string
	MultiaddrsLst []string
}
type GFp2pStatus struct {
	RendezvousSymbolStr string
	ProtocolIDstr       string

	BootstrapPeers []GFp2pPeerInfo
}

type GFp2pStatusServerCh chan GFp2pGetStatusMsg
type GFp2pGetStatusMsg struct {
	responseCh chan GFp2pStatus
}

//-------------------------------------------------
func statusServer(pConfig GFp2pConfig,
	pRuntimeSys *gf_core.RuntimeSys) GFp2pStatusServerCh {
	
	statusMngrCh := make(chan GFp2pGetStatusMsg, 10)
	go func() {
		for {
			select {
			case getStatusMsg := <-statusMngrCh:
				status := getStatus(pConfig)
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
func getStatus(pConfig GFp2pConfig) *GFp2pStatus {

	bootstrapPeers      := pConfig.BootstrapPeers
	bootstrapPeersSerialized := serializePeersInfo(bootstrapPeers)

	status := &GFp2pStatus{
		RendezvousSymbolStr: pConfig.RendezvousSymbolStr,
		ProtocolIDstr:       pConfig.ProtocolIDstr,
		BootstrapPeers:      bootstrapPeersSerialized,
	}
	return status
}

//-------------------------------------------------
func serializePeersInfo(pPeersInfoLst GFp2pAddrLst) []GFp2pPeerInfo {

	peersLst := []GFp2pPeerInfo{}
	for _, peerAddr := range pPeersInfoLst {

		peerInfo, _  := peer.AddrInfoFromP2pAddr(peerAddr)
		peerID       := peerInfo.ID
		peerAddrsLst := peerInfo.Addrs

		peerAddrsSerialized := []string{}
		for _, a := range peerAddrsLst {
			peerAddrsSerialized = append(peerAddrsSerialized, a.String())
		}
		gfPeerInfo := GFp2pPeerInfo{
			IDstr:         string(peerID),
			MultiaddrsLst: peerAddrsSerialized,
		}

		peersLst = append(peersLst, gfPeerInfo)
	}
	return peersLst
}


