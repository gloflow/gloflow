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

package gf_images_lib

import (
	"fmt"
	"testing"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
)

//---------------------------------------------------

func TestJobsUpdates(p_test *testing.T) {

	fmt.Println("TEST__JOBS_UPDATES ==============================================")



	TstartService()



	//-------------
	// HTTP
	test__job_updates__via_http(test__image_url_str,
		test__origin_page_url_str,
		test__image_client_type_str,
		test__gf_images_service_port_str,
		runtimeSys)

	// IN_PROCESS
	test__job_updates__in_process(test__image_url_str,
		test__image_flows_names_lst,
		test__image_client_type_str,
		test__origin_page_url_str,
		jobs_mngr,
		runtimeSys)

	//-------------
}

//---------------------------------------------------

func test__job_updates__via_http(pTestImageURLstr string,
	pTestOriginPageURLstr         string,
	p_test__image_client_type_str string,
	p_test_image_service_port_str string,
	pRuntimeSys                   *gf_core.RuntimeSys) {
	
	test__input_images_urls_lst                := []string{pTestImageURLstr,}
	test__input_images_origin_pages_urls_lst   := []string{pTestOriginPageURLstr,}
	test__image_service_host_port_str          := fmt.Sprintf("localhost:%s", p_test_image_service_port_str)
	runningJobIDstr, imagesOutputsLst, gfErr := Client__dispatch_process_extern_images(test__input_images_urls_lst,
		test__input_images_origin_pages_urls_lst,
		p_test__image_client_type_str,
		test__image_service_host_port_str,
		pRuntimeSys)
	
	if gfErr != nil {
		panic(gfErr.Error)
	}
	

	fmt.Println(runningJobIDstr)
	fmt.Println(imagesOutputsLst)


}

//---------------------------------------------------

func test__job_updates__in_process(pTestImageURLstr string,
	p_test__image_flows_names_lst []string,
	p_test__image_client_type_str string,
	pTestOriginPageURLstr   string,
	pJobsMngr                     gf_images_jobs.Jobs_mngr,
	pRuntimeSys                   *gf_core.RuntimeSys) {

	imagesToProcessLst := []gf_images_jobs.GFimageExternToProcess{
		gf_images_jobs.GFimageExternToProcess{
			Source_url_str:      pTestImageURLstr,
			Origin_page_url_str: pTestOriginPageURLstr,
		},
	}

	runningJob, output, gfErr := gf_images_jobs.Client__run_extern_imgs(p_test__image_client_type_str,
		imagesToProcessLst,
		p_test__image_flows_names_lst,
		pJobsMngr,
		pRuntimeSys)
	if gfErr != nil {
		panic(gfErr.Error)
	}

	fmt.Println(runningJob)
	spew.Dump(output)

	T__test_image_job__updates(runningJob.Id_str, pJobsMngr, pRuntimeSys)
}