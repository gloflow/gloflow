$( document ).ready(function() {
    main();
});



//---------------------------------------------------
function main() {




    $("body").append("<div>gf_eth_monitor</div>");
    $("body").append(`<div>
        <div>block #</div>
        <input id="block_num" value="2000000"></input>
    </div>`);





    $("input#block_num").on('keypress',(e)=>{
        // "enter" pressed
        if(e.which == 13) {
            const block_int = $("input#block_num").val();

            

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


    
}

//---------------------------------------------------
function render__block_from_workers(p_block_int,
    p_block_from_workers_map,
    p_miners_map) {
                            
    const block_element = $(`<div class="block">
        <div class="block_metadata">
            

            <div class="miners">

            </div>
            
        </div>
        

    </div>`);
    $("body").append(block_element);



    if (p_miners_map != undefined) {
        Object.entries(p_miners_map).forEach(e=> {

            const miner_map             = e[1];
            const miner_name_str        = miner_map["name_str"];
            const miner_address_hex_str = miner_map["address_hex_str"];

            $(block_element).find(".miners").append(`<div class="miner_info">
                miner: <span class="miner_name">${miner_name_str}</span>
            </div>`);

        });
    }


    //---------------------------------------------------
    function render__block(p_block_int,
        p_worker_host_str,
        p_block_map,
        p_block_parent_element) {

        const gas_used_int      = p_block_map["gas_used_int"];
        const gas_limit_int     = p_block_map["gas_limit_int"];
        const coinbase_addr_str = p_block_map["coinbase_addr_str"];

        const hash_str        = p_block_map["hash_str"];
        const parent_hash_str = p_block_map["parent_hash_str"];
        const time_int        = p_block_map["time_int"];
        


        const creation_time_f      = parseFloat(time_int);
		const creation_date        = new Date(creation_time_f*1000);
        const timeago_str          = $.timeago(creation_date);
            

        const block_from_worker__element = $(`<div class="block_from_worker">
            <div>block #        <span class="block_num">${p_block_int} </span><a href="https://etherscan.io/block/${p_block_int}" target="_blank">etherscan.io</a></div>
            <div>hash:          <span class="block_hash">${hash_str}</span></div>
            <div>parent hash:   <span class="block_parent_hash">${parent_hash_str}</span></div>
            <div>time:          <span class="block_time">${creation_date} </span><span class="block_timeago">(${timeago_str})</span></div>
            <div>gas used:      <span>${gas_used_int}</span></div>
            <div>gas limit:     <span>${gas_limit_int}</span></div>
            <div>worker host:   <span>${p_worker_host_str}</span></div>
            <div class="coinbase_addr">coinbase addr: <span>${coinbase_addr_str} </span>(<a href="https://etherscan.io/address/${coinbase_addr_str}" target="_blank">etherscan.io</a>)</div>
            
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

            const tx_map           = e[1];
            const tx_hash_str      = tx_map["hash_str"];
            const tx_from_addr_str = tx_map["from_addr_str"];
            const tx_to_addr_str   = tx_map["to_addr_str"];
            const tx_value_eth_f   = tx_map["value_eth_f"];
            const tx_gas_used_int  = tx_map["gas_used_int"];
            const tx_gas_price_int = tx_map["gas_price_int"];
            const tx_nonce_int     = tx_map["nonce_int"];
            const tx_size_f        = tx_map["size_f"];
            const tx_cost_int      = tx_map["cost_int"];

            const tx_element = $(`<div class="tx">
                <div class="tx_hash">hash           - <span>${tx_hash_str} </span><a href="https://etherscan.io/tx/${tx_hash_str}" target="_blank">etherscan.io</a></div>
                <div class="source_destination">
                    <div class="to_addr">From       - <span>${tx_from_addr_str} </span>(<a href="https://etherscan.io/address/${tx_from_addr_str}" target="_blank">etherscan.io</a>)</div>
                    <div class="from_addr">To       - <span>${tx_to_addr_str} </span>(<a href="https://etherscan.io/address/${tx_to_addr_str}" target="_blank">etherscan.io</a>)</div>
                </div>
                <div class="tx_value">value         - <span>${tx_value_eth_f}</span>eth</div>
                <div class="tx_gas_used">gas used   - <span>${tx_gas_used_int}</span></div>
                <div class="tx_gas_price">gas price - <span>${tx_gas_price_int}</span></div>
                <div class="tx_nonce">nonce         - <span>${tx_nonce_int}</span></div>
                <div class="tx_size">size           - <span>${tx_size_f}</span></div>
                <div class="tx_cost">cost           - <span>${tx_cost_int}</span></div>
            </div>`);

            
            $(txs_element).find(".txs_list").append(tx_element);


            if (tx_gas_used_int > 21000) {

                $(tx_element).find(".tx_gas_used span").addClass("not_just_value_transfer");

            }
        });

        return txs_element;
    }

    //---------------------------------------------------



    Object.entries(p_block_from_workers_map).forEach(e=> {

        const worker_host_str = e[0];
        const block_map       = e[1];


        console.log(block_map)




        render__block(p_block_int, worker_host_str, block_map, block_element);


        
    });

}



//---------------------------------------------------
function http__get_block(p_block_num_int,
	p_on_complete_fun,
	p_on_error_fun) {

	const url_str = "/gfethm/v1/block?b="+p_block_num_int;


	//-------------------------
	// HTTP AJAX
	$.get(url_str,
		function(p_data_map) {
            console.log('response received');
            
			console.log('data_map["status_str"] - '+p_data_map["status_str"]);
			
			if (p_data_map["status_str"] == "OK") {

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