FROM golang:1.21.3 as build_sales-api
ENV CGO_ENABLED 0
ARG BUILD_REF

RUN mkdir /services

# COPY go.* /services/ , 这行命令为什么就不行了呢... 真是奇怪了，奇怪
COPY . /service/
WORKDIR /service/app/services/sales-api

# RUN go mod download
RUN go build -o sales-api -ldflags "-X main.build=local"

# Run the Go Binary in Alpine.
FROM alpine:3.18
ARG BUILD_DATE
ARG BUILD_REF

# 少了这个，会导致 CreateContainerError , 找不到用户
RUN addgroup -g 1000 -S sales && \
    adduser -u 1000 -h /services -G sales -S sales

# COPY --from=build_sales-api --chown=sales:sales /services /services/
COPY --from=build_sales-api --chown=sales:sales /service/app/services/sales-api/ /service/sales-api

WORKDIR /service/sales-api
USER sales

# Docker 启动二进制文件的命令，必须要和 go.mod 的 moudles 名字进行对应起来。
# module github.com/yourusername/basic-a， 这个的后缀，打包出来后，是 basic-a 的二进制文件～
RUN chmod +x ./sales-api

# 这个二进制文件的名称，依赖与哪个文件和哪些命令？
CMD ["./sales-api"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="sales-api" \
      org.opencontainers.image.authors="William Kennedy <bill@ardanlabs.com>" \
      org.opencontainers.image.source="https://github.com/ardanlabs/service/tree/master/app/services/sales-api" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Ardan Labs"


