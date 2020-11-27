FROM golang:1.13.7-alpine3.11 as builder

ENV GO111MODULE on
COPY download_kubeedge.sh /
COPY . /go/src/apulis/AIArtsBackend
RUN umask 022
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN mkdir -p /etc/kubeedge
RUN apk update
RUN apk --no-cache add curl git build-base
RUN chmod +x /download_kubeedge.sh;
RUN /download_kubeedge.sh arm64
RUN /download_kubeedge.sh amd64

RUN git clone git@apulis-gitlab.apulis.cn:apulis/apulisedge.git /
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN cd /; swag init; GO111MODULE=${GO111MODULE} go build -o /go/bin/ApulisEdge;


FROM alpine:3.11 as server
RUN apk --no-cache add ca-certificates libdrm
WORKDIR /root/
COPY --from=0 /go/bin/ApulisEdge .

CMD ["./ApulisEdge"]