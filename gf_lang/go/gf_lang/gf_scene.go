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

package gf_lang

import (
    "math"
    "strconv"
    "strings"
    "encoding/json"
)

//-------------------------------------------------
// SCENE - serializable scene description, produced by running a gf_lang program
//         with the GFsceneRecorder extern API. It is the contract between the
//         interpreter and renderers (the headless wgpu renderer in gf_lang/rust).
//         Schema docs: gf_lang/docs/scene.md
//
// all positions in the scene are final world-space values (coord origin stack
// and rotation pivots are already applied by the recorder, mirroring ts/engine/gf_engine.ts).

const GFsceneVersionInt = 1

type GFscene struct {
    VersionInt int                        `json:"version"`
    Background []float64                  `json:"background,omitempty"` // [r,g,b] 0-1, nil == renderer default (black)
    Camera     GFsceneCamera              `json:"camera"`
    Shaders    map[string]*GFsceneShader  `json:"shaders"`
    Entities   []*GFsceneEntity           `json:"entities"`
    TruncatedBool bool                    `json:"truncated"` // true if the entities cap was hit
}

type GFsceneCamera struct {
    Position   [3]float64 `json:"position"`
    Target     [3]float64 `json:"target"`
    FOVdegF    float64    `json:"fov_deg"`
    NearF      float64    `json:"near"`
    FarF       float64    `json:"far"`
}

type GFsceneShader struct {
    Uniforms   []GFsceneShaderUniform `json:"uniforms"`
    VertexStr  string                 `json:"vertex"`
    FragmentStr string                `json:"fragment"`
}

type GFsceneShaderUniform struct {
    NameStr string      `json:"name"`
    TypeStr string      `json:"type"`
    Default interface{} `json:"default"`
}

type GFsceneEntity struct {
    KindStr   string      `json:"kind"` // cube|sphere|line
    Position  [3]float64  `json:"position"`
    Rotation  [3]float64  `json:"rotation"` // euler XYZ, radians
    Scale     [3]float64  `json:"scale"`
    Color     [3]float64  `json:"color"`
    Wireframe bool        `json:"wireframe"`
    Material  GFsceneMaterial `json:"material"`
    Points    [][3]float64 `json:"points,omitempty"` // lines only (world-space)

    // recorded only, renderers may ignore it (no animation support yet)
    Animation *GFsceneAnimation `json:"animation,omitempty"`
}

type GFsceneMaterial struct {
    TypeStr   string                 `json:"type"` // basic|shader
    ShaderStr string                 `json:"shader,omitempty"`
    Uniforms  map[string]interface{} `json:"uniforms,omitempty"` // snapshot at entity creation time
}

type GFsceneAnimation struct {
    Props         []GFsceneAnimationProp `json:"props"`
    DurationSecF  float64                `json:"duration_sec"`
    RepeatBool    bool                   `json:"repeat"`
}

type GFsceneAnimationProp struct {
    NameStr string  `json:"name"`
    StartF  float64 `json:"start"`
    EndF    float64 `json:"end"`
}

//-------------------------------------------------
// RECORDER

type GFsceneRecorder struct {
    Scene           *GFscene
    MaxEntitiesInt  int // 0 == unlimited

    coordOriginsLst [][3]float64
    pivotsLst       [][3]float64
    wireframeBool   bool
    shaderActiveBool bool
    shaderActiveStr  string
    uniformsMap     map[string]map[string]interface{} // shader_name -> uniform_name -> current value
    linePointsLst   [][3]float64
    animation       *GFsceneAnimation
}

func SceneRecorderNew(pMaxEntitiesInt int) *GFsceneRecorder {
    return &GFsceneRecorder{
        Scene: &GFscene{
            VersionInt: GFsceneVersionInt,
            // defaults mirror ts/engine/gf_engine.ts init(). the browser engine calls
            // camera.lookAt(0,5000,0) but TrackballControls.update() immediately
            // re-targets the camera to the controls target (world origin).
            Camera: GFsceneCamera{
                Position: [3]float64{10, 10, 10},
                Target:   [3]float64{0, 0, 0},
                FOVdegF:  75,
                NearF:    0.1,
                FarF:     10000,
            },
            Shaders:  map[string]*GFsceneShader{},
            Entities: []*GFsceneEntity{},
        },
        MaxEntitiesInt:  pMaxEntitiesInt,
        coordOriginsLst: [][3]float64{{0, 0, 0}},
        pivotsLst:       [][3]float64{{0, 0, 0}},
        uniformsMap:     map[string]map[string]interface{}{},
    }
}

// SceneRun - runs a program with a recorder and returns the resulting scene.
func SceneRun(pProgramASTlst GFexpr, pMaxEntitiesInt int) (*GFscene, error) {
    recorder := SceneRecorderNew(pMaxEntitiesInt)
    _, _, err := Run(pProgramASTlst, recorder.ExternAPI())
    if err != nil {
        return nil, err
    }
    return recorder.Scene, nil
}

