package gf_github

import (
	"fmt"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/davecgh/go-spew/spew"
)

//--------------------------------------------------------------------
// GITHUB_ACTIONS_RUN_WORKFLOW

// Run a Github Actions workflow on a target repository and branch.
// Workflow must be marked in github_actions to be runnable by "workflow_dispatch" events
// (in the "on" section)
// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#workflow_dispatch
// pWorkflowIDorFileNameStr - is either a workflow ID or a workflow definition file (*.yaml)
func ActionsRunWorkflow(pRepoOwnerAndNameStr string,
	pWorkflowIDorFileNameStr string,
	pBranchNameStr           string,
	pGithubBearerTokenStr    string,
	pRuntimeSys              *gf_core.RuntimeSys) {

	// https://docs.github.com/en/rest/actions/workflows#create-a-workflow-dispatch-event
	urlStr := fmt.Sprintf("https://api.github.com/repos/%s/actions/workflows/%s/dispatches",
		pRepoOwnerAndNameStr,
		pWorkflowIDorFileNameStr)

	dataMap := map[string]interface{}{
		"ref": pBranchNameStr,
	}
	dataLst, _  := json.Marshal(dataMap)

	_, body, errs := gorequest.New().
		Post(urlStr).
		Set("accept", "application/vnd.github+json").
		Set("authorization", fmt.Sprintf("Bearer %s", pGithubBearerTokenStr)).
		Send(string(dataLst)).
		End()
	if len(errs) > 0 {
		err   := errs[0]
		gfErr := gf_core.ErrorCreate("failed to run a workflow in github actions via REST API",
			"http_client_req_error",
			map[string]interface{}{
				"repo_owner_and_name_str":      pRepoOwnerAndNameStr,
				"workflow_id_or_file_name_str": pWorkflowIDorFileNameStr,
				"branch_name_str":              pBranchNameStr,
				"url_str":                      urlStr,
			},
			err, "gf_ops_lib", pRuntimeSys)
		return "", nil, gfErr
	}

	rMap := map[string]interface{}{}
	err  := json.Unmarshal([]byte(body), &rMap)
	if err != nil {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("failed to parse json response from gf_images_client start_job HTTP REST API - %s", url_str), 
			"json_unmarshal_error",
			map[string]interface{}{
				"url_str": url_str,
				"body":    body,
			},
			j_err, "gf_images_lib", pRuntimeSys)
		return "", nil, gfErr
	}
}

//--------------------------------------------------------------------
// Get IP's from which github servers are expected to send requests.
func GetIPs(pRuntimeSys *gf_core.RuntimeSys) ([]string, *gf_core.GFerror) {

	urlStr := fmt.Sprintf("https://api.github.com/meta")

	_, body, errs := gorequest.New().
		Get(urlStr).
		Set("accept", "application/json").
		End()

	if errs != nil {
		err    := errs[0]
		gfErr := gf_core.ErrorCreate("github meta HTTP REST API request failed",
			"http_client_req_error",
			map[string]interface{}{
				"url_str": urlStr,
			},
			err, "gf_ops_lib", pRuntimeSys)
		return nil, gfErr
	}

	rMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(body), &rMap)
	if err != nil {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("failed to parse json response from github HTTP REST API"), 
			"json_unmarshal_error",
			map[string]interface{}{
				"url_str": urlStr,
				"body":    body,
			},
			err, "gf_ops_lib", pRuntimeSys)
		return nil, gfErr
	}

	// spew.Dump(rMap)

	// ADD!! - return IP's of other services as well, not just Github Actions
	githubActionsIPsUncastedLst := rMap["actions"].([]interface{})
	githubActionsIPsLst := []string{}
	for _, ip := range githubActionsIPsUncastedLst {
		githubActionsIPsLst = append(githubActionsIPsLst, ip.(string))
	}

	return githubActionsIPsLst, nil
}