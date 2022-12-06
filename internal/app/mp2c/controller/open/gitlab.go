package open

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mio/internal/pkg/service/system"
	glbtyp "mio/pkg/gitlab/types"
)

var DefaultGitlabController = GitlabController{}

type GitlabController struct {
}

func (GitlabController) WebHook(ctx *gin.Context) (gin.H, error) {
	event := ctx.GetHeader("X-Gitlab-Event")
	var ev glbtyp.EventHookType
	switch event {
	case "Deployment Hook":
		ev = glbtyp.EventDeploymentHook
	case "Merge Request Hook":
		ev = glbtyp.EventMergeRequestHook
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
