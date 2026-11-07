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

// syntax highlighting for gf_lang code.
// gf_lang programs are JSON with "//" comments, so a regex tokenizer is all that is needed.

// words (as they appear inside strings) that get their own color
const WORD_CLASSES = new Map<string, string>();
const add_words = (p_class_str :string, p_words_str :string)=> {
	for (const w of p_words_str.split(" ")) WORD_CLASSES.set(w, p_class_str);
};
add_words("kw",    "lang_v set rule shader if repeat max uniform animate return print push pop list make map len rand start * uniforms vertex fragment");
add_words("shape", "cube sphere line");
add_words("prop",  "x y z rx ry rz sx sy sz color color-background material material_prop scale coord_origin rotation_pivot current_pos iters_max wireframe shader_uniform");

// groups: 1 comment, 2 string (possibly unterminated), 3 number, 4 true|false|null, 5 bracket, anything else: 1 char
const TOKEN_RE = /(\/\/[^\n]*)|("(?:[^"\\\n]|\\.)*"?)|(-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?)|\b(true|false|null)\b|([\[\]])|[\s\S]/g;

//-------------------------------------------------
export function highlight(p_code_str :string) :string {

	return p_code_str.replace(TOKEN_RE, (p_match_str, comment, str, num, lit, bracket)=>{

		var class_str = null;
		if      (comment) class_str = "comment";
		else if (str)     class_str = string_class(str);
		else if (num)     class_str = "number";
		else if (lit)     class_str = "literal";
		else if (bracket) class_str = "bracket";

		const escaped_str = p_match_str.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
		return class_str ? `<span class="${class_str}">${escaped_str}</span>` : escaped_str;
	});
}

//-------------------------------------------------
function string_class(p_str :string) :string {

	const inner_str = p_str.slice(1, -1);
	if (inner_str.startsWith("$")) return "var";
	return WORD_CLASSES.get(inner_str) || "string";
}

//-------------------------------------------------
// removes the "//" comments (the same tokenizer is used, so "//" inside strings is left alone)
export function strip_comments(p_code_str :string) :string {

	return p_code_str.replace(TOKEN_RE, (p_match_str, comment)=> comment ? "" : p_match_str);
}

//-------------------------------------------------
// .gfl (indentation syntax) highlighting.
// groups: 1 comment, 2 """ string, 3 string, 4 number, 5 true|false, 6 word, 7 punctuation, anything else: 1 char
const GFL_TOKEN_RE = /(\/\/[^\n]*)|("""[\s\S]*?(?:"""|$))|("(?:[^"\\\n]|\\.)*"?)|(-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?(?![\w$#-]))|\b(true|false)\b|([A-Za-z_$#][\w$#-]*)|([=|:(){}])|[\s\S]/g;

export function highlight_gfl(p_code_str :string) :string {

	return p_code_str.replace(GFL_TOKEN_RE, (p_match_str, comment, triple, str, num, lit, word, punct)=>{

		var class_str = null;
		if      (comment)       class_str = "comment";
		else if (triple || str) class_str = "string";
		else if (num)           class_str = "number";
		else if (lit)           class_str = "literal";
		else if (punct)         class_str = "bracket";
		else if (word) {
			if (word.startsWith("$")) class_str = "var";
			else                      class_str = WORD_CLASSES.get(word) || null;
		}

		const escaped_str = p_match_str.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
		return class_str ? `<span class="${class_str}">${escaped_str}</span>` : escaped_str;
	});
}
