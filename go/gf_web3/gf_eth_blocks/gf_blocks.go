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

package gf_eth_blocks

import (
	"fmt"
	"strings"
	"context"
	"time"
	"math/big"
	"sync"
	"github.com/getsentry/sentry-go"
	"github.com/mitchellh/mapstructure"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_web3/gf_eth_core"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_web3/gf_eth_tx"
	"github.com/gloflow/gloflow-web3-monitor/go/gf_web3/gf_eth_contract"
	// "go.mongodb.org/mongo-driver/bson"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// BLOCK__INTERNAL - internal representation of the block, with fields
//                   that are not visible to the external public users.
type GF_eth__block__int struct {
	DB_id                 string    `mapstructure:"db_id"                 json:"db_id"                 bson:"_id"`
	Creation_time__unix_f float64   `mapstructure:"creation_time__unix_f" json:"creation_time__unix_f" bson:"creation_time__unix_f"`

	Hash_str          string        `mapstructure:"hash_str"          json:"hash_str"          bson:"hash_str"`
	Parent_hash_str   string        `mapstructure:"parent_hash_str"   json:"parent_hash_str"   bson:"parent_hash_str"`
	Block_num_uint    uint64        `mapstructure:"block_num_int"     json:"block_num_int"     bson:"block_num_uint"`
	Gas_used_uint     uint64        `mapstructure:"gas_used_uint"     json:"gas_used_uint"     bson:"gas_used_uint"`
	Gas_limit_uint    uint64        `mapstructure:"gas_limit_uint"    json:"gas_limit_uint"    bson:"gas_limit_uint"`
	Coinbase_addr_str string        `mapstructure:"coinbase_addr_str" json:"coinbase_addr_str" bson:"coinbase_addr_str"`

	// txs_lst       - not stored in DB. blocks are stored in their own table/collection, separate from TX's themselves (which are in their own).
	// tx_hashes_lst - these are stored in the DB, and used to lookup TX's that are in their own collection.
	Txs_lst        []*gf_eth_tx.GF_eth__tx `mapstructure:"txs_lst"           json:"txs_lst"            bson:"-"`
	Txs_hashes_lst []string                `mapstructure:"txs_hashes_lst"    json:"txs_hashes_lst"     bson:"txs_hashes_lst"`

	Time_uint uint64 `mapstructure:"time_int" json:"time_uint"`
	
	// Block string `mapstructure:"block" json:"block"` // *eth_types.Block `json:"block"`
}

//-------------------------------------------------
// metrics that are continuously calculated

func Init_continuous_metrics(p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) {
	go func() {
		for {
			//---------------------
			// GET_BLOCKS_COUNTS
			blocks_count_int, gf_err := DB__get_count(p_metrics, p_runtime)
			if gf_err != nil {
				time.Sleep(60 * time.Second) // SLEEP
				continue
			}
			p_metrics.Block__db_count__gauge.Set(float64(blocks_count_int))

			//---------------------
			time.Sleep(60 * time.Second) // SLEEP
		}
	}()
}

//-------------------------------------------------
/*// BLOCKS__GET_AND_PERSIST_BULK

func Eth_blocks__get_and_persist_bulk__pipeline(p_block_start_uint uint64,
	p_block_end_uint      uint64,
	// p_get_worker_hosts_fn func(context.Context, *GF_runtime) []string,
	// p_abis_defs_map       map[string]*GF_eth__abi,
	// p_ctx                 context.Context,
	p_metrics             *gf_eth_core.GF_metrics,
	p_runtime             *gf_eth_core.GF_runtime) []*gf_core.GF_error {

	gf_eth_indexer.Client__index_block_range(p_block_start_uint,
		p_block_end_uint,
		p_runtime.Indexer_cmds_ch)

	return gf_errs_lst
}*/

//-------------------------------------------------
func Index__pipeline(p_block_uint uint64,
	p_get_worker_hosts_fn func(context.Context, *gf_eth_core.GF_runtime) []string,
	p_abis_defs_map       map[string]*gf_eth_contract.GF_eth__abi,
	p_ctx                 context.Context,
	p_metrics             *gf_eth_core.GF_metrics,
	p_runtime             *gf_eth_core.GF_runtime) (uint64, *gf_core.GF_error) {


	//---------------------
	// GET_BLOCK_FROM_WORKER
	// gets the same block from all the workers that it gets, and the resulting maps
	// are key-ed by worker_host.
	block_from_workers_map, _, gf_err := Get_from_workers__pipeline(p_block_uint,
		p_get_worker_hosts_fn,
		p_abis_defs_map,
		p_ctx,
		p_metrics,
		p_runtime)

	if gf_err != nil {
		return 0, gf_err
	}

	// IMPORTANT!! - for now just get the block from the first worker_host,
	//               regardless of how many workers are registered.
	var gf_block *GF_eth__block__int
	for worker_host_str := range block_from_workers_map {
		gf_block = block_from_workers_map[worker_host_str]
		break
	}

	//---------------------
	// DB_WRITE_BULK__BLOCK

	gf_err = DB__write_bulk([]*GF_eth__block__int{gf_block,},
		p_ctx,
		p_metrics,
		p_runtime)
	if gf_err != nil {
		return 0, gf_err
	}

	// METRICS
	if p_metrics != nil {
		p_metrics.Block__indexed_num__counter.Inc()
	}

	//---------------------
	
	txs_num_int := uint64(len(gf_block.Txs_lst))

	// some blocks (especially early ones) dont have any transactions in them
	if txs_num_int > 0 {

		//---------------------
		// DB_WRITE_BULK__TXS

		gf_err = gf_eth_tx.DB__write_bulk(gf_block.Txs_lst,
			p_ctx,
			p_metrics,
			p_runtime)
		if gf_err != nil {
			return 0, gf_err
		}

		// METRICS
		if p_metrics != nil {
			for _, _ = range gf_block.Txs_lst {
				p_metrics.Tx__indexed_num__counter.Inc()
			}
		}
		
		//---------------------
		// TRACES
		tx_hashes_lst := []string{}
		for _, tx := range gf_block.Txs_lst {
			tx_hashes_lst = append(tx_hashes_lst, tx.Hash_str)
		}

		worker_inspector_host_port_str := p_get_worker_hosts_fn(p_ctx, p_runtime)[0]

		// DB_WRITE
		gf_err, _ = gf_eth_tx.Trace__get_and_persist_bulk(tx_hashes_lst,
			worker_inspector_host_port_str,
			p_ctx,
			p_metrics,
			p_runtime)
		if gf_err != nil {
			return 0, gf_err
		}

		//---------------------
	}
	return txs_num_int, nil
}

//-------------------------------------------------
func Get_from_workers__pipeline(p_block_uint uint64,
	p_get_worker_hosts_fn func(context.Context, *gf_eth_core.GF_runtime) []string,
	p_abis_defs_map       map[string]*gf_eth_contract.GF_eth__abi,
	p_ctx                 context.Context,
	p_metrics             *gf_eth_core.GF_metrics,
	p_runtime             *gf_eth_core.GF_runtime) (map[string]*GF_eth__block__int, map[string]*gf_eth_core.GF_eth__miner__int, *gf_core.GF_error) {


	//---------------------
	// GET_WORKER_HOSTS

	// span__get_worker_hosts := sentry.StartSpan(p_ctx, "get_worker_hosts")
	// span__get_worker_hosts.SetTag("workers_aws_discovery", fmt.Sprint(p_runtime.Config.Workers_aws_discovery_bool))
	// var workers_inspectors_hosts_lst []string
	// if p_runtime.Config.Workers_aws_discovery_bool {
	// 	workers_inspectors_hosts_lst = p_get_worker_hosts_fn()
	// } else {
	// 	workers_inspectors_hosts_str := p_runtime.Config.Workers_hosts_str
	// 	workers_inspectors_hosts_lst = strings.Split(workers_inspectors_hosts_str, ",")
	// }
	// span__get_worker_hosts.Finish()

	workers_inspectors_hosts_lst := p_get_worker_hosts_fn(p_ctx, p_runtime)

	//---------------------
	// GET_BLOCKS__FROM_WORKERS_INSPECTORS__ALL

	span := sentry.StartSpan(p_ctx, "get_blocks__workers_inspectors__all")
	defer span.Finish()

	block_from_workers_map   := map[string]*GF_eth__block__int{}
	gf_errs_from_workers_map := map[string]*gf_core.GF_error{}

	for _, host_port_str := range workers_inspectors_hosts_lst {

		ctx := span.Context()

		// GET_BLOCK__FROM_WORKER
		gf_block, gf_err := Get_block__from_worker_inspector(p_block_uint,
			host_port_str,
			ctx,
			p_runtime.RuntimeSys)

		if gf_err != nil {
			gf_errs_from_workers_map[host_port_str] = gf_err
			
			// mark a block coming from this worker_inspector host as nil,
			// and continue processing other hosts. 
			// a particular host may fail to return a particular block for various reasons,
			// it might not have synced to that block. 
			block_from_workers_map[host_port_str] = nil
			continue
		}


		
		/*//---------------------
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

		//---------------------*/

		
		gf_err = gf_eth_tx.Enrich_from_block(gf_block.Txs_lst,
			p_abis_defs_map,
			ctx,
			p_metrics,
			p_runtime)
		if gf_err != nil {
			gf_errs_from_workers_map[host_port_str] = gf_err
			
			// mark a block coming from this worker_inspector host as nil,
			// and continue processing other hosts. 
			// a particular host may fail to return a particular block for various reasons,
			// it might not have synced to that block. 
			block_from_workers_map[host_port_str] = nil
			continue
		}

		

		block_from_workers_map[host_port_str] = gf_block
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

	miners_map, gf_err := gf_eth_core.Eth_miners__db__get_info(block_miner_addr_hex_str,
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
// GET_BLOCK__FROM_WORKER_INSPECTOR
func Get_block__from_worker_inspector(p_block_uint uint64,
	p_host_port_str string,
	p_ctx           context.Context,
	p_runtime_sys   *gf_core.Runtime_sys) (*GF_eth__block__int, *gf_core.GF_error) {



	url_str := fmt.Sprintf("http://%s/gfethm_worker_inspect/v1/blocks?b=%d",
		p_host_port_str,
		p_block_uint)

	//-----------------------
	// SPAN
	span_name_str    := fmt.Sprintf("worker_inspector__get_block:%s", p_host_port_str)
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
// GET_BLOCK__PIPELINE
func Get__pipeline(p_block_num_uint uint64,
	pEthRPCclient *ethclient.Client,
	p_ctx         context.Context,
	p_py_plugins  *gf_eth_core.GF_py_plugins,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_eth__block__int, *gf_core.GF_error) {

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

	header, err := pEthRPCclient.HeaderByNumber(span__get_header.Context(), new(big.Int).SetUint64(p_block_num_uint))
	if err != nil {
		error_defs_map := gf_eth_core.Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get block Header by number, from eth json-rpc API",
			"eth_rpc__get_header",
			map[string]interface{}{"block_num": p_block_num_uint,},
			err, "gf_eth_monitor_lib", error_defs_map, 1, p_runtime_sys)
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

	block, err := pEthRPCclient.BlockByNumber(span__get_block.Context(), new(big.Int).SetUint64(p_block_num_uint))
	if err != nil {

		error_defs_map := gf_eth_core.Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get block by number, from eth json-rpc API",
			"eth_rpc__get_block",
			map[string]interface{}{"block_num": p_block_num_uint,},
			err, "gf_eth_monitor_lib", error_defs_map, 1, p_runtime_sys)
		return nil, gf_err
	}

	span__get_block.Finish()

	//------------------
	// GET_TRANSACTIONS
	
	span__get_txs := sentry.StartSpan(p_ctx, "eth_rpc__get_txs")
	defer span__get_txs.Finish() // in case a panic happens before the main .Finish() for this span
	
	txs_lst, _ := Get__txs_pipeline(block,
		pEthRPCclient,
		span__get_txs.Context(),
		p_py_plugins,
		p_runtime_sys)

	span__get_txs.Finish()
	
	//------------------


	txs_hashes_lst := []string{}
	for _, tx := range txs_lst {
		txs_hashes_lst = append(txs_hashes_lst, tx.Hash_str)
	}

	block_hash_hex_str := block.Hash().Hex()

	gf_block := &GF_eth__block__int{
		Hash_str:          block_hash_hex_str,
		Parent_hash_str:   block.ParentHash().Hex(),
		Block_num_uint:    block.Number().Uint64(),
		Gas_used_uint:     block.GasUsed(),
		Gas_limit_uint:    block.GasLimit(),
		Coinbase_addr_str: strings.ToLower(block.Coinbase().Hex()),

		Txs_lst:        txs_lst,
		Txs_hashes_lst: txs_hashes_lst,

		Time_uint: block.Time(),
		// Block: spew.Sdump(block),
	}

	//------------------
	// IMPORTANT!! - its critical for the hashing of TX struct to get signature be done before the
	//               creation_time__unix_f attribute is set, since that always changes and would affect the hash.
	//               
	db_id_hex_str     := gf_core.Hash_val_sha256(gf_block)
	creationTimeUNIXf := float64(time.Now().UnixNano()) / 1_000_000_000.0
	
	/*obj_id_str, err := primitive.ObjectIDFromHex(db_id_hex_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to decode Block struct hash hex signature to create Mongodb ObjectID",
			"decode_hex",
			map[string]interface{}{"block_hash_hex_str": block_hash_hex_str, },
			err, "gf_eth_monitor_core", p_runtime_sys)
		return nil, gf_err
	}*/
	gf_block.DB_id                 = db_id_hex_str // obj_id_str
	gf_block.Creation_time__unix_f = creationTimeUNIXf

	//------------------

	return gf_block, nil
}

//-------------------------------------------------
func Get__txs_pipeline(p_block *types.Block,
	pEthRPCclient *ethclient.Client,
	p_ctx         context.Context,
	p_py_plugins  *gf_eth_core.GF_py_plugins,
	p_runtime_sys *gf_core.Runtime_sys) ([]*gf_eth_tx.GF_eth__tx, []*gf_core.GF_error) {



	



	var wg sync.WaitGroup


	txs_lst     := make([]*gf_eth_tx.GF_eth__tx, len(p_block.Transactions()))
	gf_errs_lst := make([]*gf_core.GF_error,     len(p_block.Transactions()))
	for i, tx := range p_block.Transactions() {
		tx_index_int := uint(i)

		wg.Add(1)

		// IMPORTANT!! - load each transaction in parallel in its own goroutine.
		//               for blocks with lots of transactions it may timeout
		//               while processing all transactions (it might take >60s) sequentially.
		go func(p_tx_index_int uint, p_tx *types.Transaction) {
			
			defer wg.Done()

			gf_tx, gf_err := gf_eth_tx.Load(p_tx,
				p_tx_index_int,
				p_block.Hash(),
				p_block.NumberU64(), // p_block_num_uint,
				p_ctx,
				pEthRPCclient,
				p_py_plugins,
				p_runtime_sys)

			if gf_err != nil {

				// if its an error continue processing the next transaction,
				// and store nil for the current transaction.
				// important that txs_lst has a length of the true number of txs.
				gf_errs_lst[p_tx_index_int] = gf_err
			}

			txs_lst[p_tx_index_int] = gf_tx
		}(tx_index_int, tx)
	}

	
	wg.Wait()

	return txs_lst, gf_errs_lst



}