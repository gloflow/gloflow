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

package gf_rl_apps

import (
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_ml_lib/gf_rl"
)

//-------------------------------------------------
func Init(pRuntimeSys *gf_core.Runtime_sys) {


	// ACTIONS_DEFS
	// actions defined for the image picking application of RL
	agentActionsDefsLst := []gf_rl.GFrlActionDef{
		// crawler fetches a new link on the same domain
		{
			NameStr:      "open_page_link_on_same_domain",
			ProbabilityF: 0.3,
		},

		// crawler fetches a new link on a different domain
		{
			NameStr:      "open_page_link_on_different_domain",
			ProbabilityF: 0.3,
		},

		// crawler opens the previous page (goes back)
		{
			NameStr:      "rewind_to_previous_page",
			ProbabilityF: 0.3,
		},
	}

	//-------------------------------------------------
	// REWARD_FUNCTION
	rewardFun := func(pState *gf_rl.GFstate, pActionStr gf_rl.GFaction) int {

		
		return 0
	}

	//-------------------------------------------------

	rlHyperparams := gf_rl.GFrlHyperparams{

		// probability (0.0-1.0) for epsilon-greedy strategy
		// (exploration/exploatation action probability)
		EpsilonF: 0.5,

		AlphaLearningRateF: 0.1,
	}

	appInfo := &gf_rl.GFappInfo {
		NameStr:        "gf_crawl_agent",
		ActionsDefsLst: agentActionsDefsLst,
		RewardFun:      rewardFun,
		Hyperparams:    rlHyperparams,
	}

	gf_rl.Init(appInfo,
		pRuntimeSys)
}