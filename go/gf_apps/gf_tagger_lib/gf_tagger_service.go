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
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
)

//-------------------------------------------------

func InitService(pTemplatesPathsMap map[string]string,
	pImagesJobsMngr gf_images_jobs_core.JobsMngr,
	pHTTPmux        *http.ServeMux,
	pRuntimeSys     *gf_core.RuntimeSys) {
	
	//------------------------
	// DB_INDEXES
	gfErr := DBmongoIndexInit(pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}
	
	//------------------------
	// STATIC FILES SERVING
	urlBaseStr      := "/tags"
	localDirPathStr := "./static"
	gf_core.HTTPinitStaticServingWithMux(urlBaseStr,
		localDirPathStr,
		pHTTPmux,
		pRuntimeSys)

	//------------------------
	
	gfErr = initHandlers(pTemplatesPathsMap,
		pImagesJobsMngr,
		pHTTPmux,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	//------------------------
}

//-------------------------------------------------

func RunService(pPortStr string,
	p_mongodb_host_str    string,
	p_mongodb_db_name_str string,
	p_init_done_ch        chan bool,
	pLogFun               func(string, string)) {

	pLogFun("INFO", "")
	pLogFun("INFO", " >>>>>>>>>>> STARTING GF_TAGGER SERVICE")
	pLogFun("INFO", "")
	
	runtimeSys := &gf_core.RuntimeSys{
		ServiceNameStr: "gf_tagger",
		LogFun:         pLogFun,
	}

	mongo_db, _, gfErr := gf_core.MongoConnectNew(p_mongodb_host_str, p_mongodb_db_name_str, nil, runtimeSys)
	if gfErr != nil {
		panic(-1)
	}
	runtimeSys.Mongo_db   = mongo_db 
	runtimeSys.Mongo_coll = mongo_db.Collection("data_symphony")

	//----------------------
	http_mux := http.NewServeMux()

	templates_dir_paths_map := map[string]string{
		"gf_tag_objects": "./templates/gf_tag_objects/gf_tag_objects.html",
	}

	// FIX!! - jobs_mngr shouldnt be used here. when gf_tagger service is run in a separate
	//         process from gf_images service, jobs_mngr can only be reaeched via HTTP or some other
	//         transport mechanism (not via Go messages as a goroutine).
	var jobs_mngr gf_images_jobs_core.JobsMngr

	InitService(templates_dir_paths_map,
		jobs_mngr,
		http_mux,
		runtimeSys)

	//----------------------
	// IMPORTANT!! - signal to user that server in this goroutine is ready to start listening 
	if p_init_done_ch != nil {
		p_init_done_ch <- true
	}

	//----------------------

	pLogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	pLogFun("INFO", "STARTING HTTP SERVER - PORT - "+pPortStr)
	pLogFun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	err := http.ListenAndServe(":"+pPortStr, nil)
	if err != nil {
		msg_str := "cant start listening on port - "+pPortStr
		pLogFun("ERROR", msg_str)
		panic(msg_str)
	}
}