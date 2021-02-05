/*
GloFlow application and media management/publishing platform
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
	"os" 
	"time"
	"strings"
	"encoding/json"
	"os/exec"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/globalsign/mgo"
)


//-------------------------------------------------
func Mongo__insert(p_record interface{},
	p_coll_name_str string,
	p_ctx           *context.Context,
	p_runtime_sys   *Runtime_sys) *Gf_error {

	_, err := p_runtime_sys.Mongo_db.Collection(p_coll_name_str).InsertOne(*p_ctx, p_record)
	if err != nil {
		gf_err := Mongo__handle_error("failed to insert a new record into the DB",
			"mongodb_insert_error",
			map[string]interface{}{"coll_name_str": p_coll_name_str,},
			err, "gf_core", p_runtime_sys)
		
		return gf_err
	}
	return nil
}

//-------------------------------------------------
func Mongo__ensure_index(p_indexes_keys_lst [][]string, 
	p_coll_name_str string,
	p_runtime_sys   *Runtime_sys) []*Gf_error {

	gf_errs_lst := []*Gf_error{}
	for _, index_keys_lst := range p_indexes_keys_lst {
		doc_type__index := mgo.Index{
			Key:        index_keys_lst, 
			Unique:     false, // index must necessarily contain only a single document per Key
			DropDups:   false, // documents with the same key as a previously indexed one will be dropped rather than an error returned.
			Background: true,  // other connections will be allowed to proceed using the collection without the index while it's being built
			Sparse:     true,  // only documents containing the provided Key fields will be included in the index
		}
	
		err := p_runtime_sys.Mongodb_db.C(p_coll_name_str).EnsureIndex(doc_type__index)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "duplicate key error index") {
				continue // ignore, index already exists
			} else {
				gf_err := Mongo__handle_error(fmt.Sprintf("failed to create db index on fields - %s", index_keys_lst), 
					"mongodb_ensure_index_error", nil, err, "gf_core", p_runtime_sys)
				gf_errs_lst = append(gf_errs_lst, gf_err)
			}
		}
	}
	return gf_errs_lst
}

//-------------------------------------------------
func Mongo__connect_new(p_mongo_server_url_str string,
	p_db_name_str string,
	p_runtime_sys *Runtime_sys) (*mongo.Database, *Gf_error) {

	connect_timeout_in_sec_int := 3
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(connect_timeout_in_sec_int) * time.Second)

	mongo_options := options.Client().ApplyURI(p_mongo_server_url_str)
	mongo_client, err := mongo.Connect(ctx, mongo_options)
	if err != nil {

		gf_err := Error__create("failed to connect to a MongoDB server at target url",
			"mongodb_connect_error",
			map[string]interface{}{
				"mongo_server_url_str": p_mongo_server_url_str,
			}, err, "gf_core", p_runtime_sys)
		return nil, gf_err
	}

	// test new connection
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(connect_timeout_in_sec_int) * time.Second)
	err = mongo_client.Ping(ctx, readpref.Primary())
	if err != nil {
		gf_err := Error__create("failed to ping a MongoDB server at target url",
			"mongodb_ping_error",
			map[string]interface{}{
				"mongo_server_url_str": p_mongo_server_url_str,
			}, err, "gf_core", p_runtime_sys)
		return nil, gf_err
	}

	mongo_db := mongo_client.Database(p_db_name_str)

	return mongo_db, nil
}

//-------------------------------------------------
func Mongo__connect(p_mongodb_host_str string,
	p_mongodb_db_name_str string,
	p_log_fun             func(string, string)) *mgo.Database {
	p_log_fun("FUN_ENTER", "gf_mongodb.Mongo__connect()")
	p_log_fun("INFO",      fmt.Sprintf("p_mongodb_host_str    - %s", p_mongodb_host_str))
	p_log_fun("INFO",      fmt.Sprintf("p_mongodb_db_name_str - %s", p_mongodb_db_name_str))
	
	session, err := mgo.DialWithTimeout(p_mongodb_host_str, time.Second * 90)
	if err != nil {
		panic(err)
	}

	//--------------------
	// IMPORTANT!! - writes are waited for to confirm them.
	// 	   			 in unsafe mode writes are fire-and-forget with no error checking. 
	//               this mode is faster, since no confirmation is expected.
	session.SetSafe(&mgo.Safe{})

	// Monotonic consistency - will read from a slave if possible, for better load distribution.
	//                         once the first write happens the connection is switched to the master.
	session.SetMode(mgo.Monotonic, true)

	//--------------------

	db := session.DB(p_mongodb_db_name_str)
	return db
}

//--------------------------------------------------------------------
func Mongo__handle_error(p_user_msg_str string,
	p_error_type_str     string,
	p_error_data_map     map[string]interface{},
	p_error              error,
	p_subsystem_name_str string,
	p_runtime_sys        *Runtime_sys) *Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_mongodb.Mongo__handle_error()")

	gf_error := Error__create_with_hook(p_user_msg_str,
		p_error_type_str,
		p_error_data_map,
		p_error,
		p_subsystem_name_str,
		func(p_gf_error *Gf_error) map[string]interface{} {

			gf_error_str := fmt.Sprint(p_gf_error)

			// IMPORTANT!! - "mgo" had behavior where after the connection was reset by mongod server,
			//               it (mgo) wouldnt reconnect to that server. so this hack is applied where the entire service
			//               is restarted so that a fresh DB connection is established
			
			if strings.Contains(gf_error_str, "connection reset by peer") {
				os.Exit(1)
			}

			// IMPORTANT!! - "mgo" specific error, where after a broken connection gets re-established the session object
			//               is still kept in an error state. to get out of this error state session.Refresh() can be called,
			//               but this might lead to inconsistencies if a signle session object is shared among multiple go-routines
			//               that might be in the middle of queries.
			//               conservative approach is taken here and error recovery is not attempted, instead the whole service
			//               is restarted. 
			if gf_error_str == "EOF" || gf_error_str == "Closed explicitly" {
				os.Exit(1)
			}

			return nil
		},
		p_runtime_sys)
	return gf_error
}

//--------------------------------------------------------------------
func Mongo__start(p_mongodb_bin_path_str string,
	p_mongodb_port_str          int,
	p_mongodb_data_dir_path_str string,
	p_mongodb_log_file_path_str string,
	p_sudo_bool                 bool,
	p_log_fun                   func(string,string)) error {
	p_log_fun("FUN_ENTER", "gf_mongodb.Mongo__start()")
	p_log_fun("INFO",      "p_mongodb_data_dir_path_str - "+p_mongodb_data_dir_path_str)
	p_log_fun("INFO",      "p_mongodb_log_file_path_str - "+p_mongodb_log_file_path_str)

	if _, err := os.Stat(p_mongodb_log_file_path_str); os.IsNotExist(err) {
		p_log_fun("ERROR", fmt.Sprintf("supplied log_file path is not a file - %s", p_mongodb_log_file_path_str))
		return err
	}

	p_log_fun("INFO", "-----------------------------------------")
	p_log_fun("INFO", "--------- STARTING - MONGODB ------------")
	p_log_fun("INFO", "-----------------------------------------")
	p_log_fun("INFO", "p_mongodb_bin_path_str      - "+p_mongodb_bin_path_str)
	p_log_fun("INFO", "p_mongodb_data_dir_path_str - "+p_mongodb_data_dir_path_str)
	p_log_fun("INFO", "p_mongodb_log_file_path_str - "+p_mongodb_log_file_path_str)

	args_lst := []string{
		"--fork",            //start the server as a deamon
		fmt.Sprintf("--dbpath %s",  p_mongodb_data_dir_path_str),
		fmt.Sprintf("--logpath %s", p_mongodb_log_file_path_str),

		"--port "+fmt.Sprint(p_mongodb_port_str),
		"--rest",            //turn on REST http API interface
		"--httpinterface",
		"--journal",         //turn journaling/durability on
		"--oplogSize 128",
	}
	
	var cmd *exec.Cmd
	if p_sudo_bool {
		new_args_lst := []string{p_mongodb_bin_path_str,}
		new_args_lst  = append(new_args_lst, args_lst...)

		cmd = exec.Command("sudo", new_args_lst...)
	} else {
		//cmd = exec.Command("/usr/bin/mongod") //fmt.Sprintf("'%s'",strings.Join(args_lst," ")),"&")
		cmd = exec.Command(p_mongodb_bin_path_str, args_lst...)
	}

	p_log_fun("INFO", "cmd - "+strings.Join(cmd.Args, " "))
	cmd.Start()

	return nil
}

//-------------------------------------------------
func Mongo__get_rs_members_info(p_mongodb_primary_host_str string,
	p_log_fun func(string, string)) ([]map[string]interface{}, error) {
	// p_log_fun("FUN_ENTER", "gf_mongodb.Mongo__get_rs_members_info()")
	// p_log_fun("INFO",      p_mongodb_primary_host_str)

	mongo_client__cmd_str := fmt.Sprintf("mongo --host %s --quiet --eval 'JSON.stringify(rs.status())'", p_mongodb_primary_host_str)

	out, err := exec.Command("sh", "-c",mongo_client__cmd_str).Output()

	//---------------
	// JSON
	var i map[string]interface{}
    err = json.Unmarshal([]byte(out), &i)
    if err != nil {
    	return nil, err
	}
	
    //---------------

	rs_members_lst := i["members"].([]map[string]interface{})
	var rs_members_info_lst []map[string]interface{}

	for _, m := range rs_members_lst {

		member_info_map := map[string]interface{} {
			"host_port_str": m["name"].(string),
			"state_str":     m["stateStr"].(string),
			"uptime_int":    m["uptime"].(int),
		}

		rs_members_info_lst = append(rs_members_info_lst, member_info_map)
	}

	return rs_members_info_lst, nil
}