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
	"os/exec"
	"context"
	"crypto/tls"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	// "encoding/json"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
// TRANSACTIONS
//-------------------------------------------------
// RUN
func MongoTXrun(p_tx_fun func() *GFerror,
	pMetaMap       map[string]interface{}, // data describing the DB write op
	p_mongo_client *mongo.Client,
	pCtx           context.Context,
	pRuntimeSys    *RuntimeSys) (mongo.Session, *GFerror) {
	
	// TX_INIT
	txSession, txOptions, gf_err := MongoTXinit(p_mongo_client,
		pMetaMap,
		pRuntimeSys)
	if gf_err != nil {
		return nil, gf_err
	}

	// TX_RUN
	err := mongo.WithSession(pCtx, txSession, func(pDBsessionCtx mongo.SessionContext) error {

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

		if err_abort := txSession.AbortTransaction(pCtx); err_abort != nil {
            gf_err := MongoHandleError("failed to execute Mongodb session",
				"mongodb_session_abort_error",
				pMetaMap,
				err, "gf_core", pRuntimeSys)
			return nil, gf_err
        }

		gf_err := MongoHandleError("failed to execute Mongodb session",
			"mongodb_session_error",
			pMetaMap,
			err, "gf_core", pRuntimeSys)
		return nil, gf_err		
	}

	return txSession, nil
}

//-------------------------------------------------
// INIT
func MongoTXinit(p_mongo_client *mongo.Client,
	pMetaMap    map[string]interface{}, // data describing the DB write op
	pRuntimeSys *RuntimeSys) (mongo.Session, *options.TransactionOptions, *GFerror) {



	wc := writeconcern.New(writeconcern.WMajority())
    rc := readconcern.Snapshot()

    tx_options := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

    session, err := p_mongo_client.StartSession()
    if err != nil {
		gf_err := MongoHandleError("failed to start a Mongo session",
			"mongodb_start_session_error",
			pMetaMap,
			err, "gf_core", pRuntimeSys)
		return nil, nil, gf_err		
    }

	return session, tx_options, nil
}

//-------------------------------------------------
// OPS
//-------------------------------------------------
// MONGO_COUNT
func MongoCount(pQuery bson.M,
	pMetaMap    map[string]interface{}, // data describing the DB write op
	p_coll      *mongo.Collection,
	pCtx        context.Context,
	pRuntimeSys *RuntimeSys) (int64, *GFerror) {

	// FIX!! - externalize this max_time value to some config.
	opts := options.Count().SetMaxTime(5 * time.Second)

	count_int, err := p_coll.CountDocuments(pCtx, pQuery, opts)
	if err != nil {
		gf_err := MongoHandleError("failed to count number of particular docs in DB",
			"mongodb_count_error",
			pMetaMap,
			err, "gf_core", pRuntimeSys)
		return 0, gf_err
	}
	return count_int, nil
}

//-------------------------------------------------
// MONGO_FIND_LATEST
func MongoFindLatest(pQuery bson.M,
	p_time_field_name_str string,
	pMetaMap              map[string]interface{}, // data describing the DB write op
	p_coll                *mongo.Collection,
	pCtx                  context.Context,
	pRuntimeSys           *RuntimeSys) (map[string]interface{}, *GFerror) {
	

	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{p_time_field_name_str: -1})
	find_opts.SetLimit(1)
	
	cursor, gf_err := MongoFind(pQuery,
		find_opts,
		pMetaMap,
		p_coll,
		pCtx,
		pRuntimeSys)

	if gf_err != nil {
		return nil, gf_err
	}

	// no result
	if cursor == nil {
		return nil, nil
	}

	var records_lst []map[string]interface{}
	err := cursor.All(pCtx, &records_lst)
	if err != nil {
		gf_err := MongoHandleError("failed to load DB results from DB cursor",
			"mongodb_cursor_all",
			pMetaMap,
			err, "gf_core", pRuntimeSys)
		return nil, gf_err
	}
	
	// get latest key
	record := records_lst[0]

	return record, nil
}

