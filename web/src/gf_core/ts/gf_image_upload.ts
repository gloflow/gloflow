/*
GloFlow application and media management/publishing platform
Copyright (C) 2025 Ivan Trajkovic

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

//-------------------------------------------------
export function init(p_flow_name_str :string,
	p_target_full_host_str :string,
	p_on_upload_fun        :Function) {

	// console.log("UPLOAD INITIALIZED")
	document.onpaste = function(p_paste_event) {

		console.log("paste event...");

		// const items = (p_paste_event.clipboardData || p_paste_event.originalEvent.clipboardData).items;
		const items = p_paste_event.clipboardData?.items;
		if (items) {

		
			// console.log(p_paste_event.clipboardData);
			// console.log(p_paste_event.originalEvent.clipboardData);
			console.log("pasted content", JSON.stringify(items)); // will give you the mime types

			for (const item of items) {

				console.log("item", item);

				if (item.kind === 'file') {

					const blob   = item.getAsFile();
					const reader = new FileReader();

					reader.onload = function(p_event) {
						
						console.log("data loaded");
						
						// result attribute contains the data as a data: URL representing
						// the file's data as a base64 encoded string.
						const img_data = p_event.target?.result;
						const img_data_str: string = (img_data instanceof ArrayBuffer 
                            ? new TextDecoder("utf-8").decode(img_data) 
                            : img_data) || "";
						
						// the beginning of data URL strings example:
						// "data:image/png;base64".
						// IMPORTANT!! - it seems all images are of "png" format when pasted in.
						const image_format_str = img_data_str.split(";")[0].replace("data:image/", "")
						console.log(`image_format_str - ${image_format_str}`);
						
						//-------------------------------------------------
						// VIEW_IMAGE
						gf_upload__view_img(img_data_str,
							p_flow_name_str,

							//-------------------------------------------------
							// UPLOAD_ACTIVATE_FUN
							async (p_image_name_str :string,
							p_flows_names_str        :string,
							p_on_upload_complete_fun :Function)=>{

								// UPLOAD_IMAGE
								const upload_gf_image_id_str = await gf_upload__run(p_image_name_str,
									img_data_str,
									image_format_str,
									p_flows_names_str,
									p_target_full_host_str);

								p_on_upload_complete_fun();
								p_on_upload_fun(upload_gf_image_id_str);
							});

							//-------------------------------------------------

						//-------------------------------------------------
					};

					// reader.readAsBinaryString(blob);
					reader.readAsDataURL(blob);
				}
			}
		}
	}
}

//-------------------------------------------------
function gf_upload__view_img(p_img_data_str :string,
	p_flow_name_str       :string,
	p_upload_activate_fun :Function) {
	
	//-------------------------------------------------
	function get_image_dialog() {

		//-------------------------------
		// first image
		if ($("#upload_image_dialog").length == 0) {

			/*
			const img_dialog = $(`
				<div id="upload_image_dialog" class="gf_center">
					<div id="background"></div>

					<div id="upload_images_detail">
						<div id="upload_images">

							<div id="upload_images_panel">
								
								<div id="images">
									
								</div>

								<div id="upload_image_flow_name_input" class="gf_center">
									<input placeholder="flow name" value="${p_flow_name_str}"></input>
								</div>

							</div>

						</div>
						<div id="upload_btn">upload image</div>
					</div>

				</div>`);
			*/

			const dialog = $(`
				<div id="upload_image_dialog" class="gf_center">
					<div id="background"></div>

					<div id="upload_image_flow_name_input" class="gf_center">
						<input placeholder="flow name" value="${p_flow_name_str}"></input>
					</div>
					
					<div id="upload_images">


					</div>

					
					<div id="upload_btn">upload image</div>

				</div>
				`);

			$("body").append(dialog);

			const new_img_id_int    = 1;
			const new_image_element = append_image(new_img_id_int, p_img_data_str);


			$(dialog).find("#upload_images").append(new_image_element)
			// $(dialog).find("#upload_images_panel #images").append(new_image_element);

			return dialog;
		}

		//-------------------------------
		// MULTI_IMAGE_UPLOAD - additional images upload at a time
		else {
			const dialog         = $("#upload_image_dialog");
			const new_img_id_int = parseInt($(dialog).find(".target_upload_image").last().attr("id")) + 1 // increment by one from last elements id/index
			
			/*
			$(dialog).find("#upload_images_panel #images").append(`
				<div class="image">
					<div class"remove_btn">
						<img src="https://assetspub.gloflow.com/assets/gf_sys/gf_close_btn.svg" draggable="false"></img>
					</div>
					<img id="${new_img_id_int}" class="target_upload_image" src="${p_img_data_str}"></img>
					<div class="upload_image_name_input">
						<input placeholder="image name"></input>
					</div>
				</div>
			`);
			*/

			const new_image_element = append_image(new_img_id_int, p_img_data_str);
			$(dialog).find("#upload_images").append(new_image_element);

			$(dialog).find("#upload_btn").text("upload images");
			
			return dialog;
		}

		//-------------------------------

		//-------------------------------------------------
		function append_image(p_new_image_dom_id_int :number, p_img_data_str :string) :HTMLElement {

			const image_element = $(`
				<div class="image">
					<div class="remove_btn">
						<img src="https://assetspub.gloflow.com/assets/gf_sys/gf_close_btn.svg" draggable="false"></img>
					</div>
					
					<img id="${p_new_image_dom_id_int}" class="target_upload_image" src="${p_img_data_str}"></img>

					<div class="upload_image_name_input">
						<input placeholder="image name"></input>
					</div>
				</div>
			`)[0];

			$(image_element).find(".remove_btn").on("click", ()=>{
				$(image_element).remove();

				// if there are no more images for upload, remove the whole dialog
				if ($(img_dialog).find(".image").length == 0) {
					$(img_dialog).remove();
					$("body").css("overflow-y", "visible"); // turn-on scroll
				}
			});

			return image_element;
		}

		//-------------------------------------------------
	}

	//-------------------------------------------------
	const img_dialog = get_image_dialog();
	

	//-------------------------------------------------
	function position_image_view() {

		// position the upload_image_dialog in view if the user scrolled
		const scroll_position_f = $(document).scrollTop();
		$("#upload_image_dialog").css("top", `${scroll_position_f}px`);
	}

	//-------------------------------------------------

	// ?????
	const this_image = $(img_dialog).find("#upload_images_panel").last();
	$(this_image).find("img").on("load", ()=>{

		$("body").css("overflow-y", "hidden"); // turn-off scroll
		position_image_view();
	});

	/*
	// reposition image_view on resize
	$(window).resize(function () {
		position_image_view();
	});
	*/

	//-------------------------------------------------
	// UPLOAD
	function upload() {

		const image_name_str    = $(this_image).find("#upload_image_name_input input").val();
		var image_flow_name_str = $(img_dialog).find("#upload_image_flow_name_input input").val();

		// if no image flow name was supplied then use the default flow
		if (image_flow_name_str.length == 0) {
			image_flow_name_str = p_flow_name_str;
		} 
		else {

			// LOCAL_STORAGE
			localStorage.setItem("gf:upload_flow_name_str", image_flow_name_str);
		}
		
		p_upload_activate_fun(image_name_str, image_flow_name_str, ()=>{

			// REMOVE_UPLOAD_DIALOG - when upload_activate function completes, remove the dialog
			$(img_dialog).remove();
			$("body").css("overflow-y", "visible"); // turn-on scroll
		});
	}

	//-------------------------------------------------

	$("body").keyup((p_event)=>{

		// ENTER_KEY
		if (p_event.which == 13) {

			console.log("enter key pressed...");

			// IMPORTANT!! - first time "enter" is pressed to upload an image, we want to unregister
			//               "enter" as the upload-activation key. 
			//               we also do this before upload begins, since it might last for some time
			//               and we dont want the user to keep pressing the enter button.
			$(this).unbind(p_event);

			upload();			
		}
		
		// ESC_KEY
		if (p_event.which == 27) {

			console.log("esc key pressed...");

			$(this).unbind(p_event);

			// REMOVE_DIALOG
			$("#upload_image_dialog").remove();
			$("body").css("overflow-y", "visible"); // turn-on scroll
		}
	});

	var uploading_in_progress_bool = false;
	$("#upload_image_dialog #upload_btn").on('click', ()=>{
		if (!uploading_in_progress_bool) {
			uploading_in_progress_bool = true;
			upload();
		}
	});

	$("#upload_image_dialog #background").on('click', ()=>{
		// REMOVE_UPLOAD_DIALOG - remove dialog on click on background
		$(img_dialog).remove();
		$("body").css("overflow-y", "visible"); // turn-on scroll
	})
}

