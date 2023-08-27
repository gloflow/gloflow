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

package gf_eth_tx

import (
	"fmt"
	"context"
	"time"
	"strings"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_stats/gf_stats_lib"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GF_eth__tx_trace struct {
	DB_id                 string                     `mapstructure:"db_id"                 json:"db_id"                 bson:"_id"`
	Creation_time__unix_f float64                    `mapstructure:"creation_time__unix_f" json:"creation_time__unix_f" bson:"creation_time__unix_f"`
	Tx_hash_str           string                     `mapstructure:"tx_hash_str"           json:"tx_hash_str"        bson:"tx_hash_str"`
	Gas_used_uint         uint64                     `mapstructure:"gas_used_uint"         json:"gas_used_uint"      bson:"gas_used_uint"`
	Value_returned_str    string                     `mapstructure:"value_returned_str"    json:"value_returned_str" bson:"value_returned_str"`	
	Failed_bool           bool                       `mapstructure:"failed_bool"           json:"failed_bool"        bson:"failed_bool"`
	Opcodes_lst           []*GF_eth__tx_trace_opcode `mapstructure:"opcodes_lst"           json:"opcodes_lst"        bson:"opcodes_lst"`
}

type GF_eth__tx_trace_opcode struct {
	Op_str            string            `mapstructure:"op_str"             json:"op_str"             bson:"op_str"`
	Pc_int            uint              `mapstructure:"pc_int"             json:"pc_int"             bson:"pc_int"`             // program counter
	Gas_cost_uint      uint             `mapstructure:"gas_cost_uint"      json:"gas_cost_uint"      bson:"gas_cost_uint"`
	Gas_remaining_uint uint64           `mapstructure:"gas_remaining_uint" json:"gas_remaining_uint" bson:"gas_remaining_uint"` // decreasing count of how much gas is left before this Op executes
	Stack_lst         []string          `mapstructure:"stack_lst"          json:"stack_lst"          bson:"stack_lst"`
	Memory_lst        []string          `mapstructure:"memory_lst"         json:"memory_lst"         bson:"memory_lst"`
	Storage_map       map[string]string `mapstructure:"storage_map"        json:"storage_map"        bson:"storage_map"`
}

//-------------------------------------------------

func Trace__get_and_persist_bulk(p_tx_hashes_lst []string,
	p_worker_inspector_host_port_str string,
	pCtx                             context.Context,
	p_metrics                        *gf_eth_core.GF_metrics,
	p_runtime                        *gf_eth_core.GF_runtime) (*gf_core.GFerror, []*gf_core.GFerror) {


	// IMPORTANT!! - these are "secondary" errors, that are not the primary one that causes the
	//               function to fail and return. these secondary errors are from getting
	//               traces from worker_inspector, and are considered recoverable and the iteration
	//               over all TX's continues.
	gf_errs__get_tx_trace_lst := []*gf_core.GFerror{}

	txs_traces_lst := []*GF_eth__tx_trace{}
	for _, tx_hash_str := range p_tx_hashes_lst {

		
		


		// GET_TRACE - WORKER_INSPECTOR
		gf_tx_trace, gfErr := Trace__get_from_worker_inspector(tx_hash_str,
			p_worker_inspector_host_port_str,
			pCtx,
			p_runtime.RuntimeSys)

		if gfErr != nil {
			gf_errs__get_tx_trace_lst = append(gf_errs__get_tx_trace_lst, gfErr)
			txs_traces_lst            = append(txs_traces_lst, nil)
			continue
		} else {
			gf_errs__get_tx_trace_lst = append(gf_errs__get_tx_trace_lst, nil)
		}

		txs_traces_lst = append(txs_traces_lst, gf_tx_trace)
	}

	// DB_WRITE_BULK
	gfErr := DBmongoTraceWriteBulk(txs_traces_lst,
		pCtx,
		p_metrics,
		p_runtime)
	if gfErr != nil {
		return gfErr, gf_errs__get_tx_trace_lst
	}

	return nil, gf_errs__get_tx_trace_lst
}

//-------------------------------------------------

