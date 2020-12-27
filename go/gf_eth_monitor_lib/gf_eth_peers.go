package gf_eth_monitor_lib

import (
	"log"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
// GF_ETH_PEER__NEW_LIFECYCLE
type GF_eth_peer__new_lifecycle struct {
	T_str              string  `bson:"t"` // "peer_new_lifecycle"
	Peer_name_str      string  `bson:"peer_name_str"` 
	Peer_enode_id_str  string  `bson:"peer_enode_id_str"`
	Peer_remote_ip_str string  `bson:"peer_remote_ip_str"` 
	Node_public_ip_str string  `bson:"node_public_ip_str"`
	Event_time_unix_f  float64 `bson:"event_time_unix_f"`
}

//-------------------------------------------------
// GET_PIPELINE
func eth_peers__get_pipeline(p_runtime *GF_runtime) []string {



	coll_name_str := "gf_eth_peers"


	ctx := context.Background()
	q   := bson.M{"t": "peer_new_lifecycle", }

	cur, err := p_runtime.Mongodb_db.Collection(coll_name_str).Find(ctx, q)
	if err != nil { log.Fatal(err) }
	defer cur.Close(ctx)



	peer_names_lst := []string{}
	for cur.Next(ctx) {

		var peer_lifecycle GF_eth_peer__new_lifecycle // bson.M
		err := cur.Decode(&peer_lifecycle)
		if err != nil { log.Fatal(err) }
		



		peer_names_lst = append(peer_names_lst, peer_lifecycle.Peer_name_str)


	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}



	return peer_names_lst
}

//-------------------------------------------------
// DB_WRITE
func eth_peers__db_write(p_peer_new_lifecycle *GF_eth_peer__new_lifecycle,
	p_metrics *GF_metrics,
	p_runtime *GF_runtime) *gf_core.Gf_error {



	coll_name_str := "gf_eth_peers"
	_, err := p_runtime.Mongodb_db.Collection(coll_name_str).InsertOne(context.Background(), p_peer_new_lifecycle)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to insert a new Peer lifecycle into the DB",
			"mongodb_insert_error",
			map[string]interface{}{"peer_name_str": p_peer_new_lifecycle.Peer_name_str,},
			err, "gf_eth_monitor_lib", p_runtime.Runtime_sys)
		return gf_err
	}


	// METRICS
	if p_metrics != nil {
		p_metrics.counter__db_writes_num__new_peer_lifecycle.Inc()
	}

	return nil



}