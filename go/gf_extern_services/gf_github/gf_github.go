package gf_github

import (
	"fmt"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_core"
	// "github.com/davecgh/go-spew/spew"
)

//--------------------------------------------------------------------

type GFissue struct {
	TitleStr  string
	BodyStr   string
	UrlStr    string
	NumberInt int
	StateStr  string
	Labels    []GFissueLabel
	MilestoneTitleStr string
	MilestoneUrlStr   string
}

type GFissueLabel struct {
	NameStr     string
	ColorHexStr string
}

//--------------------------------------------------------------------
// GET_ISSUES

// get a github issue associated with a particular repository
func GetIssues(pRepoOwnerAndNameStr string,
	pGithubBearerTokenStr string,
	pRuntimeSys           *gf_core.RuntimeSys) ([]GFissue, *gf_core.GFerror) {

	// https://docs.github.com/en/rest/issues/issues#list-repository-issues
	urlStr := fmt.Sprintf("https://api.github.com/repos/%s/issues", pRepoOwnerAndNameStr)

	_, body, errs := gorequest.New().
		Get(urlStr).
		Set("accept", "application/vnd.github+json").
		Set("authorization", fmt.Sprintf("Bearer %s", pGithubBearerTokenStr)).
		// Send(string(dataLst)).
		End()
	if len(errs) > 0 {
		err   := errs[0]
		gfErr := gf_core.ErrorCreate("failed to get repository issues in github via REST API",
			"http_client_req_error",
			map[string]interface{}{
				"repo_owner_and_name_str": pRepoOwnerAndNameStr,
				"url_str":                 urlStr,
			},
			err, "gf_github", pRuntimeSys)
		return nil, gfErr
	}

	// fmt.Println(body)

	rLst := []interface{}{}
	err := json.Unmarshal([]byte(body), &rLst)
	if err != nil {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("failed to parse json response from github HTTP REST API"), 
			"json_unmarshal_error",
			map[string]interface{}{
				"url_str": urlStr,
				"body":    body,
			},
			err, "gf_github", pRuntimeSys)
		return nil, gfErr
	}

	// spew.Dump(rLst)

	gfIssuesLst := []GFissue{}
	for _, issue := range rLst {

		issueMap := issue.(map[string]interface{})
		urlStr := issueMap["url"].(string)
		numberInt := int(issueMap["number"].(float64))
		stateStr := issueMap["state"].(string)
		titleStr := issueMap["title"].(string)
		bodyStr := issueMap["body"].(string)

		gfIssueLabelsLst := []GFissueLabel{}
		for _, label := range issueMap["labels"].([]interface{}) {
			labelMap      := label.(map[string]interface{})
			labelNameStr  := labelMap["name"].(string)
			labelColorStr := labelMap["color"].(string)

			gfIssueLabel := GFissueLabel{
				NameStr:     labelNameStr,
				ColorHexStr: labelColorStr,
			}
			gfIssueLabelsLst = append(gfIssueLabelsLst, gfIssueLabel)
		}

		

		gfIssue := GFissue{
			TitleStr:  titleStr,
			BodyStr:   bodyStr,
			UrlStr:    urlStr,
			NumberInt: numberInt,
			StateStr:  stateStr,
			Labels:    gfIssueLabelsLst,
		}

		if issueMap["milestone"] != nil {
			milestoneMap      := issueMap["milestone"].(map[string]interface{})
			milestoneTitleStr := milestoneMap["title"].(string)
			milestoneUrlStr   := milestoneMap["url"].(string)

			gfIssue.MilestoneTitleStr = milestoneTitleStr
			gfIssue.MilestoneUrlStr   = milestoneUrlStr
		}

		gfIssuesLst = append(gfIssuesLst, gfIssue)
	}

	return gfIssuesLst, nil
}

//--------------------------------------------------------------------
// RUN_ACTIONS_WORKFLOW

// Run a Github Actions workflow on a target repository and branch.
// Workflow must be marked in github_actions to be runnable by "workflow_dispatch" events
// (in the "on" section)
// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#workflow_dispatch
// pWorkflowIDorFileNameStr - is either a workflow ID or a workflow definition file (*.yaml)
func RunActionsWorkflow(pRepoOwnerAndNameStr string,
	pWorkflowIDorFileNameStr string,
	pBranchNameStr           string,
	pGithubBearerTokenStr    string,
	pRuntimeSys              *gf_core.RuntimeSys) *gf_core.GFerror {

	// https://docs.github.com/en/rest/actions/workflows#create-a-workflow-dispatch-event
	urlStr := fmt.Sprintf("https://api.github.com/repos/%s/actions/workflows/%s/dispatches",
		pRepoOwnerAndNameStr,
		pWorkflowIDorFileNameStr)

	dataMap := map[string]interface{}{
		"ref": pBranchNameStr,
	}
	dataLst, _  := json.Marshal(dataMap)

	_, _, errs := gorequest.New().
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
			err, "gf_github", pRuntimeSys)
		return gfErr
	}

	return nil
}

//--------------------------------------------------------------------
// GET_IPS

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
			err, "gf_github", pRuntimeSys)
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
			err, "gf_github", pRuntimeSys)
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