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
	"context"
	"time"
	"strings"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)

//-------------------------------------------------
type GF_eth__tx_trace struct {
	Tx_hash_str        string                     `json:"tx_hash_str"`
	Gas_used_uint      uint64                     `json:"gas_used_uint"`
	Value_returned_str string                     `json:"value_returned_str"`	
	Failed_bool        bool                       `json:"failed_bool"`
	Opcodes_lst        []*GF_eth__tx_trace_opcode `json:"opcodes_lst"`
}

type GF_eth__tx_trace_opcode struct {
	Op_str            string            `json:"op_str"`
	Pc_int            uint              `json:"pc_int"`             // program counter
	Gas_cost_uint      uint             `json:"gas_cost_uint"`
	Gas_remaining_uint uint64           `json:"gas_remaining_uint"` // decreasing count of how much gas is left before this Op executes
	Stack_lst         []string          `json:"stack_lst"`
	Memory_lst        []string          `json:"memory_lst"`
	Storage_map       map[string]string `json:"storage_map"`
}

//-------------------------------------------------
func Eth_tx_trace__get_and_persist_bulk(p_tx_hashes_lst []string,
	p_worker_inspector_host_port_str string,
	p_ctx                            context.Context,
	p_metrics                        *GF_metrics,
	p_runtime                        *GF_runtime) (*gf_core.Gf_error, []*gf_core.Gf_error) {


	// IMPORTANT!! - these are "secondary" errors, that are not the primary one that causes the
	//               function to fail and return. these secondary errors are from getting
	//               traces from worker_inspector, and are considered recoverable and the iteration
	//               over all TX's continues.
	gf_errs__get_tx_trace_lst := []*gf_core.Gf_error{}

	txs_traces_lst := []*GF_eth__tx_trace{}
	for _, tx_hash_str := range p_tx_hashes_lst {




		// GET_TRACE - WORKER_INSPECTOR
		gf_tx_trace, gf_err := Eth_tx_trace__get_from_worker_inspector(tx_hash_str,
			p_worker_inspector_host_port_str,
			p_ctx,
			p_runtime.Runtime_sys)

		if gf_err != nil {
			gf_errs__get_tx_trace_lst = append(gf_errs__get_tx_trace_lst, gf_err)
			txs_traces_lst            = append(txs_traces_lst, nil)
			continue
		} else {
			gf_errs__get_tx_trace_lst = append(gf_errs__get_tx_trace_lst, nil)
		}

		txs_traces_lst = append(txs_traces_lst, gf_tx_trace)
	}

	// DB_WRITE_BULK
	gf_err := Eth_tx_trace__db__write_bulk(txs_traces_lst,
		p_ctx,
		p_metrics,
		p_runtime)
	if gf_err != nil {
		return gf_err, gf_errs__get_tx_trace_lst
	}

	return nil, gf_errs__get_tx_trace_lst
}

//-------------------------------------------------
func Eth_tx_trace__db__write_bulk(p_txs_traces_lst []*GF_eth__tx_trace,
	p_ctx     context.Context,
	p_metrics *GF_metrics,
	p_runtime *GF_runtime) *gf_core.Gf_error {

	coll_name_str := "gf_eth_txs_traces"

	records_lst    := []interface{}{}
	txs_hashes_lst := []string{}
	for _, tx := range p_txs_traces_lst {
		records_lst    = append(records_lst, interface{}(tx))
		txs_hashes_lst = append(txs_hashes_lst, tx.Tx_hash_str)
	}

	gf_err := gf_core.Mongo__insert_bulk(records_lst,
		coll_name_str,
		map[string]interface{}{
			"txs_hashes_lst":     txs_hashes_lst,
			"caller_err_msg_str": "failed to bulk insert Eth txs_traces into DB",
		},
		p_ctx, p_runtime.Runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	return nil
}

//-------------------------------------------------
func Eth_tx_trace__plot(p_tx_id_hex_str string,
	p_get_hosts_fn func(context.Context, *GF_runtime) []string,
	p_ctx          context.Context,
	p_py_plugins   *GF_py_plugins,
	p_metrics      *GF_metrics,
	p_runtime      *GF_runtime) (string, *gf_core.Gf_error) {



	//-----------------------
	// WORKER_INSPECTOR__HOST_PORT
	host_port_str      := p_get_hosts_fn(p_ctx, p_runtime)[0]
	start_time__unix_f := float64(time.Now().UnixNano()) / 1000000000.0

	// GET_TRACE
	gf_tx_trace, gf_err := Eth_tx_trace__get_from_worker_inspector(p_tx_id_hex_str,
		host_port_str,
		p_ctx,
		p_runtime.Runtime_sys)

	end_time__unix_f := float64(time.Now().UnixNano())/1000000000.0

	// METRICS
	if p_metrics != nil {
		delta_time__unix_f := end_time__unix_f - start_time__unix_f
		p_metrics.Tx_trace__worker_inspector_durration__gauge.Set(delta_time__unix_f)
	}

	if gf_err != nil {
		return "", gf_err
	}

	//-----------------------
	// PY_PLUGIN__PLOT

	start_time__unix_f = float64(time.Now().UnixNano()) / 1000000000.0
	plot_svg_str, gf_err := py__run_plugin__plot_tx_trace(p_tx_id_hex_str,
		gf_tx_trace,
		p_py_plugins,
		p_runtime.Runtime_sys)
	end_time__unix_f = float64(time.Now().UnixNano())/1000000000.0
	
	// METRICS
	if p_metrics != nil {
		delta_time__unix_f := end_time__unix_f - start_time__unix_f
		p_metrics.Tx_trace__py_plugin__plot_durration__gauge.Set(delta_time__unix_f)
	}

	if gf_err != nil {
		return "", gf_err
	}

	//-----------------------

	return plot_svg_str, nil
}

//-------------------------------------------------
// GET_FROM_WORKER_INSPECTOR
func Eth_tx_trace__get_from_worker_inspector(p_tx_hash_str string,
	p_host_port_str string,
	p_ctx           context.Context,
	p_runtime_sys   *gf_core.Runtime_sys) (*GF_eth__tx_trace, *gf_core.Gf_error) {

	url_str := fmt.Sprintf("http://%s/gfethm_worker_inspect/v1/tx/trace?tx=%s",
	p_host_port_str,
		p_tx_hash_str)

	//-----------------------
	// SPAN
	span_name_str      := fmt.Sprintf("worker_inspector__get_tx_trace:%s", p_host_port_str)
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

	return gf_tx_trace, nil
}

//-------------------------------------------------
func Eth_tx_trace__get(p_tx_hash_str string,
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