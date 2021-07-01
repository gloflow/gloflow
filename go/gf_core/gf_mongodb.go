/*
MIT License

Copyright (c) 2019 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
	"crypto/tls"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// TRANSACTIONS
//-------------------------------------------------
// RUN - CamelCase name, for use by external projects.

func MongoTXrun(p_tx_fun func() *Gf_error,
	p_meta_map     map[string]interface{}, // data describing the DB write op
	p_mongo_client *mongo.Client,
	p_ctx          context.Context,
	p_runtime_sys  *Runtime_sys) (mongo.Session, *Gf_error) {
	
	return Mongo__tx_run(p_tx_fun,
		p_meta_map,
		p_mongo_client,
		p_ctx,
		p_runtime_sys)
}

//-------------------------------------------------
// RUN
func Mongo__tx_run(p_tx_fun func() *Gf_error,
	p_meta_map     map[string]interface{}, // data describing the DB write op
	p_mongo_client *mongo.Client,
	p_ctx          context.Context,
	p_runtime_sys  *Runtime_sys) (mongo.Session, *Gf_error) {
	
	// TX_INIT
	txSession, txOptions, gf_err := Mongo__tx_init(p_mongo_client,
		p_meta_map,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	// TX_RUN
	err := mongo.WithSession(p_ctx, txSession, func(pDBsessionCtx mongo.SessionContext) error {

		// TX_START
		if err := txSession.StartTransaction(txOptions); err != nil {
            return err
        }

		// DB_TX_FUN
		gf_err := p_tx_fun()
		if gf_err != nil {
			return gf_err.Error
		}

		// TX_COMMIT
		if err := txSession.CommitTransaction(pDBsessionCtx); err != nil {
            return err
        }

		return nil

	})

	if err != nil {

		if err_abort := txSession.AbortTransaction(p_ctx); err_abort != nil {
            gf_err := Mongo__handle_error("failed to execute Mongodb session",
				"mongodb_session_abort_error",
				p_meta_map,
				err, "gf_core", p_runtime_sys)
			return nil, gf_err
        }

		gf_err := Mongo__handle_error("failed to execute Mongodb session",
			"mongodb_session_error",
			p_meta_map,
			err, "gf_core", p_runtime_sys)
		return nil, gf_err		
	}

	return txSession, nil
}

//-------------------------------------------------
// INIT
func Mongo__tx_init(p_mongo_client *mongo.Client,
	p_meta_map    map[string]interface{}, // data describing the DB write op
	p_runtime_sys *Runtime_sys) (mongo.Session, *options.TransactionOptions, *Gf_error) {



	wc := writeconcern.New(writeconcern.WMajority())
    rc := readconcern.Snapshot()

    tx_options := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

    session, err := p_mongo_client.StartSession()
    if err != nil {
		gf_err := Mongo__handle_error("failed to start a Mongo session",
			"mongodb_start_session_error",
			p_meta_map,
			err, "gf_core", p_runtime_sys)
		return nil, nil, gf_err		
    }

	return session, tx_options, nil
}

//-------------------------------------------------
// OPS
//-------------------------------------------------
// MONGO_COUNT
func MongoCount(p_query bson.M,
	p_meta_map            map[string]interface{},
	p_coll                *mongo.Collection,
	p_ctx                 context.Context,
	p_runtime_sys         *Runtime_sys) (int64, *Gf_error) {

	return Mongo__count(p_query, p_meta_map, p_coll, p_ctx, p_runtime_sys)
}

func Mongo__count(p_query bson.M,
	p_meta_map            map[string]interface{}, // data describing the DB write op
	p_coll                *mongo.Collection,
	p_ctx                 context.Context,
	p_runtime_sys         *Runtime_sys) (int64, *Gf_error) {

	// FIX!! - externalize this max_time value to some config.
	opts := options.Count().SetMaxTime(5 * time.Second)

	count_int, err := p_coll.CountDocuments(p_ctx, p_query, opts)
	if err != nil {
		gf_err := Mongo__handle_error("failed to count number of particular docs in DB",
			"mongodb_count_error",
			p_meta_map,
			err, "gf_core", p_runtime_sys)
		return 0, gf_err
	}
	return count_int, nil
}

//-------------------------------------------------
// MONGO_FIND_LATEST
func MongoFindLatest(p_query bson.M,
	p_time_field_name_str string,
	p_meta_map            map[string]interface{}, // data describing the DB write op
	p_coll                *mongo.Collection,
	p_ctx                 context.Context,
	p_runtime_sys         *Runtime_sys) (map[string]interface{}, *Gf_error) {

	return Mongo__find_latest(p_query,
		p_time_field_name_str,
		p_meta_map,
		p_coll,
		p_ctx,
		p_runtime_sys)
}

//-------------------------------------------------
// MONGO_FIND_LATEST
func Mongo__find_latest(p_query bson.M,
	p_time_field_name_str string,
	p_meta_map            map[string]interface{}, // data describing the DB write op
	p_coll                *mongo.Collection,
	p_ctx                 context.Context,
	p_runtime_sys         *Runtime_sys) (map[string]interface{}, *Gf_error) {
	

	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{p_time_field_name_str: -1})
	find_opts.SetLimit(1)
	
	cursor, gf_err := Mongo__find(p_query,
		find_opts,
		p_meta_map,
		p_coll,
		p_ctx,
		p_runtime_sys)

	if gf_err != nil {
		return nil, gf_err
	}

	// no result
	if cursor == nil {
		return nil, nil
	}

	var records_lst []map[string]interface{}
	err := cursor.All(p_ctx, &records_lst)
	if err != nil {
		gf_err := Mongo__handle_error("failed to load DB results from DB cursor",
			"mongodb_cursor_all",
			p_meta_map,
			err, "gf_core", p_runtime_sys)
		return nil, gf_err
	}
	
	// get latest key
	record := records_lst[0]

	return record, nil
}

//-------------------------------------------------
// FIND

func MongoFind(p_query bson.M,
	p_opts        *options.FindOptions,
	p_meta_map    map[string]interface{}, // data describing the DB write op
	p_coll        *mongo.Collection,
	p_ctx         context.Context,
	p_runtime_sys *Runtime_sys) (*mongo.Cursor, *Gf_error) {

	return Mongo__find(p_query,
		p_opts,
		p_meta_map,
		p_coll,
		p_ctx,
		p_runtime_sys)
}

//-------------------------------------------------
func Mongo__find(p_query bson.M,
	p_opts        *options.FindOptions,
	p_meta_map    map[string]interface{}, // data describing the DB write op
	p_coll        *mongo.Collection,
	p_ctx         context.Context,
	p_runtime_sys *Runtime_sys) (*mongo.Cursor, *Gf_error) {

	cur, err := p_coll.Find(p_ctx, p_query, p_opts)
	if err != nil {
		
		// NO_DOCUMENTS
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		gf_err := Mongo__handle_error("failed to find records in DB",
			"mongodb_find_error",
			p_meta_map,
			err, "gf_core", p_runtime_sys)
		return nil, gf_err
	}


	// defer cur.Close(p_ctx)
	return cur, nil
}

//-------------------------------------------------
// UPSERT
func Mongo__upsert(p_query bson.M,
	p_record      interface{},
	p_meta_map    map[string]interface{}, // data describing the DB write op
	p_coll        *mongo.Collection,
	p_ctx         context.Context,
	p_runtime_sys *Runtime_sys) *Gf_error {


	_, err := p_coll.UpdateOne(p_ctx, p_query, bson.M{"$set": p_record,},
		options.Update().SetUpsert(true))
	if err != nil {

		// NO_DOCUMENTS
		if err == mongo.ErrNoDocuments {
			return nil
		}

		gf_err := Mongo__handle_error("failed to update/upsert document in the DB",
			"mongodb_update_error",
			p_meta_map,
			err, "gf_core", p_runtime_sys)
		return gf_err
	}

	return nil
}

//-------------------------------------------------
// INSERT_BULK
func Mongo__insert_bulk(p_ids_lst []string,
	p_record_lst    []interface{},
	p_coll_name_str string,
	p_meta_map      map[string]interface{}, // data describing the DB write op 
	p_ctx           context.Context,
	p_runtime_sys   *Runtime_sys) *Gf_error {

	models := []mongo.WriteModel{}
	for i, id_str := range p_ids_lst {

		replacement_doc := p_record_lst[i]

		// FIX!! - "$set" - replaces existing doc with _id with this new one. 
		//                  but if ID is some sort of hash of the document (as is in a few GF apps)
		//                  than the contents of those documents are the same as well (their _id/hashes are the same),
		//                  so the DB update with a replacement doc is redundant. 
		//                  fix this special (but frequent in GF) case.
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.D{{"_id", id_str}}).
			SetUpdate(bson.M{"$set": replacement_doc,}).
			SetUpsert(true) // upsert=true - insert new document if the _id doesnt exist
		models = append(models, model)
	}

	opts := options.BulkWrite().SetOrdered(false)

	r, err := p_runtime_sys.Mongo_db.Collection(p_coll_name_str).BulkWrite(p_ctx, models, opts)
	if err != nil {
		gf_err := Mongo__handle_error("failed to bulk write new documents into the DB",
			"mongodb_write_bulk_error",
			p_meta_map,
			err, "gf_core", p_runtime_sys)
		return gf_err
	}

	if r.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
	}

	return nil
}

//-------------------------------------------------
// INSERT
func Mongo__insert(p_record interface{},
	p_coll_name_str string,
	p_meta_map      map[string]interface{}, // data describing the DB write op 
	p_ctx           context.Context,
	p_runtime_sys   *Runtime_sys) *Gf_error {

	_, err := p_runtime_sys.Mongo_db.Collection(p_coll_name_str).InsertOne(p_ctx, p_record)
	if err != nil {
		p_meta_map["coll_name_str"] = p_coll_name_str
		gf_err := Mongo__handle_error("failed to insert a new record into the DB",
			"mongodb_insert_error",
			p_meta_map,
			err, "gf_core", p_runtime_sys)
		
		return gf_err
	}
	return nil
}

//-------------------------------------------------
// ENSURE_INDEX
func Mongo__ensure_index(p_indexes_keys_lst [][]string, 
	p_coll_name_str string,
	p_runtime_sys   *Runtime_sys) ([]string, *Gf_error) {

	models := []mongo.IndexModel{}
	for _, index_keys_lst := range p_indexes_keys_lst {
		
		keys_bson := bson.D{}
		for _, k := range index_keys_lst {
			keys_bson = append(keys_bson, bson.E{k, 1})
		}

		model := mongo.IndexModel{
		
			Keys:    keys_bson, // bson.D{{"name", 1}, {"email", 1}},
			Options: options.Index().
				SetUnique(false).    // index must necessarily contain only a single document per Key
				SetBackground(true). // other connections will be allowed to proceed using the collection without the index while it's being built
				SetSparse(true),     // only documents containing the provided Key fields will be included in the index
		}

		models = append(models, model)
	}

	// CREATE_INDEX
	ctx := context.Background()
	opts := options.CreateIndexes().SetMaxTime(600 * time.Second)

	indexes_names_lst, err := p_runtime_sys.Mongo_db.Collection(p_coll_name_str).Indexes().CreateMany(ctx, models, opts)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "duplicate key error index") {
			return []string{}, nil
		} else {
			gf_err := Mongo__handle_error(fmt.Sprintf("failed to create db indexes on fields"), 
				"mongodb_ensure_index_error",
				map[string]interface{}{"indexes_keys_lst": p_indexes_keys_lst,},
				err, "gf_core", p_runtime_sys)
			return nil, gf_err
		}
	}

	return indexes_names_lst, nil

	/*gf_errs_lst := []*Gf_error{}
	for _, index_keys_lst := range p_indexes_keys_lst {
		
		doc_type__index := mgo.Index{
			Key:        index_keys_lst, 
			Unique:     false, // index must necessarily contain only a single document per Key
			DropDups:   false, // documents with the same key as a previously indexed one will be dropped rather than an error returned.
			Background: true,  // other connections will be allowed to proceed using the collection without the index while it's being built
			Sparse:     true,  // only documents containing the provided Key fields will be included in the index
		}
		
		if p_runtime_sys.Mongo_db.Collection(p_coll_name_str) != nil {
			err := p_runtime_sys.Mongo_db.Collection(p_coll_name_str).EnsureIndex(doc_type__index)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "duplicate key error index") {
					continue // ignore, index already exists
				} else {
					gf_err := Mongo__handle_error(fmt.Sprintf("failed to create db index on fields - %s", index_keys_lst), 
						"mongodb_ensure_index_error", nil, err, "gf_core", p_runtime_sys)
					gf_errs_lst = append(gf_errs_lst, gf_err)
				}
			}
		} else {
			fmt.Printf("mongodb collection %s doesnt exist\n", p_coll_name_str)
		}
	}*/
}

