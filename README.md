## 1、概要

> ​	使用`docker-image`+`docker-compose`搭建的本地rocket-mq集群环境，使用的是`2broker(6node)+2name-server`搭建的，同时还有学习rocket-mq的一些文章：[doc](./doc)

`rocket-mq` 版本：`4.8.0`

`go-client`版本：`go get -u -v github.com/apache/rocketmq-client-go/v2@v2.0.0`

rocker-mq 官方文档：[文档链接](https://github.com/apache/rocketmq/tree/release-4.8.0/docs/cn)

学习文档：

- [源码环境搭建](./doc/源码环境搭建.md)
- [CommitLog&&ConsumerQueue&&IndexFile&&延时队列源码分析](./doc/CommitLog和ConsumerQueue和IndexFile分析.md)
- [broker学习](./doc/Broker学习.md)
- [name-server学习](./doc/name-server学习.md)
- [Rocket-MQ与Kafka的对比](./doc/rocket-mq与kafka之间的对比.md)
- [各种集群模式的优缺点](./doc/各种集群模式的优缺点.md)

## 2、特点

1、支持通过本地docker环境搭建rocket-mq 集群

2、支持横向拓展多个节点，便于查看

3、提供了WEB-UI的 [rocketmq-console](https://github.com/apache/rocketmq-externals/tree/master/rocketmq-console) 的启动支持

## 3、问题：

1、rocketmq 的 shell脚本的问题，主要原因是 `rocketMQ`的启动脚本shell的不规范问题，可以看 [https://github.com/apache/rocketmq/issues/2655](https://github.com/apache/rocketmq/issues/2655) 

2、启动时切记要修改JVM参数，不然本地集群启动起来瞬间爆炸，单台Broker内存启动为8G，可以通过环境变量`JAVA_OPT_EXT` 控制JVM启动参数，可能docker容器内存不足，直接被kill掉进程。

3、切记电脑的总内存分配给Docker的不要过于小，下面是我机器的docker配置信息

```shell
➜  ~ docker info
Containers: 55
 Running: 0
 Paused: 0
 Stopped: 55
Images: 323
Server Version: 18.09.2
## .......
OSType: linux
Architecture: x86_64
CPUs: 4
Total Memory: 5.818GiB
Name: linuxkit-025000000001
ID: 5JBH:7VOE:3R3W:6G4I:6NXD:AXGN:MI2Z:DNUI:7BZ5:KRT6:NC5T:AJGE
Docker Root Dir: /var/lib/docker
Debug Mode (client): false
Debug Mode (server): true
 File Descriptors: 24
 Goroutines: 50
 System Time: 2021-02-07T08:18:25.936909812Z
 EventsListeners: 2
## .......
```

## 4、项目目录

```shell
➜  rocket_mq git:(master) ✗ tree -L 1
.
├── Makefile ## 脚本
├── broker-01 ## 1broker，3节点
├── broker-02
├── broker-03
├── docker-compose.yml # docker-compose 启动脚本
├── image # rocket-mq的本地镜像搭建
├── nameserver-01 ## name-server两个节点
└── nameserver-02
```

通过修改本地文件可以横向拓展多个节点

## 5、启动

### 1、broker 配置

以Broker-01来说，配置文件在 `broker-01/conf/broker.conf`，所以启动参数为`mqbroker -c conf/broker.conf`，切记至少一个broker需要有三个节点，**如果broker两个节点的话如果down掉一个broker无法使用，现象是无法选举出新的master** , `DLeger`的好处是能够和`kafka`一样自动选择leader，不需要手动指定`brokerId=0`为master。

```properties
####################### broker 配置 ###########################
# 集群概念不清晰目前
brokerClusterName = RaftCluster
# broker的名称
brokerName=RaftNode00
# 如果采用dLeger模式搭建，默认为-1，所以不配置
# brokerId=-1
# 本地网卡IP，（如果是内网可以不设置这个，由于我们需要本地测试需要写你的本机eth0的网卡，切记要修改）
brokerIP1=192.168.1.4
# 默认监听10911端口
listenPort=10911
# name-server配置
namesrvAddr=nameserver-01:9876;nameserver-02:9876
# 不允许自动创建topic
autoCreateTopicEnable=false

################# DLeger配置 ##########################
# 是否开启DLegerCommitLog
enableDLegerCommitLog=true
# 这个最好和brokerName名称一致
dLegerGroup=RaftNode00
dLegerPeers=n0-broker-01:40911;n1-broker-02:40911
dLegerSelfId=n0
sendMessageThreadPoolNums=4
```

全部配置可以看: [broker](./doc/broker)

### 2、nameserver 配置

由于name-server 配置不是特别多，走默认配置即可（由于我们是单机器部署一个实例，所以不需要考虑环境冲突问题，也不推荐name-server和broker放在一台机器的做法）

默认配置可以看：[name-server](./doc/nameserver)

### 3、`dLeger`模式部署

本文采用的是 `dLeger`模式，具体的模式好坏可以看 ,[https://github.com/apache/rocketmq/blob/release-4.8.0/docs/cn/operation.md](https://github.com/apache/rocketmq/blob/release-4.8.0/docs/cn/operation.md) , 如果要高可用，必须要求每个broker至少有三台节点（由它的选举机制决定），具体down机模拟看后面

整体架构是：`2n+1m+2s`，这个集群可以随意拓展

| 机器节点      | 角色       | broker-name | 备注                                 |
| ------------- | ---------- | ----------- | ------------------------------------ |
| nameserver-01 | nameserver |             |                                      |
| nameserver-02 | nameserver |             |                                      |
| broker-01     | broker     | RaftNode00  | 通过选择选master，不需要指定brokerId |
| broker-02     | broker     | RaftNode00  |                                      |
| broker-03     | broker     | RaftNode00  |                                      |

```shell
➜  rocket_mq docker exec -it a629ee26b7a9  mqadmin clusterList -n "nameserver-01:9876;nameserver-02:9876"
RocketMQLog:WARN No appenders could be found for logger (io.netty.util.internal.PlatformDependent0).
RocketMQLog:WARN Please initialize the logger system properly.
#Cluster Name     #Broker Name            #BID  #Addr                  #Version                #InTPS(LOAD)       #OutTPS(LOAD) #PCWait(ms) #Hour #SPACE
RaftCluster       RaftNode00              0     172.15.64.10:10911     V4_8_0                   0.00(0,0ms)         0.00(0,0ms)          0 447918.62 -1.0000
RaftCluster       RaftNode00              2     172.15.64.10:10912     V4_8_0                   0.00(0,0ms)         0.00(0,0ms)          0 447918.62 -1.0000
RaftCluster       RaftNode00              3     172.15.64.10:10913     V4_8_0                   0.00(0,0ms)         0.00(0,0ms)          0 447918.62 -1.0000
```

![image-20210205144906904](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-5/9232ec4fc66c4d01bd615a7d78cfe371.png)

### 4、其他部署模式的优缺点

官方文档：[https://github.com/apache/rocketmq/blob/master/docs/cn/operation.md](https://github.com/apache/rocketmq/blob/master/docs/cn/operation.md)

其他文档：[集群模式](./doc/cluster)

### 5、启动

1、进入到 [docker目录](./docker) ，然后进入到 [镜像目录](./docker/image)执行`make` 即可

2、进入到, `make run` 或者 `docker-compose --compatibility up -d` , 由于做了资源限制，所以需要使用`--compatibility` 参数

```shell
➜  rocket_mq make run                
docker-compose --compatibility up -d
Creating network "rocket_mq_default" with the default driver
Creating rocket_mq_broker-01_1        ... done
Creating rocket_mq_nameserver-02_1    ... done
Creating rocket_mq_broker-03_1        ... done
Creating rocket_mq_rocketmq-console_1 ... done
Creating rocket_mq_broker-02_1        ... done
Creating rocket_mq_nameserver-01_1    ... done
```

2、查看broker节点

```shell
➜  rocket_mq docker ps
CONTAINER ID        IMAGE                                   COMMAND                  CREATED             STATUS              PORTS                                           NAMES
a629ee26b7a9        rocketmq:v4.8.0                         "mqbroker -c conf/br…"   2 minutes ago       Up 2 minutes        9876/tcp, 10911/tcp, 0.0.0.0:10912->10912/tcp   rocket_mq_broker-02_1
6dbbec743077        rocketmq:v4.8.0                         "mqnamesrv"              2 minutes ago       Up 2 minutes        10911/tcp, 0.0.0.0:9871->9876/tcp               rocket_mq_nameserver-01_1
e9ffbeba13bd        rocketmq:v4.8.0                         "mqbroker -c conf/br…"   2 minutes ago       Up 2 minutes        9876/tcp, 0.0.0.0:10911->10911/tcp              rocket_mq_broker-01_1
bad3443fd2b7        rocketmq:v4.8.0                         "mqnamesrv"              2 minutes ago       Up 2 minutes        10911/tcp, 0.0.0.0:9872->9876/tcp               rocket_mq_nameserver-02_1
6cb53d20f162        apacherocketmq/rocketmq-console:2.0.0   "sh -c 'java $JAVA_O…"   2 minutes ago       Up 2 minutes        0.0.0.0:8080->8080/tcp                          rocket_mq_rocketmq-console_1
10ae582153eb        rocketmq:v4.8.0                         "mqbroker -c conf/br…"   2 minutes ago       Up 2 minutes        9876/tcp, 10911/tcp, 0.0.0.0:10913->10913/tcp   rocket_mq_broker-03_1
```

3、查看docker节点的资源占用情况

```shell
➜  rocket_mq git:(master) ✗ docker stats --no-stream
CONTAINER ID        NAME                           CPU %               MEM USAGE / LIMIT   MEM %               NET I/O             BLOCK I/O           PIDS
a629ee26b7a9        rocket_mq_broker-02_1          34.21%              373.6MiB / 512MiB   72.96%              1.66MB / 1.16MB     41.6MB / 8.95MB     158
6dbbec743077        rocket_mq_nameserver-01_1      0.28%               131.9MiB / 256MiB   51.52%              843kB / 419kB       5.63MB / 8.96MB     41
e9ffbeba13bd        rocket_mq_broker-01_1          35.53%              390.1MiB / 512MiB   76.20%              1.49MB / 3.67MB     9.45MB / 8.95MB     160
bad3443fd2b7        rocket_mq_nameserver-02_1      0.30%               149.4MiB / 256MiB   58.38%              860kB / 442kB       30.5MB / 8.96MB     41
6cb53d20f162        rocket_mq_rocketmq-console_1   0.56%               240.6MiB / 256MiB   93.97%              1.45MB / 924kB      124MB / 46.4MB      43
10ae582153eb        rocket_mq_broker-03_1          38.44%              401.4MiB / 512MiB   78.40%              1.65MB / 1.15MB     2.95MB / 8.95MB     158
```

4、查看总占用的内存情况(这个是6个docker容器的内存使用情况)

```shell
➜  rocket_mq git:(master) ✗ docker stats --no-stream | awk '{print $4}' | sed '1d' | awk '{a+=$1}END{printf "%sM\n",a}'
1676.9M
```
