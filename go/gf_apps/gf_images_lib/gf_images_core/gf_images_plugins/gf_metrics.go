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

package gf_images_plugins

import (
	"github.com/prometheus/client_golang/prometheus"
)

//-------------------------------------------------

type GFmetrics struct {
	PyPluginsExecDurationGauge prometheus.Gauge
}

//-------------------------------------------------

func MetricsCreate(pNamespaceStr string) *GFmetrics {

	// PY_PLUGINS_EXEC_DURATION__GAUGE
	pyPluginsExecDurationGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: pNamespaceStr,
			Name: "py_plugins__exec_duration",
			Help: "duration in seconds for how long the Py plugin runs",
		})
	prometheus.MustRegister(pyPluginsExecDurationGauge)


	metrics := &GFmetrics{
		PyPluginsExecDurationGauge: pyPluginsExecDurationGauge,
	}
	return metrics
}