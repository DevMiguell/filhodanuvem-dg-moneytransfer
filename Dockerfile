FROM golang:1.18.2-alpine3.16 as base
RUN apk update 
WORKDIR /src/moneytransfer
ADD . . 
RUN go mod download
RUN go build -o moneytransfer ./cmd/api

FROM alpine:3.16 as binary
WORKDIR /src/app
COPY --from=base /src/moneytransfer/moneytransfer .
EXPOSE 3000
CMD ["/src/app/moneytransfer"]