/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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

package gf_core

import (
	"fmt"
	"crypto/md5"
	"encoding/hex"
)

//---------------------------------------------------

type GF_ID string

//---------------------------------------------------
// CREATES_ID

func IDcreate(pUniqueValsForIDlst []string,
	pUNIXtimeF float64) GF_ID {
	
	h := md5.New()

	h.Write([]byte(fmt.Sprint(pUNIXtimeF)))

	for _, v := range pUniqueValsForIDlst {
		h.Write([]byte(v))
	}

	sum    := h.Sum(nil)
	hexStr := hex.EncodeToString(sum)
	
	gfIDstr := GF_ID(hexStr)
	return gfIDstr
}