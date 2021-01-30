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
	"strings"
	log "github.com/sirupsen/logrus"
	"github.com/getsentry/sentry-go"
	"github.com/ethereum/go-ethereum/ethclient"
	// eth_types "github.com/ethereum/go-ethereum/core/types"
	// eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_monitor_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// GET_BLOCK
func Eth_rpc__get_block__pipeline(p_block_num_int uint64,
	p_eth_rpc_client *ethclient.Client,
	p_ctx            context.Context,
	p_py_plugins     *gf_eth_monitor_core.GF_py_plugins,
	p_runtime_sys    *gf_core.Runtime_sys) (*gf_eth_monitor_core.GF_eth__block__int, *gf_core.Gf_error) {

	//------------------
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

	span__get_header := sentry.StartSpan(p_ctx, "eth_rpc__get_header")
	defer span__get_header.Finish() // in case a panic happens before the main .Finish() for this span

	header, err := p_eth_rpc_client.HeaderByNumber(span__get_header.Context(), new(big.Int).SetUint64(p_block_num_int))
	if err != nil {
		error_defs_map := gf_eth_monitor_core.Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get block Header by number, from eth json-rpc API",
			"eth_rpc__get_header",
			map[string]interface{}{"block_num": p_block_num_int,},
			err, "gf_eth_monitor_lib", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}
	fmt.Println(header)

	span__get_header.Finish()

	//------------------
	// BLOCK_BY_NUMBER

	/*
	type Block

    func NewBlock(header *Header, txs []*Transaction, uncles []*Header, receipts []*Receipt, hasher Hasher) *Block
    func NewBlockWithHeader(header *Header) *Block
    func (b *Block) Bloom() Bloom
    func (b *Block) Body() *Body
    func (b *Block) Coinbase() common.Address
    func (b *Block) DecodeRLP(s *rlp.Stream) error
    func (b *Block) DeprecatedTd() *big.Int
    func (b *Block) Difficulty() *big.Int
    func (b *Block) EncodeRLP(w io.Writer) error
    func (b *Block) Extra() []byte
    func (b *Block) GasLimit() uint64
    func (b *Block) GasUsed() uint64
    func (b *Block) Hash() common.Hash
    func (b *Block) Header() *Header
    func (b *Block) MixDigest() common.Hash
    func (b *Block) Nonce() uint64
    func (b *Block) Number() *big.Int
    func (b *Block) NumberU64() uint64
    func (b *Block) ParentHash() common.Hash
    func (b *Block) ReceiptHash() common.Hash
    func (b *Block) Root() common.Hash
    func (b *Block) SanityCheck() error
    func (b *Block) Size() common.StorageSize
    func (b *Block) Time() uint64
    func (b *Block) Transaction(hash common.Hash) *Transaction
    func (b *Block) Transactions() Transactions
    func (b *Block) TxHash() common.Hash
    func (b *Block) UncleHash() common.Hash
    func (b *Block) Uncles() []*Header
    func (b *Block) WithBody(transactions []*Transaction, uncles []*Header) *Block
    func (b *Block) WithSeal(header *Header) *Block
	*/
	span__get_block := sentry.StartSpan(p_ctx, "eth_rpc__get_block")
	defer span__get_block.Finish() // in case a panic happens before the main .Finish() for this span

	block, err := p_eth_rpc_client.BlockByNumber(span__get_block.Context(), new(big.Int).SetUint64(p_block_num_int))
	if err != nil {

		error_defs_map := gf_eth_monitor_core.Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get block by number, from eth json-rpc API",
			"eth_rpc__get_block",
			map[string]interface{}{"block_num": p_block_num_int,},
			err, "gf_eth_monitor_lib", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}

	span__get_block.Finish()

	//------------------
	// GET_TRANSACTIONS
	span__get_txs := sentry.StartSpan(p_ctx, "eth_rpc__get_txs")
	defer span__get_txs.Finish() // in case a panic happens before the main .Finish() for this span

	txs_lst := []*gf_eth_monitor_core.GF_eth__tx{}
	for i, tx := range block.Transactions() {

		tx_index_int := uint(i)
		gf_tx, gf_err := gf_eth_monitor_core.Eth_tx__get(tx,
			tx_index_int,
			block.Hash(),
			span__get_txs.Context(),
			p_eth_rpc_client,
			p_py_plugins,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		txs_lst = append(txs_lst, gf_tx)
	}

	span__get_txs.Finish()

	//------------------
	gf_block := &gf_eth_monitor_core.GF_eth__block__int{
		Hash_str:          block.Hash().Hex(),
		Parent_hash_str:   block.ParentHash().Hex(),
		Block_num_int:     block.Number().Uint64(),
		Gas_used_int:      block.GasUsed(),
		Gas_limit_int:     block.GasLimit(),
		Coinbase_addr_str: strings.ToLower(block.Coinbase().Hex()),
		Txs_lst:           txs_lst,
		Time:              block.Time(),
		Block:             spew.Sdump(block),
	}

	return gf_block, nil
}

//-------------------------------------------------
// INIT
func Eth_rpc__init(p_host_str string,
	p_geth_port_int int,
	p_runtime_sys   *gf_core.Runtime_sys) (*ethclient.Client, *gf_core.Gf_error) {

	

	url_str := fmt.Sprintf("http://%s:%d", p_host_str, p_geth_port_int)

	client, err := ethclient.Dial(url_str)
    if err != nil {
		log.Fatal(err)

		log.WithFields(log.Fields{
			"url_str":   url_str,
			"geth_host": p_host_str,
			"port":      p_geth_port_int,
			"err":       err}).Fatal("failed to connect json-rpc connect to Eth node")
		
			
		error_defs_map := gf_eth_monitor_core.Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to connect to Eth rpc-json API in gf_eth_monitor",
			"eth_rpc__dial",
			map[string]interface{}{"host": p_host_str,},
			err, "gf_eth_monitor_lib", error_defs_map, p_runtime_sys)
		return nil, gf_err
    }

	log.WithFields(log.Fields{"host": p_host_str, "port": p_geth_port_int}).Info("Connected to Ethereum node")
	

	return client, nil
}