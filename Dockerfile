FROM node:25.2.1 AS web_builder

# 安装 pnpm
RUN npm install -g pnpm

WORKDIR /web
COPY web/package.json web/pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile --shamefully-hoist

COPY web/ ./
RUN pnpm run build:prod

FROM golang:1.24.3-alpine3.20 AS builder

ENV GO111MODULE=on
WORKDIR /go/release
ADD . .

COPY --from=web_builder /web/dist ./static/secure

RUN set -x \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o bepusdt ./main

FROM alpine:3.20

ENV TZ=Asia/Shanghai

# 安装所需的依赖
RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /go/release/bepusdt /usr/local/bin/bepusdt

# 设置时区
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

EXPOSE 8080
ENTRYPOINT ["bepusdt"]
CMD ["start"]
