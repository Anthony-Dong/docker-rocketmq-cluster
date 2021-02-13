## 1、集群

### 1、6node-2broker的rocket-mq集群

> ​	由于`DLeger`集群搭建要求每个`broker`至少三节点

1、docker节点

```shell
➜  ~ docker ps
CONTAINER ID        IMAGE                                   COMMAND                  CREATED             STATUS              PORTS                                           NAMES
8f25ff7b82f8        apacherocketmq/rocketmq-console:2.0.0   "sh -c 'java $JAVA_O…"   About an hour ago   Up About an hour    0.0.0.0:8080->8080/tcp                          docker_rocketmq-console_1
9556a05f85be        rocketmq:v4.8.0                         "mqbroker -c conf/br…"   About an hour ago   Up About an hour    9876/tcp, 10911/tcp, 0.0.0.0:10914->10914/tcp   docker_broker-04_1
288c93824863        rocketmq:v4.8.0                         "mqbroker -c conf/br…"   About an hour ago   Up About an hour    9876/tcp, 10911/tcp, 0.0.0.0:10913->10913/tcp   docker_broker-03_1
f30df1853117        rocketmq:v4.8.0                         "mqbroker -c conf/br…"   About an hour ago   Up About an hour    9876/tcp, 10911/tcp, 0.0.0.0:10915->10915/tcp   docker_broker-05_1
4cb52312a08b        rocketmq:v4.8.0                         "mqbroker -c conf/br…"   About an hour ago   Up About an hour    9876/tcp, 10911/tcp, 0.0.0.0:10916->10916/tcp   docker_broker-06_1
0bc0c79ca379        rocketmq:v4.8.0                         "mqbroker -c conf/br…"   About an hour ago   Up About an hour    9876/tcp, 0.0.0.0:10911->10911/tcp              docker_broker-01_1
79159b7b19fc        rocketmq:v4.8.0                         "mqnamesrv"              About an hour ago   Up About an hour    10911/tcp, 0.0.0.0:9871->9876/tcp               docker_nameserver-01_1
b89376852c3c        rocketmq:v4.8.0                         "mqbroker -c conf/br…"   About an hour ago   Up About an hour    9876/tcp, 10911/tcp, 0.0.0.0:10912->10912/tcp   docker_broker-02_1
32f3d9bd8027        rocketmq:v4.8.0                         "mqnamesrv"              About an hour ago   Up About an hour    10911/tcp, 0.0.0.0:9872->9876/tcp               docker_nameserver-02_1
```

2、机器信息

> ​	rocket-mq 的broker使用了大量的线程进行`空转(for + sleep)`，导致CPU占用率偏高。

```shell
➜  ~ docker stats --no-stream
CONTAINER ID        NAME                        CPU %               MEM USAGE / LIMIT   MEM %               NET I/O             BLOCK I/O           PIDS
8f25ff7b82f8        docker_rocketmq-console_1   0.20%               247.2MiB / 256MiB   96.55%              5.16MB / 3.54MB     106MB / 112MB       36
9556a05f85be        docker_broker-04_1          31.62%              401.7MiB / 512MiB   78.45%              11.4MB / 5.75MB     0B / 8.95MB         159
288c93824863        docker_broker-03_1          33.67%              421.7MiB / 512MiB   82.35%              17.1MB / 29.2MB     0B / 8.95MB         223
f30df1853117        docker_broker-05_1          31.21%              405.7MiB / 512MiB   79.23%              11.2MB / 4.37MB     0B / 8.95MB         161
4cb52312a08b        docker_broker-06_1          35.09%              419.7MiB / 512MiB   81.98%              17.8MB / 30.2MB     430kB / 8.95MB      222
0bc0c79ca379        docker_broker-01_1          31.59%              410MiB / 512MiB     80.07%              11.1MB / 5.69MB     28.7kB / 9.11MB     162
79159b7b19fc        docker_nameserver-01_1      0.42%               183.2MiB / 256MiB   71.56%              2.34MB / 1.23MB     0B / 8.96MB         41
b89376852c3c        docker_broker-02_1          31.09%              405.4MiB / 512MiB   79.17%              10.9MB / 4.27MB     0B / 8.95MB         161
32f3d9bd8027        docker_nameserver-02_1      0.22%               193.6MiB / 256MiB   75.61%              2.45MB / 1.52MB     0B / 8.96MB         41
```

