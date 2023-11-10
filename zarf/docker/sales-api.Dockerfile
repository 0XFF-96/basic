FROM golang:1.21.3 as build_sales-api
ENV CGO_ENABLED 0
ARG BUILD_REF

RUN mkdir /service

COPY . /service/
WORKDIR /service/app/services/sales-api

ENV GOOS=linux
RUN go build -o sales-api-e -ldflags "-X main.build=local"

FROM alpine:3.18
ARG BUILD_DATE
ARG BUILD_REF

RUN addgroup -g 1000 -S sales && \
    adduser -u 1000 -h /service -G sales -S sales


# 存储 key 的相关 folder ~
COPY --from=build_sales-api --chown=sales:sales /service/zarf/keys/. /service/zarf/keys/.
COPY --from=build_sales-api --chown=sales:sales /service/app/services/sales-api/ /service/sales-api

WORKDIR /service
USER sales

RUN chmod +x ./sales-api/sales-api-e

CMD ["./sales-api/sales-api-e"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="sales-api" \
      org.opencontainers.image.authors="William Kennedy <bill@ardanlabs.com>" \
      org.opencontainers.image.source="https://github.com/ardanlabs/service/tree/master/app/services/sales-api" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Ardan Labs"
