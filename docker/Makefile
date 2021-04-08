.PHONY: run image run-signle stop delete clean down-rocketmq config help 

## mac默认走en0网卡
ETH0 := $(shell if [ "${uname}"=="Darwin" ]; then echo en0; else echo eth0;fi)
IP := $(shell ifconfig $(ETH0) | grep 'inet' | grep -v 'inet6' | awk '{print $$2}')
## 如果是mac 请装一下 gsed , brew install gnu-sed;
SED := $(shell if [ "${uname}"=="Darwin" ]; then echo gsed; else echo sed;fi)

all: run

image: ## 初始化镜像
	$(MAKE) -C image

run: ## 启动docker集群容器,默认是2broker(6node,2m4s),2nameserver,1manager
	docker-compose --compatibility up -d

run-signle: ## 启动docker集群容器,默认是1broker(3node,1m2s),2nameserver,1manager
	docker-compose --compatibility up -d nameserver-01 nameserver-02 rocketmq-console broker-01 broker-02 broker-03

stop: ## 关闭容器
	docker-compose --compatibility stop
delete: ## 删除容器
	docker-compose --compatibility down

clean:delete ## 清除本地volume的文件和容器
	$(RM) -r ./broker-01/store ./broker-01/rocketmqlogs \
	./broker-02/store ./broker-02/rocketmqlogs \
	./broker-03/store ./broker-03/rocketmqlogs \
	./broker-04/store ./broker-04/rocketmqlogs \
	./broker-05/store ./broker-05/rocketmqlogs \
	./broker-06/store ./broker-06/rocketmqlogs \
	./nameserver-01/rocketmqlogs \
	./nameserver-02/rocketmqlogs
down-rocketmq: ## 下载rocketmq
	if [ ! -f rocketmq-all-4.8.0-bin-release.zip ]; then wget https://archive.apache.org/dist/rocketmq/4.8.0/rocketmq-all-4.8.0-bin-release.zip; fi

config: ## 配置docker容器内rocket-mq的IP地址可以使得宿主机服务可以访问到
	@echo "当前主机IP $(IP)"
	## 匹配正则 ^brokerIP1=.*$, 替换 sed 's/^brokerIP1=.*$/brokerIP1=127.0.0.1/g
	for item in $(shell find . -name 'broker.conf' | grep -v 'nameserver'); do echo 替换文本$$item; $(SED) -i 's/^brokerIP1=.*$$/brokerIP1=$(IP)/g' $$item; done

help: ## 帮助
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)	