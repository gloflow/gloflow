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