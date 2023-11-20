PROTO_DIR=proto

proto: ${PROTO_DIR}/profile.proto
	protoc --go_out=. --go_opt=paths=source_relative ${PROTO_DIR}/profile.proto

.PHONY: proto