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

package gf_images_jobs

import (
	"fmt"
	"time"
	"github.com/globalsign/mgo/bson"
	"gf_core"
	"apps/gf_images_lib/gf_images_utils"
)
//-------------------------------------------------
type Job_msg struct {
	job_id_str            string
	client_type_str       string 
	cmd_str               string //"start_job"|"get_running_job_ids"
	msg_response_ch       chan interface{}
	job_updates_ch        chan *Job_update_msg
	images_to_process_lst []Image_to_process
	flows_names_lst       []string
}
type Image_to_process struct {
	Source_url_str      string `bson:"source_url_str"`
	Origin_page_url_str string `bson:"origin_page_url_str"`
}

type Job_update_msg struct {
	Type_str             string        `json:"type_str"`
	Image_id_str         string        `json:"image_id_str"`
	Image_source_url_str string        `json:"image_source_url_str"`
	Err_str              string        `json:"err_str,omitempty"`  //if the update indicates an error, this is its value
	Image_thumbs         *gf_images_utils.Gf_image_thumbs `json:"-"`
}

type Running_job struct {
	Id                    bson.ObjectId        `bson:"_id,omitempty"`
	Id_str                string               `bson:"id_str"`
	T_str                 string               `bson:"t"`
	Client_type_str       string               `bson:"client_type_str"`
	Status_str            string               `bson:"status_str"` //"running"|"complete"
	Start_time_f          float64              `bson:"start_time_f"`
	End_time_f            float64              `bson:"end_time_f"`
	Images_to_process_lst []Image_to_process   `bson:"images_to_process_lst"`
	Errors_lst            []Job_Error          `bson:"errors_lst"`
	job_updates_ch        chan *Job_update_msg `bson:"-"`
}

type Job_Error struct {
	Type_str             string `bson:"type_str"`  //"fetcher_error"|"transformer_error"
	Error_str            string `bson:"error_str"` //serialization of the golang error
	Image_source_url_str string `bson:"image_source_url_str"`
}

