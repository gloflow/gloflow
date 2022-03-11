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

package gf_identity_lib

import (
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func usersCreateID(p_user_identifier_str string,
	p_creation_unix_time_f float64) gf_core.GF_ID {

	fields_for_id_lst := []string{
		p_user_identifier_str,
	}
	gf_id_str := gf_core.ID__create(fields_for_id_lst,
		p_creation_unix_time_f)

	return gf_id_str
}