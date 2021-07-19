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
)

//-------------------------------------------------
type GF_metrics struct {
	Cmd__start_job_local_imgs__count     prometheus.Counter
	Cmd__start_job_transform_imgs__count prometheus.Counter
	Cmd__start_job_uploaded_imgs__count  prometheus.Counter
	Cmd__start_job_extern_imgs__count    prometheus.Counter
}

//-------------------------------------------------
func Metrics__create() *GF_metrics {



	// CMD__START_JOB_LOCAL_IMAGES__COUNT
	cmd__start_job_local_imgs__count := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("gf_images_jobs__cmd__start_job_local_imgs__count"),
		Help: "job command start_job_local_imgs #",
	})
	prometheus.MustRegister(cmd__start_job_local_imgs__count)



	// CMD__START_JOB_LOCAL_IMAGES__COUNT
	cmd__start_job_transform_imgs__count := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("gf_images_jobs__cmd__start_job_transform_imgs__count"),
		Help: "job command start_job_transform_imgs #",
	})
	prometheus.MustRegister(cmd__start_job_transform_imgs__count)



	// CMD__START_JOB_UPLOAD_IMGS__COUNT
	cmd__start_job_uploaded_imgs__count := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("gf_images_jobs__cmd__start_job_uploaded_imgs__count"),
		Help: "job command start_job_uploaded_imgs #",
	})
	prometheus.MustRegister(cmd__start_job_uploaded_imgs__count)



	// CMD__START_JOB_EXTERN_IMGS__COUNT
	cmd__start_job_extern_imgs__count := prometheus.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("gf_images_jobs__cmd__start_job_extern_imgs__count__count"),
		Help: "job command start_job_extern_imgs #",
	})
	prometheus.MustRegister(cmd__start_job_extern_imgs__count)




	metrics := &GF_metrics{
		Cmd__start_job_local_imgs__count:     cmd__start_job_local_imgs__count,
		Cmd__start_job_transform_imgs__count: cmd__start_job_transform_imgs__count,
		Cmd__start_job_uploaded_imgs__count:  cmd__start_job_uploaded_imgs__count,
		Cmd__start_job_extern_imgs__count:    cmd__start_job_extern_imgs__count,
	}
	return metrics
}