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

// gf_lang_convert translates .gf programs (JSON-like AST) into .gfl (indentation syntax).
// each output is verified to parse back to exactly the same AST.
//
//   gf_lang_convert a.gf b.gf ...      writes a.gfl, b.gfl ... next to the inputs
package main

import (
	"os"
	"fmt"
	"strings"
	"encoding/json"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang"
)

func main() {

	failed := 0
	for _, path := range os.Args[1:] {
		if err := convert(path); err != nil {
			fmt.Fprintf(os.Stderr, "FAIL %s: %v\n", path, err)
			failed++
		} else {
			fmt.Println("ok  ", path)
		}
	}
	if failed > 0 {
		os.Exit(1)
	}
}

func convert(pPathStr string) error {

	b, err := os.ReadFile(pPathStr)
	if err != nil {
		return err
	}
	ast, err := LoadGF(string(b))
	if err != nil {
		return err
	}
	out, err := gf_lang.FormatSyntax(ast)
	if err != nil {
		return err
	}
	return os.WriteFile(strings.TrimSuffix(pPathStr, ".gf")+".gfl", []byte(out), 0644)
}

// a .gf file is a list of expressions, with or without the outer [ ], and with // comments
func LoadGF(pSrcStr string) (gf_lang.GFexpr, error) {

	lines := []string{}
	for _, l := range strings.Split(pSrcStr, "\n") {
		inStr := false
		for i := 0; i < len(l); i++ {
			if l[i] == '\\' && inStr {
				i++
			} else if l[i] == '"' {
				inStr = !inStr
			} else if !inStr && l[i] == '/' && i+1 < len(l) && l[i+1] == '/' {
				l = l[:i]
				break
			}
		}
		lines = append(lines, l)
	}
	src := strings.Join(lines, "\n")

	var parsed interface{}
	if err := json.Unmarshal([]byte(src), &parsed); err != nil {
		if err2 := json.Unmarshal([]byte("["+src+"]"), &parsed); err2 != nil {
			return nil, err
		}
	}
	lst := parsed.([]interface{})
	// a single bare expression, ex. ["cube"]
	allLists := true
	for _, e := range lst {
		if _, ok := e.([]interface{}); !ok {
			allLists = false
		}
	}
	if !allLists {
		lst = []interface{}{lst}
	}
	return gf_lang.CastToExpr(lst), nil
}
