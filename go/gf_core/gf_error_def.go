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

//-------------------------------------------------

type ErrorDef struct {
	DescrStr string
}

//-------------------------------------------------

func errorGetDefs() map[string]ErrorDef {

	errorDefsMap := map[string]ErrorDef{

		//---------------
		"panic_error": ErrorDef{
			DescrStr: "a golang panic was caught with recover()",
		},

		"generic_error": ErrorDef{
			DescrStr: "generic error occured, check error info for more details",
		},

		//---------------
		"int_parse_error": ErrorDef{
			DescrStr: "failed to parse an integer string",
		},
		"url_parse_error": ErrorDef{
			DescrStr: "failed to parse a url with url.Parse()",
		},
		"url_unescape_error": ErrorDef{
			DescrStr: "failed to unescape a url with url.QueryUnescape()",
		},
		"io_reader_error": ErrorDef{
			DescrStr: "failed to read bytes using a reader",
		},

		//---------------
		// IDENTITY
		"user_incorrect": ErrorDef{
			DescrStr: "incorrect user specified",
		},
		"policy__op_denied": ErrorDef{
			DescrStr: "access policy has denied the operation",
		},
		"auth_missing_cookie": ErrorDef{
			DescrStr: "missing auth cookie",
		},

		//---------------
		// DATA_VERIFICATION
		"verify__invalid_value_error": ErrorDef{
			DescrStr: "data failed verification, not an expected value",
		},
		"verify__value_not_integer_error": ErrorDef{
			DescrStr: "data failed verification, the supplied value is not an integer",
		},
		"verify__value_too_many_error": ErrorDef{
			DescrStr: "data failed verification, the supplied too many values",
		},
		"verify__missing_key_error": ErrorDef{
			DescrStr: "data failed verification, the needed key is missing",
		},
		"verify__invalid_key_value_error": ErrorDef{
			DescrStr: "data failed verification, the key does not have the expected value",
		},
		"verify__input_data_missing_in_req_error": ErrorDef{
			DescrStr: "data failed verification, input is missing in request",
		},

		// length
		"verify__string_too_short_error": ErrorDef{
			DescrStr: "data failed verification, the string is too short",
		},
		"verify__string_too_long_error": ErrorDef{
			DescrStr: "data failed verification, the string is too long",
		},
		"verify__string_not_correct_length_error": ErrorDef{
			DescrStr: "data failed verification, the string is too long",
		},

		"verify__invalid_image_extension_error": ErrorDef{
			DescrStr: "an unsupported image file extension was encountered",
		},
		"verify__invalid_query_string_encoding_error": ErrorDef{
			DescrStr: "string is not a valid query-string encoding",
		},
		"verify__invalid_image_nsfv_error": ErrorDef{
			DescrStr: "image NSFV verification failed",
		},

		// VALIDATOR - used primarily when validating using validator lib,
		//             where struct tags contain directives on how to validate individual fields.
		"verify__invalid_input_struct_error": ErrorDef{
			DescrStr: "input struct is invalid",
		},

		"verify__sess_data_missing_in_req": ErrorDef{
			DescrStr: "session data missing in http request",
		},
		
		//---------------
		// FILESYSTEM
		"file_open_error": ErrorDef{
			DescrStr: "os.Open() failed to open a file - package (os)",
		},
		"file_create_error": ErrorDef{
			DescrStr: "os.Create() failed to create a file - package (os)",
		},
		"file_read_error": ErrorDef{
			DescrStr: "f.Read()/ioutil.ReadFile() failed to read file - package (os/ioutil)",
		},
		"file_remove_error": ErrorDef{
			DescrStr: "os.Remove() failed to remove a file - package (os)",
		},
		"file_write_error": ErrorDef{
			DescrStr: "f.Write() failed to write to a file - package (os)",
		},
		"file_sync_error": ErrorDef{
			DescrStr: "f.Sync() failed to sync a file to the FS - package (os)",
		},
		"file_missing_error": ErrorDef{
			DescrStr: "file doesnt exist in the FS",
		},
		"file_buffer_copy_error": ErrorDef{
			DescrStr: "using a file as a source/target of a buffer copy failed - (io.Copy(),etc.)",
		},
		"file_stat_error": ErrorDef{
			DescrStr: "getting info on a file via a stat() system call (golang API or CLI) failed - (os.Stat())",
		},
		"dir_list_error": ErrorDef{
			DescrStr: "failed to list contents of a dir in the FS",
		},

		//---------------
		// CLI
		"cli_run_error": ErrorDef{
			DescrStr: "failed to run a CLI command from Go",
		},

		//---------------
		// ENCODE/DECODE
		// JSON
		"json_decode_error": ErrorDef{
			DescrStr: "json.Unmarshal() failed to parse byte array - package (encoding/json)",
		},
		"json_encode_error": ErrorDef{
			DescrStr: "json.Marshal() failed to parse byte array - package (encoding/json)",
		},
		// YAML
		"yaml_decode_error": ErrorDef{
			DescrStr: "yaml.Unmarshal() failed to parse byte array - package (gopkg.in/yaml.v2)",
		},
		// BASE64
		"base64_decoding_error": ErrorDef{
			DescrStr: "base64.StdEncoding.DecodeString() failed - package (encoding/base64)",
		},
		// HEX
		"decode_hex": ErrorDef{
			DescrStr: "failed to decode hex string",
		},
		// MAPSTRUCT
		"mapstruct_decode": ErrorDef{
			DescrStr: "failed to decode a map into a struct using mapstructure lib",
		},

		//---------------
		// HTTP
		"http_client_req_error": ErrorDef{
			DescrStr:"failed to execute a http_client request",
		},
		"http_client_req_status_error": ErrorDef{
			DescrStr:"http_client received a non 2xx/3xx HTTP status code",
		},
		"http_server_flusher_not_supported_error": ErrorDef{
			DescrStr:"http_server not supporting http.Flusher (probably for SSE support,etc.)",
		},
		"http_client_gf_status_error": ErrorDef{
			DescrStr:"http_client received a non-OK GF error",
		},
		"http_cookie": ErrorDef{
			DescrStr:"failed to handle a http cookie",
		},
		"html_parse_error": ErrorDef{
			DescrStr: "parsing of a HTML document failed",
		},
		
		//---------------
		"rpc_context_value_missing": ErrorDef{
			DescrStr: "expected key is missing from a gf_rpc context",
		},

		//---------------
		// WEBSOCKETS
		"ws_connection_init_error": ErrorDef{
			DescrStr: "websocket client failed to connect to a url",
		},

		//---------------
		// UDP
		"udp_open_socket_error": ErrorDef{
			DescrStr: "failed to open UDP socket listening on a port",
		},
		"udp_write_packge_to_socket_error": ErrorDef{
			DescrStr: "failed to write package to UDP socket",
		},
		
		//---------------
		// IMAGES
		"image_decoding_error": ErrorDef{
			DescrStr: "image.Decode() failed to decode image data - package (image)",
		},
		"image_decoding_config_error": ErrorDef{
			DescrStr: "image.DecodeConfig() failed to decode image data - package (image,image/png,image/jpeg,etc.)",
		},
		"jpeg_decoding_error": ErrorDef{
			DescrStr: "jpeg.Decode() failed to decode image data - package (image/jpeg)",
		},
		"png_decoding_error": ErrorDef{
			DescrStr: "png.Decode() failed to decode image data - package (image/png)",
		},
		"png_encoding_error": ErrorDef{
			DescrStr: "png.Encode() failed to encode image data - package (image/png)",
		},
		"gif_decoding_frames_error": ErrorDef{
			DescrStr: "gif.DecodeAll() failed to decode GIF frames - package (image/gif)",
		},

		//---------------
		// MONGODB
		"mongodb_connect_error": ErrorDef{
			DescrStr: "failed to connect to a mongodb host - package (go.mongodb.org/mongo-driver)",
		},
		"mongodb_ping_error": ErrorDef{
			DescrStr: "failed to ping a mongodb host - package (go.mongodb.org/mongo-driver)",
		},		
		"mongodb_find_error": ErrorDef{
			DescrStr:"c.Find() failed to find a mongodb document",
		},
		"mongodb_count_error": ErrorDef{
			DescrStr:"Count of documents failed in mongodb",
		},
		"mongodb_not_found_error": ErrorDef{
			DescrStr:"target document not found in mongodb",
		},
		"mongodb_insert_error": ErrorDef{
			DescrStr:"failed to insert/create new mongodb document",
		},
		"mongodb_write_bulk_error": ErrorDef{
			DescrStr:"failed to bulk write new mongodb documents",
		},
		"mongodb_update_error": ErrorDef{
			DescrStr:"failed to update a mongodb document",
		},
		"mongodb_delete_error": ErrorDef{
			DescrStr:"failed to delete a mongodb document",
		},
		"mongodb_aggregation_error": ErrorDef{
			DescrStr:"failed to run a aggregation pipeline in mongodb",
		},
		"mongodb_ensure_index_error": ErrorDef{
			DescrStr:"c.EnsureIndex() failed to create a mongodb index",
		},
		"mongodb_cursor_decode": ErrorDef{
			DescrStr:"failed to decode value from the mongodb results Cursor",
		},
		"mongodb_cursor_all": ErrorDef{
			DescrStr:"failed to get all values from the mongodb results Cursor",
		},
		"mongodb_session_error": ErrorDef{
			DescrStr:"failed to execute mongodb session",
		},
		"mongodb_start_session_error": ErrorDef{
			DescrStr:"failed to start a new mongodb session",
		},
		"mongodb_session_abort_error": ErrorDef{
			DescrStr:"failed to abort a mongodb session",
		},
		"mongodb_get_collection_names_error": ErrorDef{
			DescrStr:"failed to get all mongodb collection names",
		},
		
		//---------------
		// SQL
		"sql_failed_to_connect": ErrorDef{
			DescrStr:"failed to establish an initial connection to an SQL server",
		},
		"sql_table_creation": ErrorDef{
			DescrStr:"failed to create an SQL table",
		},
		"sql_row_insert": ErrorDef{
			DescrStr:"failed to insert a new row into a SQL table",
		},
		"sql_transaction_begin": ErrorDef{
			DescrStr:"failed to begin a transaction to an SQL server",
		},
		"sql_transaction_commit": ErrorDef{
			DescrStr:"failed to commit a transaction to an SQL server",
		},
		"sql_query_execute": ErrorDef{
			DescrStr:"failed to execute a query in a SQL server",
		},
		"sql_row_scan": ErrorDef{
			DescrStr:"failed to scan a row of results of a SQL query",
		},
		"sql_prepare_statement": ErrorDef{
			DescrStr:"failed to prepare an SQL statement with Prepare()",
		},
		"sql_generic_error": ErrorDef{
			DescrStr:"generic SQL error",
		},
		
		//---------------
		// REDIS
		"redis_cmd": ErrorDef{
			DescrStr:"failed to execute a Redis command on a redis server",
		},

		//---------------
		// ELASTICSEARCH
		"elasticsearch_get_client": ErrorDef{
			DescrStr:"c.NewClient() failed to get elasticsearch client - package (elastic)",
		},
		"elasticsearch_ping": ErrorDef{
			DescrStr:"c.Ping() failed to ping elasticsearch server from client - package (elastic)",
		},
		"elasticsearch_add_to_index": ErrorDef{
			DescrStr:"c.Index() failed to add a record to the index - package (elastic)",
		},
		"elasticsearch_query_index": ErrorDef{
			DescrStr:"c.Search() failed issue a query - package (elastic)",
		},

		//---------------
		// TEMPLATES
		"template_create_error": ErrorDef{
			DescrStr:"template.New() failed to create/load a template - package (text/template)",
		},
		"template_render_error": ErrorDef{
			DescrStr:"template.Execute() failed to render a template - package (text/template)",
		},
		
		//---------------
		// AWS
		"aws_general_error": ErrorDef{
			DescrStr: "AWS general error",
		},
		"aws_session_create": ErrorDef{
			DescrStr: "AWS failed to create new API session",
		},
		"aws_client_v2_create": ErrorDef{
			DescrStr: "AWS failed to create new API V2 client",
		},

		// EC2
		"aws_ec2_instances_describe_error": ErrorDef{
			DescrStr: "failed to describe EC2 instances",
		},
		"aws_ec2_autoscaling_scale_error": ErrorDef{
			DescrStr: "failed to change the count of instances in a EC2 autoscaling group",
		},
		"aws_ec2_autoscaling_describe_error": ErrorDef{
			DescrStr: "failed to get info on a EC2 autoscaling group",
		},
		
		// ECS
		"aws_ecs_service_update_error": ErrorDef{
			DescrStr: "failed to update an AWS ECS service",
		},

		// SECRETS_MNGR
		"aws_secretsmngr_create_secret_value_error": ErrorDef{
			DescrStr: "failed to create secret in AWS SECRETS_MANAGER service",
		},
		"aws_secretsmngr_get_secret_value_error": ErrorDef{
			DescrStr: "failed to get secret value from AWS SECRETS_MANAGER service",
		},
		
		// SQS
		"aws_sqs_queue_create_error": ErrorDef{
			DescrStr: "failed to create SQS queue",
		},
		"aws_sqs_queue_get_url_error": ErrorDef{
			DescrStr: "failed to get a URL of a SQS queue",
		},
		"aws_sqs_queue_send_msg_error": ErrorDef{
			DescrStr: "failed to send a message to a SQS queue",
		},
		"aws_sqs_queue_receive_msg_error": ErrorDef{
			DescrStr: "failed to receive a message from a SQS queue",
		},
		"aws_sqs_queue_delete_msg_error": ErrorDef{
			DescrStr: "failed to delete a message from a SQS queue",
		},

		// SES
		"aws_ses_service_send_message_error": ErrorDef{
			DescrStr: "failed to send an email message via SES",
		},
		"aws_ses_service_verify_address_error": ErrorDef{
			DescrStr: "failed to verify email address via SES",
		},

		// S3
		"s3_credentials_error": ErrorDef{
			DescrStr: "S3 credentials operation failed",
		},
		"s3_file_upload_error": ErrorDef{
			DescrStr: "failed to upload a file to S3 bucket",
		},
		"s3_file_upload_url_presign_error": ErrorDef{
			DescrStr: "failed to get a presigned URL for uploading a file to S3 bucket",
		},
		"s3_file_copy_error": ErrorDef{
			DescrStr: "failed to copy a file within S3",
		},
		"s3_file_download_error": ErrorDef{
			DescrStr: "failed to download a file from S3 to a local FS",
		},

		// IAM
		"iam_error": ErrorDef{
			DescrStr: "failed to interact with the IAM AWS service",
		},

		//---------------
		// LIBRARY_ERROR
		"library_error": ErrorDef{
			DescrStr: "third-party library has failed",
		},

		//---------------
		// CRYPTO
		"crypto_jwt_sign_token_error": ErrorDef{
			DescrStr: "failed to crypto-sign JWT token",
		},
		"crypto_jwt_verify_token_error": ErrorDef{
			DescrStr: "failed to crypto-verify JWT token",
		},
		"crypto_jwt_verify_token_invalid_error": ErrorDef{
			DescrStr: "JWT token is invalid",
		},
		"crypto_jwt_parse_token_error": ErrorDef{
			DescrStr: "JWT token failed to parse",
		},
		"crypto_ec_recover_pubkey": ErrorDef{
			DescrStr: "failed to recovery Pubkey fro signature",
		},
		"crypto_hex_decode": ErrorDef{
			DescrStr: "failed to decodes a hex string with 0x prefix",
		},
		"crypto_cert_ca_parse": ErrorDef{
			DescrStr: "failed to parse cert CA",
		},
		"crypto_signature_eth_last_byte_invalid_value": ErrorDef{
			DescrStr: "last byte of an ethereum signature does not have the proper V value",
		},

		// PEM
		"crypto_pem_decode": ErrorDef{
			DescrStr: "failed to decode a PEM (private key cryptographic info)",
		},
		"crypto_x509_parse": ErrorDef{
			DescrStr: "failed to parse x509 info",
		},

		//---------------
		// GF_LANG

		"gf_lang_program_run_failed": ErrorDef{
			DescrStr: "failed to execute a gf_lang program",
		},
		
		//---------------
	}
	return errorDefsMap
}