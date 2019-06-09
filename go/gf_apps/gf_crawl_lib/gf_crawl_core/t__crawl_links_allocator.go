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

package gf_crawl_core


import (
	//"fmt"
	"testing"
	//"github.com/stretchr/testify/assert"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__links_allocator(p_test *testing.T) {

	//-------------------
	test__crawler_name_str       := "test_crawler"
	test__block_size_int         := 100
	runtime_sys, _ := T__init()
	//-------------------



	//CREATE_ALLOCATOR
	allocator, gf_err := Link_alloc__create(test__crawler_name_str, runtime_sys)
	if gf_err != nil {

		panic(gf_err.Error)

	}
	spew.Dump(allocator)


	//CREATE_ALLOCATOR_LINKS_BLOCK
	alloc_block, gf_err := Link_alloc__create_links_block(allocator.Id_str, test__crawler_name_str, test__block_size_int, runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}
	spew.Dump(alloc_block)
}

	