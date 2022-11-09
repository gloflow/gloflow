/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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
	"os"
	"strconv"
)

//-------------------------------------------------
func getConfig() (int, int) {

	//---------------------
	// P2P_PORT
	var p2pPortInt int
	p2pPortStr, ok := os.LookupEnv("GF_P2P_PORT")
	if ok {
		var err error
		p2pPortInt, err = strconv.Atoi(p2pPortStr)
		if err != nil {
			panic(err)
		}	
	} else {
		// start on a random port
		p2pPortInt = 0
	}
	
	//---------------------
	// HTP_PORT
	var httpPortInt int
	httpPortStr, ok := os.LookupEnv("GF_HTTP_PORT")
	if ok {
		var err error
		httpPortInt, err = strconv.Atoi(httpPortStr)
		if err != nil {
			panic(err)
		}	
	} else {
		// start on a random port
		httpPortInt = 3000
	}
	
	//---------------------

	return p2pPortInt, httpPortInt
}