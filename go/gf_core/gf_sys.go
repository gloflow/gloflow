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
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/go-playground/validator"
)

//-------------------------------------------------
// RUNTIME_SYS
type Runtime_sys struct {
	Service_name_str string
	Log_fun          func(string, string)
	Mongo_db         *mongo.Database
	Mongo_coll       *mongo.Collection // main mongodb collection to use when none is specified
	Debug_bool       bool              // if debug mode is enabled (some places will print extra info in debug mode)

	// ERRORS
	Errors_send_to_mongodb_bool bool // if errors should be persisted to Mongodb
	Errors_send_to_sentry_bool  bool // if errors should be sent to Sentry service

	Names_prefix_str string

	Validator *validator.Validate

	External_plugins *External_plugins
}

// PLUGINS
type External_plugins struct {
	EventCallback        func(string, map[string]interface{}, *Runtime_sys)         *GF_error
	SecretCreateCallback func(string, map[string]interface{}, string, *Runtime_sys) *GF_error
	SecretGetCallback    func(string, *Runtime_sys) (map[string]interface{}, *GF_error)
}