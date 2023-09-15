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

package gf_tagger_core

import (
	"text/template"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------

type GFtemplates struct {
	TagObjects                     *template.Template
	TagObjectsSubtemplatesNamesLst []string

	Bookmarks                     *template.Template
	BookmarksSubtemplatesNamesLst []string
}

//-------------------------------------------------

func TemplatesLoad(pTemplatesPathsMap map[string]string,
	pRuntimeSys *gf_core.RuntimeSys) (*GFtemplates, *gf_core.GFerror) {

	mainTemplateFilepathStr := pTemplatesPathsMap["gf_tag_objects"]
	tagObjectsTmpl, subtemplatesNamesLst, gfErr := gf_core.TemplatesLoad(mainTemplateFilepathStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}



	bookmarksTemplateFilepathStr := pTemplatesPathsMap["gf_bookmarks"]
	bookmarksTmpl, bookmarksSubtemplatesNamesLst, gfErr := gf_core.TemplatesLoad(bookmarksTemplateFilepathStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	gfTemplates := &GFtemplates{
		TagObjects:                     tagObjectsTmpl,
		TagObjectsSubtemplatesNamesLst: subtemplatesNamesLst,

		Bookmarks:                     bookmarksTmpl,
		BookmarksSubtemplatesNamesLst: bookmarksSubtemplatesNamesLst,
	}
	return gfTemplates, nil
}