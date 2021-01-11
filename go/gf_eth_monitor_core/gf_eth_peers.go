/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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

package gf_eth_monitor_core

import (
	"log"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
// GF_ETH_PEER__NEW_LIFECYCLE
type GF_eth_peer__new_lifecycle struct {
	T_str              string  `bson:"t"` // "peer_new_lifecycle"
	V_str              string  `bson:"v"` // version - "v0",...
	Peer_name_str      string  `bson:"peer_name_str"` 
	Peer_enode_id_str  string  `bson:"peer_enode_id_str"`
	Peer_remote_ip_str string  `bson:"peer_remote_ip_str"` 
	Node_public_ip_str string  `bson:"node_public_ip_str"`
	Event_time_unix_f  float64 `bson:"event_time_unix_f"`
}

type GF_eth_peer__db_aggregate__name_group struct {
	Name_str             string   `bson:"_id"                  json:"name_str"`
	Peers_remote_ips_lst []string `bson:"peers_remote_ips_lst" json:"peers_remote_ips_lst"`
	Count_int            int      `bson:"count_int"            json:"count_int"`
}

//-------------------------------------------------
// metrics that are continuously calculated

func Eth_peers__init_continuous_metrics(p_metrics *GF_metrics,
	p_runtime *GF_runtime) {

	go func() {


		for {
			
			peer_names_groups_lst := Eth_peers__get_pipeline(p_metrics, p_runtime)

			unique_peer_names_num_int := len(peer_names_groups_lst)

			p_metrics.Gauge__peers_unique_names_num.Set(float64(unique_peer_names_num_int))

			// SLEEP
			time.Sleep(60 * time.Second)
		}
	}()
}

//-------------------------------------------------
// GET_PIPELINE
func Eth_peers__get_pipeline(p_metrics *GF_metrics,
	p_runtime *GF_runtime) []*GF_eth_peer__db_aggregate__name_group {



	coll_name_str := "gf_eth_peers"
	coll := p_runtime.Mongodb_db.Collection(coll_name_str)

	ctx := context.Background()




	// Start Aggregation Example 1
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{{"t", "peer_new_lifecycle"}}},
		},
		{
			{"$group", bson.D{
				{"_id",                  "$peer_name_str"},
				{"peers_remote_ips_lst", bson.D{{"$push", "$peer_remote_ip_str"}}},
				{"count_int",            bson.D{{"$sum", 1}}},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	
		// METRICS
		if p_metrics != nil {
			p_metrics.Counter__errs_num.Inc()
		}
	}
	defer cursor.Close(ctx)



	peer_names_groups_lst := []*GF_eth_peer__db_aggregate__name_group{}
	for cursor.Next(ctx) {

		var peer_name_group GF_eth_peer__db_aggregate__name_group
		err := cursor.Decode(&peer_name_group)
		if err != nil {
			log.Fatal(err)
		}
	
		peer_names_groups_lst = append(peer_names_groups_lst, &peer_name_group)
	}

	return peer_names_groups_lst

	/*q := bson.M{"t": "peer_new_lifecycle", }

	cur, err := p_runtime.Mongodb_db.Collection(coll_name_str).Find(ctx, q)
	if err != nil {
		log.Fatal(err)
	
		// METRICS
		if p_metrics != nil {
			p_metrics.counter__errs_num.Inc()
		}
	}
	defer cur.Close(ctx)

	peer_names_lst := []string{}
	for cur.Next(ctx) {

		var peer_lifecycle GF_eth_peer__new_lifecycle
		err := cur.Decode(&peer_lifecycle)
		if err != nil {
			log.Fatal(err)
		}
	
		peer_names_lst = append(peer_names_lst, peer_lifecycle.Peer_name_str)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)

		// METRICS
		if p_metrics != nil {
			p_metrics.counter__errs_num.Inc()
		}
	}
	
	return peer_names_lst*/
}

//-------------------------------------------------
// DB_WRITE
func Eth_peers__db_write(p_peer_new_lifecycle *GF_eth_peer__new_lifecycle,
	p_metrics *GF_metrics,
	p_runtime *GF_runtime) *gf_core.Gf_error {



	coll_name_str := "gf_eth_peers"
	_, err := p_runtime.Mongodb_db.Collection(coll_name_str).InsertOne(context.Background(), p_peer_new_lifecycle)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to insert a new Peer lifecycle into the DB",
			"mongodb_insert_error",
			map[string]interface{}{"peer_name_str": p_peer_new_lifecycle.Peer_name_str,},
			err, "gf_eth_monitor_lib", p_runtime.Runtime_sys)

		// METRICS
		if p_metrics != nil {
			p_metrics.Counter__errs_num.Inc()
		}
		
		return gf_err
	}


	// METRICS
	if p_metrics != nil {
		p_metrics.Counter__db_writes_num__new_peer_lifecycle.Inc()
	}

	return nil
}