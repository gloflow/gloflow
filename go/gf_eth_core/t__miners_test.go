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

package gf_eth_core

import (
	"fmt"
	"testing"
	"context"
	"github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__miners(p_test *testing.T) {

	fmt.Println("TEST__MINERS ==============================================")


	ctx := context.Background()
	runtime, _ := t__get_runtime(p_test)

	// ethermine
	miner_addr_str := "0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8"

	miners_map, gf_err := Eth_miners__db__get_info(miner_addr_str, nil, ctx, runtime)
	if gf_err != nil {
		p_test.Fail()
	}


	spew.Dump(miners_map)
}