run:
	go run app/api/main.go

readmsg:
	grpcurl -plaintext -d '{"name": "hello"}' localhost:50051 chat.Chat.ReadMsg

chat-gen:
	cd business/chat; \
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	chat.proto

migrate:
	go run app/admin/main.go migrate

seed:
	go run app/admin/main.go seed