//-------------------------------------------------
async function gf_upload__run(p_image_name_str :string,
	p_image_data_str       :string,
	p_image_format_str     :string,
	p_flows_names_str      :string,
	p_target_full_host_str :string) {
	return new Promise(async function(p_resolve_fun, p_reject_fun) {

		console.log(`UPLOAD_IMAGE - ${p_image_name_str} - ${p_image_format_str}`);

		const upload_start_f = performance.now();

		// UPLOAD__SEND_INIT
		const upload_map = await gf_upload__send_init(p_image_name_str,
			p_image_data_str,
			p_image_format_str,
			p_flows_names_str,
			p_target_full_host_str);
		
		const upload_gf_image_id_str = upload_map["upload_gf_image_id_str"];
		const presigned_url_str      = upload_map["presigned_url_str"];

		// UPLOAD_TO_S3
		const upload_transfer_duration_sec_f = await gf_upload__s3_put(presigned_url_str,
			p_image_data_str);

		// UPLOAD__SEND_COMPLETE
		await gf_upload__send_complete(upload_gf_image_id_str,
			p_flows_names_str,
			p_target_full_host_str);

		// UPLOAD__SEND_METRICS
		const upload_end_f = performance.now();
		const upload_duration_sec_f = upload_end_f - upload_start_f;

		await gf_upload__send_metrics(upload_duration_sec_f,
			upload_transfer_duration_sec_f,
			upload_gf_image_id_str,
			p_target_full_host_str);

		p_resolve_fun(upload_gf_image_id_str);
	});
}

