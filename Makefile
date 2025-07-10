SRC_DIR=./schema
DST_DIR=./proto
BINARIES = server api client

.PHONY = build

build:
	@for bin in $(BINARIES); do \
		go build -o bin/$$bin ./cmd/$$bin; \
	done

gen-proto:
	protoc -I=${SRC_DIR} --go_out=${DST_DIR} --go_opt=paths=source_relative \
	--go-grpc_out=${DST_DIR} --go-grpc_opt=paths=source_relative \
	${SRC_DIR}/definition.proto