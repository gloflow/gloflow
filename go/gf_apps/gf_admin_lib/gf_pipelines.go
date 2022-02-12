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

package gf_admin_lib

import (
	"text/template"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib"
)

//------------------------------------------------
func Pipeline__mfa_confirm(p_extern_htop_value_str string,
	p_secret_key_base32_str string,
	p_ctx                   context.Context,
	p_runtime_sys           *gf_core.Runtime_sys) (bool, *gf_core.GF_error) {




	htop_value_str, gf_err := gf_identity_lib.TOTP_generate_value(p_secret_key_base32_str,
		p_runtime_sys)
	if gf_err != nil {
		return false, gf_err
	}




	if p_extern_htop_value_str == htop_value_str {
		return true, nil
	} else {
		return false, nil
	}
	return false, nil
}

//------------------------------------------------
func Pipeline__render_login(p_mfa_confirm_bool bool,
	p_tmpl                   *template.Template,
	p_subtemplates_names_lst []string,
	p_ctx                    context.Context,
	p_runtime_sys            *gf_core.Runtime_sys) (string, *gf_core.GF_error) {

	template_rendered_str, gf_err := view__render_template_login(p_mfa_confirm_bool,
		p_tmpl,
		p_subtemplates_names_lst,
		p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}

	return template_rendered_str, nil
}

//------------------------------------------------
func Pipeline__render_dashboard(p_tmpl *template.Template,
	p_subtemplates_names_lst []string,
	p_ctx                    context.Context,
	p_runtime_sys            *gf_core.Runtime_sys) (string, *gf_core.GF_error) {

	template_rendered_str, gf_err := view__render_template_dashboard(p_tmpl,
		p_subtemplates_names_lst,
		p_runtime_sys)
	if gf_err != nil {
		return "", gf_err
	}

	return template_rendered_str, nil
}