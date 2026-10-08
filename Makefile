LOCAL_BIN:=$(CURDIR)/bin
DOCKER_USER ?= user99430e

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

generate: generate-note-api

generate-note-api:
	mkdir -p pkg/note_v1
	protoc --proto_path api/note_v1 \
	--go_out=pkg/note_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/note_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/note_v1/note.proto

server-up:
	go run cmd/grpc_server/main.go
	
client:
	go run cmd/grpc_client/main.go	


build-a:
	GOOS=linux GOARCH=amd64 go build -o bin/grpc_server ./cmd/grpc_server

docker-push:
	docker login
	docker build -t $(DOCKER_USER)/test-go-chat-server:latest .
	docker push $(DOCKER_USER)/test-go-chat-server:latest

#docker-push:
#	docker push $(DOCKER_USER)/go-chat-server:latest

to-server:
	scp ./bin/grpc_server user200@83.222.26.90:/var/www/test_server
	
