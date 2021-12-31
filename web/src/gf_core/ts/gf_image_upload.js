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

//-------------------------------------------------
var gf_upload__init = gf_upload__init;
function gf_upload__init(p_target_full_host_str) {

	// console.log("UPLOAD INITIALIZED")
	document.onpaste = function(p_paste_event) {



		const items = (p_paste_event.clipboardData || p_paste_event.originalEvent.clipboardData).items;


		console.log("paste");
		// console.log(p_paste_event.clipboardData);
		// console.log(p_paste_event.originalEvent.clipboardData);
		console.log(JSON.stringify(items)); // will give you the mime types

		for (index in items) {
			const item = items[index];
			if (item.kind === 'file') {

				const blob   = item.getAsFile();
				const reader = new FileReader();

				reader.onload = function(p_event) {
					
					console.log("data loaded");
					
					// result attribute contains the data as a data: URL representing
					// the file's data as a base64 encoded string.
					const img_data_str = p_event.target.result;
					
					
					// the beginning of data URL strings example:
					// "data:image/png;base64".
					// IMPORTANT!! - it seems all images are of "png" format when pasted in.
					const image_format_str = img_data_str.split(";")[0].replace("data:image/", "")
					console.log(`image_format_str - ${image_format_str}`);

					// VIEW_IMAGE
					gf_upload__view_img(img_data_str,

						// UPLOAD_ACTIVATE_FUN
						(p_image_name_str,
						p_flows_names_str,
						p_on_upload_complete_fun)=>{


							// UPLOAD_IMAGE
							gf_upload__run(p_image_name_str,
								img_data_str,
								image_format_str,
								p_flows_names_str,
								p_target_full_host_str,
								()=>{
									p_on_upload_complete_fun();
								});
						});
				};

				// reader.readAsBinaryString(blob);
				reader.readAsDataURL(blob);
			}
		}
	}
}

