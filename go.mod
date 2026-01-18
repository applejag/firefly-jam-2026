module github.com/applejag/firefly-jam-2026

go 1.25.5

require (
	github.com/aperturerobotics/protobuf-go-lite v0.11.0
	github.com/firefly-zero/firefly-go v0.10.0
	github.com/orsinium-labs/tinymath v1.0.0
)

require google.golang.org/protobuf v1.36.11 // indirect

tool (
	github.com/aperturerobotics/protobuf-go-lite/cmd/protoc-gen-go-lite
	google.golang.org/protobuf/cmd/protoc-gen-go
)
