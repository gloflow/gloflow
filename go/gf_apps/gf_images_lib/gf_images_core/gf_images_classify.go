/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

package gf_images_core

import (
	"fmt"
	"time"
	"context"
	"strings"
	"github.com/gloflow/gloflow/go/gf_core"
)


//---------------------------------------------------

func RunClassify(pImagesIDsLst []GFimageID,
	pPyDirPathStr string,
	pMetrics      *GFmetrics,
	pCtx          context.Context,
	pRuntimeSys   *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {


	// PY
	classesLst, gfErr := RunPyClassify(pImagesIDsLst,
		pPyDirPathStr,
		pMetrics,
		pCtx,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}



	fmt.Println(classesLst)




	return classesLst, nil
}

//---------------------------------------------------
// PY

func RunPyClassify(pImagesIDsLst []GFimageID,
	pPyDirPathStr string,
	pMetrics      *GFmetrics,
	pCtx          context.Context,
	pRuntimeSys   *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {

	

	imagesIDsLst := []string{}
	for _, imageID := range pImagesIDsLst {
		imagesIDsLst = append(imagesIDsLst, string(imageID))
	}



	pyPathStr := fmt.Sprintf("%s/gf_images_classify.py", pPyDirPathStr)
	argsLst := []string{
		fmt.Sprintf("-images_ids=%s", strings.Join(imagesIDsLst, ",")),
	}
	stdoutPrefixStr := "GF_OUT:"
	inputStdinStr   := ""

	//---------------------
	// ENV
	envMap := map[string]string{
		"SENTRY_DSN": pRuntimeSys.SentryDSNstr,
		"SENTRY_ENV": pRuntimeSys.EnvStr,
	}

	//---------------------

	runStartUNIXtimeF := float64(time.Now().UnixNano())/1000000000.0

	
	// PY_RUN
	outputsLst, gfErr := gf_core.CLIpyRun(pyPathStr,
		argsLst,
		&inputStdinStr,
		envMap,
		stdoutPrefixStr,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	runEndUNIXtimeF   := float64(time.Now().UnixNano())/1000000000.0
	runDurrationSecsF := runEndUNIXtimeF - runStartUNIXtimeF
	
	if pMetrics != nil {
		pMetrics.ImageClassifyPyExecDurationGauge.Set(runDurrationSecsF)
	}

	fmt.Println(outputsLst)



	// OUTPUTS
	classesLst := outputsLst[0]["classes_lst"].([]interface{})
	classesCastLst := []string{}
	for _, class := range classesLst {
		classesCastLst = append(classesCastLst, class.(string))
	}

	return classesCastLst, nil
}

//-------------------------------------------------