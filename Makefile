PROTO_DIR=pb

protobuf: ${PROTO_DIR}/profile.proto
	protoc --go_out=. --go_opt=paths=source_relative ${PROTO_DIR}/profile.proto

.PHONY: protobuf