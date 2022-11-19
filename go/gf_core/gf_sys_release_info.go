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

package gf_core

//-------------------------------------------------

type SysReleaseInfo struct {
	Name_str        string
	Version_str     string
    Description_str string
    Git_commit_str  string // indicates this is usually pasted in by CI systems
}

var GitCommitSHAstr = "GF_COMMIT_SHA"

//-------------------------------------------------

func GetSysReleseInfo(pRuntimeSys *RuntimeSys) SysReleaseInfo {

	r := SysReleaseInfo{
		Name_str:        "nation_genesis",
		Version_str:     "0.10.0.0", // currently deployed version
        Description_str: `
0.10.0 (nation_genesis):
    - added first version of a user system (gf_identity)
    - adding first version of p2p system based on libp2p
    - initial version of gf_home, personal control panel for users
    - gf_admin added, for admin login and control of a particular GF server/node
0.9.0 (solo_learner):
    - new gf_solo service has been introduced, that compiles all other services into a single binary.
    - first additions of ML functionality. Addition of Rust code for collage composition and efficient building of .tfrecords
    from images. addition of basic ML models in Py. 
    - various improvements to gf_core, around error handling, etc. improvements to gf_rpc_lib.      
0.8.0 (precious):
    - OPEN_SOURCING!! GF has been open sourced. gf_publisher, gf_landing_page, gf_tagger have also been refactored to use the new error handling structure,
    thats in place for gf_images and gf_analytics (gf_crawl_lib).
0.7.4:
    - first prototype of the most basic image_editor (image filters/cropping) added. Still needs polish and integration of the UI into the
    rest of the system.
0.7.3.1:
    - big refactor of how statistics are structured in all services, to standardize on gf_stats accross the whole system. 
    all stats are now accessible via gf_analytics only, not in custom dashboards of each of the services themselves
    (this way expensive stats calculations are only run in the gf_analytics service machines, not on machines of other services
    that are expected to efficiently handle real-time requests).
0.7.3:
    - big refactor of error handling, a system wide gf_error code is now used. not all packages migrated yet,
    but core ones (gf_images,gf_crawl_lib,gf_analytics) are.
0.7.2:
    - small improvements in adding crawled_images to flows in the crawler dashboard.
0.7.1:
    - gf_crawl can move images discovered in HTML pages into image flows. This is done via UI's in the crawler
    dashboard at the moment.
0.7.0 (sparkly):
    - first round of GIF viewer/posting features. GIF's are now clickable/playable in the flows_browser, and can
    be added via the gf_chrome_ext.
0.6.0 (robust):
    - migrated system to multi-node cluster, separating key services to their own nodes. new cloud provider (GCP).
0.5.0 (raise):
    - addition of the "gf_crawl" application, which crawls the image-web for images from select sites.
    these images are then accessible in their own crawled/discovered image flow.
    "gf_crawl" functionality is working as a part of the gf_analytics service for now.
0.4.4:
    - addition of the "gf_domains" application. UI for it is publicly accessible, which for now gives basic stats
    an all domains discovered in lings in the GF DB.
0.4.3:
    - addition of the "flows" functionality to the "gf_images" application.
0.4.0:
    - rewrite of the entire backend in Golang, and the front-end in Typescript. 
    major effort, and took over a year (if not longer) to complete, due to real-life coming in the way and 
    tech experimentation happening.
0.2:
    - rewrote all services and most of the front-end in Dart, very few minor improvements added.
0.1:
    - initial Python backend and gf_image/gf_post data-model, using mongodb, most basic css styling.
    first introduction of a Chrome browser extension, for creating posts only (by adding images to them), no image flows.
    (~2010)`,
                    
        // IMPORTANT!! - in CI systems this line is searched for and the git commit hash is pasted in.
        Git_commit_str: GitCommitSHAstr,
	}

	return r
}