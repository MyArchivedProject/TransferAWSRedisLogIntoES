FROM golang:1.15.5-alpine as builder
WORKDIR /opt/data1
ADD . /opt/data1
RUN go build -ldflags "-s -w" -o goaws && ls 

FROM alpine:3.12 as deployer
WORKDIR /opt/data1
COPY --from=builder /opt/data1/goaws /opt/data1/
RUN chmod +x goaws
CMD ["./goaws"]
