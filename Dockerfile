FROM golang:1.25 AS builder

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y && mkdir -p /configs/

COPY --from=builder /src/bin /app

WORKDIR /app
RUN mv xiaomiao-home-system server

EXPOSE 8080
EXPOSE 8090

CMD ["./server", "-conf", "/configs"]
