raft-service-proto:
	protoc \
	--go_out=raft/raftpc --go_opt=paths=source_relative \
	--go-grpc_out=raft/raftpc --go-grpc_opt=paths=source_relative \
    --proto_path=raft/raftpc/proto \
	raft_service.proto

raft-proto:
	protoc \
	--go_out=raft/raftpc --go_opt=paths=source_relative \
    --go-grpc_out=raft/raftpc --go-grpc_opt=paths=source_relative \
    --proto_path=raft/raftpc/proto \
	raft.proto

raftpc-proto: raft-service-proto raft-proto