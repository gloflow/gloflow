/*
GloFlow application and media management/publishing platform
Copyright (C) 2020 Ivan Trajkovic

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
	"strings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//-------------------------------------------------
// reads config argument either from the CLI or from a config (file or ENV vars)
func ConfigGetArg(pArgNameStr string, pCmd *cobra.Command) string {

	argValStr := viper.GetString(pArgNameStr)
	if argValStr == "" {
		argValStr, _ = pCmd.Flags().GetString(pArgNameStr)
	}
	if argValStr == "" {
		argValStr = viper.GetString("GF_" + strings.ToUpper(pArgNameStr))
	}
	return argValStr
}