//-------------------------------------------------
function gf_upload__view_img(p_img_data_str,
	p_upload_activate_fun) {
	
	//-----------------
	// FLOW_NAME
	// get old value from localStorage if it exists, if it doesnt use the default
	const previous_flow_name_str = localStorage.getItem("gf:upload_flow_name_str");
	var default_flow_name_str    = "general";
	if (previous_flow_name_str != null) {
		default_flow_name_str = previous_flow_name_str;
	}

	//-----------------

	//-------------------------------------------------
	function get_image_dialog() {

		// first image
		if ($("#upload_image_dialog").length == 0) {

			const img_dialog = $(`
				<div id="upload_image_dialog">
					<div id="background"></div>

					<div id="upload_images_detail">
						<div id="upload_images">

							<!-- IMAGE_PANEL -->
							<div class="upload_image_panel">
								<img id="1" class="target_upload_image" src='${p_img_data_str}'></img>
								<div id="upload_image_name_input">
									<input placeholder="image name"></input>
								</div>
								<div id="upload_image_flow_name_input">
									<input placeholder="flow name" value="${default_flow_name_str}"></input>
								</div>
							</div>
							<!-- ----------- -->

						</div>
						<div id="upload_btn">upload image</div>
					</div>

				</div>`);
			$("body").append(img_dialog);
			return img_dialog;
		}

		// additional images (multi-image upload)
		else {
			const img_dialog     = $("#upload_image_dialog");
			const new_img_id_int = parseInt($(img_dialog).find(".target_upload_image").last().attr("id")) + 1 // increment by one from last elements id/index
			
			$(img_dialog).find("#upload_images").append(`
				<!-- IMAGE_PANEL -->
				<div class="upload_image_panel">
					<img id="${new_img_id_int}" class="target_upload_image" src="${p_img_data_str}"></img>
					<div id="upload_image_name_input">
						<input placeholder="image name"></input>
					</div>
					<div id="upload_image_flow_name_input">
						<input placeholder="flow name"></input>
					</div>
				</div>
				<!-- ----------- -->`);
			
			$(img_dialog).find("#upload_btn").text("upload images");
			return img_dialog;
		}
	}

	//-------------------------------------------------
	const img_dialog = get_image_dialog();
	

	//-------------------------------------------------
	function position_image_view() {

		// position the upload_image_dialog in view if the user scrolled
		const scroll_position_f = $(document).scrollTop();
		$("#upload_image_dialog").css("top", `${scroll_position_f}px`);

		// REPOSITION IMAGES_DETAIL
		const images_detail = $(img_dialog).find("#upload_images_detail");

		const image_x = Math.max(0, (($(window).width() - $(images_detail).outerWidth()) / 2));
		const image_y = Math.max(0, (($(window).height() - $(images_detail).outerHeight()) / 2));

		$(images_detail).css("left", image_x+"px");
		$(images_detail).css("top",  image_y+"px");
	}

	//-------------------------------------------------
	const this_image = $(img_dialog).find(".upload_image_panel").last();
	$(this_image).find("img").on("load", ()=>{

		$("body").css("overflow-y", "hidden"); // turn-off scroll
		position_image_view();
	});

	// reposition image_view on resize
	$(window).resize(function () {
		position_image_view();
	});

	//-------------------------------------------------
	function upload() {
		const image_name_str    = $(this_image).find("#upload_image_name_input input").val();
		var image_flow_name_str = $(this_image).find("#upload_image_flow_name_input input").val();

		// if no image flow name was supplied then use the default flow ("general")
		if (image_flow_name_str.length == 0) {
			image_flow_name_str = "general";
		} 
		else {
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
			// IMPORTANT!! - first time "enter" is pressed to upload an image, we want to unregister
			//               "enter" as the upload-activation key. 
			//               we also do this before upload begins, since it might last for some time
			//               and we dont want the user to keep pressing the enter button.
			$(this).unbind(p_event);

			upload();			
		}
		
		// ESC_KEY
		if (p_event.which == 27) {
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
function gf_upload__run(p_image_name_str,
	p_image_data_str,
	p_image_format_str,
	p_flows_names_str,
	p_target_full_host_str,
	p_on_complete_fun) {
	console.log(`UPLOAD_IMAGE - ${p_image_name_str} - ${p_image_format_str}`);

	// UPLOAD__SEND_INIT
	gf_upload__send_init(p_image_name_str,
		p_image_data_str,
		p_image_format_str,
		p_flows_names_str,
		p_target_full_host_str,
		(p_upload_gf_image_id_str, p_presigned_url_str)=>{

			// UPLOAD_TO_S3
			gf_upload__s3_put(p_presigned_url_str,
				p_image_data_str,
				()=>{

					// UPLOAD__SEND_COMPLETE
					gf_upload__send_complete(p_upload_gf_image_id_str, 
						p_target_full_host_str,
						()=>{
							p_on_complete_fun();
						});
				});
		});
}

//-------------------------------------------------
function gf_upload__send_init(p_image_name_str,
	p_image_data_str,
	p_image_format_str,
	p_flows_names_str,
	p_target_full_host_str,
	p_on_complete_fun) {

	// UPLOAD_INIT
	const url_str = `${p_target_full_host_str}/images/v1/upload_init?imgf=${p_image_format_str}&imgn=${p_image_name_str}&f=${p_flows_names_str}&ct=browser`;
	$.ajax({
		method: "GET",
		"url":  url_str,
		//-------------------------------------------------
		"success": (p_data_map) => {

			console.log("upload initialized...")
			console.log(p_data_map);

			const upload_gf_image_id_str = p_data_map["data"]["upload_info_map"]["upload_gf_image_id_str"];
			const presigned_url_str      = p_data_map["data"]["upload_info_map"]["presigned_url_str"];

			console.log(`upload_gf_image_id - ${upload_gf_image_id_str}`);
			console.log(`presigned_url      - ${presigned_url_str}`);


			p_on_complete_fun(upload_gf_image_id_str,
				presigned_url_str);
		}

		//-------------------------------------------------
	});
}

//-------------------------------------------------
function gf_upload__s3_put(p_presigned_url_str,
	p_image_data_str,
	p_on_complete_fun) {

	const image_data_clean_str = p_image_data_str.replace("data:image/png;base64,", "");
	const image_data           = gf_base64_to_blob(image_data_clean_str, "image/png");

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
			p_on_complete_fun();
		}
	})
}

//-------------------------------------------------
function gf_upload__send_complete(p_upload_gf_image_id_str,
	p_target_full_host_str,
	p_on_complete_fun) {

	console.log("AWS S3 PUT upload done...")
	const url_str = `${p_target_full_host_str}/images/v1/upload_complete?imgid=${p_upload_gf_image_id_str}`;

	$.ajax({
		method: "POST",
		"url":  url_str,
		//-------------------------------------------------
		"success": (p_data_map) => {

			console.log("upload complete...")
			console.log(p_data_map);
			p_on_complete_fun();
		}

		//-------------------------------------------------
	})
}

//-------------------------------------------------
function gf_base64_to_blob(p_base64_str, p_mime_type_str)  {

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