run-api:
	cd grpc-api; go run main.go

run-client:
	cd grpc-client; npm start

proto-gen-api:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	grpc-api/protos/sensor.proto; \
	mv grpc-api/protos/sensor_grpc.pb.go grpc-api/server/sensorpb/sensor_grpc.pb.go; \
	mv grpc-api/protos/sensor.pb.go grpc-api/server/sensorpb/sensor.pb.go

proto-gen-client:
	protoc --js_out=import_style=commonjs,binary:./grpc-client/src/sensorpb \
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./grpc-client/src/sensorpb \
		grpc-api/protos/sensor.proto 

tidy:
	cd grpc-api; go mod tidy

build-docker: 
	docker build -t grpc-medium-envoy:1.0 ./envoy

run-docker:
	docker run --rm -p 9901:9901 -p 10000:10000 grpc-medium-envoy:1.0 

npm-install:
	cd grpc-client; npm install