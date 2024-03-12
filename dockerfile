FROM golang:1.19-alpine3.16 as builder

LABEL name="synapsis-test"
LABEL version="1.0.0"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary

ENTRYPOINT ["/app/binary"]