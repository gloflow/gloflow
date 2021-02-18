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
	"math"
	"math/big"
	"context"
	"encoding/base64"
	// "encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/ethereum/go-ethereum/ethclient"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type GF_eth__tx struct {
	Hash_str       string                `json:"hash_str"         bson:"hash_str"`
	Index_int      uint64                `json:"index_int"        bson:"index_int"` // position of the transaction in the block
	From_addr_str  string                `json:"from_addr_str"    bson:"from_addr_str"`
	To_addr_str    string                `json:"to_addr_str"      bson:"to_addr_str"`
	Value_eth_f    float64               `json:"value_eth_f"      bson:"value_eth_f"`
	Data_bytes_lst []byte                `json:"-"                bson:"-"`
	Data_b64_str   string                `json:"data_b64_str"     bson:"data_b64_str"`
	Gas_used_int   uint64                `json:"gas_used_int"     bson:"gas_used_int"`
	Gas_price_int  uint64                `json:"gas_price_int"    bson:"gas_price_int"`
	Nonce_int      uint64                `json:"nonce_int"        bson:"nonce_int"`
	Size_f         float64               `json:"size_f"           bson:"size_f"`
	Cost_int       uint64                `json:"cost_int"         bson:"cost_int"`
	Contract_new   *GF_eth__contract_new `json:"contract_new_map" bson:"contract_new_map"`
	Logs_lst       []*GF_eth__log        `json:"logs_lst"         bson:"logs_lst"`
}

type GF_eth__tx_trace struct {
	Gas_used_int       uint64                     `json:"gas_used_int"`
	Value_returned_str string                     `json:"value_returned_str"`	
	Failed_bool        bool                       `json:"failed_bool"`
	Opcodes_lst        []*GF_eth__tx_trace_opcode `json:"opcodes_lst"`
}

type GF_eth__tx_trace_opcode struct {
	Op_str            string   `json:"op_str"`
	Pc_int            uint     `json:"pc_int"`            // program counter
	Gas_cost_int      uint     `json:"gas_cost_int"`
	Gas_remaining_int uint64   `json:"gas_remaining_int"` // decreasing count of how much gas is left before this Op executes
	Stack_lst         []string          `json:"stack_lst"`
	Memory_lst        []string          `json:"memory_lst"`
	Storage_map       map[string]string `json:"storage_map"`
}

// eth_types.Log
type GF_eth__log struct {
	Address_str  string   `json:"address_str"  bson:"address_str"`  // address of the contract that generated the log
	Topics_lst   []string `json:"topics_lst"   bson:"topics_lst"`   // list of topics provided by the contract
	Data_hex_str string   `json:"data_hex_str" bson:"data_hex_str"` // supplied by contract, usually ABI-encoded
}

//-------------------------------------------------
func Eth_tx__plot_trace() {



}


