








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
}