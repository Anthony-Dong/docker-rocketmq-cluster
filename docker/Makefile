.PHONY: run stop delete clean down-rocketmq

run:
	docker-compose --compatibility up -d

run-signle:
	docker-compose --compatibility up -d nameserver-01 nameserver-02 rocketmq-console broker-01 broker-02 broker-03

stop:
	docker-compose --compatibility stop
delete:
	docker-compose --compatibility down
clean:delete
	$(RM) -r ./broker-01/store ./broker-01/rocketmqlogs \
	./broker-02/store ./broker-02/rocketmqlogs \
	./broker-03/store ./broker-03/rocketmqlogs \
	./broker-04/store ./broker-04/rocketmqlogs \
	./broker-05/store ./broker-05/rocketmqlogs \
	./broker-06/store ./broker-06/rocketmqlogs \
	./nameserver-01/rocketmqlogs \
	./nameserver-02/rocketmqlogs
down-rocketmq:
	if [ ! -f rocketmq-all-4.8.0-bin-release.zip ]; then wget https://archive.apache.org/dist/rocketmq/4.8.0/rocketmq-all-4.8.0-bin-release.zip; fi
	