//-------------------------------------------------
func Eth_tx__get_trace__from_worker_inspector(p_tx_hash_str string,
	p_host_str    string,
	p_port_int    uint,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_eth__tx_trace, *gf_core.Gf_error) {

	url_str := fmt.Sprintf("http://%s:%d/gfethm_worker_inspect/v1/tx/trace?tx=%s",
		p_host_str,
		p_port_int,
		p_tx_hash_str)

	//-----------------------
	// SPAN
	span_name_str    := fmt.Sprintf("worker_inspector__get_tx_trace:%s", p_host_str)
	span__get_tx_trace := sentry.StartSpan(p_ctx, span_name_str)
	
	// adding tracing ID as a header, to allow for distributed tracing, correlating transactions
	// across services.
	sentry_trace_id_str := span__get_tx_trace.ToSentryTrace()
	headers_map         := map[string]string{"sentry-trace": sentry_trace_id_str,}
		
	// GF_RPC_CLIENT
	data_map, gf_err := gf_rpc_lib.Client__request(url_str, headers_map, p_ctx, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	span__get_tx_trace.Finish()

	//-----------------------

	trace_map  := data_map["trace_map"].(map[string]interface{})
	result_map := trace_map["result"].(map[string]interface{})

	gf_opcodes_lst := []*GF_eth__tx_trace_opcode{}
	for _, op := range result_map["structLogs"].([]interface{}) {

		op_map := op.(map[string]interface{})

		stack_lst := []string{}
		for _, s := range op_map["stack"].([]interface{}) {
			stack_lst = append(stack_lst, s.(string))
		}

		memory_lst := []string{}
		for _, s := range op_map["memory"].([]interface{}) {
			memory_lst = append(memory_lst, s.(string))
		}

		storage_map := map[string]string{}
		for k, v := range op_map["storage"].(map[string]interface{}) {
			storage_map[k] = v.(string)
		}

		gf_opcode := &GF_eth__tx_trace_opcode{
			Op_str:            strings.TrimSpace(op_map["op"].(string)),
			Pc_int:            uint(op_map["pc"].(float64)),
			Gas_cost_int:      uint(op_map["gasCost"].(float64)),
			Gas_remaining_int: uint64(op_map["gas"].(float64)),
			Stack_lst:         stack_lst,
			Memory_lst:        memory_lst,
			Storage_map:       storage_map,
		}

		gf_opcodes_lst = append(gf_opcodes_lst, gf_opcode)
	}

	gf_tx_trace := &GF_eth__tx_trace{
		Gas_used_int:       uint64(result_map["gas"].(float64)),
		Value_returned_str: result_map["returnValue"].(string),
		Failed_bool:        result_map["failed"].(bool),
		Opcodes_lst:        gf_opcodes_lst,
	}

	return gf_tx_trace, nil
}

//-------------------------------------------------
func Eth_tx__get_trace(p_tx_hash_str string,
	p_eth_rpc_host_str string,
	p_runtime_sys      *gf_core.Runtime_sys) (map[string]interface{}, *gf_core.Gf_error) {

	// IMPORTANT!! - transaction tracing is not exposed as a function in the golang ehtclient, as explained
	//               by the authors, because it is a geth specific function and ethclient is suppose to be a 
	//               generic implementation of a client for the standard ethereum RPC API.
	input_str := fmt.Sprintf(`{
		"id":     1,
		"method": "debug_traceTransaction",
		"params": ["%s", {
			"disableStack":   false,
			"disableMemory":  false,
			"disableStorage": false
		}]
	}`, p_tx_hash_str)

	
	output_map, gf_err := Eth_rpc__call(input_str,
		p_eth_rpc_host_str,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	return output_map, nil
}

//-------------------------------------------------
func eth_tx__enrich_from_block(p_gf_block *GF_eth__block__int,
	p_abis_map map[string]*GF_eth__abi,
	p_ctx      context.Context,
	p_metrics  *GF_metrics,
	p_runtime  *GF_runtime) *gf_core.Gf_error {
	


	for _, tx := range p_gf_block.Txs_lst {


		// IMPORTANT!! - if worker_inspector encounters an error while loading
		//               the transaction, or whichever mechanisms the transaction comes from,
		//               it is marked as nil. so skip.
		if tx == nil {
			continue
		}

		//------------------
		// NEW_CONTRACT - check if its a new_contract transaction.
		//                if it is then load/decode contract information
		if tx.Contract_new != nil {

			// TEMPORARY!! - we just assume the new contract has a erc20 ABI, for testing purposes.
			//               generalize a way to specify an ABI to use to decode new contracts.
			gf_abi := p_abis_map["erc20"]
			gf_err := Eth_contract__enrich(gf_abi, p_ctx, p_metrics, p_runtime)
			if gf_err != nil {
				return gf_err
			}
		}

		//------------------
		// LOGS - check if the transaction has any logs
		if len(tx.Logs_lst) > 0 {


		
		}

		//------------------
	}
	return nil
}

//-------------------------------------------------
func Eth_tx__load(p_tx *eth_types.Transaction,
	p_tx_index_int   uint,
	p_block_hash     eth_common.Hash,
	p_ctx            context.Context,
	p_eth_rpc_client *ethclient.Client,
	p_py_plugins     *GF_py_plugins,
	p_runtime_sys    *gf_core.Runtime_sys) (*GF_eth__tx, *gf_core.Gf_error) {

	
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

		error_defs_map := Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get transaction recepit via json-rpc  in gf_eth_monitor",
			"eth_rpc__get_tx_receipt",
			map[string]interface{}{"tx_hash_hex": tx_hash_hex_str,},
			err, "gf_eth_monitor_core", error_defs_map, p_runtime_sys)
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
		error_defs_map := Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get transaction via json-rpc in gf_eth_monitor",
			"eth_rpc__get_tx",
			map[string]interface{}{"tx_hash_hex": tx_hash_hex_str,},
			err, "gf_eth_monitor_core", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}

	span__get_tx.Finish()




	

	//------------------
	// GET_LOGS
	span__parse_tx_logs := sentry.StartSpan(p_ctx, "eth_rpc__parse_tx_logs")
	defer span__parse_tx_logs.Finish() // in case a panic happens before the main .Finish() for this span

	logs, gf_err := Eth_tx__get_logs(tx_receipt,
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
		error_defs_map := Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to get transaction via json-rpc in gf_eth_monitor",
			"eth_rpc__get_tx_sender",
			map[string]interface{}{"tx_hash_hex": tx_hash_hex_str,},
			err, "gf_eth_monitor_core", error_defs_map, p_runtime_sys)
		return nil, gf_err
	}

	sender_addr_str := strings.ToLower(sender_addr.Hex())

	//------------------
	// TX_TO
	var to_str        string
	var contract__new *GF_eth__contract_new

	// NEW CONTRACT - if To() is nil and the ContractAddress is set,
	//                then its a new contract creation transaction.
	if tx.To() == nil && tx_receipt.ContractAddress.Hex() != "" {

		
		to_str = "new_contract"
		new_contract_addr_str := strings.ToLower(tx_receipt.ContractAddress.Hex())




		//------------------
		// NEW_CONTRACT_CODE

		block_num_int := tx_receipt.BlockNumber.Uint64()
		contract, gf_err := Eth_contract__get_via_rpc(new_contract_addr_str,
			block_num_int,
			p_ctx,
			p_eth_rpc_client,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		contract__new = contract

		//------------------



		// PY_PLUGIN - get info on a new contract

		gf_err = py__run_plugin__get_contract_info(new_contract_addr_str,
			p_py_plugins,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}


		//------------------

	// SELF_TRANSACTION
	// CHECK!! - make sure this is a sufficient condition to mark this transaction as "self".
	} else if tx.To() == nil {


		// SELF - these are "cancelation transactions", meant to cancel TX's in mempool,
		//        that for some reason (usually too little gas) are not getting
		//        mined into a block. these cancelation tx's have the same nonce as the
		//        tx's that are not getting mined, and have the same address in 
		//        From and To ("self"). these transaction have a gas price set to an appropriate
		//        level to get mined into a block, and have no action effectively canceling
		//        the old tx.
		to_str = "self"

	// REGULAR
	} else {
		to_str = strings.ToLower(tx.To().Hex())
	}



	//------------------
	// IMPORTANT!! - tx value is in Wei units, to convert to Eth we divide by 10^18.
	//               its important to do the conversion in Go process space, instead of transmit
	//               Wei value via JSON; JSON/JS support floats of 32bit size only, and would lose precision
	//               in transmission.
	tx_value_f        := new(big.Float).SetInt(tx.Value())
	tx_value_eth_f, _ := new(big.Float).Quo(tx_value_f, big.NewFloat(math.Pow(10, 18))).Float64()

	// TX_DATA
	data_bytes_lst := tx.Data()
	data_b64_str   := base64.StdEncoding.EncodeToString(data_bytes_lst) // base64

	gas_used_int := tx_receipt.GasUsed

	gf_tx := &GF_eth__tx{
		Hash_str:       tx_receipt.TxHash.Hex(),
		Index_int:      uint64(tx_receipt.TransactionIndex),
		From_addr_str:  sender_addr_str,
		To_addr_str:    to_str,
		Value_eth_f:    tx_value_eth_f, // tx.Value().Uint64(),

		// DATA
		Data_bytes_lst: data_bytes_lst,
		Data_b64_str:   data_b64_str,

		Gas_used_int:  gas_used_int,
		Gas_price_int: tx.GasPrice().Uint64(),
		Nonce_int:     tx.Nonce(),
		Size_f:        float64(tx.Size()),
		Cost_int:      tx.Cost().Uint64(),
		Logs_lst:      logs,

		Contract_new:  contract__new,
	}

	return gf_tx, nil
}

