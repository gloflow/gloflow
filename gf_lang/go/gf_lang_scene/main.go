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

// gf_lang_scene - runs a gf_lang program with the scene recorder and writes the
//                 scene graph as JSON. Optionally pipes it to the headless wgpu
//                 renderer (gf_lang_render, see gf_lang/rust) to produce a PNG.
//
//   gf_lang_scene -in prog.gf -out scene.json
//   gf_lang_scene -in prog.gf -render out.png [-renderer path/to/gf_lang_render] [-width 1024 -height 768] [-fit]

package main

import (
    "os"
    "fmt"
    "flag"
    "regexp"
    "strings"
    "os/exec"
    "bytes"
    "encoding/json"
    "path/filepath"
    "github.com/gloflow/gloflow/gf_lang/go/gf_lang"
)

//-------------------------------------------------

func parseProgram(pPathStr string) (gf_lang.GFexpr, error) {

    codeBytes, err := os.ReadFile(pPathStr)
    if err != nil {
        return nil, err
    }

    // same comment handling as gf_lang_server
    commentRegex := regexp.MustCompile(`(?m)(.*)//.*$`)
    cleanStr := ""
    for _, lineStr := range strings.Split(string(codeBytes), "\n") {
        noCommentsStr := commentRegex.ReplaceAllString(lineStr, "$1")
        if strings.TrimSpace(noCommentsStr) != "" {
            cleanStr += noCommentsStr + "\n"
        }
    }

    var code []interface{}
    if err := json.Unmarshal([]byte(cleanStr), &code); err != nil {
        return nil, fmt.Errorf("failed to parse program JSON: %w", err)
    }
    return gf_lang.CastToExpr(code), nil
}

//-------------------------------------------------

func main() {

    inStr      := flag.String("in", "", "gf_lang program file (required)")
    outStr     := flag.String("out", "", "scene JSON output file (default: stdout, unless -render is set)")
    renderStr  := flag.String("render", "", "render the scene to this PNG file via the gf_lang_render binary")
    rendererStr := flag.String("renderer", "", "path to the gf_lang_render binary (default: next to this binary, then rust/build)")
    widthInt   := flag.Int("width", 1024, "render width")
    heightInt  := flag.Int("height", 768, "render height")
    fitBool    := flag.Bool("fit", false, "frame all entities in the render (static camera, no trackball controls)")
    maxEntInt  := flag.Int("max-entities", 1000000, "cap on the number of entities in the scene (0 == unlimited)")
    flag.Parse()

    if *inStr == "" {
        flag.Usage()
        os.Exit(2)
    }

    programASTlst, err := parseProgram(*inStr)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }

    // the interpreter prints debug info to stdout, so redirect it to stderr
    // while running, to keep stdout clean for the scene JSON.
    stdout := os.Stdout
    os.Stdout = os.Stderr
    scene, err := gf_lang.SceneRun(programASTlst, *maxEntInt)
    os.Stdout = stdout
    if err != nil {
        fmt.Fprintln(os.Stderr, "error: program run failed:", err)
        os.Exit(1)
    }

    sceneJSON, err := scene.JSON()
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }

    //------------------------
    // JSON output
    if *outStr != "" {
        if err := os.WriteFile(*outStr, sceneJSON, 0644); err != nil {
            fmt.Fprintln(os.Stderr, "error:", err)
            os.Exit(1)
        }
    } else if *renderStr == "" {
        os.Stdout.Write(sceneJSON)
        fmt.Println()
    }

    //------------------------
    // RENDER
    if *renderStr != "" {
        binStr := findRenderer(*rendererStr)
        args := []string{
            "--out", *renderStr,
            "--width", fmt.Sprint(*widthInt),
            "--height", fmt.Sprint(*heightInt),
        }
        if *fitBool {
            args = append(args, "--fit")
        }
        cmd := exec.Command(binStr, args...)
        cmd.Stdin  = bytes.NewReader(sceneJSON)
        cmd.Stdout = os.Stderr
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
            fmt.Fprintln(os.Stderr, "error: render failed:", err)
            os.Exit(1)
        }
        fmt.Fprintln(os.Stderr, "rendered", *renderStr)
    }
}

//-------------------------------------------------

func findRenderer(pExplicitStr string) string {
    if pExplicitStr != "" {
        return pExplicitStr
    }
    exe, err := os.Executable()
    if err == nil {
        dir := filepath.Dir(exe)
        for _, c := range []string{
            filepath.Join(dir, "gf_lang_render"),
            filepath.Join(dir, "..", "..", "rust", "build", "gf_lang_render"),
        } {
            if _, err := os.Stat(c); err == nil {
                return c
            }
        }
    }
    return "gf_lang_render" // PATH
}