//called "expected" because jobs are long-running processes, and they might fail at various stages
//of their processing. in that case some of these result values will be satisfied, others will not.
type Job_Expected_Output struct {
	Image_id_str                      string `json:"image_id_str"`
	Image_source_url_str              string `json:"image_source_url_str"`
	Thumbnail_small_relative_url_str  string `json:"thumbnail_small_relative_url_str"`
	Thumbnail_medium_relative_url_str string `json:"thumbnail_medium_relative_url_str"`
	Thumbnail_large_relative_url_str  string `json:"thumbnail_large_relative_url_str"`
}
//-------------------------------------------------
//CLIENT
//-------------------------------------------------
func Start_job(p_client_type_str string,
			p_images_to_process_lst []Image_to_process,
			p_flows_names_lst       []string,
			p_jobs_mngr_ch          chan Job_msg,
			p_runtime_sys           *gf_core.Runtime_sys) (*Running_job,[]*Job_Expected_Output,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_jobs.Start_job()")
	p_runtime_sys.Log_fun("INFO"     ,"p_images_to_process_lst - "+fmt.Sprint(p_images_to_process_lst))

	job_cmd_str      := "start_job"
	job_start_time_f := float64(time.Now().UnixNano())/1000000000.0
	job_id_str       := fmt.Sprintf("job:%f",job_start_time_f)
	job_updates_ch   := make(chan *Job_update_msg,10)

	job_msg := Job_msg{
		job_id_str           :job_id_str,
		client_type_str      :p_client_type_str,
		cmd_str              :job_cmd_str,
		job_updates_ch       :job_updates_ch,
		images_to_process_lst:p_images_to_process_lst,
		flows_names_lst      :p_flows_names_lst,
	}

	p_jobs_mngr_ch <- job_msg
	//-----------------
	//CREATE RUNNING_JOB
	running_job := &Running_job{
		Id_str               :job_id_str,
		T_str                :"img_running_job",
		Client_type_str      :p_client_type_str,
		Status_str           :"running",
		Start_time_f         :job_start_time_f,
		Images_to_process_lst:p_images_to_process_lst,
		job_updates_ch       :job_updates_ch,
	}

	db_err := p_runtime_sys.Mongodb_coll.Insert(running_job)
	if db_err != nil {
		gf_err := gf_core.Error__create("failed to create a Running_job record in the DB",
			"mongodb_insert_error",
			&map[string]interface{}{
				"client_type_str":      p_client_type_str,
				"images_to_process_lst":p_images_to_process_lst,
				"flows_names_lst":      p_flows_names_lst,
			},
			db_err,"gf_images_jobs",p_runtime_sys)
		return nil,nil,gf_err
	}
	//-----------------
	//CREATE JOB_EXPECTED_OUTPUT

	job_expected_outputs_lst := []*Job_Expected_Output{}

	for _,image_to_process := range p_images_to_process_lst {

		img_source_url_str := image_to_process.Source_url_str
		p_runtime_sys.Log_fun("INFO","img_source_url_str - "+fmt.Sprint(img_source_url_str))

		//--------------
		//IMAGE_ID
		image_id_str,i_gf_err := gf_images_utils.Image__create_id_from_url(img_source_url_str,p_runtime_sys)
		if i_gf_err != nil {
			return nil,nil,i_gf_err
		}
		//--------------
		//GET FILE_FORMAT
		normalized_ext_str,gf_err := gf_images_utils.Get_image_ext_from_url(img_source_url_str,p_runtime_sys)
		
		//FIX!! - it should not fail the whole job if one image is invalid,
		//        it should continue and just mark that image with an error.
		if gf_err != nil {
			return nil,nil,gf_err
		}
		//--------------

		output := &Job_Expected_Output{
			Image_id_str:                     image_id_str,
			Image_source_url_str:             img_source_url_str,
			Thumbnail_small_relative_url_str :fmt.Sprintf("/images/d/thumbnails/%s_thumb_small.%s" ,image_id_str,normalized_ext_str),
			Thumbnail_medium_relative_url_str:fmt.Sprintf("/images/d/thumbnails/%s_thumb_medium.%s",image_id_str,normalized_ext_str),
			Thumbnail_large_relative_url_str :fmt.Sprintf("/images/d/thumbnails/%s_thumb_large.%s" ,image_id_str,normalized_ext_str),
		}
		job_expected_outputs_lst = append(job_expected_outputs_lst,output)
	}
	//-----------------

	return running_job,job_expected_outputs_lst,nil
}
//-------------------------------------------------
func get_running_job_update_ch(p_job_id_str string,
			p_jobs_mngr_ch chan Job_msg,
			p_runtime_sys  *gf_core.Runtime_sys) chan *Job_update_msg {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_jobs.get_running_job_update_ch()")

	msg_response_ch := make(chan interface{})
	defer close(msg_response_ch)

	job_cmd_str := "get_running_job_update_ch"
	job_msg     := Job_msg{
		job_id_str     :p_job_id_str,
		cmd_str        :job_cmd_str,
		msg_response_ch:msg_response_ch,
	}

	p_jobs_mngr_ch <- job_msg

	response         := <-msg_response_ch
	job_updates_ch,_ := response.(chan *Job_update_msg)

	return job_updates_ch
}
//-------------------------------------------------
//SERVER
//-------------------------------------------------
func Jobs_mngr__init(p_images_store_local_dir_path_str string,
				p_images_thumbnails_store_local_dir_path_str string,
				p_s3_bucket_name_str                         string,
				p_s3_info                                    *gf_core.Gf_s3_info,
				p_runtime_sys                                *gf_core.Runtime_sys) chan Job_msg {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_images_jobs.Jobs_mngr__init()")

	jobs_mngr_ch := make(chan Job_msg,100)

	//IMPORTANT!! - start jobs_mngr as an independent goroutine of the HTTP handlers at
	//              service initialization time
	go func() {
		
		running_jobs_map := map[string]chan *Job_update_msg{}

		//listen to messages
		for {
			job_msg := <- jobs_mngr_ch

			//IMPORTANT!! - only one job is processed per jobs_mngr.
			//              Scaling is done with multiple jobs_mngr's (exp. per-core)           
			switch job_msg.cmd_str {

				//------------------------
				case "start_job":

					job_id_str                  := job_msg.job_id_str
					running_jobs_map[job_id_str] = job_msg.job_updates_ch

					run_job_gf_err := jobs_mngr__run_job(job_id_str,
											job_msg.client_type_str,
											job_msg.images_to_process_lst,
											job_msg.flows_names_lst,
											job_msg.job_updates_ch,
											p_images_store_local_dir_path_str,
											p_images_thumbnails_store_local_dir_path_str,
											p_s3_bucket_name_str,
											p_s3_info,
											p_runtime_sys)
					if run_job_gf_err != nil {
						continue
					}
					//--------------------
					//MARK_JOB_AS_COMPLETE

					job_end_time_f := float64(time.Now().UnixNano())/1000000000.0
					update_err     := p_runtime_sys.Mongodb_coll.Update(bson.M{"t":"img_running_job","id_str":job_id_str,},
														bson.M{
															"$set":bson.M{
																"status_str":"complete",
																"end_time_f":job_end_time_f,
															},
														},)
					if update_err != nil {
						_ = gf_core.Error__create("failed to update an img_running_job in the DB, as complete and its end_time",
							"mongodb_update_error",
							&map[string]interface{}{
								"job_id_str":    job_id_str,
								"job_end_time_f":job_end_time_f,
							},
							update_err,"gf_images_jobs",p_runtime_sys)
					}
					//--------------------

					delete(running_jobs_map,job_id_str) //remove running job from lookup, since its complete
					close(job_msg.job_updates_ch)
				//------------------------
				case "get_running_job_update_ch":

					job_id_str := job_msg.job_id_str

					if _,ok := running_jobs_map[job_id_str]; ok {

						job_updates_ch := running_jobs_map[job_id_str]
						job_msg.msg_response_ch <- job_updates_ch

					} else {
						job_msg.msg_response_ch <- nil
					}
				//------------------------
			} 
		}
	}()

	return jobs_mngr_ch
}