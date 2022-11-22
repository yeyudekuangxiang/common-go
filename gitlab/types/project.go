package types

// Project
/*
{
        "id": 294,
        "name": "mp2c-go",
        "description": "绿喵后端服务 - go",
        "web_url": "https://gitlab.miotech.com/miotech-application/backend/mp2c-go",
        "avatar_url": null,
        "git_ssh_url": "git@gitlab.miotech.com:miotech-application/backend/mp2c-go.git",
        "git_http_url": "https://gitlab.miotech.com/miotech-application/backend/mp2c-go.git",
        "namespace": "backend",
        "visibility_level": 0,
        "path_with_namespace": "miotech-application/backend/mp2c-go",
        "default_branch": "master",
        "ci_config_path": null,
        "homepage": "https://gitlab.miotech.com/miotech-application/backend/mp2c-go",
        "url": "git@gitlab.miotech.com:miotech-application/backend/mp2c-go.git",
        "ssh_url": "git@gitlab.miotech.com:miotech-application/backend/mp2c-go.git",
        "http_url": "https://gitlab.miotech.com/miotech-application/backend/mp2c-go.git"
    }
*/
type Project struct {
	//项目id
	Id int `json:"id"`
	//项目名称
	Name string `json:"name"`
	//项目描述
	Description string `json:"description"`
	//项目地址
	WebUrl string `json:"web_url"`
	//项目头像?
	AvatarUrl interface{} `json:"avatar_url"`
	//git ssh地址
	GitSshUrl string `json:"git_ssh_url"`
	//git http地址
	GitHttpUrl string `json:"git_http_url"`
	//所属分组
	Namespace       string `json:"namespace"`
	VisibilityLevel int    `json:"visibility_level"`
	//完整路径
	PathWithNamespace string `json:"path_with_namespace"`
	//默认分支
	DefaultBranch string      `json:"default_branch"`
	CiConfigPath  interface{} `json:"ci_config_path"`
	Homepage      string      `json:"homepage"`
	Url           string      `json:"url"`
	SshUrl        string      `json:"ssh_url"`
	HttpUrl       string      `json:"http_url"`
}
