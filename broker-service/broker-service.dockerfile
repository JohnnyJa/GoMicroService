#base go image
#FROM golang:1.22-alpine as builder
#
#RUN mkdir /app
#
#COPY . /app
#
#WORKDIR /app
#
#RUN CGO_ENABLED=0 go build -o broker-app ./cmd/api
#
#RUN chmod +x /app/broker-app

#Build a tiny broker image
FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

CMD ["/app/brokerApp"]