//-------------------------------------------------
function gf_upload__send_init(p_image_name_str :string,
	p_image_data_str       :string,
	p_image_format_str     :string,
	p_flows_names_str      :string,
	p_target_full_host_str :string) :Promise<any> {
	
	return new Promise(function(p_resolve_fun, p_reject_fun) {

		const client_type_str = "browser";

		// auth_r=0 - dont redirect on auth fail, just return status
		const url_str = `${p_target_full_host_str}/v1/images/upload_init?imgf=${p_image_format_str}&imgn=${p_image_name_str}&f=${p_flows_names_str}&ct=${client_type_str}&auth_r=0`;
		$.get(url_str)
			.done((p_data_map) => {

				console.log("upload initialized...")
				console.log(p_data_map);

				if (p_data_map["status"] == "OK") {
					const upload_gf_image_id_str = p_data_map["data"]["upload_info_map"]["upload_gf_image_id_str"];
					const presigned_url_str      = p_data_map["data"]["upload_info_map"]["presigned_url_str"];

					console.log(`upload_gf_image_id - ${upload_gf_image_id_str}`);
					console.log(`presigned_url      - ${presigned_url_str}`);

					p_resolve_fun({
						"upload_gf_image_id_str": upload_gf_image_id_str,
						"presigned_url_str":      presigned_url_str,
					});
				}
				else {
					p_reject_fun(p_data_map["data"]);
				}
			})
			.fail(function(jqXHR, textStatus, errorThrown) {
				p_reject_fun(textStatus+" - "+errorThrown);
			});
	});
}

