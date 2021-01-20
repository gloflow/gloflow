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
                function(p_block_from_workers_map) {
                    

                    
                        
                        
                    const block_element = $(`<div class="block">
                        <div class="block_metadata">
                            <a href="https://etherscan.io/block/${block_int}" target="_blank">etherscan.io</a> 
                        </div>


                    </div>`);
                    $("body").append(block_element);



        
                    Object.entries(p_block_from_workers_map).forEach(e=> {
        
                        const worker_host_str = e[0];
                        const block_map       = e[1];
        
        
                        console.log(block_map)
        
        
        
        
                        const gas_used_int      = block_map["gas_used_int"];
                        const gas_limit_int     = block_map["gas_limit_int"];
                        const coinbase_addr_str = block_map["coinbase_addr_str"];
        
        
        
                        $(block_element).append(`<div class="block_from_worker">
                            <div>block #        <span class="block_num">${block_int}</span></div>
                            <div>worker host:   <span>${worker_host_str}</span></div>
                            <div>gas used:      <span>${gas_used_int}</span></div>
                            <div>gas limit:     <span>${gas_limit_int}</span></div>
                            <div>coinbase addr: <span>${coinbase_addr_str}</span></div>
                        </div>`);
                    });
        
        
                    
                },
                function() {});
        }
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
			
			if (p_data_map["status_str"] == 'OK') {

				const block_from_workers_map = p_data_map['data']["block_from_workers_map"];
				p_on_complete_fun(block_from_workers_map);
			}
			else {
				p_on_error_fun(p_data_map["data"]);
			}
        });
        
	//-------------------------	
}