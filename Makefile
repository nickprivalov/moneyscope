.PHONY: proto

proto:
	@echo "Generating Go code from Protobufs..."
	@mkdir -p backend/gen/go
	protoc --proto_path=backend/proto \
		--go_out=backend/gen/go --go_opt=paths=source_relative \
		--go-grpc_out=backend/gen/go --go-grpc_opt=paths=source_relative \
		backend/proto/common/v1/*.proto \
		backend/proto/ledger/v1/*.proto
	@echo "Done!"
