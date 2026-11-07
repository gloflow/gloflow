package gf_lang

import (
	"encoding/json"
	"os"
	"reflect"
	"path/filepath"
	"strings"
	"testing"
)

func loadJSONexample(t *testing.T, pPathStr string) GFexpr {
	b, err := os.ReadFile(pPathStr)
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, l := range splitLines(string(b)) {
		lines = append(lines, gfSynStripComment(l))
	}
	var v []interface{}
	if err := json.Unmarshal([]byte(joinLines(lines)), &v); err != nil {
		if err := json.Unmarshal([]byte("["+joinLines(lines)+"]"), &v); err != nil {
			t.Fatal(err)
		}
	}
	return CastToExpr(v)
}

func splitLines(s string) []string { return strings.Split(s, "\n") }
func joinLines(l []string) string  { return strings.Join(l, "\n") }

func TestSyntaxMatchesJSON(t *testing.T) {
	for _, name := range []string{"first_scene", "second_scene_rules_test", "third_scene_multirule_test"} {
		src, err := os.ReadFile("../../examples/" + name + ".gfl")
		if err != nil {
			t.Fatal(err)
		}
		got, err := ParseSyntax(string(src))
		if err != nil {
			t.Fatal(err)
		}
		want := loadJSONexample(t, "../../examples/"+name+".gf")
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s mismatch\n got: %v\nwant: %v", name, got, want)
		}
	}
}

func TestSyntaxConstructs(t *testing.T) {
	src := "rule R3 max=30:\n    if $i == 4:\n        R4\n    sx=(-0.001 * $i)\n    x=rand(1, 5) cube\n"
	got, err := ParseSyntax(src)
	if err != nil {
		t.Fatal(err)
	}
	var want []interface{}
	json.Unmarshal([]byte(`[["rule","R3",[["iters_max",30]],[["if",["==","$i",4],[["R4"]]],[["sx",["*",-0.001,"$i"]]],[["x",["rand",1,5]],"cube"]]]]`), &want)
	if !reflect.DeepEqual(got, CastToExpr(want)) {
		t.Errorf("got: %v\nwant: %v", got, CastToExpr(want))
	}
}

func TestSyntaxShader(t *testing.T) {
	src, _ := os.ReadFile("../../examples/shader_scene.gfl")
	got, err := ParseSyntax(string(src))
	if err != nil {
		t.Fatal(err)
	}
	sh := got[len(got)-1].(GFexpr)
	if sh[0] != "shader" || sh[2].(GFexpr)[0] != "uniforms" || sh[3].(GFexpr)[0] != "vertex" {
		t.Errorf("bad shader: %v", sh)
	}
}

func TestSyntaxErrors(t *testing.T) {
	for _, s := range []string{"\tcube", "x=", "repeat 3\n", "(1 2", "a = "} {
		if _, err := ParseSyntax(s); err == nil {
			t.Errorf("expected error for %q", s)
		}
	}
}

func TestSyntaxChainAndAtoms(t *testing.T) {
	src := "x=9 z=2 | y=3 | cube\ncolor = rgb(1 0 0) |\n    x=1 sphere\n\"R1\"\n"
	got, err := ParseSyntax(src)
	if err != nil {
		t.Fatal(err)
	}
	var want []interface{}
	json.Unmarshal([]byte(`[[["x",9,"z",2],["y",3],"cube"],[["set","color",["rgb",1,0,0]],["x",1],"sphere"],"R1"]`), &want)
	if !reflect.DeepEqual(got, CastToExpr(want)) {
		t.Errorf("got: %v", got)
	}
}

// every .gf example translates to .gfl and back to the identical AST
func TestSyntaxFormatAllExamples(t *testing.T) {
	for _, pattern := range []string{"../../examples/*.gf", "../../examples/primary/*.gf"} {
		files, _ := filepath.Glob(pattern)
		for _, f := range files {
			ast := loadJSONexample(t, f)
			out, err := FormatSyntax(ast)
			if err != nil {
				t.Errorf("%s: %v", f, err)
				continue
			}
			want, err := os.ReadFile(f + "l")
			if err != nil || string(want) != out {
				t.Errorf("%sl is out of date, regenerate with gf_lang_convert", f)
			}
		}
	}
}