func DBmongoTraceWriteBulk(p_txs_traces_lst []*GF_eth__tx_trace,
	pCtx      context.Context,
	p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) *gf_core.GFerror {

	coll_name_str := "gf_eth_txs_traces"

	filterDocsByFieldsLst := []map[string]string{}
	recordsLst            := []interface{}{}
	txsHashesLst          := []string{}
	
	for _, tx := range p_txs_traces_lst {

		filterDocsByFieldsLst = append(filterDocsByFieldsLst,
			map[string]string{"_id": tx.DB_id,})

		recordsLst   = append(recordsLst, interface{}(tx))
		txsHashesLst = append(txsHashesLst, tx.Tx_hash_str)
	}

	_, gfErr := gf_core.MongoUpsertBulk(filterDocsByFieldsLst, recordsLst,
		coll_name_str,
		map[string]interface{}{
			"txs_hashes_lst":     txsHashesLst,
			"caller_err_msg_str": "failed to bulk insert Eth txs_traces into DB",
		},
		pCtx, p_runtime.RuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//-------------------------------------------------

func Trace__plot(p_tx_id_hex_str string,
	p_get_hosts_fn func(context.Context, *gf_eth_core.GF_runtime) []string,
	pCtx           context.Context,
	p_py_plugins   *gf_eth_core.GF_py_plugins,
	p_metrics      *gf_eth_core.GF_metrics,
	p_runtime      *gf_eth_core.GF_runtime) (string, *gf_core.GFerror) {



	//-----------------------
	// WORKER_INSPECTOR__HOST_PORT
	host_port_str      := p_get_hosts_fn(pCtx, p_runtime)[0]
	start_time__unix_f := float64(time.Now().UnixNano()) / 1000000000.0

	// GET_TRACE
	gf_tx_trace, gfErr := Trace__get_from_worker_inspector(p_tx_id_hex_str,
		host_port_str,
		pCtx,
		p_runtime.RuntimeSys)

	end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

	// METRICS
	if p_metrics != nil {
		delta_time__unix_f := end_time__unix_f - start_time__unix_f
		p_metrics.Tx_trace__worker_inspector_durration__gauge.Set(delta_time__unix_f)
	}

	if gfErr != nil {
		return "", gfErr
	}

	//-----------------------
	// PY_PLUGIN__PLOT

	start_time__unix_f = float64(time.Now().UnixNano()) / 1000000000.0
	plot_svg_str, gfErr := py__run_plugin__plot_tx_trace(p_tx_id_hex_str,
		gf_tx_trace,
		p_py_plugins,
		p_runtime.RuntimeSys)
	end_time__unix_f = float64(time.Now().UnixNano())/1000000000.0
	
	// METRICS
	if p_metrics != nil {
		delta_time__unix_f := end_time__unix_f - start_time__unix_f
		p_metrics.Tx_trace__py_plugin__plot_durration__gauge.Set(delta_time__unix_f)
	}

	if gfErr != nil {
		return "", gfErr
	}

	//-----------------------

	return plot_svg_str, nil
}

//-------------------------------------------------
// GET_FROM_WORKER_INSPECTOR

func Trace__get_from_worker_inspector(p_tx_hash_str string,
	p_host_port_str string,
	pCtx            context.Context,
	p_RuntimeSys   *gf_core.RuntimeSys) (*GF_eth__tx_trace, *gf_core.GFerror) {

	url_str := fmt.Sprintf("http://%s/gfethm_worker_inspect/v1/tx/trace?tx=%s",
	p_host_port_str,
		p_tx_hash_str)

	//-----------------------
	// SPAN
	span_name_str      := fmt.Sprintf("worker_inspector__get_tx_trace:%s", p_host_port_str)
	span__get_tx_trace := sentry.StartSpan(pCtx, span_name_str)
	
	// adding tracing ID as a header, to allow for distributed tracing, correlating transactions
	// across services.
	sentry_trace_id_str := span__get_tx_trace.ToSentryTrace()
	headers_map         := map[string]string{"sentry-trace": sentry_trace_id_str,}
		
	// GF_RPC_CLIENT
	data_map, gfErr := gf_rpc_lib.ClientRequest(url_str, headers_map, pCtx, p_RuntimeSys)
	if gfErr != nil {
		return nil, gfErr
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
		

		// fmt.Println("------------------")
		// spew.Dump(op_map)

		storage_map := map[string]string{}
		if _, ok := op_map["storage"]; ok {
			for k, v := range op_map["storage"].(map[string]interface{}) {
				storage_map[k] = v.(string)
			}
		}
		
		gf_opcode := &GF_eth__tx_trace_opcode{
			Op_str:             strings.TrimSpace(op_map["op"].(string)),
			Pc_int:             uint(op_map["pc"].(float64)),
			Gas_cost_uint:      uint(op_map["gasCost"].(float64)),
			Gas_remaining_uint: uint64(op_map["gas"].(float64)),
			Stack_lst:          stack_lst,
			Memory_lst:         memory_lst,
			Storage_map:        storage_map,
		}
		
		gf_opcodes_lst = append(gf_opcodes_lst, gf_opcode)
	}


	
	
	gf_tx_trace := &GF_eth__tx_trace{
		Tx_hash_str:        p_tx_hash_str,
		Gas_used_uint:      uint64(result_map["gas"].(float64)),
		Value_returned_str: result_map["returnValue"].(string),
		Failed_bool:        result_map["failed"].(bool),
		Opcodes_lst:        gf_opcodes_lst,
	}


	//------------------
	// IMPORTANT!! - its critical for the hashing of TX struct to get signature be done before the
	//               creation_time__unix_f attribute is set, since that always changes and would affect the hash.
	//               
	db_id_hex_str         := gf_core.HashValSha256(gf_tx_trace)
	creation_time__unix_f := float64(time.Now().UnixNano()) / 1_000_000_000.0
	
	/*obj_id_str, err := primitive.ObjectIDFromHex(db_id_hex_str)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to decode Tx_trace struct hash hex signature to create Mongodb ObjectID",
			"decode_hex",
			map[string]interface{}{"tx_hash_str": p_tx_hash_str, },
			err, "gf_eth_monitor_core", p_RuntimeSys)
		return nil, gfErr
	}*/
	gf_tx_trace.DB_id                 = db_id_hex_str // obj_id_str
	gf_tx_trace.Creation_time__unix_f = creation_time__unix_f

	//------------------



	return gf_tx_trace, nil
}

//-------------------------------------------------

func Trace__get(p_tx_hash_str string,
	p_eth_rpc_host_str string,
	pCtx               context.Context,
	p_RuntimeSys       *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {

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

	
	output_map, gfErr := gf_eth_core.Eth_rpc__call(input_str,
		p_eth_rpc_host_str,
		map[string]interface{}{
			"tx_hash_str": p_tx_hash_str,
		},
		pCtx,
		p_RuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	return output_map, nil
}

//-------------------------------------------------
// metrics that are continuously calculated

func Trace__init_continuous_metrics(p_metrics *gf_eth_core.GF_metrics,
	p_runtime *gf_eth_core.GF_runtime) *gf_core.GFerror {
	
	ctx := context.Background()
	coll_name_str := "gf_eth_txs_traces"
	
	//---------------------
	// COLL_EXISTS_CHECK - if collection doesnt exist (yet) in the DB then dont
	//                     begin collection metrics on it (it will cause errors).
	coll_exists_bool, gfErr := gf_core.MongoCollExists(coll_name_str, ctx, p_runtime.RuntimeSys)
	if gfErr != nil {
		return gfErr
	}
	if coll_exists_bool {
		// FIX!! - this shouldnt just exit for users that dont have the needed collection yet.
		//         with enough usage they will trigger a code-path that will create a collection.
		//         so the system should detect that somehow, and yet not do this collection
		//         Mongo__coll_exists() for every iteration.
		return nil
	}
	
	//---------------------

	go func() {
		
		for {	
			//---------------------
			// GET_BLOCKS_COUNTS
			// blocks_count_int, gfErr := Eth_tx_trace__db__get_count(p_metrics, p_runtime)

			db_coll_stats, gfErr := gf_stats_lib.Db_stats__coll("gf_eth_txs_traces", ctx, p_runtime.RuntimeSys)
			blocks_count_int := db_coll_stats.Docs_count_int
			
			if gfErr != nil {
				time.Sleep(60 * time.Second) // SLEEP
				continue
			}
			p_metrics.Block__db_count__gauge.Set(float64(blocks_count_int))

			//---------------------

			time.Sleep(60 * time.Second) // SLEEP
		}
	}()

	return nil
}

/*//-------------------------------------------------
func Eth_tx_trace__db__get_count(p_metrics *GF_metrics,
	p_runtime *GF_runtime) (int64, *gf_core.GFerror) {

	coll_name_str := "gf_eth_txs_traces"
	coll := p_runtime.RuntimeSys.Mongo_db.Collection(coll_name_str)

	ctx := context.Background()
	
	count_int, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {

		// METRICS
		if p_metrics != nil {p_metrics.Errs_num__counter.Inc()}

		gfErr := gf_core.MongoHandleError("failed to DB count Transactions Trace",
			"mongodb_count_error",
			map[string]interface{}{},
			err, "gf_eth_monitor_core", p_runtime.RuntimeSys)
		return 0, gfErr
	}

	return count_int, nil
}*/