.PHONY: proto-generate

# Generates Go code from all .proto files using the 'buf' tool tracked in go.mod.
proto-generate:
	@echo "Generating protobuf and gRPC code with buf..."
	# Using 'go run' ensures we use the specific version of buf defined in our
	# go.mod, making the build reproducible and independent of global installs.
	go run github.com/bufbuild/buf/cmd/buf generate
	@echo "Done."