//--------------------------------------------------------------------
// UTILS
//--------------------------------------------------------------------
// CONNECT_NEW
func Mongo__connect_new(p_mongo_server_url_str string,
	p_db_name_str       string,
	p_tls_custom_config *tls.Config,
	p_runtime_sys       *Runtime_sys) (*mongo.Database, *mongo.Client, *Gf_error) {

	connect_timeout_in_sec_int := 3
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(connect_timeout_in_sec_int) * time.Second)

	mongo_options := options.Client().ApplyURI(p_mongo_server_url_str)

	// TLS
	if p_tls_custom_config != nil {
		mongo_options.SetTLSConfig(p_tls_custom_config)
	}

	mongo_client, err := mongo.Connect(ctx, mongo_options)
	if err != nil {

		gf_err := Error__create("failed to connect to a MongoDB server at target url",
			"mongodb_connect_error",
			map[string]interface{}{
				// "mongo_server_url_str": p_mongo_server_url_str,
			}, err, "gf_core", p_runtime_sys)
		return nil, nil, gf_err
	}

	// test new connection
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(connect_timeout_in_sec_int) * time.Second)
	err = mongo_client.Ping(ctx, readpref.Primary())
	if err != nil {
		gf_err := Error__create("failed to ping a MongoDB server at target url",
			"mongodb_ping_error",
			map[string]interface{}{
				// "mongo_server_url_str": p_mongo_server_url_str,
			}, err, "gf_core", p_runtime_sys)
		return nil, nil, gf_err
	}

	mongo_db := mongo_client.Database(p_db_name_str)

	return mongo_db, mongo_client, nil
}

//--------------------------------------------------------------------
/*func Mongo__connect(p_mongodb_host_str string,
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
}*/

//--------------------------------------------------------------------
// HANDLE_ERROR
func Mongo__handle_error(p_user_msg_str string,
	p_error_type_str     string,
	p_error_data_map     map[string]interface{},
	p_error              error,
	p_subsystem_name_str string,
	p_runtime_sys        *Runtime_sys) *Gf_error {
	// p_runtime_sys.Log_fun("FUN_ENTER", "gf_mongodb.Mongo__handle_error()")

	gf_err := Error__create_with_hook(p_user_msg_str,
		p_error_type_str,
		p_error_data_map,
		p_error,
		p_subsystem_name_str,
		func(p_gf_err *Gf_error) map[string]interface{} {

			gf_error_str := fmt.Sprint(p_gf_err)

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
	return gf_err
}

//--------------------------------------------------------------------
// START
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

//--------------------------------------------------------------------
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