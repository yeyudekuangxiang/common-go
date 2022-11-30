Name: {{.serviceName}}
Host: {{.host}}
Port: {{.port}}
Debug: true
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
    Addr:
    Password:
    DB:
    PoolSize: 50
  # redis客户端
  RedisConf:
    Host:
    Password:
    DB:
  # 诸葛埋点客户端
  ZhugeConf:
    AppKey:
    AppSecret:
  # rabbitmq
  Rabbitmq:
    Url:
JwtAuth:
  AccessSecret: