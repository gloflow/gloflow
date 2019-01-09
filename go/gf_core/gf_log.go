/*
GloFlow media management/publishing system
Copyright (C) 2019 Ivan Trajkovic

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
	"time"
	"strconv"
	"github.com/fatih/color"
)
//-------------------------------------------------
func Init_log_fun() func(string,string) {

	green  := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red    := color.New(color.FgRed).SprintFunc()

	log_fun := func(p_g string,p_m string) {
		t_str := strconv.FormatFloat(float64(time.Now().UnixNano())/1000000000.0,'f',10,64)

		if p_g == "FUN_ENTER" {
			fmt.Printf(t_str+":"+yellow(p_g)+":"+p_m+"\n")
		} else if p_g == "INFO" {
			fmt.Printf(t_str+":"+green(p_g)+":"+green(p_m)+"\n")
		} else if p_g == "ERROR" {
			fmt.Printf(t_str+":"+red(p_g)+":"+p_m+"\n")
		}
	}
	return log_fun
}