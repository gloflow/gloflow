/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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
	"fmt"
	"time"
	"github.com/getsentry/sentry-go"
	// "github.com/influxdata/influxdb-client-go/v2"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------------------

type GF_runtime struct {
	Config     *GF_config
	Py_plugins *GF_py_plugins
	RuntimeSys *gf_core.RuntimeSys
	// Influxdb_client *influxdb2.Client
}

type GF_config struct {

	// PORTS
	Port_str         string `mapstructure:"port"`
	Port_metrics_str string `mapstructure:"port_metrics"`

	// MONGODB - this is the dedicated mongodb DB
	Mongodb_host_str    string `mapstructure:"mongodb_host"`
	Mongodb_db_name_str string `mapstructure:"mongodb_db_name"`

	// AWS_SQS
	AWS_SQS_queue_str string `mapstructure:"aws_sqs_queue"`

	// WORKERS_AWS_DISCOVERY - if AWS API's should be used for workers discovery,
	//                         or if the hardcoded workers_host config is used.
	Workers_aws_discovery_bool bool `mapstructure:"workers_aws_discovery"`

	// WORKERS_HOSTS - list of "," separated hosts, that are used by gf_eth_monitor__masters
	//                 to reach a worker_inspector service running on each worker.
	Workers_hosts_str string `mapstructure:"workers_hosts"`

	// SENTRY_ENDPOINT
	Sentry_endpoint_str string `mapstructure:"sentry_endpoint"`

	// EVENTS - flag to turn on/off event consumption and processing from queues. 
	//          mostly used for debugging and testing.
	Events_consume_bool bool `mapstructure:"events_consume"`

	// PY_PLUGINS
	Py_plugins_dir_path_str string `mapstructure:"py_plugins_dir_path"`

	// EXTERN_SERVICES
	AlchemyAPIkeyStr string `mapstructure:"alchemy_api_key"`
	OpenSeaAPIkeyStr string `mapstructure:"opensea_api_key"`

	// // INFLUXDB
	// Influxdb_host_str    string `mapstructure:"influxdb_host"`
	// Influxdb_db_name_str string `mapstructure:"influxdb_db_name"`
}

//-------------------------------------------------

func SentryInit(pSentryEndpointURIstr string) {

	//-------------
	// SENTRY
	sentrySampleRateDefaultF := 0.5
	sentryTracingRateForHandlersMap := map[string]float64{
		"GET /gfethm/v1/block/index":   1.0,
		"GET /gfethm/v1/tx/trace/plot": 1.0,
		"GET /gfethm/v1/block":         1.0,
		"GET /gfethm/v1/miner":         1.0,
		"GET /gfethm/v1/peers":         1.0,
	}
	err := gf_core.ErrorInitSentry(pSentryEndpointURIstr,
		sentryTracingRateForHandlersMap,
		sentrySampleRateDefaultF)
	if err != nil {
		panic(err)
	}

	defer sentry.Flush(2 * time.Second)
}

//-------------------------------------------------

func RuntimeGet(pConfig *GF_config,
	pRuntimeSys *gf_core.RuntimeSys) (*GF_runtime, error) {

	//--------------------
	// MONGODB
	mongodb_host_str := pConfig.Mongodb_host_str
	mongodb_url_str  := fmt.Sprintf("mongodb://%s", mongodb_host_str)
	fmt.Printf("mongodb_host - %s\n", mongodb_host_str)

	mongodb_db, _, gfErr := gf_core.MongoConnectNew(mongodb_url_str,
		pConfig.Mongodb_db_name_str,
		nil,
		pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr.Error
	}
	pRuntimeSys.Mongo_db = mongodb_db

	fmt.Printf("mongodb connected...\n")

	//--------------------
	// // INFLUXDB
	// influxdb_host_str := pConfig.Influxdb_host_str
	// influxdb_client   := influxdb__init(influxdb_host_str)
	// fmt.Printf("influxdb connected...\n")

	// PY_PLUGINS
	plugins_info := &GF_py_plugins{
		Base_dir_path_str: pConfig.Py_plugins_dir_path_str,
	}

	//--------------------
	// RUNTIME
	runtime := &GF_runtime{
		Config:     pConfig,
		Py_plugins: plugins_info,
		RuntimeSys: pRuntimeSys,
		// Influxdb_client: influxdb_client,
	}

	//--------------------
	return runtime, nil
}

//-------------------------------------------------

/*// INFLUXDB
func influxdb__init(p_influxdb_host_str string) *influxdb2.Client {

	fmt.Println("influxdb get client...")
	client := influxdb2.NewClient(p_influxdb_host_str, "my-token")
	return &client
}*/