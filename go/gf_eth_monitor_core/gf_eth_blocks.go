/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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

package gf_eth_monitor_core

import (
	"fmt"
	"strings"
	"context"
	"math/big"
	"github.com/getsentry/sentry-go"
	"github.com/mitchellh/mapstructure"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// BLOCK__INTERNAL - internal representation of the block, with fields
//                   that are not visible to the external public users.
type GF_eth__block__int struct {
	Hash_str          string        `mapstructure:"hash_str"          json:"hash_str"`
	Parent_hash_str   string        `mapstructure:"parent_hash_str"   json:"parent_hash_str"`
	Block_num_int     uint64        `mapstructure:"block_num_int"     json:"block_num_int"`
	Gas_used_int      uint64        `mapstructure:"gas_used_int"      json:"gas_used_int"`
	Gas_limit_int     uint64        `mapstructure:"gas_limit_int"     json:"gas_limit_int"`
	Coinbase_addr_str string        `mapstructure:"coinbase_addr_str" json:"coinbase_addr_str"`
	Txs_lst           []*GF_eth__tx `mapstructure:"txs_lst"           json:"txs_lst"`
	Time              uint64        `mapstructure:"time_int"          json:"time_int"`
	Block             string        `mapstructure:"block"             json:"block"` // *eth_types.Block `json:"block"`
}

//-------------------------------------------------
func Eth_blocks__get_block_from_workers__pipeline(p_block_int uint64,
	p_get_hosts_fn func() []string,
	
	p_ctx     context.Context,
	p_metrics *GF_metrics,
	p_runtime *GF_runtime) (map[string]*GF_eth__block__int, map[string]*GF_eth__miner__int, *gf_core.Gf_error) {

	worker_inspector__port_int := uint(2000)

	//---------------------
	// GET_WORKER_HOSTS
	span__get_worker_hosts := sentry.StartSpan(p_ctx, "get_worker_hosts")
	span__get_worker_hosts.SetTag("workers_aws_discovery", fmt.Sprint(p_runtime.Config.Workers_aws_discovery_bool))

	var workers_inspectors_hosts_lst []string
	if p_runtime.Config.Workers_aws_discovery_bool {
		workers_inspectors_hosts_lst = p_get_hosts_fn()
	} else {
		workers_inspectors_hosts_str := p_runtime.Config.Workers_hosts_str
		workers_inspectors_hosts_lst = strings.Split(workers_inspectors_hosts_str, ",")
	}



	span__get_worker_hosts.Finish()

	//---------------------
	// GET_BLOCKS__FROM_WORKERS_INSPECTORS__ALL

	span := sentry.StartSpan(p_ctx, "get_blocks__workers_inspectors__all")
	defer span.Finish()

	block_from_workers_map   := map[string]*GF_eth__block__int{}
	gf_errs_from_workers_map := map[string]*gf_core.Gf_error{}

	for _, host_str := range workers_inspectors_hosts_lst {

		ctx := span.Context()

		// GET_BLOCK__FROM_WORKER
		gf_block, gf_err := eth_blocks__worker_inspector__get_block(p_block_int,
			host_str,
			worker_inspector__port_int,
			ctx,
			p_runtime.Runtime_sys)

		if gf_err != nil {
			gf_errs_from_workers_map[host_str] = gf_err
			
			// mark a block coming from this worker_inspector host as nil,
			// and continue processing other hosts. 
			// a particular host may fail to return a particular block for various reasons,
			// it might not have synced to that block. 
			block_from_workers_map[host_str] = nil
			continue
		}


		

		//---------------------


		// TEMPORARY!! - move getting of abis_map out of this function.
		// DB_GET
		abi_type_str := "erc20"
		abis_lst, gf_err := Eth_contract__db__get_abi(abi_type_str, p_ctx, p_metrics, p_runtime)
		if gf_err != nil {
			return nil, nil, gf_err
		}
		abis_map := map[string]*GF_eth__abi{
			"erc20": abis_lst[0],
		}


		//---------------------

		gf_err = eth_tx__enrich_from_block(gf_block,
			abis_map,
			ctx,
			p_metrics,
			p_runtime)
		if gf_err != nil {
			gf_errs_from_workers_map[host_str] = gf_err
			
			// mark a block coming from this worker_inspector host as nil,
			// and continue processing other hosts. 
			// a particular host may fail to return a particular block for various reasons,
			// it might not have synced to that block. 
			block_from_workers_map[host_str] = nil
			continue
		}

		

		block_from_workers_map[host_str] = gf_block
	}

	span.Finish()

	//---------------------
	// GET_MINERS - that own this address, potentially multiple records for the same address

	// get coinbase address from the block comming from the first worker_inspector
	var block_miner_addr_hex_str string
	for _, gf_block := range block_from_workers_map {
		
		// if worker failed to return a block, it will be set to nil, so go to the 
		// next one from which a coinbase could be acquired.
		if gf_block != nil {
			block_miner_addr_hex_str = gf_block.Coinbase_addr_str
			break
		}
	}

	miners_map, gf_err := Eth_miners__db__get_info(block_miner_addr_hex_str,
		p_metrics,
		p_ctx,
		p_runtime)
	if gf_err != nil {
		return nil, nil, gf_err
	}

	//---------------------

	return block_from_workers_map, miners_map, nil
}

//-------------------------------------------------
func eth_blocks__worker_inspector__get_block(p_block_int uint64,
	p_host_str    string,
	p_port_int    uint,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_eth__block__int, *gf_core.Gf_error) {

	url_str := fmt.Sprintf("http://%s:%d/gfethm_worker_inspect/v1/blocks?b=%d", p_host_str, p_port_int, p_block_int)

	//-----------------------
	// SPAN
	span_name_str    := fmt.Sprintf("worker_inspector__get_block:%s", p_host_str)
	span__get_blocks := sentry.StartSpan(p_ctx, span_name_str)
	
	// adding tracing ID as a header, to allow for distributed tracing, correlating transactions
	// across services.
	sentry_trace_id_str := span__get_blocks.ToSentryTrace()
	headers_map         := map[string]string{"sentry-trace": sentry_trace_id_str,}
		
	// GF_RPC_CLIENT
	data_map, gf_err := gf_rpc_lib.Client__request(url_str, headers_map, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	span__get_blocks.Finish()

	//-----------------------

	block_map := data_map["block_map"].(map[string]interface{})


	// DECODE_TO_STRUCT
	var gf_block GF_eth__block__int
	err := mapstructure.Decode(block_map, &gf_block)
	if err != nil {
		gf_err := gf_core.Error__create("failed to load response block_map into a GF_eth__block__int struct",
			"mapstruct__decode",
			map[string]interface{}{
				"url_str":   url_str,
				"block_map": block_map,
			},
			err, "gf_eth_monitor_core", p_runtime_sys)
		return nil, gf_err
	}

	return &gf_block, nil
}

//-------------------------------------------------
// GET_BLOCK
func Eth_blocks__get_block__pipeline(p_block_num_int uint64,
	p_eth_rpc_client *ethclient.Client,
	p_ctx            context.Context,
	p_py_plugins     *GF_py_plugins,
	p_runtime_sys    *gf_core.Runtime_sys) (*GF_eth__block__int, *gf_core.Gf_error) {

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
		error_defs_map := Error__get_defs()
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

		error_defs_map := Error__get_defs()
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

	txs_lst := []*GF_eth__tx{}
	for i, tx := range block.Transactions() {

		tx_index_int := uint(i)
		gf_tx, gf_err := Eth_tx__load(tx,
			tx_index_int,
			block.Hash(),
			span__get_txs.Context(),
			p_eth_rpc_client,
			p_py_plugins,
			p_runtime_sys)
		if gf_err != nil {

			// if its an error continue processing the next transaction,
			// and store nil for this one.
			// return nil, gf_err
		}

		txs_lst = append(txs_lst, gf_tx)
	}

	span__get_txs.Finish()

	//------------------
	gf_block := &GF_eth__block__int{
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