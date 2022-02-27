#syntax=docker/dockerfile:1

FROM golang:1.17.7 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /glowMarkt
 
FROM alpine:3.14.3
WORKDIR /
COPY --from=build /glowMarkt /glowMarkt
ENTRYPOINT [ "/glowMarkt" ]