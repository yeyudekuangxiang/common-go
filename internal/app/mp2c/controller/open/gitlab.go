package open

import (
	"github.com/gin-gonic/gin"
	glbtyp "gitlab.miotech.com/miotech-application/backend/common-go/gitlab/types"
	"io/ioutil"
	"mio/internal/pkg/service/system"
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
