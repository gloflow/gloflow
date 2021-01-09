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

package gf_eth_monitor_lib

//-------------------------------------------------------------
type GF_config struct {

	// PORTS
	Port_str         string `mapstructure:"port"`
	Port_metrics_str string `mapstructure:"port_metrics"`

	// MONGODB - this is the dedicated mongodb DB
	Mongodb_host_str    string `mapstructure:"mongodb_host"`
	Mongodb_db_name_str string `mapstructure:"mongodb_db_name"`

	// INFLUXDB
	Influxdb_host_str    string `mapstructure:"influxdb_host"`
	Influxdb_db_name_str string `mapstructure:"influxdb_db_name"`

	// AWS_SQS
	AWS_SQS_queue_str string `mapstructure:"aws_sqs_queue"`

	// WORKERS_INSPECTORS_HOSTS - list of "," separated hosts, that are used by gf_eth_monitor__masters
	//                            to reach a worker_inspector service running on each worker.
	Workers_inspectors_hosts_str string `mapstructure:"workers_inspectors_hosts"`
}