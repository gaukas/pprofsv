PROTO_DIR=pb

test: 
	go test ./...

protobuf: ${PROTO_DIR}/profile.proto
	protoc --go_out=. --go_opt=paths=source_relative ${PROTO_DIR}/profile.proto

.PHONY: protobuf test 
