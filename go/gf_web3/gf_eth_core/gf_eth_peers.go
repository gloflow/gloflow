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

package gf_eth_core

import (
	"time"
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

type GF_eth_peer__db_aggregate__name_group_counts struct {
	Name_str  string `bson:"_id"       json:"name_str"`
	Count_int int    `bson:"count_int" json:"count_int"`
}

//-------------------------------------------------
// metrics that are continuously calculated

func Eth_peers__init_continuous_metrics(p_metrics *GF_metrics,
	p_runtime *GF_runtime) {

	go func() {
		for {
			//---------------------
			// GET_PEERS_COUNTS
			peer_names_groups_counts_lst, gfErr := DBmongoPeersGetCount(p_metrics, p_runtime)
			if gfErr != nil {
				time.Sleep(60 * time.Second) // SLEEP
				continue
			}

			//---------------------
			unique_peer_names_num_int := len(peer_names_groups_counts_lst)
			p_metrics.Peers__unique_names_num__gauge.Set(float64(unique_peer_names_num_int))

			time.Sleep(60 * time.Second) // SLEEP
		}
	}()
}