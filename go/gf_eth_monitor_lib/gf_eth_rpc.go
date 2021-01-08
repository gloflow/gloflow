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

package gf_eth_monitor_lib

import (
	"fmt"
	"math/big"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/ethclient"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// GET_BLOCK
func Eth_rpc__get_block(p_block_num_int int64,
	p_eth_rpc_client *ethclient.Client,
	p_ctx            context.Context,
	p_runtime_sys    *gf_core.Runtime_sys) (*GF_eth__block, *gf_core.Gf_error) {

	
	/*
	type Header struct {
		ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
		UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
		Coinbase    common.Address `json:"miner"            gencodec:"required"`
		Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
		TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
		ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
		Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
		Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
		Number      *big.Int       `json:"number"           gencodec:"required"`
		GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
		GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
		Time        uint64         `json:"timestamp"        gencodec:"required"`
		Extra       []byte         `json:"extraData"        gencodec:"required"`
		MixDigest   common.Hash    `json:"mixHash"`
		Nonce       BlockNonce     `json:"nonce"`
	}

	// Time - the unix timestamp for when the block was collated

	*/
	header, err := p_eth_rpc_client.HeaderByNumber(p_ctx, big.NewInt(p_block_num_int))
	if err != nil {
		panic(err)
	}
	fmt.Println(header)
	


	block, err := p_eth_rpc_client.BlockByNumber(p_ctx, big.NewInt(p_block_num_int))
	if err != nil {

		error_defs_map := error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to read apps__info from YAML file in Cmonkeyd",
			"file_read_error",
			map[string]interface{}{"block_num": p_block_num_int,},
			err, "gf_eth_monitor_lib", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}


	txs_lst := []*GF_eth__tx{}
	for _, t := range block.Transactions() {

		/*
		type Transaction

		func NewContractCreation(nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction
		func NewTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction
		func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error)
		func (tx *Transaction) AsMessage(s Signer) (Message, error)
		func (tx *Transaction) ChainId() *big.Int
		func (tx *Transaction) CheckNonce() bool
		func (tx *Transaction) Cost() *big.Int
		func (tx *Transaction) Data() []byte
		func (tx *Transaction) DecodeRLP(s *rlp.Stream) error
		func (tx *Transaction) EncodeRLP(w io.Writer) error
		func (tx *Transaction) Gas() uint64
		func (tx *Transaction) GasPrice() *big.Int
		func (tx *Transaction) GasPriceCmp(other *Transaction) int
		func (tx *Transaction) GasPriceIntCmp(other *big.Int) int
		func (tx *Transaction) Hash() common.Hash
		func (tx *Transaction) MarshalJSON() ([]byte, error)
		func (tx *Transaction) Nonce() uint64
		func (tx *Transaction) Protected() bool
		func (tx *Transaction) RawSignatureValues() (v, r, s *big.Int)
		func (tx *Transaction) Size() common.StorageSize
		func (tx *Transaction) To() *common.Address
		func (tx *Transaction) UnmarshalJSON(input []byte) error
		func (tx *Transaction) Value() *big.Int
		func (tx *Transaction) WithSignature(signer Signer, sig []byte) (*Transaction, error)
		*/



		t.Gas()
		t.GasPrice()
		// :common.Hash - represents the 32 byte Keccak256 hash of arbitrary data.
		//                in this case its a hash of signed transaction data.
		tx_hash := t.Hash()
		t.Value()
		
		
		address := t.To()
		if address == nil {
			// its a contract creation transaction
		}

		gf_tx, gf_err := Eth_rpc__get_tx(tx_hash, p_ctx, p_eth_rpc_client, p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		fmt.Println(gf_tx)



		txs_lst = append(txs_lst, gf_tx)
	}



	gf_block := &GF_eth__block{
		Block_num_int:     block.Number().Uint64(),
		Gas_used_int:      block.GasUsed(),
		Gas_limit_int:     block.GasLimit(),
		Coinbase_addr_str: block.Coinbase().Hex(),
		Txs_lst:           txs_lst,
		Block:             spew.Sdump(block),
	}

	

	return gf_block, nil
}

//-------------------------------------------------
func Eth_rpc__get_tx(p_tx_hash eth_common.Hash,
	p_ctx            context.Context,
	p_eth_rpc_client *ethclient.Client,
	p_runtime_sys    *gf_core.Runtime_sys) (*GF_eth__tx, *gf_core.Gf_error) {

	tx_hash_hex_str := p_tx_hash.Hex()
	tx_receipt, err := p_eth_rpc_client.TransactionReceipt(p_ctx, p_tx_hash)
	if err != nil {

		error_defs_map := error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get transaction recepit in gf_eth_monitor",
			"eth_rpc__get_tx_receipt",
			map[string]interface{}{"tx_hash_hex": tx_hash_hex_str,},
			err, "gf_eth_monitor_lib", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}

	// GET_LOGS
	logs, gf_err := Eth_rpc__get_tx_logs(tx_receipt,
		p_ctx,
		p_eth_rpc_client,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}



	gas_used_int := tx_receipt.GasUsed
	tx := &GF_eth__tx{
		Gas_used_int: gas_used_int,
		Logs:         logs,
	}

	return tx, nil
}

//-------------------------------------------------
func Eth_rpc__get_tx_logs(p_tx_receipt *eth_types.Receipt,
	p_ctx            context.Context,
	p_eth_rpc_client *ethclient.Client,
	p_runtime_sys    *gf_core.Runtime_sys) ([]*eth_types.Log, *gf_core.Gf_error) {


	
	/*
	type Log struct {
		// Consensus fields:
		// address of the contract that generated the event
		Address common.Address `json:"address" gencodec:"required"`
		// list of topics provided by the contract.
		Topics []common.Hash `json:"topics" gencodec:"required"`

		// supplied by the contract, usually ABI-encoded
		Data []byte `json:"data" gencodec:"required"`

		// Derived fields. These fields are filled in by the node
		// but not secured by consensus.
		// block in which the transaction was included
		BlockNumber uint64 `json:"blockNumber"`

		// hash of the transaction
		TxHash common.Hash `json:"transactionHash" gencodec:"required"`
		// index of the transaction in the block
		TxIndex uint `json:"transactionIndex"`
		// hash of the block in which the transaction was included
		BlockHash common.Hash `json:"blockHash"`
		// index of the log in the block
		Index uint `json:"logIndex"`

		// The Removed field is true if this log was reverted due to a chain reorganisation.
		// You must pay attention to this field if you receive logs through a filter query.
		Removed bool `json:"removed"`
	}*/
	logs_lst := []*eth_types.Log{}
	for _, l := range p_tx_receipt.Logs {

		// data__byte_lst := l.Data
		// fmt.Println(data__byte_lst)

		logs_lst = append(logs_lst, l)
	}

	return logs_lst, nil
}

//-------------------------------------------------
// INIT
func Eth_rpc__init(p_host_str string,
	p_runtime_sys *gf_core.Runtime_sys) (*ethclient.Client, *gf_core.Gf_error) {

	geth_port_int := 8545

	url_str := fmt.Sprintf("http://%s:%d", p_host_str, geth_port_int)

	client, err := ethclient.Dial(url_str)
    if err != nil {
		log.Fatal(err)

		log.WithFields(log.Fields{
			"url_str":   url_str,
			"geth_host": p_host_str,
			"port":      geth_port_int,
			"err":       err}).Fatal("failed to connect json-rpc connect to Eth node")
		
			
		error_defs_map := error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to connect to Eth rpc-json API in gf_eth_monitor",
			"eth_rpc__dial",
			map[string]interface{}{"host": p_host_str,},
			err, "gf_eth_monitor_lib", error_defs_map, p_runtime_sys)
		return nil, gf_err
    }

	log.WithFields(log.Fields{"host": p_host_str, "port": geth_port_int}).Info("Connected to Ethereum node")
	

	return client, nil
}