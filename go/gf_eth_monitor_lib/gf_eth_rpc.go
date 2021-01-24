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
	"math"
	"math/big"
	"context"
	"strings"
	log "github.com/sirupsen/logrus"
	"github.com/getsentry/sentry-go"
	"github.com/ethereum/go-ethereum/ethclient"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow-ethmonitor/go/gf_eth_monitor_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// GET_BLOCK
func Eth_rpc__get_block__pipeline(p_block_num_int uint64,
	p_eth_rpc_client *ethclient.Client,
	p_ctx            context.Context,
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
		error_defs_map := error__get_defs()
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

		error_defs_map := error__get_defs()
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
		gf_tx, gf_err := Eth_rpc__get_tx(tx,
			tx_index_int,
			block.Hash(),
			span__get_txs.Context(),
			p_eth_rpc_client,
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
func Eth_rpc__get_tx(p_tx *eth_types.Transaction,
	p_tx_index_int   uint,
	p_block_hash     eth_common.Hash,
	p_ctx            context.Context,
	p_eth_rpc_client *ethclient.Client,
	p_runtime_sys    *gf_core.Runtime_sys) (*gf_eth_monitor_core.GF_eth__tx, *gf_core.Gf_error) {

	
	tx_hash         := p_tx.Hash() // :eth_common.Hash
	tx_hash_hex_str := tx_hash.Hex()

	//------------------
	// GET_TX_RECEIPT

	/*
	type Receipt struct {
		// Consensus fields: These fields are defined by the Yellow Paper
		PostState         []byte `json:"root"`
		Status            uint64 `json:"status"`
		CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gencodec:"required"`
		Bloom             Bloom  `json:"logsBloom"         gencodec:"required"`
		Logs              []*Log `json:"logs"              gencodec:"required"`

		// Implementation fields: These fields are added by geth when processing a transaction.
		// They are stored in the chain database.
		TxHash          common.Hash    `json:"transactionHash" gencodec:"required"`
		ContractAddress common.Address `json:"contractAddress"`
		GasUsed         uint64         `json:"gasUsed" gencodec:"required"`

		// Inclusion information: These fields provide information about the inclusion of the
		// transaction corresponding to this receipt.
		BlockHash        common.Hash `json:"blockHash,omitempty"`
		BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
		TransactionIndex uint        `json:"transactionIndex"`
	}
	*/

	span__get_tx_receipt := sentry.StartSpan(p_ctx, "eth_rpc__get_tx_receipt")
	defer span__get_tx_receipt.Finish() // in case a panic happens before the main .Finish() for this span

	tx_receipt, err := p_eth_rpc_client.TransactionReceipt(span__get_tx_receipt.Context(), tx_hash)
	if err != nil {

		error_defs_map := error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get transaction recepit via json-rpc  in gf_eth_monitor",
			"eth_rpc__get_tx_receipt",
			map[string]interface{}{"tx_hash_hex": tx_hash_hex_str,},
			err, "gf_eth_monitor_lib", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}
	span__get_tx_receipt.Finish()



	//------------------
	// GET_TX

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

	/*t.Gas()
	t.GasPrice()
	// :common.Hash - represents the 32 byte Keccak256 hash of arbitrary data.
	//                in this case its a hash of signed transaction data.
	tx_hash := t.Hash()
	t.Value()
	
	
	address := t.To()
	if address == nil {
		// its a contract creation transaction
	}*/


	span__get_tx := sentry.StartSpan(p_ctx, "eth_rpc__get_tx")
	defer span__get_tx.Finish() // in case a panic happens before the main .Finish() for this span

	tx, _, err := p_eth_rpc_client.TransactionByHash(span__get_tx.Context(), tx_hash)
	if err != nil {
		error_defs_map := error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get transaction via json-rpc in gf_eth_monitor",
			"eth_rpc__get_tx",
			map[string]interface{}{"tx_hash_hex": tx_hash_hex_str,},
			err, "gf_eth_monitor_lib", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}

	span__get_tx.Finish()




	

	//------------------
	// GET_LOGS
	span__parse_tx_logs := sentry.StartSpan(p_ctx, "eth_rpc__parse_tx_logs")
	defer span__parse_tx_logs.Finish() // in case a panic happens before the main .Finish() for this span

	logs, gf_err := Eth_rpc__get_tx_logs(tx_receipt,
		span__parse_tx_logs.Context(),
		p_eth_rpc_client,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}
	span__parse_tx_logs.Finish()

	//------------------
	// GET_TX_SENDER

	sender_addr, err := p_eth_rpc_client.TransactionSender(p_ctx, tx, p_block_hash, p_tx_index_int)
	if err != nil {
		error_defs_map := error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get transaction via json-rpc in gf_eth_monitor",
			"eth_rpc__get_tx_sender",
			map[string]interface{}{"tx_hash_hex": tx_hash_hex_str,},
			err, "gf_eth_monitor_lib", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}

	//------------------

	// IMPORTANT!! - tx value is in Wei units, to convert to Eth we divide by 10^18.
	//               its important to do the conversion in Go process space, instead of transmit
	//               Wei value via JSON; JSON/JS support floats of 32bit size only, and would lose precision
	//               in transmission.
	tx_value_f        := new(big.Float).SetInt(tx.Value())
	tx_value_eth_f, _ := new(big.Float).Quo(tx_value_f, big.NewFloat(math.Pow(10, 18))).Float64()

	gas_used_int := tx_receipt.GasUsed
	gf_tx := &gf_eth_monitor_core.GF_eth__tx{
		Hash_str:      tx_receipt.TxHash.Hex(),
		Index_int:     uint64(tx_receipt.TransactionIndex),
		From_addr_str: strings.ToLower(sender_addr.Hex()),
		To_addr_str:   strings.ToLower(tx.To().Hex()),
		Value_eth_f:   tx_value_eth_f, // tx.Value().Uint64(),
		Gas_used_int:  gas_used_int,
		Gas_price_int: tx.GasPrice().Uint64(),
		Nonce_int:     tx.Nonce(),
		Size_f:        float64(tx.Size()),
		
		Cost_int:      tx.Cost().Uint64(),
		Logs:          logs,
	}

	return gf_tx, nil
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

		// COST - The base cost of logging operations is 375 gas. On top of that,
		//        every included topic costs an additional 375 gas.
		//        Finally, each byte of data costs 8 gas.

		// TOPIC - used to describe the event.
		//         can only hold a maximum of 32 bytes of data.
		//         first topic usually consists of the signature (a keccak256 hash)
		//         of the name of the event that occurred, including the types
		//         (uint256, string, etc.) of its parameters.
		//         topics should only reliably be used for data that strongly
		//         narrows down search queries (like addresses)
		topics_lst := []eth_common.Hash{}

		// l.Topics() - list of topics provided by the contract
		for _, topic := range l.Topics {
			topics_lst = append(topics_lst, topic)
		}

		// DATA - while topics are searchable, data is not.
		//        including data is a lot cheaper than including topics.
		//        supplied by the contract, usually ABI-encoded.
		data_bytes_lst := l.Data
		fmt.Println(data_bytes_lst)


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