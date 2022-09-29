FROM golang:1-alpine

ADD main.go /
WORKDIR /
RUN go build -gcflags '-N -l' -o httpserver .

ENTRYPOINT ["/httpserver"]