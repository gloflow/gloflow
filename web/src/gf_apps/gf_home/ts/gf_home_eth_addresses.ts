/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

///<reference path="../../../d/jquery.d.ts" />

import * as gf_dragndrop from "./../../../gf_core/ts/gf_dragndrop";
import * as gf_tagger_input_ui from "./../../gf_tagger/ts/gf_tagger_client/gf_tagger_input_ui";
import * as gf_utils from "gf_utils";

declare var Web3;

//--------------------------------------------------------
export async function init_observed(p_parent_element,
	p_http_api_map,
	p_assets_paths_map) {
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const component_name_str = "web3_addresses_observed";
		const address_type_str   = "observed";
		const address_chain_str  = "eth";
		const output_map = await p_http_api_map["home"]["web3_addresses_get_all_fun"](address_type_str,
			address_chain_str);
		const eth_addresses_lst = output_map["addresses_lst"];

		const container = $(`
			<div id="observed_eth_addresses">
				<div id="title">observed eth addresses</div>
			</div>`);
		$(p_parent_element).append(container);

		

		// NO_ADDRESSES - there are no initial observed addresses, so have custom UI for
		//                adding an initial address.
		if (eth_addresses_lst.length == 0) {
			init_no_address_dialog(container,
				address_type_str,
				p_http_api_map,
				p_assets_paths_map);
		}
		else {

			const initial_height_int = $(container).outerHeight();
			var total_height_int = initial_height_int;
			for (const eth_address_str of eth_addresses_lst) {

				// CREATE_ETH_ADDRESS
				const eth_address_element = create_eth_address(eth_address_str,
					address_type_str,
					container,
					
					//--------------------------------------------------------
					// p_on_new_address_btn_fun
					(p_added_address_container)=>{

						// update parent container height
						total_height_int += $(p_added_address_container).outerHeight();
						$(container).css("height", `${total_height_int}px`);
					},

					//--------------------------------------------------------
					p_http_api_map,
					p_assets_paths_map);
				
				$(container).append(eth_address_element);
				
				total_height_int += $(eth_address_element).outerHeight();
			}

			// update parent container height
			$(container).css("height", `${total_height_int}px`);
		}

		// DRAG_N_DROP
		gf_dragndrop.init(container,

			//--------------------------------------------------------
			// p_on_dnd_event_fun
			async (p_dnd_event_type_str :string, p_drag_data_map)=>{

				switch (p_dnd_event_type_str) {
					case "drag_start":
						break;

					case "drag_stop":
						// update component remotely on each drag/coord change
						await gf_utils.update_viz_component_remote(component_name_str,
							p_drag_data_map,
							p_http_api_map);
						break;
				}
			},
			
			//--------------------------------------------------------
			p_assets_paths_map);

		p_resolve_fun(container);
	});
	return p;
}

//--------------------------------------------------------
export async function init_my(p_parent_element,
	p_http_api_map,
	p_assets_paths_map) {
	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {

		const component_name_str = "web3_addresses_my";
		const address_type_str   = "my";
		const address_chain_str  = "eth";
		const output_map = await p_http_api_map["home"]["web3_addresses_get_all_fun"](address_type_str,
			address_chain_str);
		const eth_addresses_lst = output_map["addresses_lst"];

		const container = $(`
			<div id="my_eth_addresses">
				<div id="title">my eth addresses</div>
			</div>`);
		$(p_parent_element).append(container);

		// NO_ADDRESSES - there are no initial observed addresses, so have custom UI for
		//                adding an initial address.
		if (eth_addresses_lst.length == 0) {

			init_no_address_dialog(container,
				address_type_str,
				p_http_api_map,
				p_assets_paths_map);
			
		}
		else {
			const initial_height_int = $(container).outerHeight();
			var total_height_int = initial_height_int;
			for (const eth_address_str of eth_addresses_lst) {

				// CREATE_ETH_ADDRESS
				const eth_address_element = create_eth_address(eth_address_str,
					address_type_str,
					container,
					
					//--------------------------------------------------------
					// p_on_new_address_btn_fun
					(p_added_address_container)=>{

						// update parent container height
						total_height_int += $(p_added_address_container).outerHeight();
						$(container).css("height", `${total_height_int}px`);
					},

					//--------------------------------------------------------
					p_http_api_map,
					p_assets_paths_map);

				$(container).append(eth_address_element);

				total_height_int += $(eth_address_element).outerHeight();
			}

			// update parent container height
			$(container).css("height", `${total_height_int}px`);
		}

		// DRAG_N_DROP
		gf_dragndrop.init(container, 
			//--------------------------------------------------------
			// p_on_dnd_event_fun
			async (p_dnd_event_type_str :string, p_drag_data_map)=>{

				switch (p_dnd_event_type_str) {
					case "drag_start":
						break;

					case "drag_stop":

						// update component remotely on each drag/coord change
						await gf_utils.update_viz_component_remote(component_name_str,
							p_drag_data_map,
							p_http_api_map);
							
						break;
				}
			},
			
			//--------------------------------------------------------
			p_assets_paths_map);

		p_resolve_fun(container);
	});
	return p;
}

