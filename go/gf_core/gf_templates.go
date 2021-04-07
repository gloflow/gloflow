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

package gf_core

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"io/ioutil"
	"text/template"
)

//-------------------------------------------------
func Templates__load(p_main_template_filepath_str string,
	// p_templates_dir_path_str string,
	p_runtime_sys *Runtime_sys) (*template.Template, []string, *Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_templates.Templates__load()")

	template_filename_str := filepath.Base(p_main_template_filepath_str)
	template_dir_path_str := filepath.Dir(p_main_template_filepath_str)

	//---------------------
	// SUB_TEMPLATES - templates that are imported into the main template
	subtemplates_dir_path_str   := fmt.Sprintf("%s/subtemplates", template_dir_path_str)
	subtemplates_names_lst      := []string{}
	subtemplates_file_paths_lst := []string{}

	// load subtemplates if the subtemplates/ dir exists
	if _, err := os.Stat(subtemplates_dir_path_str); !os.IsNotExist(err) {

		files_lst, err := ioutil.ReadDir(subtemplates_dir_path_str)
		if err != nil {
			gf_err := Error__create("failed to parse a template",
				"dir_list_error",
				map[string]interface{}{"subtemplates_dir_path_str": subtemplates_dir_path_str,},
				err, "gf_core", p_runtime_sys)
			return nil, nil, gf_err
		}

		for _, f := range files_lst {
			filename_str := f.Name()
			if strings.HasSuffix(filename_str, ".html") {
				subtemplates_names_lst      = append(subtemplates_names_lst, strings.Split(filename_str, ".")[0])
				subtemplates_file_paths_lst = append(subtemplates_file_paths_lst, fmt.Sprintf("%s/%s", subtemplates_dir_path_str, filename_str))
			}
		}
	}

	//---------------------
	// TEMPLATES
	// main_template_path_str := fmt.Sprintf("%s/%s", p_templates_dir_path_str, p_main_template_filename_str)
	templates_paths_lst := append([]string{p_main_template_filepath_str,}, subtemplates_file_paths_lst...)

	// IMPORTANT!! - load several template files into a single template name
	main__tmpl, err := template.New(template_filename_str).ParseFiles(templates_paths_lst...)
	if err != nil {
		gf_err := Error__create("failed to parse a template",
			"template_create_error",
			map[string]interface{}{"main_template_filepath_str": p_main_template_filepath_str,},
			err, "gf_core", p_runtime_sys)
		return nil, nil, gf_err
	}
	
	//---------------------

	return main__tmpl, subtemplates_names_lst, nil
}