![image-20210213141305307](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-13/907eaa7a9c1c479cbd61225c1e956db2.png)

3、集群节点

```shell
root@0bc0c79ca379:/opt/rocketmq# mqadmin clusterList -n "nameserver-01:9876;nameserver-02:9876"
RocketMQLog:WARN No appenders could be found for logger (io.netty.util.internal.PlatformDependent0).
RocketMQLog:WARN Please initialize the logger system properly.
#Cluster Name     #Broker Name            #BID  #Addr                  #Version                #InTPS(LOAD)       #OutTPS(LOAD) #PCWait(ms) #Hour #SPACE
RaftCluster       RaftNode00              0     192.168.43.3:10913     V4_8_0                   5.08(0,0ms)         5.08(0,0ms)          0 448110.12 0.4695
RaftCluster       RaftNode00              1     192.168.43.3:10911     V4_8_0                   4.98(0,0ms)         0.00(0,0ms)          0 448110.12 0.4695
RaftCluster       RaftNode00              2     192.168.43.3:10912     V4_8_0                   5.00(0,0ms)         0.00(0,0ms)          0 448110.12 0.4695
RaftCluster       RaftNode01              0     192.168.43.3:10916     V4_8_0                   5.10(0,0ms)         5.10(0,0ms)          0 448110.12 0.4695
RaftCluster       RaftNode01              1     192.168.43.3:10914     V4_8_0                   5.00(0,0ms)         0.00(0,0ms)          0 448110.12 0.4695
RaftCluster       RaftNode01              2     192.168.43.3:10915     V4_8_0                   5.00(0,0ms)         0.00(0,0ms)          0 448110.12 0.4695
```

4、控制台信息

![image-20210213125704406](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-13/3d811a5889904ef69293d83834a66c4f.png)

## 2、topic相关

### 1、创建topic(2b-2w-2r)

> ​	在没有指定broker的情况下，默认在每个`broker`都创建`topic`

```shell
➜  ~ docker exec -it docker_broker-05_1  mqadmin updateTopic -n "nameserver-01:9876;nameserver-02:9876"  -p 6 -r 2 -w 2  -t TestTopic-2  -c RaftCluster
RocketMQLog:WARN No appenders could be found for logger (io.netty.util.internal.PlatformDependent0).
RocketMQLog:WARN Please initialize the logger system properly.
create topic to 192.168.43.3:10913 success.
create topic to 192.168.43.3:10916 success.
TopicConfig [topicName=TestTopic-2, readQueueNums=2, writeQueueNums=2, perm=RW-, topicFilterType=SINGLE_TAG, topicSysFlag=0, order=false]
```

查看属性

```shell
➜  ~ docker exec -it docker_broker-05_1  mqadmin topicStatus  -n "nameserver-01:9876;nameserver-02:9876" -t TestTopic-2
RocketMQLog:WARN No appenders could be found for logger (io.netty.util.internal.PlatformDependent0).
RocketMQLog:WARN Please initialize the logger system properly.
#Broker Name                      #QID  #Min Offset           #Max Offset             #Last Updated
RaftNode00                        0     0                     0
RaftNode00                        1     0                     0
RaftNode01                        0     0                     0
RaftNode01                        1     0                     0
```

### 2、创建topic (1b-2w-2r)

```shell
root@0bc0c79ca379:/opt/rocketmq#  mqadmin updateTopic -n "nameserver-01:9876;nameserver-02:9876"  -p 6 -r 12 -w 6  -t TestTopic-3  -c RaftCluster
RocketMQLog:WARN No appenders could be found for logger (io.netty.util.internal.PlatformDependent0).
RocketMQLog:WARN Please initialize the logger system properly.
create topic to 192.168.43.3:10913 success.
create topic to 192.168.43.3:10916 success.
TopicConfig [topicName=TestTopic-3, readQueueNums=12, writeQueueNums=6, perm=RW-, topicFilterType=SINGLE_TAG, topicSysFlag=0, order=false]
```

查看