//-------------------------------------------------
function gf_upload__s3_put(p_presigned_url_str :string,
	p_image_data_str :string) :Promise<number> {
	
	return new Promise(function(p_resolve_fun, p_reject_fun) {

		const image_data_clean_str = p_image_data_str.replace("data:image/png;base64,", "");
		const image_data           = gf_base64_to_blob(image_data_clean_str, "image/png");
		const upload_start_f = performance.now();

		// AWS_S3
		$.ajax({
			"type": "PUT",
			"url":  p_presigned_url_str,
			"data": image_data, // p_image_data_str,

			// these are the headers that were included in the S3 URL signature generated by AWS.
			// so they have to be set to the same values for the received PUT request signature (on AWS side)
			// to match.
			"headers": {
				"content-type": "image/png",
				"x-amz-acl":    "public-read",
			},

			// jqeury is not to convert the image to form data
			processData: false,
			"success": ()=>{
				const upload_end_f = performance.now();
				const upload_transfer_duration_sec_f = upload_end_f - upload_start_f;
				p_resolve_fun(upload_transfer_duration_sec_f);
			},
			"error": (jqXHR, textStatus, errorThrown)=>{
				p_reject_fun(textStatus+" - "+errorThrown);
			}
		})
	});
}

//-------------------------------------------------
function gf_upload__send_metrics(p_upload_duration_sec_f :number,
	p_upload_transfer_duration_sec_f :number,
	p_upload_gf_image_id_str         :string,
	p_target_full_host_str           :string) :Promise<any> {
	
	return new Promise(function(p_resolve_fun, p_reject_fun) {

		const client_type_str = "browser";

		// auth_r=0 - dont redirect on auth fail, just return status
		const url_str = `${p_target_full_host_str}/v1/images/upload_metrics?imgid=${p_upload_gf_image_id_str}&ct=${client_type_str}&auth_r=0`;

		const data_map = {
			"upload_client_duration_sec_f":          p_upload_duration_sec_f,
			"upload_client_transfer_duration_sec_f": p_upload_transfer_duration_sec_f,
		};
		$.ajax({
			type: "POST",
			url:  url_str,
			data: JSON.stringify(data_map),
			//-------------------------------------------------
			"success": (p_data_map) => {

				console.log("upload metrics done...")
				console.log(p_data_map);
				p_resolve_fun({});
			},

			//-------------------------------------------------
			error: (jqXHR, p_text_status_str)=>{
				p_reject_fun(p_text_status_str);
			}
		})
	});
}

//-------------------------------------------------
function gf_upload__send_complete(p_upload_gf_image_id_str :string,
	p_flows_names_str      :string,
	p_target_full_host_str :string) :Promise<any> {

	return new Promise(function(p_resolve_fun, p_reject_fun) {
		console.log("AWS S3 PUT upload done...")

		// auth_r=0 - dont redirect on auth fail, just return status
		const url_str = `${p_target_full_host_str}/v1/images/upload_complete?imgid=${p_upload_gf_image_id_str}&f=${p_flows_names_str}&auth_r=0`;

		$.ajax({
			method: "POST",
			"url":  url_str,
			//-------------------------------------------------
			"success": (p_data_map) => {

				console.log("upload complete...")
				console.log(p_data_map);
				p_resolve_fun({});
			},

			//-------------------------------------------------
			error: (jqXHR, p_text_status_str)=>{
				p_reject_fun(p_text_status_str);
			}
		})
	});
}

//-------------------------------------------------
function gf_base64_to_blob(p_base64_str :string, p_mime_type_str :string)  {

	const slice_size_int = 1024;

	// window.atob() - decodes a base-64 encoded string
	const decoded_str     = window.atob(p_base64_str);
	const byte_arrays_lst = [];
	
	// pack individual slices as uint8 byte_arrays, and then pack those into
	// a array themselves - for loading into a Blob. 
	for (var offset_i = 0, len = decoded_str.length; offset_i < len; offset_i += slice_size_int) {
		const slice = decoded_str.slice(offset_i, offset_i + slice_size_int);

		var bytes_lst = new Array(slice.length);
		for (var i = 0; i < slice.length; i++) {
			bytes_lst[i] = slice.charCodeAt(i);
		}

		// BYTE_ARRAY
		const byte_array = new Uint8Array(bytes_lst);
		byte_arrays_lst.push(byte_array);
	}

	// https://developer.mozilla.org/en-US/docs/Web/API/Blob
	const blob = new Blob(byte_arrays_lst, {type: p_mime_type_str});
	return blob;
}