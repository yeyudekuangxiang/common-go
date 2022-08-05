package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/util/httputil"
	glbtyp "mio/pkg/gitlab/types"
	"mio/pkg/wxwork"
	"strconv"
)

var DefaultGitlabService = GitlabService{}

type GitlabService struct {
}

func NewGitlabService() *GitlabService {
	return &GitlabService{}
}

const private_token = "yoQqAi__rVuZj8kRwgfh"
const base_url = "https://gitlab.miotech.com/api/v4"

// MergeBranch 合并
func (srv GitlabService) MergeBranch(projectId int, source, target string) error {

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
	body, err := httputil.PostJson(url, req)
	if err != nil {
		app.Logger.Error(err)
		return err
	}
	res := struct {
		Iid int `json:"iid"`
	}{}
	//
	err = json.Unmarshal(body, &res)
	if err != nil {
		app.Logger.Error(err)
		return err
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
	mergeBody, err := httputil.PutJson(url2, req2)
	if err != nil {
		app.Logger.Error(err)
		return err
	}

	mres := struct {
		MergeStatus string `json:"merge_status"`
	}{}
	err = json.Unmarshal(mergeBody, &mres)
	if err != nil {
		app.Logger.Error(err)
		return err
	}
	if mres.MergeStatus == "can_be_merged" {
		return nil
	} else {
		app.Logger.Error(mres.MergeStatus)
		return errors.New(mres.MergeStatus)
	}
}
func (srv GitlabService) MergeState(projectId, mergeRequestIId int) (string, error) {
	url := base_url + "/projects/" + strconv.Itoa(projectId) + "/merge_requests/" + strconv.Itoa(mergeRequestIId) + "?private_token=" + private_token

	body, err := httputil.Get(url)
	if err != nil {
		return "", err
	}
	mergeRequest := MergeRequest{}

	//fmt.Println(string(body))
	err = json.Unmarshal(body, &mergeRequest)
	if err != nil {
		return "", err
	}
	return mergeRequest.State, err
}
func (srv GitlabService) Callback(event glbtyp.EventType, body []byte) error {
	switch event {
	case glbtyp.EventTypeDeploymentHook:
		data := glbtyp.Deployment{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		return srv.deployment(data)
	}
	return nil
}
func (srv GitlabService) deployment(deployment glbtyp.Deployment) error {
	switch deployment.Status {
	case "running":
		return srv.deploymentRunning(deployment)
	case "success":
		return srv.deploymentSuccess(deployment)
	default:
		return srv.deployFailed(deployment)
	}
}
func (srv GitlabService) deploymentRunning(deployment glbtyp.Deployment) error {

	return wxwork.SendRobotMessage(config.Constants.WxWorkGitlabRobotKey, wxwork.Markdown{
		Content: fmt.Sprintf(`## 开始发布通知
**应用名称:**[%s](%s)
**应用描述:**%s
**发布版本:**%s
**发布描述:**%s
**发布时间:**%s
**发布人:**%s
**查看发布:**[%d](%s)
`, deployment.Project.Name, deployment.Project.WebUrl,
			deployment.Project.Description,
			deployment.Ref,
			deployment.CommitTitle,
			deployment.StatusChangedAt,
			deployment.User.Name,
			deployment.DeployableId, deployment.DeployableUrl),
	})
}
func (srv GitlabService) deploymentSuccess(deployment glbtyp.Deployment) error {
	return wxwork.SendRobotMessage(config.Constants.WxWorkGitlabRobotKey, wxwork.Markdown{
		Content: fmt.Sprintf(`## 发布成功通知
**应用名称:**[%s](%s)
**应用描述:**%s
**发布版本:**%s
**发布描述:**%s
**发布时间:**%s
**发布人:**%s
**查看发布:**[%d](%s)
`, deployment.Project.Name, deployment.Project.WebUrl,
			deployment.Project.Description,
			deployment.Ref,
			deployment.CommitTitle,
			deployment.StatusChangedAt,
			deployment.User.Name,
			deployment.DeployableId, deployment.DeployableUrl),
	})
}
func (srv GitlabService) deployFailed(deployment glbtyp.Deployment) error {
	return wxwork.SendRobotMessage(config.Constants.WxWorkGitlabRobotKey, wxwork.Markdown{
		Content: fmt.Sprintf(`## 发布失败通知
**应用名称:**[%s](%s)
**应用描述:**%s
**发布版本:**%s
**发布描述:**%s
**发布时间:**%s
**发布人:**%s
**查看发布:**[%d](%s)
`, deployment.Project.Name, deployment.Project.WebUrl,
			deployment.Project.Description,
			deployment.Ref,
			deployment.CommitTitle,
			deployment.StatusChangedAt,
			deployment.User.Name,
			deployment.DeployableId, deployment.DeployableUrl),
	})
}
