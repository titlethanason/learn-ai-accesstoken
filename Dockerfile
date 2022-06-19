# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY local.env ./

RUN go build -o /learn-ai-accesstoken

EXPOSE 8888

CMD [ "/learn-ai-accesstoken" ]