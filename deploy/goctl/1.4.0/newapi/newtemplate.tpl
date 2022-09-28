type Request {
  Name string `form:"name" json:"name" binding:"required" alias:"姓名"`
}

type Response {
  Message string `json:"message"`
}

@server(
	jwt: JwtAuth
)
service {{.name}}-api {
  @handler {{.handler}}Handler
  get /from/:name(Request) returns (Response)
}
