Name: {{.serviceName}}
Host: {{.host}}
Port: {{.port}}
Debug: false
#rpc service
RpcConf:
  Endpoints:
    - 127.0.0.1:2006
  NonBlock: true
JwtAuth:
  AccessSecret: