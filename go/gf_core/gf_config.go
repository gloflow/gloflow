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

//-------------------------------------------------------------

type GFconfig struct {

	// ENVIRONMENT
	EnvStr string `mapstructure:"env"`

	// DOMAINS - PRIMARY_DOMAIN - where this gf_solo instance is reachable on
	DomainBaseStr      string `mapstructure:"domain_base"`
	DomainAdminBaseStr string `mapstructure:"domain_admin_base"`

	// PORTS
	PortStr        string `mapstructure:"port"`
	PortAdminStr   string `mapstructure:"port_admin"`
	PortMetricsStr string `mapstructure:"port_metrics"`

	// SQL
	SQLuserNameStr string `mapstructure:"sql_user_name"`
	SQLpassStr     string `mapstructure:"sql_pass"`
	SQLhostStr     string `mapstructure:"sql_host"`
	SQLdbNameStr   string `mapstructure:"sql_db_name"`

	// MONGODB - this is the dedicated mongodb DB
	MongoHostStr   string `mapstructure:"mongodb_host"`
	MongoDBnameStr string `mapstructure:"mongodb_db_name"`

	// SENTRY_DSN
	SentryDSNstr string `mapstructure:"sentry_dsn"`

	// TEMPLATES
	TemplatesPathsMap map[string]string `mapstructure:"templates_paths"`

	//--------------------
	// IDENTITY
	AuthSubsystemTypeStr       string `mapstructure:"auth_subsystem_type"`
	AdminMFAsecretKeyBase32str string `mapstructure:"admin_mfa_secret_key_base32"`
	AdminEmailStr              string `mapstructure:"admin_email"`

	// DOMAIN_FOR_AUTH_COOKIES - domain/pattern that is set on the auth cookies to restrict their scope.
	DomainForAuthCookiesStr string `mapstructure:"domain_for_auth_cookies"`

	//--------------------
	// GF_IMAGES
	ImagesConfigFilePathStr string `mapstructure:"images__config_file_path"`

	//--------------------
	// GF_ANALYTICS

	AnalyticsPyStatsDirsLst []string `mapstructure:"analytics__py_stats_dirs"`
	AnalyticsRunIndexerBool bool     `mapstructure:"analytics__run_indexer"`

	//--------------------
	// ALCHEMY
	AlchemyAPIkeyStr string `mapstructure:"alchemy_api_key"`

	//--------------------
	// NEW_STORAGE_ENGINE - flag indicating if the new image storage engine should be used
	ImagesUseNewStorageEngineBool bool `mapstructure:"images_use_new_storage_engine"`

	// IPFS
	IPFSnodeHostStr string `mapstructure:"ipfs_node_host"`

	//--------------------
}

//-------------------------------------------------------------
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
