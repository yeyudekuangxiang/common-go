package system

import (
	"encoding/json"
	"fmt"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/queue/producer/wxworkpdr"
	"mio/internal/pkg/queue/types/message/wxworkqueue"
	"mio/internal/pkg/util/httputil"
	"mio/pkg/errno"
	glbtyp "mio/pkg/gitlab/types"
	"mio/pkg/wxwork"
	"regexp"
	"strconv"
	"strings"
)

var DefaultGitlabService = GitlabService{}

type GitlabService struct {
}

func NewGitlabService() *GitlabService {
	return &GitlabService{}
}

const private_token = "yoQqAi__rVuZj8kRwgfh"
const base_url = "https://gitlab.miotech.com/api/v4"

var gitlabRefMap = map[string]string{
	"^master$":                  "## 预发布版本",
	"^develop$":                 "## 测试版本",
	"^feature-.+$":              "## 开发版本 [{ref}]",
	"^hotfix-.+$":               "## 热修复版本 [{ref}]",
	"v[0-9]+\\.[0-9]+\\.[0-9]+": "# 正式版本 [{ref}]",
	"^develop-[a-zA-Z]+-(api|rpc)-v[0-9]+\\.[0-9]+\\.[0-9]+$":    "# 开发版本 [{ref}]",
	"^production-[a-zA-Z]+-(api|rpc)-v[0-9]+\\.[0-9]+\\.[0-9]+$": "# 正式版本 [{ref}]",
}

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
		return errno.ErrCommon.WithMessage(mres.MergeStatus)
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
	case "failed":
		return srv.deployFailed(deployment)
	default:
		return srv.deployCancel(deployment)
	}
}
func (srv GitlabService) formatName(ref string) string {
	result := ""
	for regStr, name := range gitlabRefMap {
		reg, err := regexp.Compile(regStr)
		if err != nil {
			panic(err)
		}
		if reg.MatchString(ref) {
			result = name
			break
		}
	}
	return strings.ReplaceAll(result, "{ref}", ref)
}
func (srv GitlabService) deploymentRunning(deployment glbtyp.Deployment) error {

	return wxworkpdr.SendRobotMessage(wxworkqueue.RobotMessage{
		Key:  config.Constants.WxWorkGitlabRobotKey,
		Type: wxwork.MsgTypeMarkdown,
		Message: wxwork.Markdown{
			Content: fmt.Sprintf(`%s 开始发布通知
**应用名称:**[%s](%s)
**应用描述:**%s
**发布分支:**%s
**发布描述:**%s
**发布时间:**%s
**发布人:**%s
**查看发布:**[%d](%s)
`, srv.formatName(deployment.Ref),
				deployment.Project.Name, deployment.Project.WebUrl,
				deployment.Project.Description,
				deployment.Ref,
				deployment.CommitTitle,
				deployment.StatusChangedAt,
				deployment.User.Name,
				deployment.DeployableId, deployment.DeployableUrl),
		},
	})
}
func (srv GitlabService) deploymentSuccess(deployment glbtyp.Deployment) error {
	return wxworkpdr.SendRobotMessage(wxworkqueue.RobotMessage{Key: config.Constants.WxWorkGitlabRobotKey, Type: wxwork.MsgTypeMarkdown, Message: wxwork.Markdown{
		Content: fmt.Sprintf(`%s 发布成功通知
**应用名称:**[%s](%s)
**应用描述:**%s
**发布分支:**%s
**发布描述:**%s
**发布时间:**%s
**发布人:**%s
**查看发布:**[%d](%s)
`, srv.formatName(deployment.Ref),
			deployment.Project.Name, deployment.Project.WebUrl,
			deployment.Project.Description,
			deployment.Ref,
			deployment.CommitTitle,
			deployment.StatusChangedAt,
			deployment.User.Name,
			deployment.DeployableId, deployment.DeployableUrl),
	}})
}
func (srv GitlabService) deployFailed(deployment glbtyp.Deployment) error {
	return wxworkpdr.SendRobotMessage(wxworkqueue.RobotMessage{Key: config.Constants.WxWorkGitlabRobotKey, Type: wxwork.MsgTypeMarkdown, Message: wxwork.Markdown{
		Content: fmt.Sprintf(`%s 发布失败通知
**应用名称:**[%s](%s)
**应用描述:**%s
**发布分支:**%s
**发布描述:**%s
**发布时间:**%s
**发布人:**%s
**查看发布:**[%d](%s)
`, srv.formatName(deployment.Ref),
			deployment.Project.Name, deployment.Project.WebUrl,
			deployment.Project.Description,
			deployment.Ref,
			deployment.CommitTitle,
			deployment.StatusChangedAt,
			deployment.User.Name,
			deployment.DeployableId, deployment.DeployableUrl),
	}})
}
func (srv GitlabService) deployCancel(deployment glbtyp.Deployment) error {
	return wxworkpdr.SendRobotMessage(wxworkqueue.RobotMessage{Key: config.Constants.WxWorkGitlabRobotKey, Type: wxwork.MsgTypeMarkdown, Message: wxwork.Markdown{
		Content: fmt.Sprintf(`%s 发布取消通知
**应用名称:**[%s](%s)
**应用描述:**%s
**发布分支:**%s
**发布描述:**%s
**发布时间:**%s
**发布人:**%s
**查看发布:**[%d](%s)
`, srv.formatName(deployment.Ref),
			deployment.Project.Name, deployment.Project.WebUrl,
			deployment.Project.Description,
			deployment.Ref,
			deployment.CommitTitle,
			deployment.StatusChangedAt,
			deployment.User.Name,
			deployment.DeployableId, deployment.DeployableUrl),
	}})
}
