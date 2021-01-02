/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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
	"fmt"
	"math/big"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/ethclient"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/davecgh/go-spew/spew"
)


//-------------------------------------------------
func eth_rpc__get_block(p_block_int int64, p_runtime *GF_runtime) *eth_types.Block {



	ctx := context.Background()
	block, err := p_runtime.Eth_rpc_client.BlockByNumber(ctx, big.NewInt(p_block_int))
	if err != nil {



		fmt.Println(err)
		panic(1)
	}




	spew.Dump(block)

	return block

}


//-------------------------------------------------
func eth_rpc__init() *ethclient.Client {

	geth_host_str := "54.147.190.195" // "127.0.0.1"
	geth_port_int := 8545

	url_str := fmt.Sprintf("http://%s:%d", geth_host_str, geth_port_int)
	client, err := ethclient.Dial(url_str)
    if err != nil {
		log.Fatal(err)

		log.WithFields(log.Fields{
			"url_str":   url_str,
			"geth_host": geth_host_str,
			"port":      geth_port_int,
			"err":       err}).Fatal("failed to connect json-rpc connect to Eth node")
		panic(1)
    }

	log.WithFields(log.Fields{"host": geth_host_str, "port": geth_port_int}).Info("Connected to Ethereum node")
	

	return client
}