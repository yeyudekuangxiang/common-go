package system

import (
	"encoding/json"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util"
	"strconv"
)

var DefaultGitlabService = GitlabService{}

type GitlabService struct {
}

const private_token = "yoQqAi__rVuZj8kRwgfh"
const base_url = "https://gitlab.miotech.com/api/v4"

// MergeBranch 合并
func (srv GitlabService) MergeBranch(projectId int, source, target string) bool {

	//创建合并
	url := base_url + "/projects/" + strconv.Itoa(projectId) + "/merge_requests?private_token=" + private_token
	req := struct {
		Id            int    `json:"id"`
		Source_branch string `json:"source_branch"`
		Target_branch string `json:"target_branch"`
		Title         string `json:"title"`
	}{
		projectId,
		source,
		target,
		"系统自动合并：" + source + " -> " + target,
	}
	body, err := util.DefaultHttp.PostJson(url, &req)
	if err != nil {
		app.Logger.Error(err)
		return false
	}
	res := struct {
		Iid int `json:"iid"`
	}{}
	//
	err = json.Unmarshal(body, &res)
	if err != nil {
		app.Logger.Error(err)
		return false
	}

	//time.Sleep(5 * time.Second)
	//合并通过
	req2 := struct {
		Id  int `json:"id"`
		Iid int `json:"merge_request_iid"`
	}{
		projectId,
		res.Iid,
	}
	url2 := base_url + "/projects/" + strconv.Itoa(projectId) + "/merge_requests/" + strconv.Itoa(res.Iid) + "/merge?private_token=" + private_token
	mergeBody, err := util.DefaultHttp.PutJson(url2, req2)
	if err != nil {
		app.Logger.Error(err)
		return false
	}

	mres := struct {
		MergeStatus string `json:"merge_status"`
	}{}
	err = json.Unmarshal(mergeBody, &mres)
	if err != nil {
		app.Logger.Error(err)
		return false
	}
	if mres.MergeStatus == "can_be_merged" {
		return true
	} else {
		app.Logger.Error(mres.MergeStatus)
		return false
	}
}