```shell
root@0bc0c79ca379:/opt/rocketmq# mqadmin topicStatus  -n "nameserver-01:9876;nameserver-02:9876" -t TestTopic-3
RocketMQLog:WARN No appenders could be found for logger (io.netty.util.internal.PlatformDependent0).
RocketMQLog:WARN Please initialize the logger system properly.
#Broker Name                      #QID  #Min Offset           #Max Offset             #Last Updated
RaftNode00                        0     0                     0
RaftNode00                        1     0                     0
RaftNode00                        2     0                     0
RaftNode00                        3     0                     0
RaftNode00                        4     0                     0
RaftNode00                        5     0                     0
RaftNode01                        0     0                     0
RaftNode01                        1     0                     0
RaftNode01                        2     0                     0
RaftNode01                        3     0                     0
RaftNode01                        4     0                     0
RaftNode01                        5     0                     0
```

### 3、关于topic读取和写

#### 1、配置为`2b-1w-1r`

对于多broker来说，其实topic是配置了多份broker，所以会出现*n的现象

![image-20210213134055048](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-13/c64707e2cea24e0cab7dc8aa6e8293e3.png)

所以对于1w-1r的topic来说，其写入的队列其实是两个

![image-20210213134207635](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-13/5ab298d87424436e951d87d6ed0194e7.png)

消费的信息，可以看到可以消费提高了能力

```shell
subscribe callback: QueueId:0, QueueOffset:1672, message:Hello RocketMQ Go Client! 3344, store_host: 192.168.43.3:10916
subscribe callback: QueueId:0, QueueOffset:1672, message:Hello RocketMQ Go Client! 3345, store_host: 192.168.43.3:10913
```

然后就是其读取的队列了，其实是两个上限，当第三者加入进去后会发现

```shell
time="2021-02-13T13:38:37+08:00" level=warning msg="[BUG] ConsumerId not in cidAll" cidAll="[192.168.43.3@4424 192.168.43.3@4451]" consumerGroup=go_client_dev_2 consumerId=192.168.43.3@4479
```

#### 2、配置为`2b-1w-2r`

对于消费来说，只能有两个消费者,同样会提示以下信息：

```shell
time="2021-02-13T13:47:31+08:00" level=warning msg="[BUG] ConsumerId not in cidAll" cidAll="[192.168.43.3@4667 192.168.43.3@4696]" consumerGroup=go_client_dev_2 consumerId=192.168.43.3@4723
```

同样的，我们可以看到消费信息，的queue其实多了两个，分别是`node1-1`和`node0-1`，但是写入未写入进入，所以读取和写入队列设置的时候最好一致。

![image-20210213134911125](/Users/fanhaodong/Library/Application Support/typora-user-images/image-20210213134911125.png)

## 3、消费者/生产者相关

### 1、查看消费客户端信息

```shell
root@f30df1853117:/opt/rocketmq# mqadmin consumerStatus -n "nameserver-01:9876;nameserver-02:9876"  -g go_client_dev
RocketMQLog:WARN No appenders could be found for logger (io.netty.util.internal.PlatformDependent0).
RocketMQLog:WARN Please initialize the logger system properly.
001  192.168.43.3@5231                        V4_5_2               1613197331645/192.168.43.3@5231
002  192.168.43.3@5289                        V4_5_2               1613197331645/192.168.43.3@5289
003  192.168.43.3@5259                        V4_5_2               1613197331645/192.168.43.3@5259
```

### 2、查看消费积压

1）可以看到消费队列是 4个，每个broker分别有两个

```shell
➜  ~ docker exec -it docker_broker-05_1  mqadmin consumerProgress  -n "nameserver-01:9876;nameserver-02:9876"  -g go_client_dev
RocketMQLog:WARN No appenders could be found for logger (io.netty.util.internal.PlatformDependent0).
RocketMQLog:WARN Please initialize the logger system properly.
#Topic                            #Broker Name                      #QID  #Broker Offset        #Consumer Offset      #Client IP           #Diff                 #LastTime
%RETRY%go_client_dev              RaftNode00                        0     0                     0                     N/A                  0                     N/A
%RETRY%go_client_dev              RaftNode01                        0     0                     0                     N/A                  0                     N/A
TestTopic-6                       RaftNode00                        0     6006                  6005                  N/A                  1                     2021-02-13 06:26:29
TestTopic-6                       RaftNode00                        1     0                     0                     N/A                  0                     N/A
TestTopic-6                       RaftNode01                        0     6006                  6005                  N/A                  1                     2021-02-13 06:26:29
TestTopic-6                       RaftNode01                        1     0                     0                     N/A                  0                     N/A

Consume TPS: 10.01
Diff Total: 2
```

