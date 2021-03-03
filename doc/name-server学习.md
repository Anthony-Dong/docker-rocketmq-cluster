## [1、基本概念](https://github.com/apache/rocketmq/blob/release-4.8.0/docs/cn/architecture.md#1-%E6%8A%80%E6%9C%AF%E6%9E%B6%E6%9E%84)

![img](https://tyut.oss-accelerate.aliyuncs.com/image/2021/3-3/f9199d8be0aa4ec3ab7db8f0f6baadf9.png)

NameServer是一个非常简单的Topic路由注册中心，其角色类似Dubbo中的zookeeper，支持Broker的动态注册与发现。主要包括两个功能：Broker管理，NameServer接受Broker集群的注册信息并且保存下来作为路由信息的基本数据。然后提供心跳检测机制，检查Broker是否还存活；路由信息管理，每个NameServer将保存关于Broker集群的整个路由信息和用于客户端查询的队列信息。然后Producer和Conumser通过NameServer就可以知道整个Broker集群的路由信息，从而进行消息的投递和消费。NameServer通常也是集群的方式部署，各实例间相互不进行信息通讯。Broker是向每一台NameServer注册自己的路由信息，所以每一个NameServer实例上面都保存一份完整的路由信息。当某个NameServer因某种原因下线了，Broker仍然可以向其它NameServer同步其路由信息，Producer,Consumer仍然可以动态感知Broker的路由的信息。

## 2、分布式

NameServer是一个几乎无状态节点，可集群部署，节点之间无任何信息同步。

## 3、关系

### 1、与producer和consumer的关系

- p和c都需要拿到 `NameServer` 来获取集群信息，作为客户端连接的参数

- Producer与NameServer集群中的其中一个节点（随机选择）建立长连接，定期从NameServer获取Topic路由信息，并向提供Topic 服务的Master建立长连接，且定时向Master发送心跳。Producer完全无状态，可集群部署。
- Consumer与NameServer集群中的其中一个节点（随机选择）建立长连接，定期从NameServer获取Topic路由信息，并向提供Topic服务的Master、Slave建立长连接，且定时向Master、Slave发送心跳。Consumer既可以从Master订阅消息，也可以从Slave订阅消息，消费者在向Master拉取消息时，Master服务器会根据拉取偏移量与最大偏移量的距离（判断是否读老消息，产生读I/O），以及从服务器是否可读等因素建议下一次是从Master还是Slave拉取。

### 2、与Broker的关系

- 启动NameServer，NameServer起来后监听端口，等待Broker、Producer、Consumer连上来，相当于一个路由控制中心。
- Broker启动，跟所有的NameServer保持长连接，定时发送心跳包。心跳包中包含当前Broker信息(IP+端口等)以及存储所有Topic信息。注册成功后，NameServer集群中就有Topic跟Broker的映射关系。
- Broker部署相对复杂，Broker分为Master与Slave，一个Master可以对应多个Slave，但是一个Slave只能对应一个Master，Master与Slave 的对应关系通过指定相同的BrokerName，不同的BrokerId 来定义，BrokerId为0表示Master，非0表示Slave。Master也可以部署多个。每个Broker与NameServer集群中的所有节点建立长连接，定时注册Topic信息到所有NameServer。 注意：当前RocketMQ版本在部署架构上支持一Master多Slave，但只有BrokerId=1的从服务器才会参与消息的读负载。

## 4、配置

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