//-------------------------------------------------
// FIND
func MongoFind(pQuery bson.M,
	p_opts      *options.FindOptions,
	pMetaMap    map[string]interface{}, // data describing the DB write op
	p_coll      *mongo.Collection,
	pCtx        context.Context,
	pRuntimeSys *RuntimeSys) (*mongo.Cursor, *GFerror) {

	cur, err := p_coll.Find(pCtx, pQuery, p_opts)
	if err != nil {
		
		// NO_DOCUMENTS
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		gf_err := MongoHandleError("failed to find records in DB",
			"mongodb_find_error",
			pMetaMap,
			err, "gf_core", pRuntimeSys)
		return nil, gf_err
	}

	// defer cur.Close(pCtx)
	return cur, nil
}

//-------------------------------------------------
func MongoDelete(pQuery bson.M,
	pCollNameStr string,
	pMetaMap     map[string]interface{}, // data describing the DB write op
	pCtx         context.Context,
	pRuntimeSys  *RuntimeSys) *GFerror {

	_, err := pRuntimeSys.Mongo_db.Collection(pCollNameStr).DeleteMany(pCtx, pQuery)
	if err != nil {
		gf_err := MongoHandleError("failed to delete documents in the DB",
			"mongodb_delete_error",
			pMetaMap,
			err, "gf_core", pRuntimeSys)
		return gf_err
	}
	return nil
}

//-------------------------------------------------
// UPSERT
func MongoUpsert(pQuery bson.M,
	pRecord     interface{},
	pMetaMap    map[string]interface{}, // data describing the DB write op
	pColl       *mongo.Collection,
	pCtx        context.Context,
	pRuntimeSys *RuntimeSys) *GFerror {

	_, err := pColl.UpdateOne(pCtx, pQuery, bson.M{"$set": pRecord,},
		options.Update().SetUpsert(true))
	if err != nil {

		// NO_DOCUMENTS
		if err == mongo.ErrNoDocuments {
			return nil
		}

		gfErr := MongoHandleError("failed to update/upsert document in the DB",
			"mongodb_update_error",
			pMetaMap,
			err, "gf_core", pRuntimeSys)
		return gfErr
	}

	return nil
}

//-------------------------------------------------
// UPSERT_BULK
func MongoUpsertBulk(pFilterDocsByFieldsLst []map[string]string,
	pRecordsLst  []interface{},
	pCollNameStr string,
	pMetaMap     map[string]interface{}, // data describing the DB write op 
	pCtx         context.Context,
	pRuntimeSys  *RuntimeSys) (int64, *GFerror) {

	models := []mongo.WriteModel{}
	for i, filterDocByFieldsMap := range pFilterDocsByFieldsLst {

		replacementDoc := pRecordsLst[i]

		// bulk filter lists to select objects to run updates(upserts) on
		filter := bson.D{}
		for k, v := range filterDocByFieldsMap {
			filter = append(filter, bson.E{k, v}) 
		}

		// FIX!! - "$set" - replaces existing doc with _id with this new one. 
		//                  but if ID is some sort of hash of the document (as is in a few GF apps)
		//                  than the contents of those documents are the same as well (their _id/hashes are the same),
		//                  so the DB update with a replacement doc is redundant. 
		//                  fix this special (but frequent in GF) case.
		model := mongo.NewUpdateOneModel().
			SetFilter(filter). // bson.D{{"id_str", IDstr}}).
			SetUpdate(bson.M{"$set": replacementDoc,}).
			SetUpsert(true) // upsert=true - insert new document if the _id doesnt exist
		models = append(models, model)
	}

	opts := options.BulkWrite().SetOrdered(false)

	r, err := pRuntimeSys.Mongo_db.Collection(pCollNameStr).BulkWrite(pCtx, models, opts)
	if err != nil {
		gf_err := MongoHandleError("failed to bulk write new documents into the DB",
			"mongodb_write_bulk_error",
			pMetaMap,
			err, "gf_core", pRuntimeSys)
		return 0, gf_err
	}

	if r.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
	}

	insertedNewDocsInt := r.InsertedCount
	return insertedNewDocsInt, nil
}

