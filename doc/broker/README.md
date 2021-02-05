
这个是启动broker的启动参数，在目录的`/root/logs/rocketmqlogs/broker.log`下面
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
