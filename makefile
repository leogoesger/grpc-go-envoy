run:
	go run app/api/main.go

readmsg:
	grpcurl -plaintext localhost:8000 chat.Chat.ReadMsg

sendmsg:
	grpcurl -plaintext -d '{"toUser": "Leo", "fromUser": "Noelle", "message": "are you well?"}' localhost:8000 chat.Chat.WriteMsg

streammsg:
	grpcurl -plaintext localhost:8000 chat.Chat.StreamLstMsg

list:
	grpcurl -plaintext localhost:8000 list

gen:
	cd business/chat; \
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	chat.proto

migrate:
	go run app/admin/main.go migrate

seed:
	go run app/admin/main.go seed

docker-build:
	cd envoy; \
	docker build -t grpc-go-envoy:1.0 .

docker-run:
	docker run -d -p 8000:8000 -p 9901:9901 grpc-go-envoy:1.0

tidy:
	go mod tidy