//--------------------------------------------------------
function init_no_address_dialog(p_container,
	p_address_type_str,
	p_http_api_map,
	p_assets_paths_map) {

	const add_initial_btn = $(`
		<div id="add_initial_btn">
			<div class="add_initial_address_btn">
				<img src="${p_assets_paths_map["gf_add_btn"]}" draggable="false"></img>
			</div>
		</div>`)
	$(p_container).append(add_initial_btn);



	var total_height_int = 0;
	$(add_initial_btn).find(".add_initial_address_btn").on("click", ()=>{

		// ADD_NEW_ADDRESS
		const added_address_container = create_eth_address_input(p_address_type_str,
			p_http_api_map,
			p_assets_paths_map);
		$(p_container).append(added_address_container);

		// update parent container height
		total_height_int += $(added_address_container).outerHeight();
		$(p_container).css("height", `${total_height_int}px`);

		// add_initial_btn is attached to the DOM (has a parent),
		// so remove it because its only used for adding initial addresses
		if ($(add_initial_btn).parent().length > 0) {
			$(add_initial_btn).remove();
		}
	});
}

//--------------------------------------------------------
function create_eth_address(p_address_str :string,
	p_address_type_str :string,
	p_parent_element,
	p_on_new_address_btn_fun,
	p_http_api_map,
	p_assets_paths_map) {
	
	const eth_address_short_start_str = `${p_address_str.substr(0, 7)}`;
	const eth_address_short_end_str   = `${p_address_str.substr(p_address_str.length-7, 7)}`;
	const eth_address_element = $(`
		<div class="eth_address">
			<div class="hex_address">${eth_address_short_start_str}...${eth_address_short_end_str}</div>
		</div>
	`);

	var added_bool = false;
	var info_container_element;
	$(eth_address_element).on("click", ()=>{

		if (!added_bool) {

			// ETH_ADDRESS - UI control that represents an ETH address and actions supported on it
			info_container_element = $(`
				<div class="info">
					<div class="index_nfts_for_owner_btn">
						get NFTs
					</div>
					<div class="add_tag_btn">
						t
					</div>
					<div class="etherscan_btn">
						<a href="https://etherscan.io/address/${p_address_str}" target="_blank">e</a>
					</div>
					<div class="add_new_address_btn">
						<img src="${p_assets_paths_map["gf_add_btn"]}" draggable="false"></img>
					</div>
					<div class="tags_container">

					</div>
				</div>`);

			//------------------------------
			// INDEX_NFTS
			$(info_container_element).find(".index_nfts_for_owner_btn").on("click", async (p_e)=>{
				p_e.stopImmediatePropagation();

				console.log("index owner NFTs");



				console.log("+++++++++++++++++++++++++++++++++++++++++++++++")


				const nfts_container_element = await index_address_nfts(p_address_str,
					p_http_api_map);


				$(eth_address_element).append(nfts_container_element as HTMLElement);


			});

			//------------------------------
			// ADD_TAG
			$(info_container_element).find(".add_tag_btn").on("click", async (p_e)=>{
				p_e.stopImmediatePropagation();

				console.log("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

				const obj_id_str   = p_address_str;
				const obj_type_str = "address"
				const obj_element  = info_container_element;


				const log_fun = (p_g, p_m)=>{
					console.log(p_g+":"+p_m);
				};
				const tagging_input_ui_element = gf_tagger_input_ui.init_tagging_input_ui_element(obj_id_str,
					obj_type_str,

					//--------------------------------------------------------
					// p_on_tags_created_fun
					(p_tags_lst)=>{

						console.log("CCCCCCCCCCCCCCCCCCCCCC")
						console.log(p_tags_lst)


						const tags_element = $(info_container_element).find(".tags_container");

						for (const tag_str of p_tags_lst) {


							const new_tag = $(`<div class="tag">
									${tag_str}
								</div>`);
							$(tags_element).append(new_tag);


							if (tag_str.startsWith("name:")) {
								const name_str = tag_str.split(":")[1];
								$(eth_address_element).find(".hex_address").text(name_str);
							}


						}


					},

					//--------------------------------------------------------
					// p_on_tag_ui_remove_fun
					()=>{

					},
					p_http_api_map,
					log_fun);

				gf_tagger_input_ui.place_tagging_input_ui_element(tagging_input_ui_element,
						obj_element,
						log_fun);
				
			});

			//------------------------------
			// ADD_NEW_ADDRESS
			$(info_container_element).find(".add_new_address_btn").on("click", (p_e)=>{
				p_e.stopImmediatePropagation();

				// ADD_NEW_ADDRESS
				const added_address_container = create_eth_address_input(p_address_type_str,
					p_http_api_map,
					p_assets_paths_map);
				$(p_parent_element).append(added_address_container);

				p_on_new_address_btn_fun(added_address_container);
			});

			//------------------------------

			$(eth_address_element).append(info_container_element);

			added_bool = true;
		}
		else {

			$(info_container_element).remove();
			info_container_element = null;
			added_bool = false;
		}


	});

	return eth_address_element;
}

//--------------------------------------------------------
function create_eth_address_input(p_address_type_str :string,
	p_http_api_map,
	p_assets_paths_map) {

	const new_address_container = $(`
		<div class="eth_address_input">
			<input class="hex_address_input"></input>
			<div class="confirm_btn">
				<img src="${p_assets_paths_map["gf_confirm_btn"]}" draggable="false"></img>
			</div>
		</div>`);


	// ADD_NEW_ADDRESS
	$(new_address_container).find(".confirm_btn").on("click", async ()=>{

		const new_address_str = $(new_address_container).find("input").val();

		if (new_address_str != "") {

			// VALIDATE
			const valid_bool = Web3.utils.isAddress(new_address_str)
			if (!valid_bool) {
				$(new_address_container).css("background-color", "red");
				return;
			}
			
			// HTTP
			const chain_str = "eth";
			const output_map = await p_http_api_map["home"]["web3_address_add_fun"](new_address_str,
				p_address_type_str,
				chain_str);
			
			console.log(output_map);
			
			// CREATE_NEW_ADDRESS
			const parent_element = $(new_address_container).parent();
			const new_address_element = create_eth_address(new_address_str,
				p_address_type_str,
				parent_element,

				//--------------------------------------------------------
				// on_new_address_btn_fun
				(p_added_address_container)=>{

				},

				//--------------------------------------------------------
				p_http_api_map,
				p_assets_paths_map);
			
			$(parent_element).append(new_address_element);
			
			// remove address input field
			$(new_address_container).remove();
		}
	});
	
	return new_address_container;
}

//--------------------------------------------------------
// INDEX_ADDRESS_NFTS

function index_address_nfts(p_address_str :string,
	p_http_api_map) {
	

	const p = new Promise(async function(p_resolve_fun, p_reject_fun) {



		const chain_str = "eth";

		// HTTP_CALL - NFT_INDEX - for a given ETH address
		const output_map = await p_http_api_map["home"]["web3_nft_index_for_address_fun"](p_address_str,
			chain_str)


		const nfts_lst = output_map["nfts_lst"];


		const nfts_container = $(`
			<div class="nfts_container">
			
			</div>`);

		if (nfts_lst != null) {
			for (const nft_map of nfts_lst) {



				console.log(">>>>", nft_map);


				view_nft(nft_map);			

			}
		}

		//--------------------------------------------------------
		function view_nft(p_nft_map) {
			const owner_address_str = p_nft_map["owner_address_str"]
			const token_id_str      = p_nft_map["token_id_str"]
			const contract_address_str = p_nft_map["contract_address_str"]
			const contract_name_str    = p_nft_map["contract_name_str"]
			const chain_str            = p_nft_map["chain_str"]

			const token_uri_raw_str = p_nft_map["token_uri_raw_str"]
			const media_uri_raw_str = p_nft_map["media_uri_raw_str"]


			const nft_element = $(`
				<div class="nft">
					<div class="owner_address">
						${owner_address_str}
					</div>

					<div class="token_id">
						${token_id_str}
					</div>

					<div class="contract_address">
						${contract_address_str}
					</div>

					<div class="contract_name">
						${contract_name_str}
					</div>

					<div class="chain">
						${chain_str}
					</div>

					<div class="token_uri_raw">
						${token_uri_raw_str}
					</div>

					<div class="media_uri_raw">
						${media_uri_raw_str}
					</div>
				</div>`);
			$(nfts_container).append(nft_element);
		}

		//--------------------------------------------------------
		p_resolve_fun(nfts_container);


	});
	return p;



	
}