2）topic信息，可以看到生产者的队列，每个broker只有一个，这是因为创建的时候设置的`1w-2r`

```shell
➜  ~ docker exec -it docker_broker-05_1  mqadmin topicStatus  -n "nameserver-01:9876;nameserver-02:9876"  -t TestTopic-6
RocketMQLog:WARN No appenders could be found for logger (io.netty.util.internal.PlatformDependent0).
RocketMQLog:WARN Please initialize the logger system properly.
#Broker Name                      #QID  #Min Offset           #Max Offset             #Last Updated
RaftNode00                        0     0                     6311                    2021-02-13 06:27:30,408
RaftNode01                        0     0                     6311                    2021-02-13 06:27:30,309
```

## 4、broker配置参数相关

[1、官方配置介绍：](https://github.com/apache/rocketmq/blob/release-4.8.0/docs/cn/best_practice.md#33-broker-%E9%85%8D%E7%BD%AE)

| 参数名                  | 默认值                    | 说明                                                         |
| ----------------------- | ------------------------- | ------------------------------------------------------------ |
| listenPort              | 10911                     | 接受客户端连接的监听端口                                     |
| namesrvAddr             | null                      | nameServer 地址                                              |
| brokerIP1               | 网卡的 InetAddress        | 当前 broker 监听的 IP                                        |
| brokerIP2               | 跟 brokerIP1 一样         | 存在主从 broker 时，如果在 broker 主节点上配置了 brokerIP2 属性，broker 从节点会连接主节点配置的 brokerIP2 进行同步 |
| brokerName              | null                      | broker 的名称                                                |
| brokerClusterName       | DefaultCluster            | 本 broker 所属的 Cluser 名称                                 |
| brokerId                | 0                         | broker id, 0 表示 master, 其他的正整数表示 slave             |
| storePathCommitLog      | $HOME/store/commitlog/    | 存储 commit log 的路径                                       |
| storePathConsumerQueue  | $HOME/store/consumequeue/ | 存储 consume queue 的路径                                    |
| mappedFileSizeCommitLog | 1024 * 1024 * 1024(1G)    | commit log 的映射文件大小                                    |
| deleteWhen              | 04                        | 在每天的什么时间删除已经超过文件保留时间的 commit log        |
| fileReservedTime        | 72                        | 以小时计算的文件保留时间                                     |
| brokerRole              | ASYNC_MASTER              | SYNC_MASTER/ASYNC_MASTER/SLAVE                               |
| flushDiskType           | ASYNC_FLUSH               | SYNC_FLUSH/ASYNC_FLUSH SYNC_FLUSH 模式下的 broker 保证在收到确认生产者之前将消息刷盘。ASYNC_FLUSH 模式下的 broker 则利用刷盘一组消息的模式，可以取得更好的性能。 |

2、broker全部配置，这个是启动broker的启动参数，在目录的`/root/logs/rocketmqlogs/broker.log`下面

```shell
2021-02-05 00:44:49 INFO main - rocketmqHome=/opt/rocketmq
2021-02-05 00:44:49 INFO main - namesrvAddr=nameserver-01:9876;nameserver-02:9876
2021-02-05 00:44:49 INFO main - brokerIP1=172.25.0.5
2021-02-05 00:44:49 INFO main - brokerIP2=172.25.0.5
2021-02-05 00:44:49 INFO main - brokerName=RaftNode00
2021-02-05 00:44:49 INFO main - brokerClusterName=RaftCluster
2021-02-05 00:44:49 INFO main - brokerId=-1
2021-02-05 00:44:49 INFO main - brokerPermission=6
2021-02-05 00:44:49 INFO main - defaultTopicQueueNums=8
2021-02-05 00:44:49 INFO main - autoCreateTopicEnable=true
2021-02-05 00:44:49 INFO main - clusterTopicEnable=true
2021-02-05 00:44:49 INFO main - brokerTopicEnable=true
2021-02-05 00:44:49 INFO main - autoCreateSubscriptionGroup=true
2021-02-05 00:44:49 INFO main - messageStorePlugIn=
2021-02-05 00:44:49 INFO main - msgTraceTopicName=RMQ_SYS_TRACE_TOPIC
2021-02-05 00:44:49 INFO main - traceTopicEnable=false
2021-02-05 00:44:49 INFO main - sendMessageThreadPoolNums=4
2021-02-05 00:44:49 INFO main - pullMessageThreadPoolNums=28
2021-02-05 00:44:49 INFO main - processReplyMessageThreadPoolNums=28
2021-02-05 00:44:49 INFO main - queryMessageThreadPoolNums=14
2021-02-05 00:44:49 INFO main - adminBrokerThreadPoolNums=16
2021-02-05 00:44:49 INFO main - clientManageThreadPoolNums=32
2021-02-05 00:44:49 INFO main - consumerManageThreadPoolNums=32
2021-02-05 00:44:49 INFO main - heartbeatThreadPoolNums=6
2021-02-05 00:44:49 INFO main - endTransactionThreadPoolNums=20
2021-02-05 00:44:49 INFO main - flushConsumerOffsetInterval=5000
2021-02-05 00:44:49 INFO main - flushConsumerOffsetHistoryInterval=60000
2021-02-05 00:44:49 INFO main - rejectTransactionMessage=false
2021-02-05 00:44:49 INFO main - fetchNamesrvAddrByAddressServer=false
2021-02-05 00:44:49 INFO main - sendThreadPoolQueueCapacity=10000
2021-02-05 00:44:49 INFO main - pullThreadPoolQueueCapacity=100000
2021-02-05 00:44:49 INFO main - replyThreadPoolQueueCapacity=10000
2021-02-05 00:44:49 INFO main - queryThreadPoolQueueCapacity=20000
2021-02-05 00:44:49 INFO main - clientManagerThreadPoolQueueCapacity=1000000
2021-02-05 00:44:49 INFO main - consumerManagerThreadPoolQueueCapacity=1000000
2021-02-05 00:44:49 INFO main - heartbeatThreadPoolQueueCapacity=50000
2021-02-05 00:44:49 INFO main - endTransactionPoolQueueCapacity=100000
2021-02-05 00:44:49 INFO main - filterServerNums=0
2021-02-05 00:44:49 INFO main - longPollingEnable=true
2021-02-05 00:44:49 INFO main - shortPollingTimeMills=1000
2021-02-05 00:44:49 INFO main - notifyConsumerIdsChangedEnable=true
2021-02-05 00:44:49 INFO main - highSpeedMode=false
2021-02-05 00:44:49 INFO main - commercialEnable=true
2021-02-05 00:44:49 INFO main - commercialTimerCount=1
2021-02-05 00:44:49 INFO main - commercialTransCount=1
2021-02-05 00:44:49 INFO main - commercialBigCount=1
2021-02-05 00:44:49 INFO main - commercialBaseCount=1
2021-02-05 00:44:49 INFO main - transferMsgByHeap=true
2021-02-05 00:44:49 INFO main - maxDelayTime=40
2021-02-05 00:44:49 INFO main - regionId=DefaultRegion
2021-02-05 00:44:49 INFO main - registerBrokerTimeoutMills=6000
2021-02-05 00:44:49 INFO main - slaveReadEnable=false
2021-02-05 00:44:49 INFO main - disableConsumeIfConsumerReadSlowly=false
2021-02-05 00:44:49 INFO main - consumerFallbehindThreshold=17179869184
2021-02-05 00:44:49 INFO main - brokerFastFailureEnable=true
2021-02-05 00:44:49 INFO main - waitTimeMillsInSendQueue=200
2021-02-05 00:44:49 INFO main - waitTimeMillsInPullQueue=5000
2021-02-05 00:44:49 INFO main - waitTimeMillsInHeartbeatQueue=31000
2021-02-05 00:44:49 INFO main - waitTimeMillsInTransactionQueue=3000
2021-02-05 00:44:49 INFO main - startAcceptSendRequestTimeStamp=0
2021-02-05 00:44:49 INFO main - traceOn=true
2021-02-05 00:44:49 INFO main - enableCalcFilterBitMap=false
2021-02-05 00:44:49 INFO main - expectConsumerNumUseFilter=32
2021-02-05 00:44:49 INFO main - maxErrorRateOfBloomFilter=20
2021-02-05 00:44:49 INFO main - filterDataCleanTimeSpan=86400000
2021-02-05 00:44:49 INFO main - filterSupportRetry=false
2021-02-05 00:44:49 INFO main - enablePropertyFilter=false
2021-02-05 00:44:49 INFO main - compressedRegister=false
2021-02-05 00:44:49 INFO main - forceRegister=true
2021-02-05 00:44:49 INFO main - registerNameServerPeriod=30000
2021-02-05 00:44:49 INFO main - transactionTimeOut=6000
2021-02-05 00:44:49 INFO main - transactionCheckMax=15
2021-02-05 00:44:49 INFO main - transactionCheckInterval=60000
2021-02-05 00:44:49 INFO main - aclEnable=false
2021-02-05 00:44:49 INFO main - storeReplyMessageEnable=true
2021-02-05 00:44:49 INFO main - autoDeleteUnusedStats=false
2021-02-05 00:44:49 INFO main - listenPort=30911
2021-02-05 00:44:49 INFO main - serverWorkerThreads=8
2021-02-05 00:44:49 INFO main - serverCallbackExecutorThreads=0
2021-02-05 00:44:49 INFO main - serverSelectorThreads=3
2021-02-05 00:44:49 INFO main - serverOnewaySemaphoreValue=256
2021-02-05 00:44:49 INFO main - serverAsyncSemaphoreValue=64
2021-02-05 00:44:49 INFO main - serverChannelMaxIdleTimeSeconds=120
2021-02-05 00:44:49 INFO main - serverSocketSndBufSize=131072
2021-02-05 00:44:49 INFO main - serverSocketRcvBufSize=131072
2021-02-05 00:44:49 INFO main - serverPooledByteBufAllocatorEnable=true
2021-02-05 00:44:49 INFO main - useEpollNativeSelector=false
2021-02-05 00:44:49 INFO main - clientWorkerThreads=4
2021-02-05 00:44:49 INFO main - clientCallbackExecutorThreads=6
2021-02-05 00:44:49 INFO main - clientOnewaySemaphoreValue=65535
2021-02-05 00:44:49 INFO main - clientAsyncSemaphoreValue=65535
2021-02-05 00:44:49 INFO main - connectTimeoutMillis=3000
2021-02-05 00:44:49 INFO main - channelNotActiveInterval=60000
2021-02-05 00:44:49 INFO main - clientChannelMaxIdleTimeSeconds=120
2021-02-05 00:44:49 INFO main - clientSocketSndBufSize=131072
2021-02-05 00:44:49 INFO main - clientSocketRcvBufSize=131072
2021-02-05 00:44:49 INFO main - clientPooledByteBufAllocatorEnable=false
2021-02-05 00:44:49 INFO main - clientCloseSocketIfTimeout=false
2021-02-05 00:44:49 INFO main - useTLS=false
2021-02-05 00:44:49 INFO main - storePathRootDir=/root/store
2021-02-05 00:44:49 INFO main - storePathCommitLog=/root/store/commitlog
2021-02-05 00:44:49 INFO main - mappedFileSizeCommitLog=1073741824
2021-02-05 00:44:49 INFO main - mappedFileSizeConsumeQueue=6000000
2021-02-05 00:44:49 INFO main - enableConsumeQueueExt=false
2021-02-05 00:44:49 INFO main - mappedFileSizeConsumeQueueExt=50331648
2021-02-05 00:44:49 INFO main - bitMapLengthConsumeQueueExt=64
2021-02-05 00:44:49 INFO main - flushIntervalCommitLog=500
2021-02-05 00:44:49 INFO main - commitIntervalCommitLog=200
2021-02-05 00:44:49 INFO main - useReentrantLockWhenPutMessage=false
2021-02-05 00:44:49 INFO main - flushCommitLogTimed=false
2021-02-05 00:44:49 INFO main - flushIntervalConsumeQueue=1000
2021-02-05 00:44:49 INFO main - cleanResourceInterval=10000
2021-02-05 00:44:49 INFO main - deleteCommitLogFilesInterval=100
2021-02-05 00:44:49 INFO main - deleteConsumeQueueFilesInterval=100
2021-02-05 00:44:49 INFO main - destroyMapedFileIntervalForcibly=120000
2021-02-05 00:44:49 INFO main - redeleteHangedFileInterval=120000
2021-02-05 00:44:49 INFO main - deleteWhen=04
2021-02-05 00:44:49 INFO main - diskMaxUsedSpaceRatio=75
2021-02-05 00:44:49 INFO main - fileReservedTime=72
2021-02-05 00:44:49 INFO main - putMsgIndexHightWater=600000
2021-02-05 00:44:49 INFO main - maxMessageSize=4194304
2021-02-05 00:44:49 INFO main - checkCRCOnRecover=true
2021-02-05 00:44:49 INFO main - flushCommitLogLeastPages=4
2021-02-05 00:44:49 INFO main - commitCommitLogLeastPages=4
2021-02-05 00:44:49 INFO main - flushLeastPagesWhenWarmMapedFile=4096
2021-02-05 00:44:49 INFO main - flushConsumeQueueLeastPages=2
2021-02-05 00:44:49 INFO main - flushCommitLogThoroughInterval=10000
2021-02-05 00:44:49 INFO main - commitCommitLogThoroughInterval=200
2021-02-05 00:44:49 INFO main - flushConsumeQueueThoroughInterval=60000
2021-02-05 00:44:49 INFO main - maxTransferBytesOnMessageInMemory=262144
2021-02-05 00:44:49 INFO main - maxTransferCountOnMessageInMemory=32
2021-02-05 00:44:49 INFO main - maxTransferBytesOnMessageInDisk=65536
2021-02-05 00:44:49 INFO main - maxTransferCountOnMessageInDisk=8
2021-02-05 00:44:49 INFO main - accessMessageInMemoryMaxRatio=40
2021-02-05 00:44:49 INFO main - messageIndexEnable=true
2021-02-05 00:44:49 INFO main - maxHashSlotNum=5000000
2021-02-05 00:44:49 INFO main - maxIndexNum=20000000
2021-02-05 00:44:49 INFO main - maxMsgsNumBatch=64
2021-02-05 00:44:49 INFO main - messageIndexSafe=false
2021-02-05 00:44:49 INFO main - haListenPort=30912
2021-02-05 00:44:49 INFO main - haSendHeartbeatInterval=5000
2021-02-05 00:44:49 INFO main - haHousekeepingInterval=20000
2021-02-05 00:44:49 INFO main - haTransferBatchSize=32768
2021-02-05 00:44:49 INFO main - haMasterAddress=
2021-02-05 00:44:49 INFO main - haSlaveFallbehindMax=268435456
2021-02-05 00:44:49 INFO main - brokerRole=ASYNC_MASTER
2021-02-05 00:44:49 INFO main - flushDiskType=ASYNC_FLUSH
2021-02-05 00:44:49 INFO main - syncFlushTimeout=5000
2021-02-05 00:44:49 INFO main - messageDelayLevel=1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
2021-02-05 00:44:49 INFO main - flushDelayOffsetInterval=10000
2021-02-05 00:44:49 INFO main - cleanFileForciblyEnable=true
2021-02-05 00:44:49 INFO main - warmMapedFileEnable=false
2021-02-05 00:44:49 INFO main - offsetCheckInSlave=false
2021-02-05 00:44:49 INFO main - debugLockEnable=false
2021-02-05 00:44:49 INFO main - duplicationEnable=false
2021-02-05 00:44:49 INFO main - diskFallRecorded=true
2021-02-05 00:44:49 INFO main - osPageCacheBusyTimeOutMills=1000
2021-02-05 00:44:49 INFO main - defaultQueryMaxNum=32
2021-02-05 00:44:49 INFO main - transientStorePoolEnable=false
2021-02-05 00:44:49 INFO main - transientStorePoolSize=5
2021-02-05 00:44:49 INFO main - fastFailIfNoBufferInStorePool=false
2021-02-05 00:44:49 INFO main - enableDLegerCommitLog=true
2021-02-05 00:44:49 INFO main - dLegerGroup=RaftNode00
2021-02-05 00:44:49 INFO main - dLegerPeers=n0-broker-01:40911;n1-broker-02:40911
2021-02-05 00:44:49 INFO main - dLegerSelfId=n0
2021-02-05 00:44:49 INFO main - preferredLeaderId=
2021-02-05 00:44:49 INFO main - isEnableBatchPush=false
```

