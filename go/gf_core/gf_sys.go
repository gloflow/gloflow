/*
GloFlow application and media management/publishing platform
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

import (
	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/globalsign/mgo"
)

//-------------------------------------------------
type Runtime_sys struct {
	Service_name_str string
	Log_fun      func(string, string)
	Mongo_db     *mongo.Database
	Mongo_coll   *mongo.Collection // main mongodb collection to use when none is specified
	Debug_bool   bool              // if debug mode is enabled (some places will print extra info in debug mode)

	// ERRORS
	Errors_send_to_mongodb_bool bool // if errors should be persisted to Mongodb
	Errors_send_to_sentry_bool  bool // if errors should be sent to Sentry service

	// Mongodb_db   *mgo.Database   // DEPRECATED!! - remove - use Mongo_db/Mongo_coll
	// Mongodb_coll *mgo.Collection // DEPRECATED!! - remove - use Mongo_db/Mongo_coll
}