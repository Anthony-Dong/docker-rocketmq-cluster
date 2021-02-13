package main

import (
	"clinet/common"
	"clinet/config"
	"context"
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"strconv"
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
	conf := config.RocketMqProducer{
		RocketMqConfig: config.RocketMqConfig{
			Host:       []string{"127.0.0.1:9871", "127.0.0.1:9872"},
			RetryTimes: 3,
			GroupName:  gname,
		},
		SendMsgTimeout: 3 * time.Second,
		Topic:          topic,
	}
	p, err := rocketmq.NewProducer(
		producer.WithSendMsgTimeout(conf.SendMsgTimeout),
		producer.WithGroupName(conf.GroupName),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: conf.AccessKey,
			SecretKey: conf.SecretKey,
		}),
		producer.WithNameServer(conf.Host),
		producer.WithRetry(conf.RetryTimes),
	)
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	defer func() {
		common.EchoError(p.Shutdown())
	}()
	count := 0
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		<-ticker.C
		msg := &primitive.Message{
			Topic: conf.Topic,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(count)),
		}
		res, err := p.SendSync(context.Background(), msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
		count++
	}
}
