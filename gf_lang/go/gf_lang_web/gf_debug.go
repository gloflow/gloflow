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

import (
	"fmt"
	"github.com/gloflow/gloflow/gf_lang/go/gf_lang"
)

//-------------------------------------------------

func debugView(pProgramsDebugLst []*gf_lang.GFprogramDebug) {

	fmt.Println(">>>>>>>>>>>>> state history:")
	i:=0
	for _, p := range pProgramsDebugLst {

		// fmt.Println(">>>>>>> program")

		for _, s := range p.StateHistoryLst {
			fmt.Printf("x %f  time %f\n", s.Xf, s.CreationUNIXtimeF)
		}
		i+=1
	}


	fmt.Println(">>>>>>>>>>>>>> program outputs:",)
	j:=0
	for _, p := range pProgramsDebugLst {

		// fmt.Println(">>>>>>> program", p.OutputLst)

		for _, o := range p.OutputLst {

			if v, ok := o.(*gf_lang.GFentityOutput); ok {
				fmt.Printf("%s x %f\n", v.TypeStr, v.Props.Xf)
				j+=1
			}
		}
	}

	fmt.Printf("states # %d\n", i)
	fmt.Printf("cubes  # %d\n", j)
}