//-------------------------------------------------
// INSERT
func MongoInsert(p_record interface{},
	pCollNameStr string,
	pMetaMap     map[string]interface{}, // data describing the DB write op 
	pCtx         context.Context,
	pRuntimeSys  *RuntimeSys) *GFerror {

	_, err := pRuntimeSys.Mongo_db.Collection(pCollNameStr).InsertOne(pCtx, p_record)
	if err != nil {
		pMetaMap["coll_name_str"] = pCollNameStr
		gfErr := MongoHandleError("failed to insert a new record into the DB",
			"mongodb_insert_error",
			pMetaMap,
			err, "gf_core", pRuntimeSys)
		
		return gfErr
	}
	return nil
}

//-------------------------------------------------
// ENSURE_INDEX
func MongoEnsureIndex(pIndexesKeysLst [][]string, 
	pIndexesNamesLst []string,
	pCollNameStr     string,
	pRuntimeSys      *RuntimeSys) *GFerror {

	models := []mongo.IndexModel{}
	for i, indexKeysLst := range pIndexesKeysLst {
		
		keysBson := bson.D{}
		for _, k := range indexKeysLst {
			keysBson = append(keysBson, bson.E{k, 1})
		}

		indexOptions := options.Index().
			SetUnique(false).    // index must necessarily contain only a single document per Key
			SetBackground(true). // other connections will be allowed to proceed using the collection without the index while it's being built
			SetSparse(true)      // only documents containing the provided Key fields will be included in the index

		indexNameStr := pIndexesNamesLst[i]
		if indexNameStr != "" {
			indexOptions.SetName(indexNameStr)
		}
		
		model := mongo.IndexModel{
		
			Keys:    keysBson, // bson.D{{"name", 1}, {"email", 1}},
			Options: indexOptions,
		}

		models = append(models, model)
	}

	// CREATE_INDEX
	ctx := context.Background()
	opts := options.CreateIndexes().SetMaxTime(600 * time.Second)

	_, err := pRuntimeSys.Mongo_db.Collection(pCollNameStr).Indexes().CreateMany(ctx, models, opts)
	if err != nil {

		errStr := fmt.Sprint(err)

		// index exists already
		if strings.Contains(errStr, "duplicate key error index") || 
			strings.Contains(errStr, "existing index has the same name as the requested index") {
			return nil

		} else {
			gfErr := MongoHandleError(fmt.Sprintf("failed to create db indexes on fields"), 
				"mongodb_ensure_index_error",
				map[string]interface{}{
					"indexes_keys_lst":  pIndexesKeysLst,
					"indexes_names_lst": pIndexesNamesLst,
				},
				err, "gf_core", pRuntimeSys)
			return gfErr
		}
	}

	return nil
}

//--------------------------------------------------------------------
// COLLECTIONS
//--------------------------------------------------------------------
func MongoCollExists(pCollNameStr string,
	pCtx        context.Context,
	pRuntimeSys *RuntimeSys) (bool, *GFerror) {

	coll_names_lst, err := pRuntimeSys.Mongo_db.ListCollectionNames(pCtx, bson.D{})
	if err != nil {
		gfErr := ErrorCreate("failed to get a list of all collection names to check if the given collection exists",
			"mongodb_get_collection_names_error",
			map[string]interface{}{
				"coll_name_str": pCollNameStr,
			}, err, "gf_core", pRuntimeSys)
		return false, gfErr
	}

	for _, name_str := range coll_names_lst {
		if name_str == pCollNameStr {
			return true, nil
		}
	}
	return false, nil
}

