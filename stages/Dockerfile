FROM golang:1.15.5-alpine as builder
WORKDIR /opt/data1
ADD . /opt/data1
go build -ldflags "-s -w" -o goaws

FROM alpine: 3.12
WORKDIR /opt/data1
COPY --from=builder /opt/data1/* /opt/data1/*
RUN chmod +x go-aws
CMD ["./goaws"]
