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

$( document ).ready(function() {
    main();
});

//---------------------------------------------------
function main() {




    $("body").append("<div class='app_label'>gf_eth_monitor</div>");
    
    
    const monitor_element = $(`<div class='monitor'>
        <div>block #</div>
        <input id="search" value="2000000"></input>

        <div id="index_block">
            <input id="block_range__start" value="2000000">
            <input id="block_range__end"   value="2000001">
            <div id="submit_btn"></div>
        </div>
    </div>`);
    
    $("body").append(monitor_element);





    $("input#search").on('keypress',(e)=>{
        // "enter" pressed
        if(e.which == 13) {
            const block_int = $("input#search").val();

            

            http__get_block(block_int,
                function(p_block_from_workers_map, p_miners_map) {
                    
                    $(".block").remove(); // remove old block display if any

                    render__block_from_workers(block_int,
                        p_block_from_workers_map,
                        p_miners_map);
                    
        
        
                    
                },
                function() {});
        }
    });




    $(".submit_btn").on('click', (e)=>{

        const block_range__start_int = $("input#block_range__start").val();
        const block_range__end_int   = $("input#block_range__end").val();


        http__index_block(block_range__start_int,
            block_range__end_int,
            function() {


                $(monitor_element).find("#index_block #submit_btn").css("background-color", "green");
            },
            function(){
                $(monitor_element).find("#index_block #submit_btn").css("background-color", "red");
            });
    });
    
}

//---------------------------------------------------
function render__block_from_workers(p_block_uint,
    p_block_from_workers_map,
    p_miners_map) {
                            
    const block_element = $(`<div class="block">
        <div class="block_metadata">
            

            <div class="miners">

            </div>
            
        </div>
        

    </div>`);
    $("body").append(block_element);

    // MINERS
    Object.entries(p_miners_map).forEach(e=> {

        const miner_map             = e[1];
        const miner_name_str        = miner_map["name_str"];
        const miner_address_hex_str = miner_map["address_hex_str"];

        $(block_element).find(".miners").append(`<div class="miner_info">
            miner: <span class="miner_name">${miner_name_str}</span>
        </div>`);

    });


    //---------------------------------------------------
    function render__block(p_block_uint,
        p_worker_host_str,
        p_block_map,
        p_block_parent_element) {

        const gas_used_uint     = p_block_map["gas_used_uint"];
        const gas_limit_uint    = p_block_map["gas_limit_uint"];
        const coinbase_addr_str = p_block_map["coinbase_addr_str"];

        const hash_str        = p_block_map["hash_str"];
        const parent_hash_str = p_block_map["parent_hash_str"];
        const time_int        = p_block_map["time_uint"];
        


        const creation_time_f = parseFloat(time_int);
		const creation_date   = new Date(creation_time_f*1000);
        const timeago_str     = $.timeago(creation_date);
            

        const block_from_worker__element = $(`<div class="block_from_worker">
            <div>block #        <span class="block_num">${p_block_uint} </span><a href="https://etherscan.io/block/${p_block_uint}" target="_blank">etherscan.io</a></div>
            <div>h:             <span class="block_hash">${hash_str}</span></div>
            <div>parent hash:   <span class="block_parent_hash">${parent_hash_str}</span></div>
            <div>time:          <span class="block_time">${creation_date} </span><span class="block_timeago">(${timeago_str})</span></div>
            <div>gas used:      <span class="gas_used">${gas_used_uint}</span></div>
            <div>gas limit:     <span class="gas_limit">${gas_limit_uint}</span></div>
            <div>worker host:   <span class="worker_host">${p_worker_host_str}</span></div>
            <div class="coinbase_addr">coinbase addr: <span>${coinbase_addr_str} (<a href="https://etherscan.io/address/${coinbase_addr_str}" target="_blank">etherscan.io</a>)</span></div>
            
        </div>`);

        $(p_block_parent_element).append(block_from_worker__element);




        const txs_lst     = p_block_map["txs_lst"];
        const txs_element = render__block_txs(txs_lst);

        $(block_from_worker__element).append(txs_element);

    }

    //---------------------------------------------------
    function render__block_txs(p_txs_lst) {

        const txs_element = $(`<div class="txs">
            <div>txs # <span>${p_txs_lst.length}</span></div>
            <div class="txs_list">

            </div>
        </div>`);

        Object.entries(p_txs_lst).forEach(e => {

            const tx_map = e[1];

            // if transaction failed to be loaded from a block on the backend,
            // it will be marked as null, and we should skip to the next transaction
            if (tx_map == null) {
                return
            }
            
            const tx_hash_str      = tx_map["hash_str"];
            const tx_from_addr_str = tx_map["from_addr_str"];
            const tx_to_addr_str   = tx_map["to_addr_str"];
            const tx_value_eth_f   = tx_map["value_eth_f"];
            const tx_gas_used_int  = tx_map["gas_used_int"];
            const tx_gas_price_int = tx_map["gas_price_int"];
            const tx_nonce_int     = tx_map["nonce_int"];
            const tx_size_f        = tx_map["size_f"];
            const tx_cost_gwei_f   = tx_map["cost_gwei_f"];

            // TRANSACTION
            const tx_element = $(`<div class="tx">
                <div class="tx_hash">hash           - <span>${tx_hash_str} </span><a href="https://etherscan.io/tx/${tx_hash_str}" target="_blank">etherscan.io</a></div>
                <div class="source_destination">
                    <div class="from_addr">From   - <span>${tx_from_addr_str} </span>(<a href="https://etherscan.io/address/${tx_from_addr_str}" target="_blank">etherscan.io</a>)</div>
                    <div class="to_addr">To       - <span>${tx_to_addr_str} </span></div>
                </div>
                
                <div class="tx_value">value         - <span>${tx_value_eth_f}</span>eth</div>
                <div class="tx_gas_used">gas used   - <span>${tx_gas_used_int}</span></div>
                <div class="tx_gas_price">gas price - <span>${tx_gas_price_int}</span></div>
                <div class="tx_nonce">nonce         - <span>${tx_nonce_int}</span></div>
                <div class="tx_size">size           - <span>${tx_size_f}</span></div>
                <div class="tx_cost">cost           - <span>${tx_cost_gwei_f}</span> gwei</div>
            </div>`);

            

            //----------------------------
            // NEW_CONTRACT
            // for new_contract transactions, where the To() returned by Eth node is nil, to_addr is set to "new_contract" by worker_inspector.
            // dont include etherescan.io validation link for those
            if (tx_to_addr_str == "new_contract") {
                const contract_new_addr_str   = tx_map["contract_new_map"]["addr_str"];
                const contract_code_bytes_lst = tx_map["contract_new_map"]["code_bytes_lst"];

                $(tx_element).find(".to_addr").append(`<a href="https://etherscan.io/address/${contract_new_addr_str}" target="_blank">etherscan.io</a>`);

                $(tx_element).append(`<div>
                    <div class="addr">${contract_new_addr_str}</div>
                    <div class="code_bytes">${contract_code_bytes_lst}</div>
                </div>`)
            }
            
            //----------------------------
            // for all transactions include the etherscan.io validation link
            else {
                $(tx_element).find(".to_addr").append(`<a href="https://etherscan.io/address/${tx_to_addr_str}" target="_blank">etherscan.io</a>`);
            }


            $(txs_element).find(".txs_list").append(tx_element);

            if (tx_gas_used_uint > 21000) {
                $(tx_element).find(".tx_gas_used span").addClass("not_just_value_transfer");




                // TX_TRACE_BUTTON - trace TX's that are not simple value transfers
                $(tx_element).append(`<div class="trace_btn">trace</div>`);
                $(tx_element).find(".trace_btn").on('click', function() {

                    // VIEW_TRACE
                    view_trace(tx_hash_str,
                        // p_on_error_fun
                        function() {
                            $(tx_element).find(".trace_btn").css("background-color", "red");
                        });
                });
            }
        });

        return txs_element;
    }

    //---------------------------------------------------

    Object.entries(p_block_from_workers_map).forEach(e=> {

        const worker_host_str = e[0];
        const block_map       = e[1];

        // NO_BLOCK - from a particular worker. so just skip it.
        if (block_map == null) {
            return
        }

        render__block(p_block_uint, worker_host_str, block_map, block_element);
    });
}

