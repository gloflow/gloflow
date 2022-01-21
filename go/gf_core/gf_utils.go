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
	"fmt"
	"time"
	"crypto/sha256"
	"math/rand"
)

//-------------------------------------------------------------
func Str_random() string {
	rand_with_seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand_int := rand_with_seed.Int()
	rand_str := fmt.Sprintf("%d", rand_int)
	return rand_str
}

//-------------------------------------------------------------
func Hash_val_sha256(p_val interface{}) string {

	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", p_val)))
	hash_hex_str := fmt.Sprintf("%x", h.Sum(nil))
	return hash_hex_str
}

//-------------------------------------------------------------
func Str_in_lst(p_str string, p_lst []string) bool {
	for _,s := range p_lst {
		if p_str == s {
			return true
		}
	}
	return false
}

//-------------------------------------------------------------
func Map_has_key(p_map interface{}, p_key_str string) bool {

	if _,ok := p_map.(map[string]interface{}); ok {
		_,ok := p_map.(map[string]interface{})[p_key_str]
		return ok
	} else if _,ok := p_map.(map[string]string); ok {
		_,ok := p_map.(map[string]string)[p_key_str]
		return ok
	}

	panic("no handler for p_map type")
	return false
}