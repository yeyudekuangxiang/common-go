package open

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/system"
	"mio/internal/pkg/util"
	"mio/pkg/gitlab"
	glbtyp "mio/pkg/gitlab/types"
)

var DefaultGitlabController = GitlabController{}

type GitlabController struct {
}

func (GitlabController) WebHook(ctx *gin.Context) (gin.H, error) {
	event := ctx.GetHeader("X-Gitlab-Event")
	var ev glbtyp.EventType
	switch event {
	case "Deployment Hook":
		ev = glbtyp.EventTypeDeploymentHook
	default:
		return nil, nil
	}
	body, err := ioutil.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if err != nil {
		return nil, err
	}
	return nil, system.NewGitlabService().Callback(ev, body)
}

// MergeRequestCallback 当创建新的合并请求、更新/合并/关闭现有合并请求或在源分支中添加提交时触发
func (GitlabController) MergeRequestCallback(data []byte) (gin.H, error) {
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