//---------------------------------------------------
function view_trace(p_tx_id_str, p_on_error_fun) {

    http__get_trace(p_tx_id_str,
        function(p_tx_trace_svg_str) {

            $("body").append(`<div id="tx_trace"></div>`)
            const draw = SVG().addTo('#tx_trace');
            draw.svg(p_tx_trace_svg_str);



            

            // IMPORTANT!! - position the plot depending on where the user scroll
            //               to, to always keep the plot in view when the user opens it.
            const current_global_scroll_y = $(window).scrollTop();
            $("body #tx_trace").css("top", `${current_global_scroll_y+200}px`);

            const svg_e = $("body #tx_trace svg")[0];
            const svg_plot_bbox = svg_e.getBBox();

            //----------------------------
            // CSS
            const svg_plot_height_int = Math.floor(svg_plot_bbox.height);
            $(svg_e).css("background-color", "white");
            $(svg_e).css("height", `${svg_plot_height_int+100}px`);

            // FIX!! - limit the possible width of svg plots (gas cost of instructions), in py_plugin
            $(svg_e).css("width", `${2000}px`);

            //----------------------------
        },
        function(){
            p_on_error_fun();
        })
}

//---------------------------------------------------
function http__index_block(p_range_start_int,
    p_range_end_int,
    p_on_complete_fun,
	p_on_error_fun) {


    const url_str = `/gfethm/v1/block/index?br=${p_range_start_int}-${p_range_end_int}`;

    //-------------------------
	// HTTP AJAX
	$.get(url_str,
		function(p_data_map) {

            console.log('response received');
			console.log(`data_map["status"] - ${p_data_map["status"]}`);
			
			if (p_data_map["status"] == "OK") {


                p_on_complete_fun();

			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
        });
    
	//-------------------------	
}

//---------------------------------------------------
function http__get_trace(p_tx_id_str,
    p_on_complete_fun,
	p_on_error_fun) {

    const url_str = `/gfethm/v1/tx/trace/plot?tx=${p_tx_id_str}`;

    //-------------------------
	// HTTP AJAX
	$.get(url_str,
		function(p_data_map) {

            console.log('response received');
			console.log(`data_map["status"] - ${p_data_map["status"]}`);
			
			if (p_data_map["status"] == "OK") {

				const tx_trace_svg_str = p_data_map["data"]["plot_svg_str"];
                p_on_complete_fun(tx_trace_svg_str);

			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
        });
    
	//-------------------------	
}

//---------------------------------------------------
function http__get_block(p_block_num_int,
	p_on_complete_fun,
	p_on_error_fun) {

	const url_str = `/gfethm/v1/block?b=${p_block_num_int}`;

	//-------------------------
	// HTTP AJAX
	$.get(url_str,
		function(p_data_map) {
            console.log('response received');
			console.log(`data_map["status"] - ${p_data_map["status"]}`);
			
			if (p_data_map["status"] == "OK") {

				const block_from_workers_map = p_data_map["data"]["block_from_workers_map"];
                const miners_map             = p_data_map["data"]["miners_map"];
                p_on_complete_fun(block_from_workers_map, miners_map);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
        });
    
	//-------------------------	
}