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

package main

import(
	"fmt"
	"strings"
	"regexp"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

func ParseProgramASTfromFile(pLocalFilePathStr string,
	pRuntimeSys *gf_core.RuntimeSys) (gf_lang.GFexpr, *gf_core.GFerror) {

	//------------------------
	// READ_FILE
	programCodeStr, gfErr := gf_core.FileRead(pLocalFilePathStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	//------------------------
	// REMOVE_COMMENTS
	commentRegex := regexp.MustCompile(`(?m)(.*)//.*$`)

	cleanJSONcodeStr := ""
	for _, lineStr := range strings.Split(programCodeStr, "\n") {

		lineNoCommentsStr := commentRegex.ReplaceAllString(lineStr, "$1")

		if strings.TrimSpace(lineNoCommentsStr) != "" {
			cleanJSONcodeStr += fmt.Sprintf("%s\n", lineNoCommentsStr)
		}
	}

	fmt.Println("clean JSON code:", cleanJSONcodeStr)

	//------------------------
	// PARSE

	code, gfErr := gf_core.ParseJSONfromString(cleanJSONcodeStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	codeUncastedLst := code.([]interface{})
	codeLst := gf_lang.CastToExpr(codeUncastedLst)

	//------------------------
    
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println(codeLst)


	return codeLst, nil
}