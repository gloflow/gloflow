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

package gf_images_core

import (
	// "github.com/prometheus/client_golang/prometheus"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_plugins"
)

//-------------------------------------------------
type GFmetrics struct {
	PyPluginsMetrics *gf_images_plugins.GFmetrics
}

//-------------------------------------------------
func MetricsCreate(pNamespaceStr string) *GFmetrics {

	pyPluginsMetrics := gf_images_plugins.MetricsCreate(pNamespaceStr)

	metrics := &GFmetrics{
		PyPluginsMetrics: pyPluginsMetrics,
	}
	return metrics
}