func (pScene *GFscene) JSON() ([]byte, error) {
    return json.Marshal(pScene)
}

//-------------------------------------------------

func (r *GFsceneRecorder) ExternAPI() GFexternAPI {
    return GFexternAPI{
        InitEngineFun: r.initEngine,
        SetStateFun:   r.setState,
        CreateCubeFun:   r.entityFun("cube"),
        CreateSphereFun: r.entityFun("sphere"),
        CreateLineFun:   r.lineFun,
        AnimateFun: r.animate,

        // not supported when recording a scene
        RPCcall: func(string, string, string, []interface{}) map[string]interface{} { return nil },
        RPCserve: func(string, []*GFrpcServerHandler, GFexternAPI) {},
    }
}

//-------------------------------------------------

func (r *GFsceneRecorder) initEngine(pShaderDefsMap map[string]*GFshaderDef) {
    for nameStr, def := range pShaderDefsMap {
        shader := &GFsceneShader{
            Uniforms:    []GFsceneShaderUniform{},
            VertexStr:   def.VertexCodeStr,
            FragmentStr: def.FragmentCodeStr,
        }
        r.uniformsMap[nameStr] = map[string]interface{}{}
        for _, u := range def.UniformsDefsLst {
            shader.Uniforms = append(shader.Uniforms, GFsceneShaderUniform{u.NameStr, u.TypeStr, u.DefaultVal})
            r.uniformsMap[nameStr][u.NameStr] = u.DefaultVal
        }
        r.Scene.Shaders[nameStr] = shader
    }
}

//-------------------------------------------------

func (r *GFsceneRecorder) setState(pChange GFstateChange) []interface{} {

    switch pChange.PropertyNameStr {

    case "coord_origin":
        if pChange.OriginTypeStr != "current_pos" {
            return nil
        }
        switch pChange.SetterTypeStr {
        case "push":
            latest := r.coordOriginsLst[len(r.coordOriginsLst)-1]
            p := realWorldPos(
                [3]float64{latest[0] + pChange.Xf, latest[1] + pChange.Yf, latest[2] + pChange.Zf},
                [3]float64{pChange.RotationXf, pChange.RotationYf, pChange.RotationZf},
                r.pivotsLst[len(r.pivotsLst)-1])
            r.coordOriginsLst = append(r.coordOriginsLst, p)
        case "pop":
            if len(r.coordOriginsLst) > 1 {
                r.coordOriginsLst = r.coordOriginsLst[:len(r.coordOriginsLst)-1]
            }
        }

    case "rotation_pivot":
        switch pChange.SetterTypeStr {
        case "push":
            if pChange.AxisTypeStr == "current_pos" {
                p := realWorldPos(
                    [3]float64{pChange.Xf, pChange.Yf, pChange.Zf},
                    [3]float64{pChange.RotationXf, pChange.RotationYf, pChange.RotationZf},
                    r.pivotsLst[len(r.pivotsLst)-1])
                r.pivotsLst = append(r.pivotsLst, p)
            }
        case "pop":
            if len(r.pivotsLst) > 1 {
                r.pivotsLst = r.pivotsLst[:len(r.pivotsLst)-1]
            }
        }

    case "color_rgb_lst":
        c := pChange.ColorRGBlst
        return []interface{}{c[0], c[1], c[2]}

    case "color_rgb_hex":
        rgb, ok := parseHexColor(pChange.ColorRGBhexStr)
        if !ok {
            return []interface{}{0.0, 0.0, 0.0}
        }
        return []interface{}{rgb[0], rgb[1], rgb[2]}

    case "color_background_rgb_lst":
        c := pChange.ColorBackgroundRGBlst
        r.Scene.Background = []float64{c[0], c[1], c[2]}

    case "color_background_rgb_hex":
        if rgb, ok := parseHexColor(pChange.ColorBackgroundHexStr); ok {
            r.Scene.Background = rgb[:]
        }

    case "material_type":
        switch pChange.MaterialTypeStr {
        case "wireframe":
            r.wireframeBool = pChange.MaterialValueBool
        case "shader":
            // like the browser engine, once a shader is set it stays active for all subsequent entities
            r.shaderActiveBool = true
            r.shaderActiveStr  = pChange.MaterialValueStr
        }

    case "material_shader_uniform":
        // browser engine applies it to the currently active shader (not to the named one)
        if r.shaderActiveBool {
            if m, ok := r.uniformsMap[r.shaderActiveStr]; ok {
                m[pChange.MaterialShaderUniformNameStr] = pChange.MaterialShaderUniformVal
            }
        }

    case "line_cmd":
        if pChange.LineCmdStr == "start" {
            r.linePointsLst = nil
        }
    }
    return nil
}

//-------------------------------------------------

func (r *GFsceneRecorder) full() bool {
    if r.MaxEntitiesInt > 0 && len(r.Scene.Entities) >= r.MaxEntitiesInt {
        r.Scene.TruncatedBool = true
        return true
    }
    return false
}

