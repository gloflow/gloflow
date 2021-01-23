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
                function(p_block_from_workers_map, p_miners_lst) {
                    
                    $(".block").remove(); // remove old block display if any

                    render__block_from_workers(block_int,
                        p_block_from_workers_map,
                        p_miners_lst);
                    
        
        
                    
                },
                function() {});
        }
    });


    
}

//---------------------------------------------------
function render__block_from_workers(p_block_int,
    p_block_from_workers_map,
    p_miners_lst) {
                            
    const block_element = $(`<div class="block">
        <div class="block_metadata">
            <a href="https://etherscan.io/block/${p_block_int}" target="_blank">etherscan.io</a>

            <div class="miners">

            </div>
            
        </div>
        

    </div>`);
    $("body").append(block_element);




    Object.entries(p_miners_lst).forEach(e=> {

        const miner_map             = e[1];
        const miner_name_str        = miner_map["name_str"];
        const miner_address_hex_str = miner_map["address_hex_str"];

        $(block_element).find(".miners").append(`<div class="miner_info">
            miner: <span class="miner_name">${miner_name_str}</span>
        </div>`);

    });


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
        


        const block_from_worker__element = $(`<div class="block_from_worker">
            <div>block #        <span class="block_num">${p_block_int}</span></div>
            <div>hash:          <span class="block_hash">${hash_str}</span>
            <div>parent hash:   <span class="block_parent_hash">${parent_hash_str}</span>
            <div>time:          <span class="block_time">${time_int}</span>
            <div>gas used:      <span>${gas_used_int}</span></div>
            <div>gas limit:     <span>${gas_limit_int}</span></div>
            <div>worker host:   <span>${p_worker_host_str}</span></div>
            <div class="coinbase_addr">coinbase addr: <span>${coinbase_addr_str}</span></div>
            
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

            const tx_map          = e[1];
            const tx_hash_str     = tx_map["hash_str"];
            const tx_gas_used_int = tx_map["gas_used_int"];

            const tx_element = $(`<div class="tx">
                <div class="etherscan_link"><a href="https://etherscan.io/address/${tx_hash_str}">etherscan.io</a></div>
                <div class="tx_hash">hash - <span>${tx_hash_str}</span></div>
                <div class="tx_gas_used">gas used - <span>${tx_gas_used_int}</span></div>
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
                const miners_lst             = p_data_map["data"]["miners_lst"];
                p_on_complete_fun(block_from_workers_map, miners_lst);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
        });
        
	//-------------------------	
}