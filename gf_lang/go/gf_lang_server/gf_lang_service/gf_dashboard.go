/*
GloFlow application and media management/publishing platform
Copyright (C) 2023 Ivan Trajkovic

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

package gf_lang_service

import (
	"os"
	"sort"
	"strings"
	"net/http"
	"path/filepath"
	"encoding/json"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
// the dev dashboard (gf_lang/dashboard) is a static web app (built into gf_lang/build/ by
// run_build_gf_lang.sh) that runs gf_lang programs in the browser (WASM interpreter + three.js).
// the server serves its files at "/" and the list of example programs at /v1/lang/examples.

type GFexample struct {
	NameStr string `json:"name_str"`
	CodeStr string `json:"code_str"`
}

//-------------------------------------------------

func initDashboard(pHTTPmux *http.ServeMux,
	pConfig     *GFconfig,
	pRuntimeSys *gf_core.RuntimeSys) {

	// STATIC_FILES
	pHTTPmux.Handle("/", http.FileServer(http.Dir(pConfig.DashboardDirStr)))

	// EXAMPLES - read on every request, so edited/new example files show up on a page reload
	pHTTPmux.HandleFunc("/v1/lang/examples", func(pResp http.ResponseWriter, pReq *http.Request) {

		examplesLst, err := loadExamples(pConfig.ExamplesDirStr)
		if err != nil {
			http.Error(pResp, err.Error(), http.StatusInternalServerError)
			return
		}
		pResp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(pResp).Encode(map[string]interface{}{"examples_lst": examplesLst})
	})
}

//-------------------------------------------------

func loadExamples(pDirPathStr string) ([]*GFexample, error) {

	filesLst, err := filepath.Glob(filepath.Join(pDirPathStr, "*.gf"))
	if err != nil {
		return nil, err
	}
	sort.Strings(filesLst)

	examplesLst := []*GFexample{}
	for _, fileStr := range filesLst {
		codeBytesLst, err := os.ReadFile(fileStr)
		if err != nil {
			return nil, err
		}
		nameStr := strings.TrimSuffix(filepath.Base(fileStr), ".gf")
		examplesLst = append(examplesLst, &GFexample{nameStr, string(codeBytesLst)})
	}
	return examplesLst, nil
}
