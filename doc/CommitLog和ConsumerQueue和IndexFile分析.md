## 1、消息的物理位置

### 1、数据文件(`CommitLog`)

​		数据文件(`CommitLog`)所在位置，由于采用的是`dledger`模式，所以会有以下目录，如果是其他模式则是在 `storePathCommitLog=/root/store/commitlog`下面

> ​	普通模式下：消息主体以及元数据的存储主体，存储Producer端写入的消息主体内容,消息内容不是定长的。单个文件大小默认1G ，文件名长度为20位，左边补零，剩余为起始偏移量，比如00000000000000000000代表了第一个文件，起始偏移量为0，文件大小为`mappedFileSizeCommitLog=1024 * 1024 * 1024`；当第一个文件写满了，第二个文件为00000000001073741824，起始物理偏移量为1073741824，以此类推。消息主要是顺序写入日志文件，当文件满了，写入下一个文件；
>
> ​	`dledger`模式下：后续再讲

```shell
# 文件名，其实位置为0，表示偏移量为0，单文件大小为1G，写满新建文件，比如写满的偏移量为1073741824，那么新文件名就是00000000001073741824
root@af03d8240e98:~/store/dledger-n0# ls -l {data,index}
data:
total 944
-rw-r--r-- 1 root root 1073741824 Feb  7 08:32 00000000000000000000
index:
total 128
-rw-r--r-- 1 root root 167772160 Feb  7 08:32 00000000000000000000
# 文件大小，单个文件大小默认1G，可以修改mappedFileSizeCommitLog参数去变更
root@af03d8240e98:~/store/dledger-n0/data# du -hs 00000000000000000000
944K	00000000000000000000
# 文件名长度20
root@af03d8240e98:~/store/dledger-n0/data# a=`ls`
root@af03d8240e98:~/store/dledger-n0/data# echo ${#a}
20
```

### 2、消费队列(`ConsumeQueue`)

> ​	消息消费队列，引入的目的主要是提高消息消费的性能，由于RocketMQ是基于主题topic的订阅模式，消息消费是针对主题进行的，如果要遍历commitlog文件中根据topic检索消息是非常低效的。Consumer即可根据ConsumeQueue来查找待消费的消息。其中，ConsumeQueue（逻辑消费队列）作为消费消息的索引，保存了指定Topic下的队列消息在CommitLog中的起始物理偏移量offset，消息大小size和消息Tag的HashCode值。consumequeue文件可以看成是基于topic的commitlog索引文件，故consumequeue文件夹的组织方式如下：topic/queue/file三层组织结构，具体存储路径为：$HOME/store/consumequeue/{topic}/{queueId}/{fileName}。同样consumequeue文件采取定长设计，每一个条目共20个字节，分别为8字节的commitlog物理偏移量、4字节的消息长度、8字节tag hashcode，单个文件由30W个条目组成，可以像数组一样随机访问每一个条目，每个ConsumeQueue文件大小约5.72M；

```shell
root@6063de159208:~/store/consumequeue/TopicTest# pwd
/root/store/consumequeue/TopicTest
root@6063de159208:~/store/consumequeue/TopicTest# ls -l
total 0
drwxr-xr-x 3 root root 96 Feb  7 11:34 0
drwxr-xr-x 3 root root 96 Feb  7 11:34 1
drwxr-xr-x 3 root root 96 Feb  7 11:34 2
drwxr-xr-x 3 root root 96 Feb  7 11:34 3
```

然后查看topic信息，可以看到有四个队列，注意这个主要是和写队列数量有关，推荐读写队列数量设置一致，不然会出现空的消费队列问题

```shell
root@6063de159208:~/store/consumequeue/TopicTest# mqadmin topicStatus -n "nameserver-01:9876;nameserver-02:9876" -t TopicTest
RocketMQLog:WARN No appenders could be found for logger (io.netty.util.internal.PlatformDependent0).
RocketMQLog:WARN Please initialize the logger system properly.
#Broker Name                      #QID  #Min Offset           #Max Offset             #Last Updated
RaftNode00                        0     0                     14                      2021-02-07 11:34:56,070
RaftNode00                        1     0                     14                      2021-02-07 11:34:55,770
RaftNode00                        2     0                     14                      2021-02-07 11:34:55,868
RaftNode00                        3     0                     14                      2021-02-07 11:34:55,972
```

文件命名是：`偏移量`

### 3、索引文件(`IndexFile`)

