package types

// Deployment
/*
{
    "object_kind": "deployment",
    "status": "success",
    "status_changed_at": "2022-08-05 14:25:15 +0800",
    "deployment_id": 5200,
    "deployable_id": 302172,
    "deployable_url": "https://github.com/yeyudekuangxiang/mp2c-go/-/jobs/302172",
    "environment": "hotfix-tagsort",
    "project": {},
    "short_sha": "91b829b3",
    "user": {},
    "user_url": "https://gitlab.miotech.com/neiljin",
    "commit_url": "https://github.com/yeyudekuangxiang/mp2c-go/-/commit/91b829b3e42ddbb1ce1e213231ae353272c72d81",
    "commit_title": "fix:修复tag排序问题",
    "ref": "hotfix-tagsort"
}
*/
type Deployment struct {
	//类型  deployment部署
	ObjectKind string `json:"object_kind"`
	//状态 success成功 running运行中
	Status string `json:"status"`
	//状态改变时间 2022-08-05 14:25:15 +0800
	StatusChangedAt string `json:"status_changed_at"`
	//deployment_id 5200
	DeploymentId int `json:"deployment_id"`
	//deployable_id 部署任务的编号 302172
	DeployableId int `json:"deployable_id"`
	//部署任务url https://github.com/yeyudekuangxiang/mp2c-go/-/jobs/302172
	DeployableUrl string `json:"deployable_url"`
	//环境名称 hotfix-tagsort
	Environment string `json:"environment"`
	//项目信息
	Project Project `json:"project"`
	//部署时commit id
	ShortSha string `json:"short_sha"`
	//用户信息
	User User `json:"user"`
	//用户主页信息
	UserUrl string `json:"user_url"`
	//提交信息地址
	CommitUrl string `json:"commit_url"`
	//提交标题信息
	CommitTitle string `json:"commit_title"`
	//部署的分支或者标签
	Ref string `json:"ref"`
}
