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

package gf_images_jobs

/*
#cgo LDFLAGS: -L../../../../rust/build/ -lgf_images_jobs
#include "../../../../rust/gf_images_jobs/gf_images_jobs.h"
*/
import "C"

// IMPORTANT!! - LDFLAGS - options for the external (GCC "ld") linker.
//               "-L"    - linker flag for directory in which to look for libs (.so/.a).
//               "-l"    - linker flag for name of the library to link.
//                         this name is a short version of the full lib name:
//                         "gf_images_jobs" name is a full name "libgf_images_jobs.so"|"libgf_images_jobs.a"
// "#include" - points to the C header files to use to load C definitions
//              used in this file in Go.

func run_job_rust() {



	job_name_str := "rust_job"
	C.c__run_job(C.CString(job_name_str))


}