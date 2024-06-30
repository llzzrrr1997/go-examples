#第一次构建镜像
#FROM golang:1.17
#
#ENV GO111MODULE=on \
#    GOPROXY=https://goproxy.cn,direct
#
#WORKDIR /service
#
#COPY . .
#
#RUN go build -o app main.go
#
#EXPOSE 8080
#
#CMD ["./app"]

#第二次构建镜像
#FROM alpine
#
#WORKDIR /service
#
#COPY app .
#
#EXPOSE 8080
#
#CMD ["./app"]

#第三次构建 多阶段构建
FROM golang:1.17 AS builder

LABEL maintainer="564952475@qq.com"

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /service

COPY . .

RUN CGO_ENABLED=0 GOOS=linux  go build -o app main.go


FROM alpine

WORKDIR /service

COPY --from=builder /service/app /service/app
COPY main.go /service

EXPOSE 8080

CMD ["./app"]