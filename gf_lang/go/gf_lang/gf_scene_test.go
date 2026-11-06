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
    "testing"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestSceneRecorder(t *testing.T) {

    // cubes with a coord origin, color, wireframe and a background
    program := GFexpr{
        GFexpr{"lang_v", 0.3},
        GFexpr{"set", "color-background", GFexpr{"rgb", 0.5, 0.25, 0.0}},
        GFexpr{"set", "material", GFexpr{"wireframe", true}},
        GFexpr{GFexpr{"set", "color", GFexpr{"rgb", 0.0, 0.0, 1.0}}, GFexpr{"x", 2.0, "y", 3.0}, "cube"},
        GFexpr{GFexpr{"x", 5.0}, "sphere"},
    }

    scene, err := SceneRun(program, 0)
    if err != nil {
        t.Fatal(err)
    }
    if len(scene.Entities) != 2 {
        t.Fatalf("expected 2 entities, got %d", len(scene.Entities))
    }
    if scene.Background == nil || scene.Background[1] != 0.25 {
        t.Fatalf("bad background %v", scene.Background)
    }

    c := scene.Entities[0]
    if c.KindStr != "cube" || !approx(c.Position[0], 2) || !approx(c.Position[1], 3) {
        t.Fatalf("bad cube %+v", c)
    }
    if c.Color != [3]float64{0, 0, 1} || !c.Wireframe || c.Material.TypeStr != "basic" {
        t.Fatalf("bad cube material/color %+v", c)
    }
    if scene.Entities[1].KindStr != "sphere" {
        t.Fatalf("expected a sphere, got %s", scene.Entities[1].KindStr)
    }

    // serializes
    if _, err := scene.JSON(); err != nil {
        t.Fatal(err)
    }
}

func TestSceneRecorderEntitiesCap(t *testing.T) {
    program := GFexpr{
        GFexpr{"lang_v", 0.3},
        GFexpr{"*", 10, GFexpr{GFexpr{"x", 1.0}, "cube"}},
    }
    scene, err := SceneRun(program, 3)
    if err != nil {
        t.Fatal(err)
    }
    if len(scene.Entities) != 3 || !scene.TruncatedBool {
        t.Fatalf("expected 3 truncated entities, got %d (truncated=%v)", len(scene.Entities), scene.TruncatedBool)
    }
}

func TestSceneRotateWorld(t *testing.T) {
    // rotating (1,0,0) 90deg around Z gives (0,1,0)
    p := rotateWorld([3]float64{1, 0, 0}, 0, 0, math.Pi/2)
    if !approx(p[0], 0) || !approx(p[1], 1) {
        t.Fatalf("bad rotation %v", p)
    }
    // around a pivot
    q := realWorldPos([3]float64{2, 0, 0}, [3]float64{0, 0, math.Pi / 2}, [3]float64{1, 0, 0})
    if !approx(q[0], 1) || !approx(q[1], 1) {
        t.Fatalf("bad pivot rotation %v", q)
    }
}

func TestParseHexColor(t *testing.T) {
    c, ok := parseHexColor("#ff8000")
    if !ok || !approx(c[0], 1) || !approx(c[2], 0) || math.Abs(c[1]-128.0/255) > 1e-9 {
        t.Fatalf("bad color %v %v", c, ok)
    }
    if _, ok := parseHexColor("nope"); ok {
        t.Fatal("expected failure")
    }
}
