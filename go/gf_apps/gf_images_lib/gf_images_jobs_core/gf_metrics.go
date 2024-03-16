// SPDX-License-Identifier: GPL-2.0
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

package gf_images_jobs_core

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_core/gf_images_plugins"
)

//-------------------------------------------------

type GFmetrics struct {
	Cmd__start_job_local_imgs__count     prometheus.Counter
	Cmd__start_job_transform_imgs__count prometheus.Counter
	Cmd__start_job_uploaded_imgs__count  prometheus.Counter
	Cmd__start_job_extern_imgs__count    prometheus.Counter
	CmdStartJobClassifyImagesCount       prometheus.Counter
	ImagesPluginsMetrics                 *gf_images_plugins.GFmetrics
}

//-------------------------------------------------

func MetricsCreate(pNamespaceStr string) *GFmetrics {

	// CMD__START_JOB_LOCAL_IMAGES__COUNT
	cmd__start_job_local_imgs__count := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: pNamespaceStr,
		Name: fmt.Sprintf("cmd__start_job_local_imgs__count"),
		Help: "job command start_job_local_imgs #",
	})
	prometheus.MustRegister(cmd__start_job_local_imgs__count)

	// CMD__START_JOB_LOCAL_IMAGES__COUNT
	cmd__start_job_transform_imgs__count := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: pNamespaceStr,
		Name: fmt.Sprintf("cmd__start_job_transform_imgs__count"),
		Help: "job command start_job_transform_imgs #",
	})
	prometheus.MustRegister(cmd__start_job_transform_imgs__count)

	// CMD__START_JOB_UPLOAD_IMGS__COUNT
	cmd__start_job_uploaded_imgs__count := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: pNamespaceStr,
		Name: fmt.Sprintf("cmd__start_job_uploaded_imgs__count"),
		Help: "job command start_job_uploaded_imgs #",
	})
	prometheus.MustRegister(cmd__start_job_uploaded_imgs__count)

	// CMD__START_JOB_EXTERN_IMGS__COUNT
	cmd__start_job_extern_imgs__count := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: pNamespaceStr,
		Name: fmt.Sprintf("cmd__start_job_extern_imgs__count"),
		Help: "job command start_job_extern_imgs #",
	})
	prometheus.MustRegister(cmd__start_job_extern_imgs__count)

	// CMD__START_JOB_CLASSIFY_IMGS__COUNT
	cmdStartJobClassifyImagesCount := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: pNamespaceStr,
		Name: fmt.Sprintf("cmd__start_job_classify_imgs__count"),
		Help: "job command start_job_classify_imgs #",
	})
	prometheus.MustRegister(cmdStartJobClassifyImagesCount)

	// IMAGES_PLUGINS
	imagesPluginsMetrics := gf_images_plugins.MetricsCreate(pNamespaceStr)

	metrics := &GFmetrics{
		Cmd__start_job_local_imgs__count:     cmd__start_job_local_imgs__count,
		Cmd__start_job_transform_imgs__count: cmd__start_job_transform_imgs__count,
		Cmd__start_job_uploaded_imgs__count:  cmd__start_job_uploaded_imgs__count,
		Cmd__start_job_extern_imgs__count:    cmd__start_job_extern_imgs__count,
		CmdStartJobClassifyImagesCount:       cmdStartJobClassifyImagesCount,
		ImagesPluginsMetrics:                 imagesPluginsMetrics,
	}
	return metrics
}