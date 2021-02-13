FROM openjdk:8-jdk as jdk

RUN apt-get update &&\
    apt-get install -y vim

WORKDIR /opt
RUN wget https://archive.apache.org/dist/rocketmq/4.8.0/rocketmq-all-4.8.0-bin-release.zip\
    && unzip rocketmq-all-4.8.0-bin-release.zip \
    && mv rocketmq-all-4.8.0-bin-release rocketmq \
    && rm -rf /opt/rocketmq/bin

WORKDIR /opt/rocketmq

ADD bin ./bin
RUN mkdir -p /root/store \
    && mkdir -p /root/store/commitlog
ENV PATH=${PATH}:/opt/rocketmq/bin

EXPOSE 9876 10911

CMD [ "/bin/bash","-c","bin/dledger/fast-try.sh start"]