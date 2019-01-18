/*
GloFlow media management/publishing system
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

type error_def struct {
	Descr_str string
}
//-------------------------------------------------
func error__get_defs() map[string]error_def {

	error_defs_map := map[string]error_def{

		//---------------
		"panic_error":error_def{
			Descr_str:"a golang panic was caught with recover()",
		},
		//---------------
		"int_parse_error":error_def{
			Descr_str:"failed to parse an integer string",
		},
		"url_parse_error":error_def{
			Descr_str:"failed to parse a url with url.Parse()",
		},
		"io_reader_error":error_def{
			Descr_str:"failed to read bytes using the io.Reader.ReadBytes()",
		},
		//---------------
		//DATA_VERIFICATION
		"verify__invalid_value_error":{
			Descr_str:"data failed verification, not an expected value",
		},
		"verify__value_not_integer_error":{
			Descr_str:"data failed verification, the supplied value is not an integer",
		},
		"verify__value_too_many_error":{
			Descr_str:"data failed verification, the supplied too many values",
		},
		"verify__missing_key_error":{
			Descr_str:"data failed verification, the needed key is missing",
		},
		"verify__invalid_key_value_error":{
			Descr_str:"data failed verification, the key does not have the expected value",
		},
		"verify__string_too_short_error":{
			Descr_str:"data failed verification, the string is too short",
		},
		"verify__string_too_long_error":{
			Descr_str:"data failed verification, the string is too long",
		},
		"verify__invalid_image_extension_error":{
			Descr_str:"an unsupported image file extension was encountered",
		},
		"verify__invalid_query_string_encoding_error":{
			Descr_str:"string is not a valid query-string encoding",
		},
		"verify__invalid_image_nsfv_error":{
			Descr_str:"image NSFV verification failed",
		},
		//---------------
		//FILESYSTEM
		"file_open_error":error_def{
			Descr_str:"os.Create() failed to create a file - package (os)",
		},
		"file_create_error":error_def{
			Descr_str:"os.Open() failed to open a file - package (os)",
		},
		"file_remove_error":error_def{
			Descr_str:"os.Remove() failed to remove a file - package (os)",
		},
		"file_write_error":error_def{
			Descr_str:"f.Write() failed to write to a file - package (os)",
		},
		"file_sync_error":error_def{
			Descr_str:"f.Sync() failed to sync a file to the FS - package (os)",
		},
		"file_missing_error":error_def{
			Descr_str:"file doesnt exist in the FS",
		},
		"file_buffer_copy_error":error_def{
			Descr_str:"using a file as a source/target of a buffer copy failed - (io.Copy(),etc.)",
		},
		"dir_list_error":error_def{
			Descr_str:"failed to list contents of a dir in the FS",
		},
		//---------------
		//ENCODE/DECODE
		"json_decode_error":error_def{
			Descr_str:"json.Unmarshal() failed to parse byte array  - package (encoding/json)",
		},
		"json_encode_error":error_def{
			Descr_str:"json.Marshal() failed to parse byte array  - package (encoding/json)",
		},
		"base64_decoding_error":error_def{
			Descr_str:"base64.StdEncoding.DecodeString() failed - package (encoding/base64)",
		},
		//---------------
		//IMAGES
		"image_decoding_error":error_def{
			Descr_str:"image.Decode() failed to decode image data - package (image)",
		},
		"image_decoding_config_error":error_def{
			Descr_str:"image.DecodeConfig() failed to decode image data - package (image,image/png,image/jpeg,etc.)",
		},
		"jpeg_decoding_error":error_def{
			Descr_str:"jpeg.Decode() failed to decode image data - package (image/jpeg)",
		},
		"png_decoding_error":error_def{
			Descr_str:"png.Decode() failed to decode image data - package (image/png)",
		},
		"png_encoding_error":error_def{
			Descr_str:"png.Encode() failed to encode image data - package (image/png)",
		},
		"gif_decoding_frames_error":error_def{
			Descr_str:"gif.DecodeAll() failed to decode GIF frames - package (image/gif)",
		},
		//---------------
		//MONGODB
		"mongodb_find_error":error_def{
			Descr_str:"c.Find() failed to find a mongodb document - package (mgo)",
		},
		"mongodb_not_found_error":error_def{
			Descr_str:"target document not found in mongodb - package (mgo)",
		},
		"mongodb_insert_error":error_def{
			Descr_str:"c.Insert() failed to insert/create new mongodb document - package (mgo)",
		},
		"mongodb_update_error":error_def{
			Descr_str:"c.Update() failed to update a mongodb document- package (mgo)",
		},
		"mongodb_delete_error":error_def{
			Descr_str:"c.Update() failed to update a mongodb document- package (mgo)",
		},
		"mongodb_aggregation_error":error_def{
			Descr_str:"pipe.All() failed to run a aggregation pipeline in mongodb - package (mgo)",
		},
		"mongodb_ensure_index_error":error_def{
			Descr_str:"c.EnsureIndex() failed to create a mongodb index - package (mgo)",
		},
		//---------------
		//ELASTICSEARCH
		"elasticsearch_get_client":error_def{
			Descr_str:"c.NewClient() failed to get elasticsearch client - package (elastic)",
		},
		"elasticsearch_ping":error_def{
			Descr_str:"c.Ping() failed to ping elasticsearch server from client - package (elastic)",
		},
		"elasticsearch_add_to_index":error_def{
			Descr_str:"c.Index() failed to add a record to the index - package (elastic)",
		},
		"elasticsearch_query_index":error_def{
			Descr_str:"c.Search() failed issue a query - package (elastic)",
		},
		//---------------
		//TEMPLATES
		"template_create_error":error_def{
			Descr_str:"template.New() failed to create/load a template - package (text/template)",
		},
		"template_render_error":error_def{
			Descr_str:"template.Execute() failed to render a template - package (text/template)",
		},
		//---------------
		//HTTP
		"http_client_req_error":{
			Descr_str:"failed to execute a http_client request",
		},
		"http_client_req_status_error":{
			Descr_str:"http_client received a non 2xx/3xx HTTP status code",
		},
		"http_server_flusher_not_supported_error":{
			Descr_str:"http_server not supporting http.Flusher (probably for SSE support,etc.)",
		},
		//---------------
		//S3
		"s3_credentials_error":{
			Descr_str:"S3 credentials operation failed",
		},
		"s3_file_upload_error":{
			Descr_str:"failed to upload a file to S3 bucket",
		},
		"s3_file_copy_error":{
			Descr_str:"failed to copy a file within S3",
		},
		//---------------
		//HTML_PARSING
		"html_parse_error":{
			Descr_str:"parsing of a HTML document failed",
		},
		//---------------
	}
	return error_defs_map
}