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
	"flag"
	"strings"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	multiaddr "github.com/multiformats/go-multiaddr"
)

//-------------------------------------------------
type GFp2pConfig struct {
	RendezvousSymbolStr string
	BootstrapPeers      GFp2pAddrLst
	ListenAddresses     GFp2pAddrLst
	ProtocolIDstr       string
}

//-------------------------------------------------
func getConfig() GFp2pConfig {

	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}
	config.RendezvousSymbolStr = "gloflow_testnet"
	config.ProtocolIDstr       = "/gf/general/0.0.1"

	return config
}

//-------------------------------------------------
func ParseFlags() (GFp2pConfig, error) {
	config := GFp2pConfig{}

	flag.Var(&config.BootstrapPeers, "peer", "Adds a peer multiaddress to the bootstrap list")
	flag.Var(&config.ListenAddresses, "listen", "Adds a multiaddress to the listen list")
	flag.Parse()

	if len(config.BootstrapPeers) == 0 {
		config.BootstrapPeers = dht.DefaultBootstrapPeers
	}

	return config, nil
}

//-------------------------------------------------
func (al *GFp2pAddrLst) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

//-------------------------------------------------
func (al *GFp2pAddrLst) Set(value string) error {
	addr, err := multiaddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}

//-------------------------------------------------
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