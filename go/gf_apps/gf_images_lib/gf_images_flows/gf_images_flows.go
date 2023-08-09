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

package gf_images_flows

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/gloflow/gloflow/go/gf_identity/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_identity/gf_policy"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_jobs_client"
)

//-------------------------------------------------
// IMPORTANT!! - image_flow's are ordered sequences of images, that the user creates and then
//               over time adds images to it... 

type GFflow struct {
	Vstr              string
	IDstr             gf_core.GF_ID
	CreationUNIXtimeF float64
	NameStr           string
	OwnerUserID       gf_core.GF_ID
	EditorUserIDs     []gf_core.GF_ID
}

type GFimageExistsCheck struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty"`
	IDstr               gf_core.GF_ID      `bson:"id_str"`
	Tstr                string             `bson:"t"`
	CreationUNIXtimeF   float64            `bson:"creation_unix_time_f"`
	ImagesExternURLsLst []string           `bson:"images_extern_urls_lst"`
}

//-------------------------------------------------

func Init(pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {


	ctx := context.Background()


	// SQL_CREATE_TABLES
	gfErr := DBsqlCreateTables(pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	// CREATE_DISCOVERED_FLOWS
	gfErr = pipelineCreateDiscoveredFlows(ctx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	return nil
}

//-------------------------------------------------
/*
TEMPORARY - this is mainly needed while flows are held as a property of images
        	and discovered there in aggregate to get the total list.
			going forward flows are held in the SQL db and this function
			migrates/creates them in SQL if they dont already exist.

			in the future this function wont be necessary, unless there's some
			need for copying of flows from DB to DB.
*/

// consistency function that discovers all flows listed in images
// under their "flows" attribute, and creates them as explicit entities
// in the main GF DB if they dont already exist.
func pipelineCreateDiscoveredFlows(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {
	
	// get all flows from the current Mongodb
	allFlowsLst, gfErr := DBgetAll(pCtx, pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//----------------------
	// IMPORTANT!! - system user "gf" is assigned as the owner of this flow
	ownerUserID := gf_core.GF_ID("gf")

	//----------------------

	// SQL
	for _, flowMap := range allFlowsLst {

		nameStr := flowMap["_id"].(string)
		fmt.Printf("flow name - %s", nameStr)

		existsBool, gfErr := DBsqlCheckFlowExists(nameStr, pRuntimeSys)
		if gfErr != nil {
			return gfErr
		}

		// create flow if it doesnt exist
		if !existsBool {

			//----------------------
			// CREATE_FLOW
			_, gfErr := Create(nameStr,
				ownerUserID,
				pCtx,
				pRuntimeSys)

			if gfErr != nil {
				return gfErr
			}

			//----------------------
		}

	}



	return nil
}

//-------------------------------------------------

func pipelineGetAll(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {

	resultsLst, gfErr := DBgetAll(pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	allFlowsLst := []map[string]interface{}{}
	for _, flowInfoMap := range resultsLst {
		flowNameStr      := flowInfoMap["_id"].(string)
		flowImgsCountInt := flowInfoMap["count_int"].(int32)

		allFlowsLst = append(allFlowsLst, map[string]interface{}{
			"flow_name_str":       flowNameStr,
			"flow_imgs_count_int": flowImgsCountInt,
		})
	}
	return allFlowsLst, nil
}

//-------------------------------------------------
// GET_PAGE__PIPELINE

func pipelineGetPage(p_req *http.Request,
	p_resp        http.ResponseWriter,
	p_ctx         context.Context,
	pRuntimeSys *gf_core.RuntimeSys) ([]*gf_images_core.GFimage, *gf_core.GFerror) {

	//--------------------
	// INPUT

	qs_map := p_req.URL.Query()

	flow_name_str := "general" // default
	if a_lst,ok := qs_map["fname"]; ok {
		flow_name_str = a_lst[0]
	}

	var err error
	page_index_int := 0 // default
	if a_lst, ok := qs_map["pg_index"]; ok {
		pg_index           := a_lst[0]
		page_index_int, err = strconv.Atoi(pg_index) // user supplied value
		
		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to parse integer pg_index query string arg",
				"int_parse_error",
				map[string]interface{}{"pg_index": pg_index,},
				err, "gf_images_lib", pRuntimeSys)
			return nil, gfErr
		}
	}

	page_size_int := 10 // default
	if a_lst, ok := qs_map["pg_size"]; ok {
		pg_size          := a_lst[0]
		page_size_int,err = strconv.Atoi(pg_size) // user supplied value
		if err != nil {
			gfErr := gf_core.ErrorCreate("failed to parse integer pg_size query string arg",
				"int_parse_error",
				map[string]interface{}{"pg_size": pg_size,},
				err, "gf_images_lib", pRuntimeSys)
			return nil, gfErr
		}
	}

	pRuntimeSys.LogFun("INFO",fmt.Sprintf("flow_name_str  - %s", flow_name_str))
	pRuntimeSys.LogFun("INFO",fmt.Sprintf("page_index_int - %d", page_index_int))
	pRuntimeSys.LogFun("INFO",fmt.Sprintf("page_size_int  - %d", page_size_int))

	//--------------------

	//--------------------
	// GET_PAGES
	cursor_start_position_int := page_index_int*page_size_int
	pages_lst, gfErr := dbGetPage(flow_name_str,
		cursor_start_position_int, // p_cursor_start_position_int
		page_size_int,             // p_elements_num_int
		p_ctx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------
	return pages_lst, nil
}

//-------------------------------------------------
// IMAGES_EXIST_CHECK

func imagesExistCheck(pImagesExternURLsLst []string,
	pFlowNameStr   string,
	pClientTypeStr string,
	pRuntimeSys  *gf_core.RuntimeSys) ([]map[string]interface{}, *gf_core.GFerror) {
	
	existing_images_lst, gfErr := dbImagesExist(pImagesExternURLsLst,
		pFlowNameStr,
		pClientTypeStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//-------------------------
	// PERSIST IMAGE_EXISTS_CHECK

	go func() {
		creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
		idStr             := gf_core.GF_ID(fmt.Sprintf("img_exists_check:%f", creationUNIXtimeF))
		
		check := GFimageExistsCheck{
			IDstr:               idStr,
			Tstr:                "img_exists_check",
			CreationUNIXtimeF:   creationUNIXtimeF,
			ImagesExternURLsLst: pImagesExternURLsLst,
		}

		ctx           := context.Background()
		coll_name_str := "gf_flows_img_exists_check" // pRuntimeSys.Mongo_coll.Name()
		_              = gf_core.MongoInsert(check,
			coll_name_str,
			map[string]interface{}{
				"images_extern_urls_lst": pImagesExternURLsLst,
				"flow_name_str":          pFlowNameStr,
				"client_type_str":        pClientTypeStr,
				"caller_err_msg_str":     "failed to insert a img_exists_check into the DB",
			},
			ctx,
			pRuntimeSys)
	}()

	//-------------------------

	return existing_images_lst, nil
}

//-------------------------------------------------
// ADD_EXTERN_IMAGE_WITH_POLICY

func AddExternImageWithPolicy(pImageExternURLstr string,
	pImageOriginPageURLstr string,
	pFlowsNamesLst         []string,
	pClientTypeStr         string,
	pUserIDstr             gf_core.GF_ID,
	pJobsMngrCh            chan gf_images_jobs_core.JobMsg,
	pCtx                   context.Context,
	pRuntimeSys            *gf_core.RuntimeSys) (*string, *string, gf_images_core.GF_image_id, *gf_core.GFerror) {

	//-------------------------
	// POLICY_VERIFY - raises error if policy rejects the op
	opStr := gf_policy.GF_POLICY_OP__FLOW_ADD_IMG
	gfErr := VerifyPolicy(opStr,
		pFlowsNamesLst,
		pUserIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gf_images_core.GFimageID(""), gfErr
	}

	//-------------------------

	runningJobIDstr, thumbnailSmallRelativeURLstr, imageIDstr, gfErr := AddExternImage(pImageExternURLstr,
		pImageOriginPageURLstr,
		pFlowsNamesLst,
		pClientTypeStr,
		pUserIDstr,
		pJobsMngrCh,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, nil, gf_images_core.GFimageID(""), gfErr
	}

	return runningJobIDstr, thumbnailSmallRelativeURLstr, imageIDstr, nil
}

//-------------------------------------------------
// ADD_EXTERN_IMAGES - BATCH

func AddExternImages(pImagesExternURLsLst []string,
	pImagesOriginPagesURLsStr []string,
	pFlowsNamesLst            []string,
	pClientTypeStr            string,
	pJobsMngrCh               chan gf_images_jobs_core.JobMsg,
	pRuntimeSys               *gf_core.RuntimeSys) (*string, []*string, []gf_images_core.GFimageID, *gf_core.GFerror) {

	//------------------
	imagesURLsToProcessLst := []gf_images_jobs_core.GFimageExternToProcess{}
	for i := 0; i < len(pImagesExternURLsLst); i++ {

		imageExternURLstr := pImagesExternURLsLst[i]
		imageOriginPageURLstr := pImagesOriginPagesURLsStr[i]
		
		image := gf_images_jobs_core.GFimageExternToProcess{
			SourceURLstr:     imageExternURLstr,
			OriginPageURLstr: imageOriginPageURLstr,
		}
		imagesURLsToProcessLst = append(imagesURLsToProcessLst, image)
	}
	
	// GF_IMAGES_JOBS_CLIENT
	runningJob, jobExpectedOutputsLst, gfErr := gf_images_jobs_client.RunExternImages(pClientTypeStr,
		imagesURLsToProcessLst,
		pFlowsNamesLst,
		pJobsMngrCh,
		pRuntimeSys)

	if gfErr != nil {
		return nil, nil, nil, gfErr
	}

	//------------------

	imagesIDsLst                   := []gf_images_core.GFimageID{}
	imagesThumbSmallRelativeURLlst := []*string{}

	for i:=0; i < len(jobExpectedOutputsLst); i++ {	
		imageIDstr               := gf_images_core.GFimageID(jobExpectedOutputsLst[i].Image_id_str)
		thumbSmallRelativeURLstr := jobExpectedOutputsLst[i].Thumbnail_small_relative_url_str


		imagesIDsLst = append(imagesIDsLst, imageIDstr)
		imagesThumbSmallRelativeURLlst = append(imagesThumbSmallRelativeURLlst, &thumbSmallRelativeURLstr)
	}

	return &runningJob.Id_str, imagesThumbSmallRelativeURLlst, imagesIDsLst, nil
}

//-------------------------------------------------
// ADD_EXTERN_IMAGE

func AddExternImage(pImageExternURLstr string,
	pImageOriginPageURLstr string,
	pFlowsNamesLst         []string,
	pClientTypeStr         string,
	pUserIDstr             gf_core.GF_ID,
	pJobsMngrCh            chan gf_images_jobs_core.JobMsg,
	pCtx                   context.Context,
	pRuntimeSys            *gf_core.RuntimeSys) (*string, *string, gf_images_core.GFimageID, *gf_core.GFerror) {

	//------------------
	// FLOWS
	// check if each flow that was specified for this new image exists,
	// and if it doesnt create it first, before processing images.

	for _, flowNameStr := range pFlowsNamesLst {
		
		// check flow exists
		existsBool, gfErr := DBsqlCheckFlowExists(flowNameStr, pRuntimeSys)
		if gfErr != nil {
			return nil, nil, gf_images_core.GFimageID(""), gfErr
		}

		// if it doesnt exist, create it... 
		if !existsBool {

			// FLOW_CREATE
			_, gfErr := Create(flowNameStr, pUserIDstr,
				pCtx,
				pRuntimeSys)
			if gfErr != nil {
				return nil, nil, gf_images_core.GFimageID(""), gfErr
			}
		}
	}

	//------------------
	imagesURLsToProcessLst := []gf_images_jobs_core.GFimageExternToProcess{
			{
				SourceURLstr:     pImageExternURLstr,
				OriginPageURLstr: pImageOriginPageURLstr,
			},
		}
	
	// GF_IMAGES_JOBS_CLIENT
	runningJob, jobExpectedOutputsLst, gfErr := gf_images_jobs_client.RunExternImages(pClientTypeStr,
		imagesURLsToProcessLst,
		pFlowsNamesLst,
		pJobsMngrCh,
		pRuntimeSys)

	if gfErr != nil {
		return nil, nil, gf_images_core.GFimageID(""), gfErr
	}

	//------------------

	imageIDstr                       := gf_images_core.GFimageID(jobExpectedOutputsLst[0].Image_id_str)
	thumbnail_small_relative_url_str := jobExpectedOutputsLst[0].Thumbnail_small_relative_url_str

	return &runningJob.Id_str, &thumbnail_small_relative_url_str, imageIDstr, nil
}

//-------------------------------------------------
// CREATE

func Create(pFlowNameStr string,
	pOwnerUserIDstr gf_core.GF_ID,
	pCtx            context.Context,
	pRuntimeSys     *gf_core.RuntimeSys) (*GFflow, *gf_core.GFerror) {

	creationUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0
	idStr             := gf_core.GF_ID(fmt.Sprintf("img_flow:%f", creationUNIXtimeF))
	



	flow := &GFflow{
		IDstr:             idStr,
		NameStr:           pFlowNameStr,
		CreationUNIXtimeF: creationUNIXtimeF,
		OwnerUserID:       pOwnerUserIDstr,
	}

	//------------------
	// DB

	// SQL
	if pRuntimeSys.SQLdb != nil {


		gfErr := DBsqlCreateFlow(idStr,
			pFlowNameStr,
			pOwnerUserIDstr,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}


	// MONGODB
	} else {
		collNameStr := pRuntimeSys.Mongo_coll.Name()
		gfErr := gf_core.MongoInsert(flow,
			collNameStr,
			map[string]interface{}{
				"images_flow_name_str": pFlowNameStr,
				"caller_err_msg_str":   "failed to insert a image Flow into the DB",
			},
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return nil, gfErr
		}
	}
	
	//------------------
	// POLICY_CREATE
	
	/*
	gfErr = gf_policy.PipelineCreate(idStr, pOwnerUserIDstr, pCtx, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}
	*/

	//------------------

	return flow, nil
}