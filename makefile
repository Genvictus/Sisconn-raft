PROTO_PATH="raft/raftpc/proto"

go-test:
	go test ./raft ./raft/transport -coverprofile out/coverage
	go tool cover -html out/coverage -o out/coverage.html

raft-service-proto:
	protoc \
	--go_out=raft/raftpc --go_opt=paths=source_relative \
	--go-grpc_out=raft/raftpc --go-grpc_opt=paths=source_relative \
    --proto_path=${PROTO_PATH} \
	raft_service.proto

raft-proto:
	protoc \
	--go_out=raft/raftpc --go_opt=paths=source_relative \
    --go-grpc_out=raft/raftpc --go-grpc_opt=paths=source_relative \
    --proto_path=${PROTO_PATH} \
	raft.proto

raftpc-proto: raft-service-proto raft-proto

raftpc-proto-client-web: raft-proto-client-web raft-service-proto-client-web

client:
	go build Sisconn-raft/cmd/client

server:
	go build Sisconn-raft/cmd/server
