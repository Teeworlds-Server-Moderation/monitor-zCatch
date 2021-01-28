FROM golang:alpine as build

LABEL maintainer "github.com/jxsl13"

WORKDIR /build
COPY *.go ./
COPY go.* ./

ENV CGO_ENABLED=0
ENV GOOS=linux 

RUN go get -d && go build -a -ldflags '-w -extldflags "-static"' -o monitor-zCatch .


FROM alpine:latest as minimal

ENV BROKER_ADDRESS=tcp://mosquitto:1883
ENV CLIENT_ID=monitor-zCatch


WORKDIR /app
COPY --from=build /build/monitor-zCatch .
VOLUME ["/data"]
ENTRYPOINT ["/app/monitor-zCatch"]