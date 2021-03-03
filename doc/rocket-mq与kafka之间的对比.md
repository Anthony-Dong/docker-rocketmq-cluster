## 1、写入性能对比

1、rocket-mq对于写入是采用的是同步和异步方式，由于业务中使用的都是同步方式，所以broker都采用的是同步方式。具体的异步同步模式可以看：[官网的生产者-最佳实践中](https://github.com/apache/rocketmq/blob/master/docs/cn/best_practice.md)

> ​	具体下面不懂的参考：[CommitLog和ConsumerQueue和IndexFile源码分析](./CommitLog和ConsumerQueue和IndexFile分析.md)

- 写入的时候每个broker，都会去顺写commit-log 日志文件，这个文件采用append的方式，所有的topic都在一个commit-log日志下(所以在Topic较多的情况下写入速度比较高)，如果大小大于设置的阈值，则会创建一个新的
- 写入 `consumer-log`采用异步的方式，开启一个线程，定期的扫描commit-log日志文件，写入topic信息 （以每个topic-queue）进行隔离
- 写入`index-file`也是采用的异步的方式，跟`consumer-log`一样（每个broker写入一个文件）

2、kafka

- kafka 的日志采用多文件模式，也就是每个topic，每个分区，都会创建日志文件，也就是写入性能会随着topic数量下降
- kafka的索引文件采用稀疏索引，索引读取方式采用 `memory map`的方式，可以降低**程序的内存占用**
- 我们的协议是建立在一个 “消息块” 的抽象基础上，合理将消息分组。 这使得网络请求将多个消息打包成一组，而不是每次发送一条消息，从而使整组消息分担网络中往返的开销。Consumer 每次获取多个大型有序的消息块，并由服务端 依次将消息块一次加载到它的日志中。
- 客户端会

## 2、读取性能对比

1、对于上面写解释来说，rocket-mq采用的是 `mmap`技术，对于读取每条消息来说，他会先去读取`consumer-log`然后查看到偏移量，然后读取 `commit-log`，对于读写基本一致的情况下(由于采用mmap方式，如果数据刚好在内存中，性能会提升很多)，读取性能会高一些，但是它对于随机读取文件的能力也会要求高一些，所以推荐使用ssd作为存储

2、kafka需要根据偏移量查找到指定的index文件，然后读取找到偏移量，由于稀疏索引的原因需要查找到指定的偏移量就需要读取日志进行顺序查找

3、kafka读直接通过 Linux的 `sendfile`机制进行读操作



## 3、源码难度和质量

1、kafka的代码写的很优秀，核心代码采用的`scala`语言作为主要开发语言，学习难度会较高，其次是代码的开源质量来说很高，本人学习的是`2.1.1`版本

2、rocket-mq 的代码质量整体来说，可以说是中等，对于一些核心实现来说可以值得学习一下，但是代码注释欠缺太多，可以看出来很多代码是阿里未开源的，但是代码难度来说可以说是偏于简单。代码冗余也稍微多一些。学习版本是`4.8.0`

3、kafka本地代码环境不好搭建，使用的`gradle`管理工具



## 4、为什么Rocket-MQ 采用MMap，kafka采用Sendfile技术

1、mmap是一种不需要内核态(page-cache) copy 到用户态的一种技术，但是它依靠内存映射可以对数据进行读取和加工

2、sendfile技术是一种 可以相对于 mmap技术来说，减少一次系统调用（局限在于无法进行文件的读取和修改操作），在 linux内核版本高于 2.4 的时候，可以减少一次 cpu->socket的拷贝，可以直接拷贝到网卡上的一种技术。更形象的说只是一个命令罢了（局限是 **输入流(读取)必须是文件**，这也就是无法进行socket读取，只能写入socket了）

1）所以这也就是 rocket-mq 采用的是 mmap技术，因为它需要在brocker发送给消费者时需要将数据进行二次的封装

2）kafka的数据采用的是客户端拆包了，直接将数据进行sendfile机制

具体关于 mmap技术和sendfile技术可以参考：[https://zhuanlan.zhihu.com/p/308054212](https://zhuanlan.zhihu.com/p/308054212)



## 5、总结

本质上来说，如何解决mq的读写性能问题，更多的层面还是从硬件着手，比如可以参考：

- [基于SSD的Kafka应用层缓存架构设计与实现](https://juejin.cn/post/6918703114061250568)
- [腾讯云CKafka冷热分离技术](https://juejin.cn/post/6844904046340341768)


## 参考

[磁盘I/O那些事](https://tech.meituan.com/2017/05/19/about-desk-io.html)

[零拷贝技术](https://zhuanlan.zhihu.com/p/308054212)

[基于SSD的Kafka应用层缓存架构设计与实现](https://juejin.cn/post/6918703114061250568)

[腾讯云CKafka]()