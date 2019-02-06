/*
GloFlow media management/publishing system
Copyright (C) 2019 Ivan Trajkovic

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

///<reference path="../../../../d/jquery.d.ts" />

declare var Caman;
//-------------------------------------------------
export function init(p_target_image_div_element, p_log_fun) {
	p_log_fun('FUN_ENTER', 'gf_image_editor.init()');

	const target_image = $(p_target_image_div_element).find('img')[0];
	var width_int  = target_image.clientWidth;
	var height_int = target_image.clientHeight;

	console.log('img width  - '+width_int);
	console.log('img height - '+height_int);


	const container = $(`
		<div class='gf_image_editor'>
			<div class='open_editor_btn'>editor</div>
		</div>`);

	$(p_target_image_div_element).append(container);

	//-------------------------------------------------
	function create_pane() {
		p_log_fun('FUN_ENTER', 'gf_image_editor.init().create_pane()');

		const editor_pane = $(`
			<div class='editor_pane'>
				<div class='close_btn'>x</div>
				<div class='save_btn'>save</div>

				<canvas class='modified_image_canvas' width="`+width_int+`" height="`+height_int+`"></canvas>

				<div class="slider_input">
					<form>
						<div>
							<input id="contrast" name="contrast" type="range" min="-100" max="100" value="0">
							<label for="contrast">contrast</label>
						</div>
						
						<div>
							<input id="brightness" name="brightness" type="range" min="-100" max="100" value="0">
							<label for="brightness">brightness</label>
						</div>
						
						<div>
							<input id="saturation" name="saturation" type="range" min="-100" max="100" value="0">
							<label for="saturation">saturation</label>
						</div>

						<div>
							<input id="sharpen" name="sharpen" type="range" min="0" max="100" value="0">
							<label for="sharpen">sharpen</label>
						</div>

						<div>
							<input id="sepia" name="sepia" type="range" min="0" max="100" value="0">
							<label for="sepia">sepia</label>
						</div>

						<div>
							<input id="noise" name="noise" type="range" min="0" max="100" value="0">
							<label for="noise">noise</label>
						</div>

						<div>
							<input id="hue" name="hue" type="range" min="0" max="100" value="0">
							<label for="hue">hue</label>
						</div>
					</form>
				</div>
			</div`);

		$(editor_pane).find('input[type=range]').change(apply_filters_fun);
		//-------------------------------------------------
		function apply_filters_fun() {
			
			const contrast   = parseInt($('#contrast').val());
			const brightness = parseInt($('#brightness').val());
			const saturation = parseInt($('#saturation').val());
			const sharpen    = parseInt($('#sharpen').val());
			const sepia      = parseInt($('#sepia').val());
			const noise      = parseInt($('#noise').val());
			const hue        = parseInt($('#hue').val());

			Caman('.editor_pane canvas', target_image, function() {
				this.revert(false);
				
				this.contrast(contrast);
				this.brightness(brightness);
				this.saturation(saturation);
				this.sharpen(sharpen);
				this.sepia(sepia);
				this.noise(noise);
				this.hue(hue);
				this.render(()=>console.log('filter applied'));
			});
		}
		//-------------------------------------------------
			
		const canvas = $(editor_pane).find('canvas')[0];
		Caman(canvas, $(target_image).attr('src'), function () {
			this.render();
		});

		//-------------
		//SAVE_MODIFIED_IMAGE
		$(editor_pane).find('.save_btn').on('click', ()=>{
			save_modified_image(editor_pane);
		});
		//-------------

		return editor_pane;
	}
	//-------------------------------------------------
	function save_modified_image(p_editor_pane) {
		p_log_fun('FUN_ENTER', 'gf_image_editor.init().save_modified_image()');

		const canvas            = $(p_editor_pane).find('.modified_image_canvas')[0];
		const canvas_base64_str = (canvas as HTMLCanvasElement).toDataURL();

		console.log(canvas_base64_str);

		http_save(canvas_base64_str)
	}
	//-------------------------------------------------

	var opened_bool = false;
	$(container).find('.open_editor_btn').on('click', ()=>{

		if (opened_bool) {
			return;
		}

		const editor_pane = create_pane()

		$(editor_pane).find('.close_btn').on('click',()=>{
			$(editor_pane).remove();
			opened_bool = false;
		});

		$(container).append(editor_pane);

		opened_bool = true;
	});

	return container;
}
//-------------------------------------------------
function http_save(p_canvas_base64_str) {

	$.ajax({
		type: "POST",
		url:  "/images/editor/save",
		data: { 
			imgBase64: p_canvas_base64_str
		}
		}).done((p_data)=>{});
}