FROM golang:1-alpine AS builder

ENV GO111MODULE=on
WORKDIR /go/release
ADD . .
RUN set -x \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o bepusdt ./main

FROM alpine:3.20

ENV TZ=Asia/Shanghai

# 安装所需的依赖
RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /go/release/bepusdt /runtime/bepusdt

# 设置时区
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

EXPOSE 8080
CMD ["/runtime/bepusdt"]
