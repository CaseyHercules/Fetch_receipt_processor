FROM golang:1.20.2

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy