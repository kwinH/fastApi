# stage 1: build src code to binary
FROM gobase as builder

ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on

COPY ./ /app/
RUN cd /app && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main .

# stage 2: use alpine as base image
FROM alpine:3.10
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add tzdata ca-certificates && \
    cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    # apk del tzdata && \
    rm -rf /var/cache/apk/*


COPY --from=builder /app/main /go/main
COPY ./config/settings.yml /go/config/settings.yml

WORKDIR /go
EXPOSE 8000

CMD ["/go/main","server","-c", "/go/config/settings.yml"]