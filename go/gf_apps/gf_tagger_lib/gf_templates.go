/*
GloFlow application and media management/publishing platform
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

package gf_tagger_lib

import (
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type gf_templates struct {
	tag_objects__tmpl                   *template.Template
	tag_objects__subtemplates_names_lst []string

	bookmarks__tmpl                   *template.Template
	bookmarks__subtemplates_names_lst []string
}

//-------------------------------------------------

func tmpl__load(p_templates_paths_map map[string]string,
	pRuntimeSys *gf_core.RuntimeSys) (*gf_templates, *gf_core.GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_templates.tmpl__load()")

	main_template_filepath_str := p_templates_paths_map["gf_tag_objects"]
	tag_objects__tmpl, subtemplates_names_lst, gf_err := gf_core.TemplatesLoad(main_template_filepath_str, pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}



	bookmarks_template_filepath_str := p_templates_paths_map["gf_bookmarks"]
	bookmarks__tmpl, bookmarks_subtemplates_names_lst, gf_err := gf_core.TemplatesLoad(bookmarks_template_filepath_str, pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}

	gf_templates := &gf_templates{
		tag_objects__tmpl:                   tag_objects__tmpl,
		tag_objects__subtemplates_names_lst: subtemplates_names_lst,

		bookmarks__tmpl:                   bookmarks__tmpl,
		bookmarks__subtemplates_names_lst: bookmarks_subtemplates_names_lst,
	}
	return gf_templates, nil
}