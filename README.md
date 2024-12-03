# Streaming Demo

- https://kafka.apache.org/quickstart

## System Info

```
docker pull apache/kafka:3.9.0
docker run -p 9092:9092 apache/kafka:3.9.0

go run producers/sysinfo/sysinfo.go

go run consumers/sysinfo/sysinfo.go
```