//--------------------------------------------------------------------
// UTILS
//--------------------------------------------------------------------
// CONNECT_NEW
func MongoConnectNew(p_mongo_server_url_str string,
	p_db_name_str       string,
	p_tls_custom_config *tls.Config,
	pRuntimeSys         *RuntimeSys) (*mongo.Database, *mongo.Client, *GFerror) {

	connect_timeout_in_sec_int := 3
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(connect_timeout_in_sec_int) * time.Second)

	mongo_options := options.Client().ApplyURI(p_mongo_server_url_str)

	// TLS
	if p_tls_custom_config != nil {
		mongo_options.SetTLSConfig(p_tls_custom_config)
	}

	mongo_client, err := mongo.Connect(ctx, mongo_options)
	if err != nil {

		gfErr := ErrorCreate("failed to connect to a MongoDB server at target url",
			"mongodb_connect_error",
			map[string]interface{}{
				// "mongo_server_url_str": p_mongo_server_url_str,
			}, err, "gf_core", pRuntimeSys)
		return nil, nil, gfErr
	}

	// test new connection
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(connect_timeout_in_sec_int) * time.Second)
	err = mongo_client.Ping(ctx, readpref.Primary())
	if err != nil {
		gfErr := ErrorCreate("failed to ping a MongoDB server at target url",
			"mongodb_ping_error",
			map[string]interface{}{
				// "mongo_server_url_str": p_mongo_server_url_str,
			}, err, "gf_core", pRuntimeSys)
		return nil, nil, gfErr
	}

	mongo_db := mongo_client.Database(p_db_name_str)

	return mongo_db, mongo_client, nil
}

//--------------------------------------------------------------------
// HANDLE_ERROR
func MongoHandleError(p_user_msg_str string,
	p_error_type_str     string,
	p_error_data_map     map[string]interface{},
	p_error              error,
	p_subsystem_name_str string,
	pRuntimeSys          *RuntimeSys) *GFerror {

	// stack/callgraph distance of this call to ErrorCreateWithHook from
	// error stack fetching and storage, so provide this level explicitly here
	// so that gf_error can pick the proper place to report as the function that
	// originated this mongo error (caller of MongoHandleError)
	skipStackFramesNumInt := 4

	gfErr := ErrorCreateWithHook(p_user_msg_str,
		p_error_type_str,
		p_error_data_map,
		p_error,
		p_subsystem_name_str,
		skipStackFramesNumInt,
		func(p_gf_err *GFerror) map[string]interface{} {

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
		pRuntimeSys)
	return gfErr
}

//--------------------------------------------------------------------
// START
func MongoStart(p_mongodb_bin_path_str string,
	p_mongodb_port_str          int,
	p_mongodb_data_dir_path_str string,
	p_mongodb_log_file_path_str string,
	p_sudo_bool                 bool,
	pLogFun                     func(string, string)) error {
	pLogFun("INFO", "p_mongodb_data_dir_path_str - "+p_mongodb_data_dir_path_str)
	pLogFun("INFO", "p_mongodb_log_file_path_str - "+p_mongodb_log_file_path_str)

	if _, err := os.Stat(p_mongodb_log_file_path_str); os.IsNotExist(err) {
		pLogFun("ERROR", fmt.Sprintf("supplied log_file path is not a file - %s", p_mongodb_log_file_path_str))
		return err
	}

	pLogFun("INFO", "-----------------------------------------")
	pLogFun("INFO", "--------- STARTING - MONGODB ------------")
	pLogFun("INFO", "-----------------------------------------")
	pLogFun("INFO", "p_mongodb_bin_path_str      - "+p_mongodb_bin_path_str)
	pLogFun("INFO", "p_mongodb_data_dir_path_str - "+p_mongodb_data_dir_path_str)
	pLogFun("INFO", "p_mongodb_log_file_path_str - "+p_mongodb_log_file_path_str)

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

	pLogFun("INFO", "cmd - "+strings.Join(cmd.Args, " "))
	cmd.Start()

	return nil
}