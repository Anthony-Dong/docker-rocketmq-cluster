## 一、RocketMQ 常用部署模式

官方的部署模式：[https://github.com/apache/rocketmq/blob/release-4.8.0/docs/cn/best_practice.md](https://github.com/apache/rocketmq/blob/release-4.8.0/docs/cn/best_practice.md)

RocketMQ 常用部署方案有以下几种：

- **单机模式**
- **多主模式**
- **双主双从/多主多从模式（异步复制）**
- **双主双从/多主多从模式（同步双写）**
- **Dledger 集群模式**

**(1)、单机模式**

这种模式就如该名单机模式一样，就是部署单个 RocketMQ Broker 来使用，一般使用这种方式在生产中风险较大，一旦 Broker 重启或者宕机时，会导致整个服务不可用，所以常常在学习、开发过程中才会使用这种模式。

优缺点分析：

- **优点：** 本地开发测试，配置简单，同步刷盘消息不会丢失。
- **缺点：** 不可靠，如果宕机会导致服务不可用。

**(2)、多主模式**

全部由 Broker Master 节点组成（即可能部署两个或者更多 Broker），生产者发送的数据会分别存入不同的 Broker 中，这样能够避免某个 Broker 一直接收处理数据从而负载过高。

优缺点分析：

- **优点：** 性能高，配置简单，单个 Master 宕机或重启维护对应用无影响，在磁盘配置为 RAID10 时，即使机器宕机不可恢复，由于 RAID10 磁盘非常可靠，消息也不会丢（异步刷盘可能会丢失少量消息，同步刷盘能保证消息不丢）。
- **缺点：** 单台服务器宕机期间，不可订阅该服务器上未被消费者消费的消息，只有机器恢复后才可恢复订阅，所以可能会影响消息的实时性。

**(3)、双主双从/多主多从模式（异步复制）**

一般会部署多个 Broker Master，同时也会为各个 Broker Master 部署一个 Broker Slave，且 Master 和 Slave 之间采用”异步复制数据”方式进行数据同步（主从同步消息会有延迟，毫秒级），这样在生产者将消息发送到 Broker Master 后不必等待数据同步到 Slave 节点，就返回成功。

优缺点分析：

- **优点：** 性能高，且磁盘损坏也不会丢失大量消息，消息实时性不会受影响，Master 宕机后，消费者仍然可以从 Slave 消费。
- **缺点：** 主备有短暂消息延迟，毫秒级，如果Master宕机，磁盘损坏情况，会丢失少量消息。

**(4)、双主双从/多主多从模式（同步双写）**

一般会部署多个 Broker Master，同时也会为各个 Broker Master 部署一个 Broker Slave，且 Master 和 Slave 之间采用”同步复制数据”方式进行数据同步，这样在生产者将消息发送到 Broker Master 后需要等待数据同步到 Slave 节点成功后，才返回成功。

优缺点分析：

- **优点：** 数据与服务都无单点故障，Master 宕机情况下，消息无延迟，服务可用性与数据可用性都非常高；
- **缺点：** 性能比异步复制模式略低（大约低10%左右），发送单个消息的 RT 会略高，且目前版本在主节点宕机后，备机不能自动切换为主机。

**(5)、Dledger 集群模式**

RocketMQ-on-DLedger Group 是指一组相同名称的 Broker，至少需要 3 个节点，通过 Raft 自动选举出一个 Leader，其余节点 作为 Follower，并在 Leader 和 Follower 之间复制数据以保证高可用。当 Master 节点出现问题宕机后也能自动容灾切换，并保证数据一致性。该模式也支持 Broker 水平扩展，即可以部署任意多个 RocketMQ-on-DLedger Group 同时对外提供服务。

优缺点分析：

- 优点：多节点（至少三个）组成集群，其中一个为 Leader 节点，其余为 Follower 节点组成高可用，能够自动容灾切换。
- 缺点：需要 RocketMQ 4.5 及以后版本才支持。

## 二、RocketMQ Dledger 集群模式简介

### 1、传统部署方式的不足

在 RocketMQ 4.5 之前的版本中，部署 RocketMQ 高可用方案一般都会采用多主多从方式，这种方式需要多个 Master 节点与实时备份 Master 节点数据的 Slave 节点，Slave 节点通过同步复制或异步复制的方式去同步 Master 节点数据。但这样的部署模式存在一定缺陷。比如故障转移方面，如果 Master 点挂了，还需要人为手动对 Master 节点进行重启或者切换，它无法自动的将 Slave 节点转换为 Master 节点。因此，我们希望能有一个新的多副本架构，去解决这个问题。

### 2、新技术解决的问题

新的多副本架构首先需要解决自动故障转移的问题，本质上来说问题关键点在于 Broker 如何自动推选主节点。这个问题的解决方案基本可以分为两种：

- 利用第三方协调服务集群完成选主，比如 Zookeeper 或者 Etcd，但是这种方案会引入了重量级外部组件，使部署变得复杂，同时也会增加运维对组件的故障诊断成本，比如在维护 RocketMQ 集群还需要维护 Zookeeper 集群，保证 Zookeeper 集群如何高可用，不仅仅如此，如果 zookeeper 集群出现故障也会影响到 RocketMQ 集群。
- 利用 raft 协议来完成一个自动选主，raft 协议相比前者的优点是不需要引入外部组件，自动选主逻辑集成到各个节点的进程中，节点之间通过通信就可以完成选主。

RocketMQ 最终选择使用 raft 协议来解决这个问题，而 DLedger 就是一个基于 raft 协议的 commitlog 存储库，也是 RocketMQ 实现新的高可用多副本架构的关键。

### 3、Dledger 简介

分布式算法中比较常常听到的是 Paxos 算法，但是由于 Paxos 算法难于理解，且实现比较苦难，所以不太受业界欢迎。然后出现新的分布式算法 Raft，其比 Paxos 更容易懂与实现，到如今在实际中运用的也已经很成熟，不同的语言都有对其的实现。Dledger 就是其中一个 Java 语言的实现，其将算法方面的内容全部抽象掉，这样开发人员只需要关系业务即可，大大降低使用难度。

> Dledger 相关资料来源于网上搜索。

### 4、DLedger 定位

![img](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-5/595ca0dd90f9494b97bfd575923b9904.png)

Raft 协议是复制状态机的实现，这种模型应用到消息系统中就会存在问题。对于消息系统来说，它本身是一个中间代理，commitlog 状态是系统最终状态，并不需要状态机再去完成一次状态构建。因此 DLedger 去掉了 raft 协议中状态机的部分，但基于raft协议保证commitlog 是一致的，并且是高可用的。

![img](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-5/5c9c3127f7994cd1965157057a37c0a1.png)

另一方面 DLedger 又是一个轻量级的 java library。它对外提供的 API 非常简单，append 和 get。Append 向 DLedger 添加数据，并且添加的数据会对应一个递增的索引，而 get 可以根据索引去获得相应的数据。因此 DLedger 是一个 append only 的日志系统。

### 5、DLedger 应用场景

![img](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-5/1b9bd25c12d44a18bbb20c5f89ad5537.png)

DLedger 其中一个应用就是在分布式消息系统中，RocketMQ 4.5 版本发布后，可以采用 RocketMQ on DLedger 方式进行部署。DLedger commitlog 代替了原来的 commitlog，使得 commitlog 拥有了选举复制能力，然后通过角色透传的方式，raft 角色透传给外部 broker 角色，leader 对应原来的 master，follower 和 candidate 对应原来的 slave。

因此 RocketMQ 的 broker 拥有了自动故障转移的能力，在一组 broker 中如果 Master 挂了，能够依靠 DLedger 自动选主能力重新选出一个 leader，然后通过角色透传变成新的 Master。

![img](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-5/10bb53c1231c4e71a1a6b2a50d27fe55.png)

DLedger 还可以构建高可用的嵌入式 KV 存储。我们把对一些数据的操作记录到 DLedger 中，然后根据数据量或者实际需求，恢复到hashmap 或者 rocksdb 中，从而构建一致的、高可用的 KV 存储系统，应用到元信息管理等场景。

### 6、RocketMQ Dledger 的方案简介

![img](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-5/76e6f941ff6d4c5a94e244e82d897feb.png)

RocketMQ-on-DLedger Group 是指一组相同名称的 Broker，组中至少需要 3 个 Broker 节点来保证集群能够运行，在 Broker 启动时候，通过 raft 算法能够自动选举出一个 Broker 为 Leader 节点，其余为 Follower 节点。这种模式下 Leader 和 Follower 之间复制数据以保证高可用，如果 Leader 节点出现问题是可以自动进行容灾切换并保证数据一致性。且不仅仅如此，该模式也支持 Broker 节点水平扩展来增加吞吐量。所以该模式将会是部署 RocketMQ 常用模式之一。





## 参考

[DLedger 模式集群部署](http://www.mydlq.club/article/97/)

[RocketMQ系列：搭建3m-3s模式的rocketmq集群](https://blog.51cto.com/14900374/2539774)