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
	// datastore "github.com/ipfs/go-datastore"
	// datastore_sync "github.com/ipfs/go-datastore/sync"

	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func dhtInit(pNode host.Host,
	pCtx context.Context) (*dht.IpfsDHT, error) {
	
	optionsLst := []dht.Option{
		// dht.ProtocolPrefix(protocol.ID(pConfig.ProtocolIDstr)),
		// dht.NamespacedValidator("v", blankValidator{}),
		
		// start the node in Server mode
		dht.Mode(dht.ModeServer),
		// dht.Mode(dht.ModeAuto),

		// how many times each k/v pair should be replicated across
		// the peers in the network.
		// dht.OptionReplicationFactor(3),

		// DisableAutoRefresh(),
	}

	// Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
	// dstore := datastore_sync.MutexWrap(datastore.NewMapDatastore())

	/*// https://github.com/libp2p/go-libp2p-kad-dht/blob/master/dht.go
	// NewDHT creates a new DHT object with the given peer as the 'local' host.
	// IpfsDHT's initialized with this function will respond to DHT requests,
	// whereas IpfsDHT's initialized with NewDHTClient will not.
	dht := dht.NewDHT(pCtx, pNode, dstore)*/

	dht, err := dht.New(pCtx, pNode, optionsLst...)
	if err != nil {
		return nil, err
	}

	// someKeyStr := "something"
	// _, err := dht.GetClosestPeers(ctx, someKeyStr)

	// tells the DHT to get into a bootstrapped state satisfying the IpfsRouter interface.
	// in the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	logger.Print("Bootstrapping the DHT")
	if err := dht.Bootstrap(pCtx); err != nil {
		return nil, err
	}

	fmt.Printf("DHT mode: %s\n", dht.Mode())


	return dht, nil
}

//-------------------------------------------------
// GET
func dhtGet(pKeyStr string,
	pDHT *dht.IpfsDHT,
	pCtx context.Context) (string, error) {
	
	val, err := pDHT.GetValue(pCtx, pKeyStr)
	if err != nil {
		return "", err
	}

	valStr := string(val)
	return valStr, nil
}

//-------------------------------------------------
// PUT
func dhtPut(pKeyStr string,
	pValStr string,
	pDHT    *dht.IpfsDHT,
	pCtx    context.Context) error {

	err := pDHT.PutValue(pCtx, "/gf/0.0.1/key_test_1", []byte("key_value1"))
	if err != nil {
		return err
	}
	
	return nil
}

//-------------------------------------------------
func dhtFindSelf(pNode host.Host,
	pDHT *dht.IpfsDHT) (string, []string) {

	// find self
	fmt.Printf("looking for self\n")

	// :peer.AddrInfo
	selfAddrInfo := pDHT.FindLocal(pNode.ID())
	
	selfPeerIDstr, selfAddrsSerializedLst :=  serializeAddrInfo(selfAddrInfo)

	return selfPeerIDstr, selfAddrsSerializedLst
}