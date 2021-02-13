package main

import (
	"clinet/common"
	"clinet/config"
	"context"
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"os"
	"time"
)

var (
	topic string
	gname string
)

func init() {
	flag.StringVar(&topic, "t", "", "topic-name")
	flag.StringVar(&gname, "g", "go_client_dev", "topic-name")
}

func main() {
	flag.Parse()

	conf := config.RocketMqConsumer{
		RocketMqConfig: config.RocketMqConfig{
			Host:       []string{"127.0.0.1:9871", "127.0.0.1:9872"},
			RetryTimes: 3,
			GroupName:  gname,
		},
		Topic: "DelayTopic-1",
	}
	con, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.GroupName),
		consumer.WithNameServer(conf.Host),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: conf.AccessKey,
			SecretKey: conf.SecretKey,
		}),
		consumer.WithPullBatchSize(10),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = con.Shutdown()
		if err != nil {
			common.EchoError(err)
		}
	}()
	err = con.Subscribe(conf.Topic, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			time.Sleep(time.Millisecond * 100)
			fmt.Printf("subscribe callback: QueueId:%v, QueueOffset:%v, message:%s, store_host: %v, cur_time: %v\n", msgs[i].Queue.QueueId, msgs[i].QueueOffset, msgs[i].Body, msgs[i].StoreHost, common.NowTimeString())
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		common.EchoError(err)
	}
	err = con.Start()
	if err != nil {
		common.EchoError(err)
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
}
