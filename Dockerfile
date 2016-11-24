FROM golang:1.6
MAINTAINER jang.jaehue@eland.co.kr

# install go package

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# add application
WORKDIR /go/src/best/p2-customer-service
ADD . /go/src/best/p2-customer-service

RUN go install

EXPOSE 9000

CMD ["/go/bin/p2-customer-service"]