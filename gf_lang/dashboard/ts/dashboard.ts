/*
GloFlow application and media management/publishing platform
Copyright (C) 2026 Ivan Trajkovic

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

// gf_lang dev dashboard - the smallest possible IDE.
// left: editable, syntax highlighted gf code (+ a picker for the programs in gf_lang/examples/).
// right: the program run by the interpreter (Go WASM) driving the three.js engine.

///<reference path="./../../../web/src/d/jquery.d.ts" />

import * as gf_engine from "../../ts/engine/gf_engine";
import * as gf_syntax from "./gf_syntax";

declare var Go;           // wasm_exec.js
declare var gf_lang_run;  // set by the WASM module

const examples_select = document.getElementById("examples") as HTMLSelectElement;
const run_btn         = document.getElementById("run")      as HTMLButtonElement;
const status_el       = document.getElementById("status");
const code_el         = document.getElementById("code")     as HTMLTextAreaElement;
const highlight_el    = document.getElementById("highlight");
const render_el       = document.getElementById("render");

var engine_api_map = null;

//-------------------------------------------------
async function main() {

	// EXAMPLES
	const response   = await fetch("/v1/lang/examples");
	const examples   = (await response.json())["examples_lst"];
	for (const e of examples) {
		examples_select.add(new Option(e["name_str"], e["code_str"]));
	}

	// WASM - the gf_lang interpreter
	set_status("loading interpreter...");
	const go     = new Go();
	const result = await WebAssembly.instantiateStreaming(fetch("gf_lang_web.wasm"), go.importObject);
	go.run(result.instance); // never returns (the Go main() blocks), so not awaited

	// UI
	examples_select.onchange = ()=> { set_code(examples_select.value); run(); };
	run_btn.onclick          = run;
	code_el.oninput          = update_highlight;
	code_el.onscroll         = sync_scroll;
	code_el.onkeydown        = (p_e)=> {
		if (p_e.key === "Enter" && (p_e.ctrlKey || p_e.metaKey)) { p_e.preventDefault(); run(); }
		if (p_e.key === "Tab") { p_e.preventDefault(); document.execCommand("insertText", false, "  "); }
	};

	if (examples.length > 0) {
		set_code(examples_select.value);
		run();
	} else {
		set_status("no examples found");
	}
}

//-------------------------------------------------
function set_code(p_code_str :string) {
	code_el.value = p_code_str;
	update_highlight();
}

// the highlighted <pre> is shown under a transparent <textarea>, so they have to show the same text
function update_highlight() {
	// trailing newline - otherwise the last (empty) line of the textarea is not scrollable in the <pre>
	highlight_el.innerHTML = gf_syntax.highlight(code_el.value) + "\n ";
	sync_scroll();
}

function sync_scroll() {
	highlight_el.scrollTop  = code_el.scrollTop;
	highlight_el.scrollLeft = code_el.scrollLeft;
}

function set_status(p_msg_str :string, p_error_bool = false) {
	status_el.textContent = p_msg_str;
	status_el.classList.toggle("error", p_error_bool);
}

//-------------------------------------------------
function run() {

	// PARSE
	var program_ast_lst;
	try {
		program_ast_lst = parse_program(code_el.value);
	} catch (e) {
		set_status("parse error: " + e.message, true);
		return;
	}

	// previous run
	if (engine_api_map) engine_api_map["dispose_fun"]();
	engine_api_map = null;

	// RUN - the interpreter calls the extern_api functions, which drive the engine.
	//       returns "" or an error message
	const error_str = gf_lang_run(program_ast_lst, extern_api_map);
	if (error_str) set_status("run error: " + error_str, true);
	else           set_status("ok");
}

//-------------------------------------------------
// a program is a list of expressions. the example files are either that list in "[ ]",
// or just the expressions separated by commas.
function parse_program(p_code_str :string) {

	const json_str = gf_syntax.strip_comments(p_code_str);
	var parsed;
	try {
		parsed = JSON.parse(json_str);
	} catch (e) {
		parsed = JSON.parse("[" + json_str + "]");
	}
	// a single bare expression, ex. ["cube"]
	if (!parsed.every(Array.isArray)) parsed = [parsed];
	return parsed;
}

//-------------------------------------------------
const extern_api_map = {
	"init_engine_fun":   (p_shader_defs_map)=> {
		engine_api_map = gf_engine.init(p_shader_defs_map, render_el);

		// the engine default is a black background, and most programs draw black cubes (so nothing would be visible).
		// a program that sets its own color-background overrides this.
		engine_api_map["set_state_fun"]({
			"property_name_str":          "color_background_rgb_lst",
			"color_background_rgb_lst": [0.9, 0.9, 0.9],
		});
	},
	"set_state_fun":     (p_state_change_map)=> engine_api_map["set_state_fun"](p_state_change_map),
	"create_cube_fun":   (...p_args)=> engine_api_map["create_cube_fun"](...p_args),
	"create_sphere_fun": (...p_args)=> engine_api_map["create_sphere_fun"](...p_args),
	"create_line_fun":   (...p_args)=> engine_api_map["create_line_fun"](...p_args),
	"animate_fun":       ()=> {},
};

main();
