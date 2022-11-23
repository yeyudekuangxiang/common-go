Name: {{.serviceName}}
Host: {{.host}}
Port: {{.port}}
Debug: false
Timeout: 10000
#rpc service
RpcConf:
  Endpoints:
    - 127.0.0.1:2006
  NonBlock: true
  Timeout: 5000
#全局客户端配置
GlobalClientConf:
  # 阿里云日志
  AliYunSlsConf:
    Endpoint:
    AccessKeyID:
    AccessKeySecret:
    Project:
    LogStore:
  # 定时延时任务
  Asynq:
    Addr: 10.2.22.171:33469
    Password: G62m50oigInC30sf
    DB: 2
    PoolSize: 50
  # redis客户端
  RedisConf:
    Host: 10.2.22.171:33469
    Password: G62m50oigInC30sf
    DB: 0
  # 诸葛打点客户端
  ZhugeConf:
    AppKey:
    AppSecret:
JwtAuth:
  AccessSecret: