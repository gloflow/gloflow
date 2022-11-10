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
)

//-------------------------------------------------
type GFp2pPeerInfo struct {
	IDstr         string
	MultiaddrsLst []string
}

//-------------------------------------------------
func serializePeersInfo(pPeersInfoLst GFp2pAddrLst) []GFp2pPeerInfo {

	peersLst := []GFp2pPeerInfo{}
	for _, peerAddr := range pPeersInfoLst {

		peerAddrInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		peerIDstr, peerAddrsSerializedLst := serializeAddrInfo(*peerAddrInfo)

		gfPeerInfo := GFp2pPeerInfo{
			IDstr:         peerIDstr,
			MultiaddrsLst: peerAddrsSerializedLst,
		}

		peersLst = append(peersLst, gfPeerInfo)
	}
	return peersLst 
}

//-------------------------------------------------
func serializeAddrInfo(pAddrInfo peer.AddrInfo) (string, []string) {

	peerIDstr := string(pAddrInfo.ID)

	addrsSerialized := []string{}
	for _, a := range pAddrInfo.Addrs {
		addrsSerialized = append(addrsSerialized, a.String())
	}

	return peerIDstr, addrsSerialized
}