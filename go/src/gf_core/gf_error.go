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
	"runtime"
	"runtime/debug"
	"github.com/globalsign/mgo/bson"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
type Gf_error struct {
	Id                   bson.ObjectId `bson:"_id,omitempty"`
	Id_str               string        `bson:"id_str"` 
	T_str                string        `bson:"t"`                    //"gf_error"
	Creation_unix_time_f float64       `bson:"creation_unix_time_f"`
	Type_str             string                  `bson:"type_str"`
	User_msg_str         string                  `bson:"user_msg_str"`
	Data_map             *map[string]interface{} `bson:"data_map"`
	Descr_str            string                  `bson:"descr_str"`
	Error                error                   `bson:"error"`
	Service_name_str     string        `bson:"service_name_str"`   
	Subsystem_name_str   string        `bson:"subsystem_name_str"`   //major portion of functionality, a particular package, or a logical group of functions
	Stack_trace_str      string        `bson:"stack_trace_str"`
	Function_name_str    string        `bson:"func_name_str"`
	File_str             string        `bson:"file_str"`
	Line_num_int         int           `bson:"line_num_int"`
}

//-------------------------------------------------
func Error__create(p_user_msg_str string,
	p_error_type_str     string,
	p_error_data_map     *map[string]interface{},
	p_error              error,
	p_subsystem_name_str string,
	p_runtime_sys        *Runtime_sys) *Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_error.Error__create()")


	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := fmt.Sprintf("%s:%s",p_error_type_str,creation_unix_time_f)
	stack_trace_str      := string(debug.Stack())

	//IMPORTANT!! - number of stack frames to skip before recording. without skipping 
	//              we would get info on this function, not its caller which is where
	//              the error occured.
	skip_stack_frames_num_int := 1

	//https://golang.org/pkg/runtime/#Caller
	program_counter,file_str,line_num_int,_ := runtime.Caller(skip_stack_frames_num_int)

	//FuncForPC - returns a *Func describing the function that contains the given program counter address
	function          := runtime.FuncForPC(program_counter)
	function_name_str := function.Name()

	//--------------------
	//VIEW
	p_runtime_sys.Log_fun("ERROR",fmt.Sprintf("gf_error created - type:%s - service:%s - subsystem:%s - func:%s - usr_msg:%s",
				p_error_type_str,
				p_runtime_sys.Service_name_str,
				p_subsystem_name_str,
				function_name_str,
				p_user_msg_str))
	//--------------------
	//ERROR_DEF

	error_defs_map := error__get_defs()
	error_def,ok   := error_defs_map[p_error_type_str]
	if !ok {
		panic(fmt.Sprintf("unknown gf_error type encountered - %s",p_error_type_str))
	}
	//--------------------


	gf_error := Gf_error{
		Id_str:              id_str,
		T_str:               "gf_error",
		Creation_unix_time_f:creation_unix_time_f,
		Type_str:            p_error_type_str,
		User_msg_str:        p_user_msg_str,
		Data_map:            p_error_data_map,
		Descr_str:           error_def.Descr_str,
		Error:               p_error,
		Service_name_str:    p_runtime_sys.Service_name_str,
		Subsystem_name_str:  p_subsystem_name_str,
		Stack_trace_str:     stack_trace_str,
		Function_name_str:   function_name_str,
		File_str:            file_str,
		Line_num_int:        line_num_int,
	}

	fmt.Println("")
	fmt.Println("  FAILED FUNCTION ------------- "+function_name_str)
	fmt.Println("")
	fmt.Println(stack_trace_str)
	fmt.Println("")
	spew.Dump(gf_error)
	fmt.Println("")
	fmt.Println("")
	//--------------------
	//PERSIST
	err := p_runtime_sys.Mongodb_coll.Insert(gf_error)
	if err != nil {

	}
	//--------------------

	return &gf_error
}