//-------------------------------------------------
func Eth_tx__enrich_logs(p_tx_logs []*GF_eth__log,
	p_abis_map map[string]*GF_eth__abi,
	p_ctx      context.Context,
	p_metrics  *GF_metrics,
	p_runtime  *GF_runtime) ([]map[string]interface{}, *gf_core.Gf_error) {
	

	



	gf_abi := p_abis_map["erc20"]
	abi, gf_err := Eth_contract__get_abi(gf_abi,
		p_ctx,
		p_metrics,
		p_runtime)
	if gf_err != nil {
		return nil, gf_err
	}






	spew.Dump(abi)
	


	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAA----------------")

	decoded_logs_lst := []map[string]interface{}{}
	for _, l := range p_tx_logs {
		


		fmt.Println(">>>>>>>>>>>>>>>>")


		fmt.Println(len(l.Topics_lst))

		

		event_1_map := map[string]interface{}{}
		event_2_map := map[string]interface{}{}
		event_3_map := map[string]interface{}{}

		abi.UnpackIntoMap(event_1_map, "Transfer", eth_common.HexToHash(l.Topics_lst[0]).Bytes())
		abi.UnpackIntoMap(event_2_map, "Transfer", eth_common.HexToHash(l.Topics_lst[1]).Bytes())
		abi.UnpackIntoMap(event_3_map, "Transfer", eth_common.HexToHash(l.Topics_lst[2]).Bytes())


		hash          := eth_common.HexToHash(l.Data_hex_str)
		log_bytes_lst := hash.Bytes()

		event_map := map[string]interface{}{}

		// UnpackIntoMap - unpacks a log into the provided map[string]interface{}.
		err := abi.UnpackIntoMap(event_map, "Transfer", log_bytes_lst)
		if err != nil {
			error_defs_map := Error__get_defs()
			gf_err := gf_core.Error__create_with_defs("failed to decode a Tx Log",
				"eth_tx_log__decode",
				map[string]interface{}{
					"address_str":  l.Address_str,
					"data_hex_str": l.Data_hex_str,
				},
				err, "gf_eth_monitor_core", error_defs_map, p_runtime.Runtime_sys)
			return nil, gf_err
		}
		
		spew.Dump(log_bytes_lst)
		spew.Dump(event_map)
		spew.Dump(event_1_map)
		spew.Dump(event_2_map)
		spew.Dump(event_3_map)


		fmt.Println(eth_common.BigToHash(event_1_map["value"].(*big.Int)).Hex())
		fmt.Println(eth_common.BigToHash(event_2_map["value"].(*big.Int)).Hex())
		fmt.Println(eth_common.BigToHash(event_3_map["value"].(*big.Int)).Hex())



		decoded_logs_lst = append(decoded_logs_lst, event_map)
	}




	return decoded_logs_lst, nil
}

//-------------------------------------------------
func Eth_tx__get_logs(p_tx_receipt *eth_types.Receipt,
	p_ctx            context.Context,
	p_eth_rpc_client *ethclient.Client,
	p_runtime_sys    *gf_core.Runtime_sys) ([]*GF_eth__log, *gf_core.Gf_error) {


	
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
	}
	*/


	logs_lst := []*GF_eth__log{} // eth_types.Log{}
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
		topics_lst := []string{}

		// l.Topics() - list of topics provided by the contract
		for _, topic := range l.Topics {
			topics_lst = append(topics_lst, topic.Hex())
		}

		// DATA - while topics are searchable, data is not.
		//        including data is a lot cheaper than including topics.
		//        supplied by the contract, usually ABI-encoded.
		data_bytes_lst := l.Data
		data_hex_str   := eth_common.BytesToHash(data_bytes_lst).Hex() // base64.StdEncoding.EncodeToString(data_bytes_lst) // base64

		tx_log := &GF_eth__log{
			Address_str:  l.Address.Hex(),
			Topics_lst:   topics_lst,
			Data_hex_str: data_hex_str,
		}

		logs_lst = append(logs_lst, tx_log)
	}

	return logs_lst, nil
}