## 1、RocketMQ又双叒叕system busy了，怎么破？

异常日志：

```shell
2021/03/01 20:13:15 producer_api.go:58: [ERROR] [trace_id=s:7fbdb2a87a8711ebac9bacde48001122] [RocketMq-Producer] sync send err, err: [TIMEOUT_CLEAN_QUEUE]broker busy, start flow control for a while, period in queue: 201ms, size of queue: 125
2021/03/01 20:13:32 producer_api.go:58: [ERROR] [trace_id=s:8a2257bc7a8711ebac9bacde48001122] [RocketMq-Producer] sync send err, err: [TIMEOUT_CLEAN_QUEUE]broker busy, start flow control for a while, period in queue: 202ms, size of queue: 133
```

如何FIX

[https://cloud.tencent.com/developer/article/1451310](https://cloud.tencent.com/developer/article/1451310)

## 2、

