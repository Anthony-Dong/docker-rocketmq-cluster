可以看name-server机器的启动日志`/root/logs/rocketmqlogs/namesrv.log`，这个是没有修改任何配置的启动参数

```shell
2021-02-05 14:34:33 INFO main - rocketmqHome=/opt/rocketmq
2021-02-05 14:34:33 INFO main - kvConfigPath=/root/namesrv/kvConfig.json
2021-02-05 14:34:33 INFO main - configStorePath=/root/namesrv/namesrv.properties
2021-02-05 14:34:33 INFO main - productEnvName=center
2021-02-05 14:34:33 INFO main - clusterTest=false
2021-02-05 14:34:33 INFO main - orderMessageEnable=false
2021-02-05 14:34:33 INFO main - listenPort=9876
2021-02-05 14:34:33 INFO main - serverWorkerThreads=8
2021-02-05 14:34:33 INFO main - serverCallbackExecutorThreads=0
2021-02-05 14:34:33 INFO main - serverSelectorThreads=3
2021-02-05 14:34:33 INFO main - serverOnewaySemaphoreValue=256
2021-02-05 14:34:33 INFO main - serverAsyncSemaphoreValue=64
2021-02-05 14:34:33 INFO main - serverChannelMaxIdleTimeSeconds=120
2021-02-05 14:34:33 INFO main - serverSocketSndBufSize=65535
2021-02-05 14:34:33 INFO main - serverSocketRcvBufSize=65535
2021-02-05 14:34:33 INFO main - serverPooledByteBufAllocatorEnable=true
2021-02-05 14:34:33 INFO main - useEpollNativeSelector=false
2021-02-05 14:34:34 INFO main - Server is running in TLS permissive mode
2021-02-05 14:34:34 INFO main - Tls config file doesn't exist, skip it
2021-02-05 14:34:34 INFO main - Log the final used tls related configuration
2021-02-05 14:34:34 INFO main - tls.test.mode.enable = true
2021-02-05 14:34:34 INFO main - tls.server.need.client.auth = none
2021-02-05 14:34:34 INFO main - tls.server.keyPath = null
2021-02-05 14:34:34 INFO main - tls.server.keyPassword = null
2021-02-05 14:34:34 INFO main - tls.server.certPath = null
2021-02-05 14:34:34 INFO main - tls.server.authClient = false
2021-02-05 14:34:34 INFO main - tls.server.trustCertPath = null
2021-02-05 14:34:34 INFO main - tls.client.keyPath = null
2021-02-05 14:34:34 INFO main - tls.client.keyPassword = null
2021-02-05 14:34:34 INFO main - tls.client.certPath = null
2021-02-05 14:34:34 INFO main - tls.client.authServer = false
2021-02-05 14:34:34 INFO main - tls.client.trustCertPath = null
```

