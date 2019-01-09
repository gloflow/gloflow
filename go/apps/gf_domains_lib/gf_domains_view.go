package gf_domains_lib

import (
	"net/http"
	"text/template"
	"gf_core"
)
//--------------------------------------------------
func domains_browser__render_template(p_domains_lst []Domain,
						p_tmpl        *template.Template,
						p_resp        http.ResponseWriter,
						p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_domains_view.domains_browser__render_template()")

	sys_release_info := gf_core.Get_sys_relese_info(p_runtime_sys)

	type tmpl_data struct {
		Domains_lst      []Domain
		Sys_release_info gf_core.Sys_release_info
	}

	err := p_tmpl.Execute(p_resp,tmpl_data{
		Domains_lst     :p_domains_lst,
		Sys_release_info:sys_release_info,
	})

	if err != nil {
		gf_err := gf_core.Error__create("failed to render the domains_browser template",
            "template_render_error",
            &map[string]interface{}{"domains_lst":p_domains_lst,},
            err,"gf_domains_lib",p_runtime_sys)
		return gf_err
	}

	return nil
}