func (r *GFsceneRecorder) entityFun(pKindStr string) GFgeometryFunc {
    return func(x, y, z, rx, ry, rz, sx, sy, sz, cr, cg, cb float64) {
        if r.full() {
            return
        }

        // position relative to the latest coord origin
        o := r.coordOriginsLst[len(r.coordOriginsLst)-1]
        pos := [3]float64{o[0] + x, o[1] + y, o[2] + z}
        rot := [3]float64{rx, ry, rz}

        var world [3]float64
        if pKindStr == "cube" {
            world = realWorldPos(pos, rot, r.pivotsLst[len(r.pivotsLst)-1])
        } else {
            // the browser engine rotates spheres around the world origin (ignores pivots)
            world = realWorldPos(pos, rot, [3]float64{0, 0, 0})
        }

        e := &GFsceneEntity{
            KindStr:   pKindStr,
            Position:  world,
            Rotation:  rot,
            Scale:     [3]float64{sx, sy, sz},
            Color:     [3]float64{cr, cg, cb},
            Wireframe: r.wireframeBool,
            Material:  GFsceneMaterial{TypeStr: "basic"},
        }
        if r.shaderActiveBool {
            uniforms := map[string]interface{}{}
            for k, v := range r.uniformsMap[r.shaderActiveStr] {
                uniforms[k] = v
            }
            e.Material = GFsceneMaterial{TypeStr: "shader", ShaderStr: r.shaderActiveStr, Uniforms: uniforms}
        }
        if pKindStr == "cube" && r.animation != nil {
            e.Animation = r.animation
        }
        r.Scene.Entities = append(r.Scene.Entities, e)
    }
}

func (r *GFsceneRecorder) lineFun(x, y, z, rx, ry, rz, sx, sy, sz, cr, cg, cb float64) {
    o := r.coordOriginsLst[len(r.coordOriginsLst)-1]
    p := realWorldPos(
        [3]float64{o[0] + x, o[1] + y, o[2] + z},
        [3]float64{rx, ry, rz},
        r.pivotsLst[len(r.pivotsLst)-1])

    // a previous point exists, so this one completes a segment
    if len(r.linePointsLst) > 0 && !r.full() {
        r.Scene.Entities = append(r.Scene.Entities, &GFsceneEntity{
            KindStr:  "line",
            Scale:    [3]float64{1, 1, 1},
            Color:    [3]float64{0, 0, 0}, // browser engine draws lines black
            Material: GFsceneMaterial{TypeStr: "basic"},
            Points:   [][3]float64{r.linePointsLst[0], p},
        })
    }
    r.linePointsLst = [][3]float64{p}
}

func (r *GFsceneRecorder) animate(pPropsLst []map[string]interface{}, pDurationSecF float64, pRepeatBool bool) {
    a := &GFsceneAnimation{DurationSecF: pDurationSecF, RepeatBool: pRepeatBool}
    for _, p := range pPropsLst {
        name, _  := p["name_str"].(string)
        start, _ := p["start_val_f"].(float64)
        end, _   := p["end_val_f"].(float64)
        a.Props = append(a.Props, GFsceneAnimationProp{name, start, end})
    }
    r.animation = a
}

//-------------------------------------------------
// MATH - mirrors ts/engine/gf_engine_utils.ts

// rotation matrix Rx*Ry*Rz applied to a point
func rotateWorld(p [3]float64, rx, ry, rz float64) [3]float64 {
    // Rz first, then Ry, then Rx (M = Rx*Ry*Rz)
    c, s := math.Cos(rz), math.Sin(rz)
    p = [3]float64{c*p[0] - s*p[1], s*p[0] + c*p[1], p[2]}
    c, s = math.Cos(ry), math.Sin(ry)
    p = [3]float64{c*p[0] + s*p[2], p[1], -s*p[0] + c*p[2]}
    c, s = math.Cos(rx), math.Sin(rx)
    p = [3]float64{p[0], c*p[1] - s*p[2], s*p[1] + c*p[2]}
    return p
}

func realWorldPos(pPos [3]float64, pRot [3]float64, pPivot [3]float64) [3]float64 {
    d := [3]float64{pPos[0] - pPivot[0], pPos[1] - pPivot[1], pPos[2] - pPivot[2]}
    r := rotateWorld(d, pRot[0], pRot[1], pRot[2])
    return [3]float64{r[0] + pPivot[0], r[1] + pPivot[1], r[2] + pPivot[2]}
}

// "#rrggbb" or "#rgb" -> 0-1 channels
func parseHexColor(pHexStr string) ([3]float64, bool) {
    h := strings.TrimPrefix(strings.TrimSpace(pHexStr), "#")
    if len(h) == 3 {
        h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
    }
    if len(h) != 6 {
        return [3]float64{}, false
    }
    v, err := strconv.ParseUint(h, 16, 32)
    if err != nil {
        return [3]float64{}, false
    }
    return [3]float64{
        float64((v>>16)&0xff) / 255,
        float64((v>>8)&0xff)  / 255,
        float64(v&0xff)       / 255,
    }, true
}
