FROM golang:1.21.3 as build_sales-api
ENV CGO_ENABLED 0
ARG BUILD_REF

RUN mkdir /services

COPY . /services/
WORKDIR /services/app/services/sales-api

ENV GOOS=linux
RUN go build -o sales-api-e -ldflags "-X main.build=local"

FROM alpine:3.18
ARG BUILD_DATE
ARG BUILD_REF

RUN addgroup -g 1000 -S sales && \
    adduser -u 1000 -h /services -G sales -S sales

COPY --from=build_sales-api --chown=sales:sales /services/app/services/sales-api/ /service/sales-api

WORKDIR /service/sales-api
USER sales

RUN chmod +x ./sales-api-e

CMD ["./sales-api-e"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="sales-api" \
      org.opencontainers.image.authors="William Kennedy <bill@ardanlabs.com>" \
      org.opencontainers.image.source="https://github.com/ardanlabs/service/tree/master/app/services/sales-api" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Ardan Labs"
