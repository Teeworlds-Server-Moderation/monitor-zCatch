FROM golang:alpine as build

LABEL maintainer "github.com/jxsl13"

WORKDIR /build
COPY *.go ./
COPY go.* ./

ENV CGO_ENABLED=0
ENV GOOS=linux 

RUN go get -d && go build -a -ldflags '-w -extldflags "-static"' -o monitor-zCatch .


FROM alpine:latest as minimal

ENV MONITOR_BROKER_ADDRESS=tcp://mosquitto:1883
ENV MONITOR_ECON_ADDRESS=localhost:9303
ENV MONITOR_ECON_PASSWORD=""



WORKDIR /app
COPY --from=build /build/monitor-zCatch .
VOLUME ["/data"]
ENTRYPOINT ["/app/monitor-zCatch"]