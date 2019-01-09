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

package gf_images_lib

import (
	"fmt"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_gif_lib"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_image_editor"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_images_jobs"
)
//-------------------------------------------------
func Run_service(p_port_str string,
	p_mongodb_host_str                           string,
	p_mongodb_db_name_str                        string,
	p_images_store_local_dir_path_str            string,
	p_images_thumbnails_store_local_dir_path_str string,
	p_images_main_s3_bucket_name_str             string,
	p_templates_dir_paths_map                    map[string]interface{},
	p_init_done_ch                               chan bool,
	p_log_fun                                    func(string,string)) {
	p_log_fun("FUN_ENTER","gf_images_service.Run_service()")

	p_log_fun("INFO","")
	p_log_fun("INFO"," >>>>>>>>>>> STARTING GF_IMAGES SERVICE")
	p_log_fun("INFO","")
	logo_str := `.           ..         .         
      .         .            .          .       .
            .         ..xxxxxxxxxx....               .       .             .
    .             MWMWMWWMWMWMWMWMWMWMWMWMW                       .
              IIIIMWMWMWMWMWMWMWMWMWMWMWMWMWMttii:        .           .
 .      IIYVVXMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWxx...         .           .
     IWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMx..
   IIWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWMWNMWMWMWMWMWMWMWMWMWMWMWMWMx..        .
    ""MWMWMWMWMWM"""""""".  .:..   ."""""MWMWMWMWMWMWMWMWMWMWMWMWMWti.
 .     ""   .   .: . :. : .  . :.  .  . . .  """"MWMWMWMWMWMWMWMWMWMWMWMWMti=
        . .   : . :   .  .'.' '....xxxxx...,'. '   ' ."""YWMWMWMWMWMWMWMWMWMW+
     ; .  .  . : . .' :  . ..XXXXXXXXXXXXXXXXXXXXx.         . YWMWMWMWMWMWMW
.    .  .  .    . .   .  ..XXXXXXXXWWWWWWWWWWWWWWWWXXXX.  .     .     """""""
        ' :  : . : .  ...XXXXXWWW"   W88N88@888888WWWWWXX.   .   .       . .
   . ' .    . :   ...XXXXXXWWW"    M88Ng8GGGG5G888^8M "WMBX.          .   ..  :
         :     ..XXXXXXXXWWW"     M88a8WWRWWWMW8oo88M   WWMX.     .    :    .
           "XXXXXXXXXXXXWW"       WN8s88WWWWW  W8@@@8M    BMBRX.         .  : :
  .       XXXXXXXX=MMWW":  .      W8N888WWWWWWWW88888W      XRBRXX.  .       .
     ....  ""XXXXXMM::::. .        W8@889WWWWWM8@8N8W      . . :RRXx.    .
         .....'''  MMM::.:.  .      W888N89999888@8W      . . ::::"RXV    .  :
 .       ..'''''      MMMm::.  .      WW888N88888WW     .  . mmMMMMMRXx
      ..' .            ""MMmm .  .       WWWWWWW   . :. :,miMM"""  : ""    .
   .                .       ""MMMMmm . .  .  .   ._,mMMMM"""  :  ' .  :
               .                  ""MMMMMMMMMMMMM""" .  : . '   .        .
          .              .     .    .                      .         .
.                                         .          .         .`
	p_log_fun("INFO",logo_str)


	//-------------
	//RUNTIME_SYS

	mongo_db := gf_core.Mongo__connect(p_mongodb_host_str,
							p_mongodb_db_name_str,
							p_log_fun )
	mongodb_coll := mongo_db.C("data_symphony")
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str:"gf_images",
		Log_fun:         p_log_fun,
		Mongodb_coll:    mongodb_coll,
	}
	//-------------
	//DB_INDEXES

	//IMPORTANT!! - make sure mongo has indexes build for relevant queries
	db_index__init(runtime_sys)
	//-------------
	//S3
	s3_info,gf_err := gf_core.S3__init(runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}
	//-------------

	jobs_mngr_ch := gf_images_jobs.Jobs_mngr__init(p_images_store_local_dir_path_str,
							p_images_thumbnails_store_local_dir_path_str,
							p_images_main_s3_bucket_name_str,
							s3_info,
							runtime_sys)
	//-------------
	//IMAGE_FLOWS
	flows__templates_dir_path_str := p_templates_dir_paths_map["flows_str"].(string)
	gf_err = Flows__init_handlers(flows__templates_dir_path_str,
							jobs_mngr_ch,
							runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}
	//-------------
	//GIF
	gif__templates_dir_path_str := p_templates_dir_paths_map["gif_str"].(string)
	gf_err = gf_gif_lib.Gif__init_handlers(gif__templates_dir_path_str,runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}
	//-------------

	gf_image_editor.Init_handlers(runtime_sys)
	//-------------
	/*gf_gif_lib.Init_img_to_gif_migration(*p_images_store_local_dir_path_str,
							*p_images_main_s3_bucket_name_str,
							s3_client,
							s3_uploader, //s3_client,
							mongodb_coll,
							p_log_fun)*/
	//-------------
	//JOBS_MANAGER
	gf_images_jobs.Jobs_mngr__init_handlers(jobs_mngr_ch,runtime_sys)
	//-------------
	//OTHER
	init_handlers(jobs_mngr_ch, runtime_sys)
	//------------------------
	//DASHBOARD SERVING
	static_files__url_base_str := "/images"
	gf_core.HTTP__init_static_serving(static_files__url_base_str,runtime_sys)
	//------------------------
	
	//----------------------
	//IMPORTANT!! - signal to user that server in this goroutine is ready to start listening 
	if p_init_done_ch != nil {
		p_init_done_ch <- true
	}
	//----------------------

	runtime_sys.Log_fun("INFO",">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	runtime_sys.Log_fun("INFO","STARTING HTTP SERVER - PORT - "+p_port_str)
	runtime_sys.Log_fun("INFO",">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	http_err := http.ListenAndServe(":"+p_port_str,nil)
	if http_err != nil {
		msg_str := "cant start listening on port - "+p_port_str
		runtime_sys.Log_fun("ERROR",msg_str)
		runtime_sys.Log_fun("ERROR",fmt.Sprint(http_err))
		
		panic(fmt.Sprint(http_err))
	}
}