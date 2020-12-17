## How to Run

1. `make npm-install`
2. `make build-docker`
3. `make run-client`
4. `make run-docker`
5. `make run-api`

## Setup

1. Download the compiler [https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)
2. Download Go plugins

```
$ export GO111MODULE=on  # Enable module mode
$ go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

3. Copy compiler to your PATH, I put it under $GOPATH/bin
4. After write the `proto` file, you will get both `_grpc.pb.go` and `.pb.go`. I wrote it into `makefile`. They will be saved in `grpc-api/sserver/sensorpb/`

Update `makefile`  
```makefile
proto-gen-api:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	grpc-api/protos/sensor.proto; \
	mv grpc-api/protos/sensor_grpc.pb.go grpc-api/server/sensorpb/sensor_grpc.pb.go; \
	mv grpc-api/protos/sensor.pb.go grpc-api/server/sensorpb/sensor.pb.go
```

```
make proto-gen-api
```

5. Follow the [blog](https://medium.com/swlh/building-a-realtime-dashboard-with-reactjs-go-grpc-and-envoy-7be155dfabfb), `server struct` needs to implement `mustEmbedUnimplementedSensorServer`. 

```go
type server struct {
	sensorpb.UnimplementedSensorServer
}
```

6. Onto client side, download plugin [https://github.com/grpc/grpc-web/releases](https://github.com/grpc/grpc-web/releases)

7. Move it under $GOPATH/bin, change its name and give execute permission

```
mv ~/Downloads/protoc-gen-grpc-web-1.x.x-x-x86_xx $GOPATH/bin/protoc-gen-grpc-web
chmod +x protoc-gen-grpc-web
```

8. Generate client files, they will be saved under `grpc-client/src/sensorpb/`

Update `makefile`  
```makefile
proto-gen-client:
	protoc --js_out=import_style=commonjs,binary:./grpc-client/src/sensorpb \
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./grpc-client/src/sensorpb \
		grpc-api/protos/sensor.proto 
```

```
make proto-gen-client
```

9. Under `root` create `envoy/envoy.yaml`. The `yaml` file from the blog didn't work, probably indentation is off. I think this is not working yet.

```yaml
admin:
  access_log_path: /tmp/admin_access.log
  address:
    socket_address: { address: 0.0.0.0, port_value: 9901 }

static_resources:
  listeners:
  - name: listener_0
    address:
      socket_address: { address: 0.0.0.0, port_value: 8080 }
    filter_chains:
    - filters:
      - name: envoy.filters.network.http_connection_manager
        typed_config:
          "@type": type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager
          codec_type: auto
          stat_prefix: ingress_http
          route_config:
            name: local_route
            virtual_hosts:
            - name: local_service
              domains: ["*"]
              routes:
              - match: { prefix: "/" }
                route:
                  cluster: greeter_service
                  max_grpc_timeout: 0s
              cors:
                allow_origin_string_match:
                - prefix: "*"
                allow_methods: GET, PUT, DELETE, POST, OPTIONS
                allow_headers: keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,custom-header-1,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout
                max_age: "1728000"
                expose_headers: custom-header-1,grpc-status,grpc-message
          http_filters:
          - name: envoy.filters.http.grpc_web
          - name: envoy.filters.http.cors
          - name: envoy.filters.http.router
  clusters:
  - name: greeter_service
    connect_timeout: 0.25s
    type: logical_dns
    http2_protocol_options: {}
    lb_policy: round_robin
    # win/mac hosts: Use address: host.docker.internal instead of address: localhost in the line below
    load_assignment:
      cluster_name: cluster_0
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: 0.0.0.0
                    port_value: 9090
```

10. npm install on client directory

```
npm install grpc-web google-protobuf --save
```



