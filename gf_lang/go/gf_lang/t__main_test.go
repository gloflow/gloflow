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
	"os"
	"fmt"
	"testing"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------

func TestMain(m *testing.M) {

	os.Setenv("GF_LOG_LEVEL", "debug")
	v := m.Run()
	os.Exit(v)
}

//---------------------------------------------------

func TestTreeExpansion(pTest *testing.T) {
	_, logNewFun := gf_core.LogsInit()
	externAPI := GetTestExternAPI()

	// program
	programASTlst := GFexpr{
		GFexpr{"lang_v", "0.0.6"},
		GFexpr{"*", 10, GFexpr{GFexpr{"x", 5.4}, GFexpr{"*", 1, GFexpr{GFexpr{"z", 1.6,}, "cube"},},},}, // 3 * {x 5.4} 1 * {z 1.6} cube
	}

	// run program
	resultsLst, programsDebugLst, err := Run(programASTlst, externAPI)
	if err != nil {
		logNewFun("ERROR", "failed to run program to test State Operations",
			map[string]interface{}{"err": err,})
		pTest.Fail()
	}

	fmt.Println("============================")
	fmt.Println("results:")
	spew.Dump(resultsLst)
	fmt.Println("-----")
	viewDebug(programASTlst, programsDebugLst)
}

//---------------------------------------------------

func TestStateOps(pTest *testing.T) {
	_, logNewFun := gf_core.LogsInit()
	externAPI := GetTestExternAPI()

	// program
	programASTlst := GFexpr{
		GFexpr{"lang_v", "0.0.6"},
		GFexpr{"set", "color-background", GFexpr{"rgb", 0.5, 0.5, 0.5,},},
		GFexpr{GFexpr{"set", "color", GFexpr{"rgb", 0, 0, 1.0},}, "cube"},                   // cube
    	GFexpr{GFexpr{"set", "color", GFexpr{"rgb", 0, 0, 1.0},}, GFexpr{"x", 2.0}, "cube"}, // {x 2} cube
	}

	// run program
	resultsLst, programsDebugLst, err := Run(programASTlst, externAPI)
	if err != nil {
		logNewFun("ERROR", "failed to run program to test State Operations",
			map[string]interface{}{"err": err,})
		pTest.Fail()
	}

	fmt.Println("============================")
	fmt.Println("results:")
	spew.Dump(resultsLst)
	fmt.Println("-----")
	viewDebug(programASTlst, programsDebugLst)
}

//---------------------------------------------------

func TestMap(pTest *testing.T) {
	_, logNewFun := gf_core.LogsInit()
	externAPI := GetTestExternAPI()

	// program
	programASTlst := GFexpr{
		GFexpr{"lang_v", "0.0.6"},

		GFexpr{
			GFexpr{"$test_map", GFexpr{"make", GFexpr{"map", 
				GFexpr{
					GFexpr{"a", "b"},
					GFexpr{"c", 2},
				}},
			}},
			GFexpr{"return", "$test_map"},
		},
	}

	// run program
	resultsLst, _, err := Run(programASTlst, externAPI)
	if err != nil {
		logNewFun("ERROR", "failed to run program to test Maps", map[string]interface{}{"err": err,})
		pTest.Fail()
	}


	spew.Dump(resultsLst)

	if len(resultsLst) != 1 {
		logNewFun("ERROR", "results list should be 1 elements", nil)
		pTest.Fail()
	}

	testMap := resultsLst[0].(map[string]interface{})
	spew.Dump(testMap)

	if testMap["a"].(string) != "b" {
		logNewFun("ERROR", "first result map key 'a' should be 'b'", nil)
		pTest.Fail()
	}

	if testMap["c"].(int) != 2 {
		logNewFun("ERROR", "first result map key 'c' should be 2", nil)
		pTest.Fail()
	}
}

//---------------------------------------------------

func TestLists(pTest *testing.T) {
	_, logNewFun := gf_core.LogsInit()
	externAPI := GetTestExternAPI()

	// program
	programASTlst := GFexpr{
		GFexpr{"lang_v", "0.0.6"},
		
		GFexpr{
			"return", GFexpr{
				GFexpr{"$test_list",        GFexpr{"make", GFexpr{"list", GFexpr{1, 2, 3,},},},},
				GFexpr{"$test_list_length", GFexpr{"len", GFexpr{"$test_list",},},},
				GFexpr{"return", "$test_list_length"},
			},
		},
	}

	// run program
	resultsLst, _, err := Run(programASTlst, externAPI)
	if err != nil {
		logNewFun("ERROR", "failed to run program to test Lists", map[string]interface{}{"err": err,})
		pTest.Fail()
	}

	if len(resultsLst) != 1 {
		logNewFun("ERROR", "results list should be 1 elements", nil)
		pTest.Fail()
	}
	if resultsLst[0] != 3 {
		logNewFun("ERROR", "first results should be equal to 3", map[string]interface{}{"result": resultsLst[0],})
		pTest.Fail()
	}
}

//---------------------------------------------------

func TestVariables(pTest *testing.T) {
	
	_, logNewFun := gf_core.LogsInit()
	externAPI := GetTestExternAPI()

	// program
	programASTlst := GFexpr{
		GFexpr{"lang_v", "0.0.6"},
		
		// testing setting of a user variable, and referencing
		// and returning it.
		GFexpr{
			"return", GFexpr{
					GFexpr{"$test_var", 10},
					GFexpr{"return", "$test_var"},
				},
		},
	}

	// run program
	resultsLst, _, err := Run(programASTlst, externAPI)
	if err != nil {
		logNewFun("ERROR", "failed to run program to test Variables", map[string]interface{}{"err": err,})
		pTest.Fail()
	}

	spew.Dump(resultsLst)

	if len(resultsLst) != 1 {
		logNewFun("ERROR", "results list should be 1 elements", nil)
		pTest.Fail()
	}
	if resultsLst[0] != 10 {
		logNewFun("ERROR", "first results should be equal to 10", map[string]interface{}{"result": resultsLst[0],})
		pTest.Fail()
	}
}

//---------------------------------------------------

func TestReturnExpression(pTest *testing.T) {
	
	_, logNewFun := gf_core.LogsInit()
	externAPI := GetTestExternAPI()

	// program
	programASTlst := GFexpr{
		GFexpr{"lang_v", "0.0.6"},
		GFexpr{"return", 10},
		GFexpr{"return", GFexpr{"*", 3, 12}},
		GFexpr{"return", GFexpr{"return", GFexpr{"*", 10, 2}}},
	}

	// run program
	resultsLst, _, err := Run(programASTlst, externAPI)
	if err != nil {
		logNewFun("ERROR", "failed to run program to test Return statements", map[string]interface{}{"err": err,})
		pTest.Fail()
	}

	spew.Dump(resultsLst)

	if len(resultsLst) != 3 {
		logNewFun("ERROR", "results list should be 3 elements", nil)
		pTest.Fail()
	}

	if resultsLst[0] != 10 {
		logNewFun("ERROR", "first results should be equal to 10", map[string]interface{}{"result": resultsLst[0],})
		pTest.Fail()
	}

	if int(resultsLst[1].(float64)) != 36 {
		logNewFun("ERROR", "first results should be equal to 36", map[string]interface{}{"result": resultsLst[1],})
		pTest.Fail()
	}

	if int(resultsLst[2].(float64)) != 20 {
		logNewFun("ERROR", "first results should be equal to 20", map[string]interface{}{"result": resultsLst[2],})
		pTest.Fail()
	}
}