> ​	IndexFile（索引文件）提供了一种可以通过key或时间区间来查询消息的方法。Index文件的存储位置是：$HOME \store\index${fileName}，文件名fileName是以创建时的时间戳命名的，固定的单个IndexFile文件大小约为400M，一个IndexFile可以保存 2000W个索引，IndexFile的底层存储设计为在文件系统中实现HashMap结构，故rocketmq的索引文件其底层实现为hash索引。

```shell
root@8dcb9a7a5a49:~/store/index# du -hs ./*
20M	./20210207091126016
```

文件名是根据：`存储的时间搓进行命名的`,

```java
String fileName =this.storePath + File.separator+ UtilAll.timeMillisToHumanString(System.currentTimeMillis());
```

## 2、源码分析

> ​	`rocket-mq`的源码是Java写的，代码难度相对太低，通过我看`rocket-mq`的代码可以看到很多代码应该是未开源的，代码注释少到了极致。

### 1、`CommitLog`实现

`Topic`最大长度为`2<<8`个字节, properties 最大长度为 `2>>16`字节，msg最大值可以设置` maxMessageSize = 1024 * 1024 * 4`,默认是4M；

其实整体来说这块并不复杂，只是顺写日志罢了，rocket-mq大量使用了`mmap`技术去实现快速的读写，并且减少内存开销，有关于 `mmap`技术的可以看，[https://www.huaweicloud.com/articles/9d78b21a2838f491ca5ae899ae7a8467.html](https://www.huaweicloud.com/articles/9d78b21a2838f491ca5ae899ae7a8467.html)

![image-20210208163320581](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-8/58566019ecd64dfda1c23f3bc504bcdb.png)



1、[https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/CommitLog.java#L787](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/CommitLog.java#L787)

这个主要逻辑就是写消息，有个比较特殊的就是，支持 delay

同步刷盘，注意可以设置刷盘的超时时间，为`syncFlushTimeout = 1000 * 5= 5s `

![img](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-8/d9cc9e53235b404483fe1bdf234c0449.png)

原理其实就是

```java
public boolean flush(final int flushLeastPages) {
    boolean result = true;
    MappedFile mappedFile = this.findMappedFileByOffset(this.flushedWhere, this.flushedWhere == 0);
    if (mappedFile != null) {
      // mapped file 刷盘指定的page 
        long tmpTimeStamp = mappedFile.getStoreTimestamp();
        int offset = mappedFile.flush(flushLeastPages);
        long where = mappedFile.getFileFromOffset() + offset;
        result = where == this.flushedWhere;
        this.flushedWhere = where;
        if (0 == flushLeastPages) {
            this.storeTimestamp = tmpTimeStamp;
        }
    }
    return result;
}

//We only append data to fileChannel or mappedByteBuffer, never both.
if (writeBuffer != null || this.fileChannel.position() != 0) {
    this.fileChannel.force(false);
} else {
    this.mappedByteBuffer.force();
}
```

1、[https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/MappedFile.java#L199](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/MappedFile.java#L199)

```java
public AppendMessageResult appendMessagesInner(final MessageExt messageExt, final AppendMessageCallback cb) {
    assert messageExt != null;
    assert cb != null;

    int currentPos = this.wrotePosition.get();

    if (currentPos < this.fileSize) {
        // 其实每次都将buffer-浅拷贝一下，然后设置position，写入消息成功position增加消息长度
        ByteBuffer byteBuffer = writeBuffer != null ? writeBuffer.slice() : this.mappedByteBuffer.slice();
        // 设置当前的写入位置 currentPos
        byteBuffer.position(currentPos);
        AppendMessageResult result;
        if (messageExt instanceof MessageExtBrokerInner) {
            // 当前文件的开始偏移量，commit log是根据物理偏移量进行命令的
            // buffer
            // 文件剩余空间
            // 消息
            result = cb.doAppend(this.getFileFromOffset(), byteBuffer, this.fileSize - currentPos, (MessageExtBrokerInner) messageExt);
        } else if (messageExt instanceof MessageExtBatch) {
            result = cb.doAppend(this.getFileFromOffset(), byteBuffer, this.fileSize - currentPos, (MessageExtBatch) messageExt);
        } else {
            return new AppendMessageResult(AppendMessageStatus.UNKNOWN_ERROR);
        }
        this.wrotePosition.addAndGet(result.getWroteBytes());
        this.storeTimestamp = result.getStoreTimestamp();
        return result;
    }
    log.error("MappedFile.appendMessage return null, wrotePosition: {} fileSize: {}", currentPos, this.fileSize);
    return new AppendMessageResult(AppendMessageStatus.UNKNOWN_ERROR);
}
```

2、[https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/CommitLog.java#L1521](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/CommitLog.java#L1521)


```java
public AppendMessageResult doAppend(final long fileFromOffset, final ByteBuffer byteBuffer, final int maxBlank,
    final MessageExtBrokerInner msgInner) {
    // STORETIMESTAMP + STOREHOSTADDRESS + OFFSET <br>

    // PHY OFFSET
    // 消息的整体偏移量，byteBuffer.position(当前文件内的物理偏移量)+fileFromOffset(文件的物理偏移量)
    long wroteOffset = fileFromOffset + byteBuffer.position();

    int sysflag = msgInner.getSysFlag();

    int bornHostLength = (sysflag & MessageSysFlag.BORNHOST_V6_FLAG) == 0 ? 4 + 4 : 16 + 4;
    int storeHostLength = (sysflag & MessageSysFlag.STOREHOSTADDRESS_V6_FLAG) == 0 ? 4 + 4 : 16 + 4;

    // 内存
    ByteBuffer bornHostHolder = ByteBuffer.allocate(bornHostLength);
    ByteBuffer storeHostHolder = ByteBuffer.allocate(storeHostLength);

    this.resetByteBuffer(storeHostHolder, storeHostLength);
    String msgId;
    if ((sysflag & MessageSysFlag.STOREHOSTADDRESS_V6_FLAG) == 0) {
        msgId = MessageDecoder.createMessageId(this.msgIdMemory, msgInner.getStoreHostBytes(storeHostHolder), wroteOffset);
    } else {
        msgId = MessageDecoder.createMessageId(this.msgIdV6Memory, msgInner.getStoreHostBytes(storeHostHolder), wroteOffset);
    }

    //
    // Record ConsumeQueue information
    keyBuilder.setLength(0);
    keyBuilder.append(msgInner.getTopic());
    keyBuilder.append('-');
    keyBuilder.append(msgInner.getQueueId());
    String key = keyBuilder.toString();

    // 当前队列的偏移量： key 格式： Topic-QueueID
    Long queueOffset = CommitLog.this.topicQueueTable.get(key);
    if (null == queueOffset) {
        queueOffset = 0L;
        CommitLog.this.topicQueueTable.put(key, queueOffset);
    }

    // Transaction messages that require special handling
    final int tranType = MessageSysFlag.getTransactionValue(msgInner.getSysFlag());
    switch (tranType) {
        // Prepared and Rollback message is not consumed, will not enter the
        // consumer queuec
        case MessageSysFlag.TRANSACTION_PREPARED_TYPE:
        case MessageSysFlag.TRANSACTION_ROLLBACK_TYPE:
            queueOffset = 0L;
            break;
        case MessageSysFlag.TRANSACTION_NOT_TYPE:
        case MessageSysFlag.TRANSACTION_COMMIT_TYPE:
        default:
            break;
    }

    /**
     * Serialize message
     */

    // 消息属性
    final byte[] propertiesData =
        msgInner.getPropertiesString() == null ? null : msgInner.getPropertiesString().getBytes(MessageDecoder.CHARSET_UTF8);

    final int propertiesLength = propertiesData == null ? 0 : propertiesData.length;

    if (propertiesLength > Short.MAX_VALUE) {
        log.warn("putMessage message properties length too long. length={}", propertiesData.length);
        return new AppendMessageResult(AppendMessageStatus.PROPERTIES_SIZE_EXCEEDED);
    }

    // topic
    final byte[] topicData = msgInner.getTopic().getBytes(MessageDecoder.CHARSET_UTF8);
    final int topicLength = topicData.length;

    // body
    final int bodyLength = msgInner.getBody() == null ? 0 : msgInner.getBody().length;

    final int msgLen = calMsgLength(msgInner.getSysFlag(), bodyLength, topicLength, propertiesLength);

    // Exceeds the maximum message
    if (msgLen > this.maxMessageSize) {
        CommitLog.log.warn("message size exceeded, msg total size: " + msgLen + ", msg body size: " + bodyLength
            + ", maxMessageSize: " + this.maxMessageSize);
        return new AppendMessageResult(AppendMessageStatus.MESSAGE_SIZE_EXCEEDED);
    }


    // 消息长度+文件最小的空白长度>文件剩余空间
    // 一些reset操作
    // 返回EOF
    // Determines whether there is sufficient free space
    if ((msgLen + END_FILE_MIN_BLANK_LENGTH) > maxBlank) {
        this.resetByteBuffer(this.msgStoreItemMemory, maxBlank);
        // 1 TOTALSIZE
        this.msgStoreItemMemory.putInt(maxBlank);
        // 2 MAGICCODE
        this.msgStoreItemMemory.putInt(CommitLog.BLANK_MAGIC_CODE);
        // 3 The remaining space may be any value
        // Here the length of the specially set maxBlank
        final long beginTimeMills = CommitLog.this.defaultMessageStore.now();
        byteBuffer.put(this.msgStoreItemMemory.array(), 0, maxBlank);
        return new AppendMessageResult(AppendMessageStatus.END_OF_FILE, wroteOffset, maxBlank, msgId, msgInner.getStoreTimestamp(),
            queueOffset, CommitLog.this.defaultMessageStore.now() - beginTimeMills);
    }

    // 重置这个buffer，设置limit为message-len
    // Initialization of storage space
    this.resetByteBuffer(msgStoreItemMemory, msgLen);
    // 1 TOTALSIZE
    this.msgStoreItemMemory.putInt(msgLen);
    // 2 MAGICCODE
    this.msgStoreItemMemory.putInt(CommitLog.MESSAGE_MAGIC_CODE);
    // 3 BODYCRC
    this.msgStoreItemMemory.putInt(msgInner.getBodyCRC());
    // 4 QUEUEID
    this.msgStoreItemMemory.putInt(msgInner.getQueueId());
    // 5 FLAG
    this.msgStoreItemMemory.putInt(msgInner.getFlag());
    // 6 QUEUEOFFSET
    this.msgStoreItemMemory.putLong(queueOffset);
    // 7 PHYSICALOFFSET
    this.msgStoreItemMemory.putLong(fileFromOffset + byteBuffer.position());
    // 8 SYSFLAG
    this.msgStoreItemMemory.putInt(msgInner.getSysFlag());
    // 9 BORNTIMESTAMP
    this.msgStoreItemMemory.putLong(msgInner.getBornTimestamp());
    // 10 BORNHOST
    this.resetByteBuffer(bornHostHolder, bornHostLength);
    this.msgStoreItemMemory.put(msgInner.getBornHostBytes(bornHostHolder));
    // 11 STORETIMESTAMP
    this.msgStoreItemMemory.putLong(msgInner.getStoreTimestamp());
    // 12 STOREHOSTADDRESS
    this.resetByteBuffer(storeHostHolder, storeHostLength);
    this.msgStoreItemMemory.put(msgInner.getStoreHostBytes(storeHostHolder));
    // 13 RECONSUMETIMES
    this.msgStoreItemMemory.putInt(msgInner.getReconsumeTimes());
    // 14 Prepared Transaction Offset
    this.msgStoreItemMemory.putLong(msgInner.getPreparedTransactionOffset());
    // 15 BODY
    this.msgStoreItemMemory.putInt(bodyLength);
    if (bodyLength > 0)
        this.msgStoreItemMemory.put(msgInner.getBody());
    // 16 TOPIC
    this.msgStoreItemMemory.put((byte) topicLength);
    this.msgStoreItemMemory.put(topicData);
    // 17 PROPERTIES
    this.msgStoreItemMemory.putShort((short) propertiesLength);
    if (propertiesLength > 0)
        this.msgStoreItemMemory.put(propertiesData);

    final long beginTimeMills = CommitLog.this.defaultMessageStore.now();
    // Write messages to the queue buffer
    byteBuffer.put(this.msgStoreItemMemory.array(), 0, msgLen);

    AppendMessageResult result = new AppendMessageResult(AppendMessageStatus.PUT_OK, wroteOffset, msgLen, msgId,
        msgInner.getStoreTimestamp(), queueOffset, CommitLog.this.defaultMessageStore.now() - beginTimeMills);

    switch (tranType) {
        case MessageSysFlag.TRANSACTION_PREPARED_TYPE:
        case MessageSysFlag.TRANSACTION_ROLLBACK_TYPE:
            break;
        case MessageSysFlag.TRANSACTION_NOT_TYPE:
        case MessageSysFlag.TRANSACTION_COMMIT_TYPE:
            // The next update ConsumeQueue information
            CommitLog.this.topicQueueTable.put(key, ++queueOffset);
            break;
        default:
            break;
    }
    return result;
}
```

### 2、`ConsumerQueue` 实现

> ​	他其实就是做一个Map文件的映射，方便高速的消费，我们知道对于读取`commitlog`来说，说句的访问效率是极低的，因为它是顺写的需要遍历，其次是随机读取的缓慢。所以需要在写入消息的时候写入这个文件。rockte-mq中是开启一个线程去写`ConsumerQueue`

1、注册dispatch（下面`indexfile`）同理，不解释，通过 [doDispatch](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/DefaultMessageStore.java#L1506) 调用 ， [doDispatch](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/DefaultMessageStore.java#L1506)  由  [doReput方法](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/DefaultMessageStore.java#L1922) 调用

```java
this.dispatcherList = new LinkedList<>();
this.dispatcherList.addLast(new CommitLogDispatcherBuildConsumeQueue());
this.dispatcherList.addLast(new CommitLogDispatcherBuildIndex());
```

2、broker启动的时候会启动`ReputMessageService` 这个线程

[run](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/DefaultMessageStore.java#L1998)

```java
@Override
public void run() {
    DefaultMessageStore.log.info(this.getServiceName() + " service started");

    while (!this.isStopped()) {
        try {
            Thread.sleep(1);
            this.doReput();
        } catch (Exception e) {
            DefaultMessageStore.log.warn(this.getServiceName() + " service has exception. ", e);
        }
    }

    DefaultMessageStore.log.info(this.getServiceName() + " service end");
}
```

3、其次就是 `doReput`的实现

[https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/DefaultMessageStore.java#L1922](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/DefaultMessageStore.java#L1922)

核心的一步是:

```java
public void putMessagePositionInfo(DispatchRequest dispatchRequest) {
  // 根据topic和qid 获取 cq
    ConsumeQueue cq = this.findConsumeQueue(dispatchRequest.getTopic(), dispatchRequest.getQueueId());
    cq.putMessagePositionInfoWrapper(dispatchRequest);
}
```

4、其次就是 [org.apache.rocketmq.store.ConsumeQueue#putMessagePositionInfoWrapper](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/ConsumeQueue.java#L379)

主要看看它的[构造器](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/ConsumeQueue.java#L47)

```java
public ConsumeQueue(
    final String topic,
    final int queueId,
    final String storePath,
    final int mappedFileSize,
    final DefaultMessageStore defaultMessageStore) {
    this.storePath = storePath;
    this.mappedFileSize = mappedFileSize;
    this.defaultMessageStore = defaultMessageStore;

    this.topic = topic;
    this.queueId = queueId;

  // 存储路径： ${stroe.dir}/topic/queueid
    String queueDir = this.storePath
        + File.separator + topic
        + File.separator + queueId;

    this.mappedFileQueue = new MappedFileQueue(queueDir, mappedFileSize, null);

  // 每行CQ_STORE_UNIT_SIZE=20字节
    this.byteBufferIndex = ByteBuffer.allocate(CQ_STORE_UNIT_SIZE);

    if (defaultMessageStore.getMessageStoreConfig().isEnableConsumeQueueExt()) {
        this.consumeQueueExt = new ConsumeQueueExt(
            topic,
            queueId,
            StorePathConfigHelper.getStorePathConsumeQueueExt(defaultMessageStore.getMessageStoreConfig().getStorePathRootDir()),
            defaultMessageStore.getMessageStoreConfig().getMappedFileSizeConsumeQueueExt(),
            defaultMessageStore.getMessageStoreConfig().getBitMapLengthConsumeQueueExt()
        );
    }
}
```

然后查看`messaage`的组成: [putMessagePositionInfo](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/ConsumeQueue.java#L425)

```java
this.byteBufferIndex.flip();
this.byteBufferIndex.limit(CQ_STORE_UNIT_SIZE);
this.byteBufferIndex.putLong(offset);// 偏移量
this.byteBufferIndex.putInt(size); // 消息的大小
this.byteBufferIndex.putLong(tagsCode); // 消息类型，MULTI_TAGS_FLAG||SINGLE_TAG
```

#### 1、如何根据偏移量进行查询

[org.apache.rocketmq.store.DefaultMessageStore#getMessage](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/DefaultMessageStore.java#L555)

1） 文件命名是以 `偏移量`进行命令的

2）然后根据偏移量进行查询 指定的文件 （二分查询，[org.apache.rocketmq.store.MappedFileQueue#findMappedFileByOffset(long, boolean)](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/MappedFileQueue.java#L462)）

3）找到文件后，然后查询偏移量信息，根据偏移量*固定步长(consumer-queue 每个消息20字节) % 文件固定长度, 就可以找到文件的物理位置( [org.apache.rocketmq.store.ConsumeQueue#getIndexBuffer](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/ConsumeQueue.java#L491))

4）读取消息即可

#### 2、如何根据时间进行查询

具体逻辑在： [org.apache.rocketmq.store.DefaultMessageStore#getOffsetInQueueByTime](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/DefaultMessageStore.java#L759)

1、先根据文件的modify时间，选择文件(所以可以依靠文件的变更时间进行确认时间，这里有个问题就是：consumer-queue是异步写的，但是实际生产时间一定是小于写入时间，也就是说一定不会出现选错文件的问题）

2、然后遍历即可，这个时间复杂程度较高

### 3、`IndexFile`核心实现

其实commitlog就是元数据文件，而consumer-queue可以看作是每个TOPIC的commitlog的索引文件，比如我们消费一条消息，知道broker，去拿topic，然后去拿指定的队列ID，可以看到只需要顺序的去读取消费队列，每个消费会告诉commitlog的物理存储索引位置，然后读取出来即可。

关于索引`index` 其实是一个特殊的hashmap，它的key是

`org.apache.rocketmq.store.index.IndexService#putKey#L223`

```java
if (req.getUniqKey() != null) {
    indexFile = putKey(indexFile, msg, buildKey(topic, req.getUniqKey()));
    if (indexFile == null) {
        log.error("putKey error commitlog {} uniqkey {}", req.getCommitLogOffset(), req.getUniqKey());
        return;
    }
}
```

继续看这俩方法

```java
private String buildKey(final String topic, final String key) {
    return topic + "#" + key;
}
```

以及这个

```java
private IndexFile putKey(IndexFile indexFile, DispatchRequest msg, String idxKey) {
    for (boolean ok = indexFile.putKey(idxKey, msg.getCommitLogOffset(), msg.getStoreTimestamp()); !ok; ) {
			// ......... 死循环的去写入，如果IndexFile不为空的话
    }
    return indexFile;
}
```

**可以看看它写入了哪些信息,这个比较核心**

```java
public boolean putKey(final String key, final long phyOffset, final long storeTimestamp) {
  	// 最大写入 maxIndexNum= 5000000个哈希槽*4
    if (this.indexHeader.getIndexCount() < this.indexNum) {
        int keyHash = indexKeyHashMethod(key);// abs hashcode 了一下
        int slotPos = keyHash % this.hashSlotNum; //获取hash槽
        int absSlotPos = IndexHeader.INDEX_HEADER_SIZE + slotPos * hashSlotSize; // 其实就是读取一下hash槽的索引位置，这个文件的文件结构是，header-hash索引-数据，所以这个是确定第三个位置，核心就是hashSlotNum是多少了，默认是maxHashSlotNum=5000000个哈希槽

        FileLock fileLock = null;

        try {

            // fileLock = this.fileChannel.lock(absSlotPos, hashSlotSize,
            // false);
          // 如果槽内没有数据或者槽内的数据比当前的递增值还要大就置空
            int slotValue = this.mappedByteBuffer.getInt(absSlotPos);
            if (slotValue <= invalidIndex || slotValue > this.indexHeader.getIndexCount()) {
                slotValue = invalidIndex;
            }

          // 这个其实就是获取 与index文件名的偏移量，因为文件名就有时间么
            long timeDiff = storeTimestamp - this.indexHeader.getBeginTimestamp();

            timeDiff = timeDiff / 1000;

            if (this.indexHeader.getBeginTimestamp() <= 0) {
                timeDiff = 0;
            } else if (timeDiff > Integer.MAX_VALUE) {
                timeDiff = Integer.MAX_VALUE;
            } else if (timeDiff < 0) {
                timeDiff = 0;
            }

          // 绝对位置= header+hash槽数量*每个槽的大小(4字节)+当前索引的自增数量*索引的大小(20字节)
            int absIndexPos =
                IndexHeader.INDEX_HEADER_SIZE + this.hashSlotNum * hashSlotSize
                    + this.indexHeader.getIndexCount() * indexSize;

            this.mappedByteBuffer.putInt(absIndexPos, keyHash);
            this.mappedByteBuffer.putLong(absIndexPos + 4, phyOffset);
            this.mappedByteBuffer.putInt(absIndexPos + 4 + 8, (int) timeDiff);
         	 // 这个写入slotValue 的值，实际上就是前继节点（当hash冲突的时候可以进行遍历，slot记录的是hash值一样的最后一个索引的值）
            this.mappedByteBuffer.putInt(absIndexPos + 4 + 8 + 4, slotValue);

          // 所以hash槽只需要方 当前索引的自增数量即可，用自增的好处是，我提前分配好hash槽的空间，后续只需要append，当大于最大hash槽，这个文件就无法写入了（这就是最一开始的判断）
            this.mappedByteBuffer.putInt(absSlotPos, this.indexHeader.getIndexCount());

            if (this.indexHeader.getIndexCount() <= 1) {
                this.indexHeader.setBeginPhyOffset(phyOffset);
                this.indexHeader.setBeginTimestamp(storeTimestamp);
            }
        		// 这个代码没有用处！！
            if (invalidIndex == slotValue) {
                this.indexHeader.incHashSlotCount();
            }
            //自增++
            this.indexHeader.incIndexCount();
            this.indexHeader.setEndPhyOffset(phyOffset);
            this.indexHeader.setEndTimestamp(storeTimestamp);
            return true;
        } catch (Exception e) {
            log.error("putKey exception, Key: " + key + " KeyHashCode: " + key.hashCode(), e);
        } finally {
            if (fileLock != null) {
                try {
                    fileLock.release();
                } catch (IOException e) {
                    log.error("Failed to release the lock", e);
                }
            }
        }
    } else {
        log.warn("Over index file capacity: index count = " + this.indexHeader.getIndexCount()
            + "; index max num = " + this.indexNum);
    }

    return false;
}
```

这个结构是一个 

![image-20210208102321898](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-8/c2f5c4dca4a2414dabc0c513aefd2c56.png)





![img](https://tyut.oss-accelerate.aliyuncs.com/image/2021/2-7/0dd8b841819d4454a90f4a7a35e10ec1.png)



## 3、延时队列

### 1、代码展示

消息只需要：

```go
msg := &primitive.Message{
  Topic: conf.Topic,
  Body:  []byte(time.Now().Format("2006-01-02 15:04:05")),
}
msg = msg.WithDelayTimeLevel(3)
```

消费的信息

```go
err = con.Subscribe(conf.Topic, consumer.MessageSelector{}, func(ctx context.Context,
  msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
  for i := range msgs {
    time.Sleep(time.Millisecond * 100)
    fmt.Printf("subscribe callback: QueueId:%v, QueueOffset:%v, message:%s, store_host: %v, cur_time: %v\n", msgs[i].Queue.QueueId, msgs[i].QueueOffset, msgs[i].Body, msgs[i].StoreHost, common.NowTimeString())
  }
  return consumer.ConsumeSuccess, nil
})
```

输出：基本可以保证level

```shell
subscribe callback: QueueId:0, QueueOffset:85, message:2021-02-13 16:10:15, store_host: 192.168.43.3:10916, cur_time: 2021-02-13 16:10:25
subscribe callback: QueueId:0, QueueOffset:85, message:2021-02-13 16:10:16, store_host: 192.168.43.3:10913, cur_time: 2021-02-13 16:10:26
```

生产的文件

```shell
root@288c93824863:~/store/consumequeue/SCHEDULE_TOPIC_XXXX# pwd
/root/store/consumequeue/SCHEDULE_TOPIC_XXXX
root@288c93824863:~/store/consumequeue/SCHEDULE_TOPIC_XXXX# ls -l
total 0
drwxr-xr-x 3 root root 96 Feb 13 08:07 2
```

### 2、实现原理

[https://cloud.tencent.com/developer/article/1581368](https://cloud.tencent.com/developer/article/1581368)

具体就是创建一个临时的Topic的消费队列，然后定期去检测，如果到期，才要放到指定的topic中和消费队列中。(这块可以理解为，我额外创建了一个Topic叫做 `SCHEDULE_TOPIC_XXXX`，然后呢只要是延时消息我就放到这个topic中，然后呢我就消费这个Topic，这个是一个定时器定期去消费，如果发现触达，我就投递到真正的Topic中)

细节就是，rocket-mq为了提高性能，并不支持任意的延时，因此它需要配置中指定延时队列的延时level: 其实就是提高读写性能

```properties
messageDelayLevel=1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
```

这可以理解为增加的读写的队列.

源码： 

1、[消息投递](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/CommitLog.java#L573)

```java
final int tranType = MessageSysFlag.getTransactionValue(msg.getSysFlag());
if (tranType == MessageSysFlag.TRANSACTION_NOT_TYPE
    || tranType == MessageSysFlag.TRANSACTION_COMMIT_TYPE) {
    // Delay Delivery
    if (msg.getDelayTimeLevel() > 0) {
        if (msg.getDelayTimeLevel() > this.defaultMessageStore.getScheduleMessageService().getMaxDelayLevel()) {
            msg.setDelayTimeLevel(this.defaultMessageStore.getScheduleMessageService().getMaxDelayLevel());
        }

        topic = TopicValidator.RMQ_SYS_SCHEDULE_TOPIC;
        queueId = ScheduleMessageService.delayLevel2QueueId(msg.getDelayTimeLevel());

        // Backup real topic, queueId
        MessageAccessor.putProperty(msg, MessageConst.PROPERTY_REAL_TOPIC, msg.getTopic());
        MessageAccessor.putProperty(msg, MessageConst.PROPERTY_REAL_QUEUE_ID, String.valueOf(msg.getQueueId()));
        // Properties
        msg.setPropertiesString(MessageDecoder.messageProperties2String(msg.getProperties()));

        // topic
        msg.setTopic(topic);

        // 队列id
        msg.setQueueId(queueId);
    }
}
```

2、[消息消费](https://github.com/apache/rocketmq/blob/release-4.8.0/store/src/main/java/org/apache/rocketmq/store/schedule/ScheduleMessageService.java#L262)

```java
for (; i < bufferCQ.getSize(); i += ConsumeQueue.CQ_STORE_UNIT_SIZE) {
// 20bit
// 消息的物理偏移量&&消息大小
long offsetPy = bufferCQ.getByteBuffer().getLong();
int sizePy = bufferCQ.getByteBuffer().getInt();
long tagsCode = bufferCQ.getByteBuffer().getLong();

if (cq.isExtAddr(tagsCode)) {
    if (cq.getExt(tagsCode, cqExtUnit)) {
        tagsCode = cqExtUnit.getTagsCode();
    } else {
        //can't find ext content.So re compute tags code.
        log.error("[BUG] can't find consume queue extend file content!addr={}, offsetPy={}, sizePy={}",
            tagsCode, offsetPy, sizePy);
        long msgStoreTime = defaultMessageStore.getCommitLog().pickupStoreTimestamp(offsetPy, sizePy);
        tagsCode = computeDeliverTimestamp(delayLevel, msgStoreTime);
    }
}

long now = System.currentTimeMillis();
long deliverTimestamp = this.correctDeliverTimestamp(now, tagsCode);

nextOffset = offset + (i / ConsumeQueue.CQ_STORE_UNIT_SIZE);

// 下发时间如果闭当前时间小，说明需要触达
long countdown = deliverTimestamp - now;

if (countdown <= 0) {
    // 消费消息，消费的是 `SCHEDULE_TOPIC_XXXX`
    MessageExt msgExt =
        ScheduleMessageService.this.defaultMessageStore.lookMessageByOffset(
            offsetPy, sizePy);

    if (msgExt != null) {
        try {
            // 获取真实消息
            MessageExtBrokerInner msgInner = this.messageTimeup(msgExt);
            if (TopicValidator.RMQ_SYS_TRANS_HALF_TOPIC.equals(msgInner.getTopic())) {
                log.error("[BUG] the real topic of schedule msg is {}, discard the msg. msg={}",
                        msgInner.getTopic(), msgInner);
                continue;
            }
            // 生产消息，落盘的commit-log中
            PutMessageResult putMessageResult =
                ScheduleMessageService.this.writeMessageStore
                    .putMessage(msgInner);

            if (putMessageResult != null
                && putMessageResult.getPutMessageStatus() == PutMessageStatus.PUT_OK) {
                continue;
            } else {
                // XXX: warn and notify me
              // 注意如果失败会投递失败，需要看日志报警！！！
                log.error(
                    "ScheduleMessageService, a message time up, but reput it failed, topic: {} msgId {}",
                    msgExt.getTopic(), msgExt.getMsgId());
                ScheduleMessageService.this.timer.schedule(
                    new DeliverDelayedMessageTimerTask(this.delayLevel,
                        nextOffset), DELAY_FOR_A_PERIOD);
                ScheduleMessageService.this.updateOffset(this.delayLevel,
                    nextOffset);
                return;
            }
        } catch (Exception e) {
						//.......
        }
    }
}
```