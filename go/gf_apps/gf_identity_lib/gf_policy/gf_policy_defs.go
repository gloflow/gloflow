/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

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

package gf_policy



const (
	GF_POLICY_OP__FLOW_GET = "gf:images:flow_get"
	
	GF_POLICY_OP__FLOW_ADD_TAG         = "gf:images:flow_add_tag"
	GF_POLICY_OP__FLOW_REMOVE_TAG      = "gf:images:flow_remove_tag"
	GF_POLICY_OP__FLOW_ADD_IMG_TAG     = "gf:images:flow_add_img_tag"
	GF_POLICY_OP__FLOW_REMOVE_IMG_TAG  = "gf:images:flow_remove_img_tag"
	GF_POLICY_OP__FLOW_ADD_IMG_NOTE    = "gf:images:flow_add_img_note"
	GF_POLICY_OP__FLOW_REMOVE_IMG_NOTE = "gf:images:flow_remove_img_note"

	GF_POLICY_OP__FLOW_ADD_IMG    = "gf:images:flow_add_img"
	GF_POLICY_OP__FLOW_REMOVE_IMG = "gf:images:flow_remove_img"
)

//---------------------------------------------------
func getDefs() map[string][]string {


	defsLst := map[string][]string{

		"viewing": []string{
			GF_POLICY_OP__FLOW_GET,
		},

		"tagging": []string{
			GF_POLICY_OP__FLOW_ADD_TAG,
			GF_POLICY_OP__FLOW_REMOVE_TAG,
			GF_POLICY_OP__FLOW_ADD_IMG_TAG,
			GF_POLICY_OP__FLOW_REMOVE_IMG_TAG,
			GF_POLICY_OP__FLOW_ADD_IMG_NOTE,
			GF_POLICY_OP__FLOW_REMOVE_IMG_NOTE,
		},

		"editing": []string{
			GF_POLICY_OP__FLOW_ADD_IMG,
			GF_POLICY_OP__FLOW_REMOVE_IMG,
		},
	}


	return defsLst

}