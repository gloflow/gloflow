/*
GloFlow application and media management/publishing platform
Copyright (C) 2026 Ivan Trajkovic

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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang"
)

//-------------------------------------------------
// the dev dashboard's "Render PNG" button: the program (already parsed in the browser,
// so it is the same AST the browser runs) is run by the Go interpreter with the scene
// recorder, and the scene is rendered to a PNG by the headless wgpu renderer
// (gf_lang/rust/gf_lang_render, see gf_lang/docs/headless_rendering.md).
//
// POST /v1/lang/render?width=N&height=N&fit=1
//   body:  {"program_ast_lst": [...], "camera": {"position": [x,y,z], "target": [x,y,z], "up": [x,y,z], "fov_deg": N}}
//          "camera" is optional: it is the view the user has in the dashboard. when given, it is used
//          as-is and "fit" is ignored. without it the renderer frames the scene (fit) or uses the scene camera.
//   reply: image/png, or text/plain error (400: bad input or program error, 500: renderer error)

const (
	renderDefaultWidthInt  = 1024
	renderDefaultHeightInt = 768
	renderMaxSideInt       = 4096
	renderMaxBodyBytesInt  = 4 * 1024 * 1024
	renderMaxEntitiesInt   = 1000000
	renderTimeout          = 60 * time.Second
	renderMaxConcurrentInt = 2 // each render is a GPU job, so only a few run at once
)

var renderSemaphore = make(chan struct{}, renderMaxConcurrentInt)

//-------------------------------------------------

func initRender(pHTTPmux *http.ServeMux,
	pConfig     *GFconfig,
	pRuntimeSys *gf_core.RuntimeSys) {

	pHTTPmux.HandleFunc("/v1/lang/render", func(pResp http.ResponseWriter, pReq *http.Request) {

		if pReq.Method != "POST" {
			http.Error(pResp, "use POST", http.StatusMethodNotAllowed)
			return
		}

		//---------------------
		// INPUT
		var input struct {
			ProgramASTlst []interface{}      `json:"program_ast_lst"`
			Camera        *renderCameraInput `json:"camera"`
		}
		if err := json.NewDecoder(io.LimitReader(pReq.Body, renderMaxBodyBytesInt)).Decode(&input); err != nil {
			http.Error(pResp, "invalid JSON body: "+err.Error(), http.StatusBadRequest)
			return
		}
		if len(input.ProgramASTlst) == 0 {
			http.Error(pResp, "program_ast_lst is empty", http.StatusBadRequest)
			return
		}

		widthInt, err := renderSizeFromQuery(pReq, "width", renderDefaultWidthInt)
		if err != nil {
			http.Error(pResp, err.Error(), http.StatusBadRequest)
			return
		}
		heightInt, err := renderSizeFromQuery(pReq, "height", renderDefaultHeightInt)
		if err != nil {
			http.Error(pResp, err.Error(), http.StatusBadRequest)
			return
		}
		// the browser has trackball controls, a static render needs the camera framed on the scene
		fitBool := pReq.URL.Query().Get("fit") != "0"

		var cameraPtr *renderCamera
		if input.Camera != nil {
			cam, err := input.Camera.validate()
			if err != nil {
				http.Error(pResp, "invalid camera: "+err.Error(), http.StatusBadRequest)
				return
			}
			cameraPtr = cam
		}

		//---------------------
		// RUN + RENDER
		select {
		case renderSemaphore <- struct{}{}:
			defer func() { <-renderSemaphore }()
		case <-pReq.Context().Done():
			return
		}

		pngBytesLst, err := renderProgramToPNG(pReq.Context(),
			input.ProgramASTlst,
			widthInt, heightInt, fitBool, cameraPtr,
			pConfig.RendererBinPathStr)
		if err != nil {
			status := http.StatusInternalServerError
			if _, isProgramErr := err.(*renderProgramError); isProgramErr {
				status = http.StatusBadRequest
			}
			pRuntimeSys.LogNewFun("ERROR", "render of a gf_lang program failed: "+err.Error(), nil)
			http.Error(pResp, err.Error(), status)
			return
		}

		pResp.Header().Set("Content-Type", "image/png")
		pResp.Write(pngBytesLst)
	})
}

//-------------------------------------------------

// renderProgramError - the program itself failed (as opposed to the renderer).
type renderProgramError struct {
	msgStr string
}

func (e *renderProgramError) Error() string { return e.msgStr }

//-------------------------------------------------

// renderCameraInput - the camera as sent by the dashboard (JSON), not yet validated.
type renderCameraInput struct {
	Position []float64 `json:"position"`
	Target   []float64 `json:"target"`
	Up       []float64 `json:"up"`
	FOVdegF  float64   `json:"fov_deg"`
}

// renderCamera - a validated camera, passed to gf_lang_render as flags.
type renderCamera struct {
	Position [3]float64
	Target   [3]float64
	Up       [3]float64
	FOVdegF  float64
}

func (pIn *renderCameraInput) validate() (*renderCamera, error) {

	vec3 := func(pName string, pLst []float64) ([3]float64, error) {
		var v [3]float64
		if len(pLst) != 3 {
			return v, fmt.Errorf("%s must have 3 numbers", pName)
		}
		for i, f := range pLst {
			if math.IsNaN(f) || math.IsInf(f, 0) {
				return v, fmt.Errorf("%s must be finite numbers", pName)
			}
			v[i] = f
		}
		return v, nil
	}

	pos, err := vec3("position", pIn.Position)
	if err != nil {
		return nil, err
	}
	target, err := vec3("target", pIn.Target)
	if err != nil {
		return nil, err
	}
	up, err := vec3("up", pIn.Up)
	if err != nil {
		return nil, err
	}
	if pIn.FOVdegF <= 0 || pIn.FOVdegF >= 180 {
		return nil, fmt.Errorf("fov_deg must be between 0 and 180")
	}

	// the same checks as gf_lang_render's CameraView::validate, so bad input is a 400 and not a renderer failure
	dir := [3]float64{target[0] - pos[0], target[1] - pos[1], target[2] - pos[2]}
	if dir == [3]float64{0, 0, 0} {
		return nil, fmt.Errorf("position and target are the same point")
	}
	if up == [3]float64{0, 0, 0} {
		return nil, fmt.Errorf("up must not be zero")
	}
	cross := [3]float64{
		dir[1]*up[2] - dir[2]*up[1],
		dir[2]*up[0] - dir[0]*up[2],
		dir[0]*up[1] - dir[1]*up[0],
	}
	if cross == [3]float64{0, 0, 0} {
		return nil, fmt.Errorf("up must not be parallel to the view direction")
	}
	return &renderCamera{Position: pos, Target: target, Up: up, FOVdegF: pIn.FOVdegF}, nil
}

//-------------------------------------------------

func renderSizeFromQuery(pReq *http.Request, pNameStr string, pDefaultInt int) (int, error) {
	valStr := pReq.URL.Query().Get(pNameStr)
	if valStr == "" {
		return pDefaultInt, nil
	}
	sizeInt, err := strconv.Atoi(valStr)
	if err != nil || sizeInt < 1 || sizeInt > renderMaxSideInt {
		return 0, fmt.Errorf("%s must be a number between 1 and %d", pNameStr, renderMaxSideInt)
	}
	return sizeInt, nil
}

//-------------------------------------------------

func renderProgramToPNG(pCtx context.Context,
	pProgramASTlst []interface{},
	pWidthInt      int,
	pHeightInt     int,
	pFitBool       bool,
	pCamera        *renderCamera,
	pRendererBinStr string) (pngBytes []byte, err error) {

	//---------------------
	// RUN - the interpreter records the scene instead of drawing it
	scene, err := runProgramScene(pProgramASTlst)
	if err != nil {
		return nil, &renderProgramError{msgStr: "program run failed: " + err.Error()}
	}

	sceneJSONbytes, err := scene.JSON()
	if err != nil {
		return nil, err
	}

	//---------------------
	// RENDER - the headless renderer runs as a separate process (see docs/headless_rendering.md)
	tmpDirStr, err := os.MkdirTemp("", "gf_lang_render_")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDirStr)
	outPathStr := filepath.Join(tmpDirStr, "out.png")

	args := []string{
		"--out",    outPathStr,
		"--width",  strconv.Itoa(pWidthInt),
		"--height", strconv.Itoa(pHeightInt),
	}
	if pCamera != nil {
		// the camera overrides fit: the view is exactly the one the user had
		args = append(args,
			"--camera-pos",    renderVec3flag(pCamera.Position),
			"--camera-target", renderVec3flag(pCamera.Target),
			"--camera-up",     renderVec3flag(pCamera.Up),
			"--fov",           strconv.FormatFloat(pCamera.FOVdegF, 'g', -1, 64),
		)
	} else if pFitBool {
		args = append(args, "--fit")
	}

	ctx, cancel := context.WithTimeout(pCtx, renderTimeout)
	defer cancel()

	var stderrBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, pRendererBinStr, args...)
	cmd.Stdin  = bytes.NewReader(sceneJSONbytes)
	cmd.Stdout = &stderrBuf
	cmd.Stderr = &stderrBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("renderer %s failed: %w: %s", pRendererBinStr, err, stderrBuf.String())
	}

	pngBytes, err = os.ReadFile(outPathStr)
	if err != nil {
		return nil, err
	}
	return pngBytes, nil
}

//-------------------------------------------------

// renderVec3flag - "x,y,z" as the gf_lang_render flag format
func renderVec3flag(pV [3]float64) string {
	return fmt.Sprintf("%s,%s,%s",
		strconv.FormatFloat(pV[0], 'g', -1, 64),
		strconv.FormatFloat(pV[1], 'g', -1, 64),
		strconv.FormatFloat(pV[2], 'g', -1, 64))
}

//-------------------------------------------------

// runProgramScene - recovers from panics in the interpreter, so a bad program
// is reported as a program error and doesn't take the server down.
func runProgramScene(pProgramASTlst []interface{}) (scene *gf_lang.GFscene, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	programASTlst := gf_lang.CastToExpr(pProgramASTlst)
	return gf_lang.SceneRun(programASTlst, renderMaxEntitiesInt)
}
