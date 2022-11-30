Name: {{.serviceName}}.rpc
Debug: true
ListenOn: 127.0.0.1:2006
Timeout: 10000
Database:
  # mysql 、postgres
  Type: postgres
  Host: 127.0.0.1
  UserName: postgres
  Password: postgres
  Database: postgres
  Port: 5432
  TablePrefix:
  #最大连接数 <=0表示不限制连接数
  MaxOpenConns: 0
  #最大空闲数 <=0表示不保留空闲连接
  MaxIdleConns: 0
  #连接可重用时间 <=0表示永远可用(单位秒)
  MaxLifetime: 0
Cache:
  - Host: 127.0.0.1:6379
    Pass: redispass

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