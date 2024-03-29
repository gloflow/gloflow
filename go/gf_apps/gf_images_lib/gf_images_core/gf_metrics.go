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
	"github.com/prometheus/client_golang/prometheus"
)

//-------------------------------------------------

type GFmetrics struct {

	// IMAGE_UPLOAD
	ImageUploadClientDurationGauge         prometheus.Gauge
	ImageUploadClientTransferDurationGauge prometheus.Gauge

	// IMAGE_CLASSIFY
	ImageClassifyPyExecDurationGauge prometheus.Gauge
}

//-------------------------------------------------

func MetricsCreate(pNamespaceStr string) *GFmetrics {

	// IMAGE_UPLOAD_CLIENT_DURATION
	imageUploadClientDurationGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: pNamespaceStr,
			Name: "client_upload__duration",
			Help: "duration in seconds (client reported) of how long it takes for the whole image upload process (in seconds)",
		})
	prometheus.MustRegister(imageUploadClientDurationGauge)

	// IMAGE_UPLOAD_CLIENT_TRANSFER_DURATION
	imageUploadClientTransferDurationGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: pNamespaceStr,
			Name: "client_upload__transfer_duration",
			Help: "duration in seconds (client reported) of how long it takes for the image upload data transfer (in seconds)",
		})
	prometheus.MustRegister(imageUploadClientTransferDurationGauge)


	// IMAGE_CLASSIFY_EXECUTION_DURATION
	imageClassifyPyExecDurationGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: pNamespaceStr,
			Name: "image_classify__py_exec_duration",
			Help: "duration in seconds of how long it takes for the image classification Py pipeline to execute (in seconds)",
		})
	prometheus.MustRegister(imageClassifyPyExecDurationGauge)

	

	metrics := &GFmetrics{
		ImageUploadClientDurationGauge:         imageUploadClientDurationGauge,
		ImageUploadClientTransferDurationGauge: imageUploadClientTransferDurationGauge,

		ImageClassifyPyExecDurationGauge: imageClassifyPyExecDurationGauge,
	}
	return metrics
}