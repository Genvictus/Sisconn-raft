raft-service-proto:
	protoc \
	--go_out=raft/transport --go_opt=paths=source_relative \
	--go-grpc_out=raft/transport --go-grpc_opt=paths=source_relative \
    --proto_path=raft/transport/proto \
	raft_service.proto

raft-proto:
	protoc \
	--go_out=raft/transport --go_opt=paths=source_relative \
    --go-grpc_out=raft/transport --go-grpc_opt=paths=source_relative \
    --proto_path=raft/transport/proto \
	raft.proto
