module firefly-jam-2026

go 1.25.5

require github.com/firefly-zero/firefly-go v0.10.0

require (
	github.com/orsinium-labs/tinymath v1.0.0
	github.com/planetscale/vtprotobuf v0.6.0
	google.golang.org/protobuf v1.36.11
)

require github.com/aperturerobotics/protobuf-go-lite v0.11.0 // indirect

tool (
	github.com/aperturerobotics/protobuf-go-lite/cmd/protoc-gen-go-lite
	github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto
	google.golang.org/protobuf/cmd/protoc-gen-go
)
