/*
GloFlow application and media management/publishing platform
Copyright (C) 2024 Ivan Trajkovic

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

package gf_mixpanel

import (
	"context"
	"testing"

	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/mixpanel/mixpanel-go"
	"github.com/stretchr/testify/assert"
)

//-------------------------------------------------------------

func TestEventSend(t *testing.T) {
	
	ctx := context.Background()




	info := &GFmixpanelInfo{
		UsernameStr:    "testuser",
		SecretStr:      "secret",
		ProjectIDstr:   "projectID",
		ProjectTokenStr: "0b9c1c4d50d55f6ae626a7c2cf66bab7", // "projectToken",
	}
	eventTypeStr := "test_event"
	eventMetaMap := map[string]interface{}{
		"key": "value",
	}
	userID := gf_core.GF_ID("userID")
	runtimeSys := &gf_core.RuntimeSys{}

	err := EventSend(eventTypeStr, eventMetaMap, info, userID, ctx, runtimeSys)
	assert.Nil(t, err, "Expected no error from EventSend")
}
