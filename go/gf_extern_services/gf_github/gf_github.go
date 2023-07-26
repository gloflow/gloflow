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
	LabelsLst []GFissueLabel
	MilestoneTitleStr string
	MilestoneUrlStr   string
	CreatedAtStr      string
}

type GFissueLabel struct {
	NameStr     string
	ColorHexStr string
}

type GFproject struct {
	NameStr string
	UrlStr string
}

type GFgithubProject struct {
	TitleStr     string
	URLstr       string
	GraphqlIDstr string
}

//-------------------------------------------------

// get basic info on github projects belonging to an organization, using the Github GraphQL API.
func GetProjects(pOrgNameStr string,
	pGithubBearerTokenStr string,
	pRuntimeSys           *gf_core.RuntimeSys) ([]GFgithubProject, *gf_core.GFerror) {

	graphQLqueryStr := fmt.Sprintf(`query {
		organization(login: "%s") {
		  	projectsV2(first: 40) {
				nodes {
			  		id
			  		title
			  		url
				}
		  	}
		}
	}`, pOrgNameStr)


	rMap, gfErr := RunGraphQLquery(graphQLqueryStr, pGithubBearerTokenStr, pRuntimeSys)
	if gfErr != nil {
		return nil, gfErr
	}

	// spew.Dump(rMap)


	projectsLst := rMap["data"].(map[string]interface{})["organization"].(map[string]interface{})["projectsV2"].(map[string]interface{})["nodes"].([]interface{})
	// spew.Dump(projectsLst)

	gfProjectsLst := []GFgithubProject{}
	for _, p := range projectsLst {
		projectMap := p.(map[string]interface{})
		projectTitleStr := projectMap["title"].(string)
		projectURLstr   := projectMap["url"].(string)
		projectGraphqlIDstr := projectMap["id"].(string)

		project := GFgithubProject{
			TitleStr: projectTitleStr,
			URLstr:   projectURLstr,
			GraphqlIDstr: projectGraphqlIDstr,
		}

		gfProjectsLst = append(gfProjectsLst, project)
	}

	return gfProjectsLst, nil
}

//--------------------------------------------------------------------

func RunGraphQLquery(pGraphQLqueryStr string,
	pGithubBearerTokenStr string,
	pRuntimeSys           *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {

	urlStr := fmt.Sprintf("https://api.github.com/graphql")
	dataMap := map[string]interface{}{
		"query": pGraphQLqueryStr,
	}
	dataLst, _ := json.Marshal(dataMap)

	_, body, errs := gorequest.New().
		Post(urlStr).
		Set("accept", "application/vnd.github+json").
		Set("authorization", fmt.Sprintf("Bearer %s", pGithubBearerTokenStr)).
		Send(string(dataLst)).
		End()
	if len(errs) > 0 {
		err   := errs[0]
		gfErr := gf_core.ErrorCreate("failed to get github projects via GraphQL",
			"http_client_req_error",
			map[string]interface{}{
				"url_str": urlStr,
			},
			err, "gf_project", pRuntimeSys)
		return nil, gfErr
	}

	rMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(body), &rMap)
	if err != nil {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("failed to parse json response from github HTTP GraphQL API"), 
			"json_decode_error",
			map[string]interface{}{
				"url_str": urlStr,
				"body":    body,
			},
			err, "gf_github", pRuntimeSys)
		return nil, gfErr
	}

	return rMap, nil
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

	fmt.Println(body)

	rLst := []interface{}{}
	err := json.Unmarshal([]byte(body), &rLst)
	if err != nil {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("failed to parse json response from github HTTP REST API"), 
			"json_decode_error",
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
		urlStr := issueMap["html_url"].(string)
		numberInt := int(issueMap["number"].(float64))
		stateStr := issueMap["state"].(string)
		titleStr := issueMap["title"].(string)
		bodyStr := issueMap["body"].(string)
		createdAtStr := issueMap["created_at"].(string)

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
			LabelsLst: gfIssueLabelsLst,
			CreatedAtStr: createdAtStr,
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
func GetIPs(pRuntimeSys *gf_core.RuntimeSys) ([]string, []string, *gf_core.GFerror) {

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
		return nil, nil, gfErr
	}

	rMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(body), &rMap)
	if err != nil {
		gfErr := gf_core.ErrorCreate(fmt.Sprintf("failed to parse json response from github HTTP REST API"), 
			"json_decode_error",
			map[string]interface{}{
				"url_str": urlStr,
				"body":    body,
			},
			err, "gf_github", pRuntimeSys)
		return nil, nil, gfErr
	}

	// spew.Dump(rMap)

	// ADD!! - return IP's of other services as well, not just Github Actions
	githubActionsIPsUncastedLst := rMap["actions"].([]interface{})
	githubWebHooksIPsUncastedLst := rMap["hooks"].([]interface{})
	
	githubActionsIPsLst := []string{}
	githubWebHooksIPsLst := []string{}
	for _, ip := range githubActionsIPsUncastedLst {
		githubActionsIPsLst = append(githubActionsIPsLst, ip.(string))
	}
	for _, ip := range githubWebHooksIPsUncastedLst {
		githubWebHooksIPsLst = append(githubWebHooksIPsLst, ip.(string))
	}

	return githubActionsIPsLst, githubWebHooksIPsLst, nil
}