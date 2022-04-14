package system

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/system"
	"mio/internal/pkg/util"
	"mio/pkg/gitlab"
	"strings"
)

var DefaultGitlabController = GitlabController{}

type GitlabController struct {
}

func (ctr GitlabController) Callback(ctx *gin.Context) (gin.H, error) {
	gitlabEvent := ctx.GetHeader("X-Gitlab-Event")
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	app.Logger.Infof("event:%s body:%s", gitlabEvent, string(body))
	_ = ioutil.WriteFile("./runtime/"+strings.ReplaceAll(gitlabEvent, " ", "")+".json", body, 0777)
	switch gitlabEvent {
	case "Merge Request Hook":
		return ctr.MergeRequestCallback(body)
	}

	return nil, nil
}

// MergeRequestCallback 当创建新的合并请求、更新/合并/关闭现有合并请求或在源分支中添加提交时触发
func (ctr GitlabController) MergeRequestCallback(data []byte) (gin.H, error) {
	form := gitlab.GitlabWebHookForm{}
	if err := json.Unmarshal(data, &form); err != nil {
		return nil, err
	}

	//合并完成
	if form.Attributes.Action == "merge" {
		//自动合并到develop
		if form.Attributes.TargetBranch == "master" && util.ArrayIsContains([]string{"hotfix-", "release-"}, form.Attributes.SourceBranch) {
			app.Logger.Infof("执行自动合并")
			_ = system.DefaultGitlabService.MergeBranch(form.Project.Id, form.Attributes.SourceBranch, "develop")
		}
	}
	return nil, nil
}
