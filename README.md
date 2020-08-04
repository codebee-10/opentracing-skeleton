# opentracing-skeleton
opentracing go、java wrapper 


### quick start

1. run grpc
```
go run src/opentrace/grpc/helloworld/server.go
```

2. run http server
```
go run src/opentrace/trace_server.go 
```

3. run http client
```
go run src/opentrace/trace_client.go 
```


### run Jaeger 

```
docker run \
-p 5775:5775/udp \
-p 16686:16686 \
-p 6831:6831/udp \
-p 6832:6832/udp \
-p 5778:5778 \
-p 14268:14268 \
jaegertracing/all-in-one:latest
```

### init trace
```

```

### spanContext 传递


### reporter 定义


### http.Handler


### rpc.Handler