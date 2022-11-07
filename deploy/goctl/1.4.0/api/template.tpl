syntax = "v1"

info (
	title: // TODO: add title
	desc: // TODO: add description
	author: "{{.gitUser}}"
	email: "{{.gitEmail}}"
)

type PingReq {

}

type PingResp {
    Pong string `json:"pong"`
}

service {{.serviceName}} {
	@handler Ping
	get /ping(PingReq) returns(PingResp)
}
