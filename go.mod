module github.com/myriadeinc/zircon

go 1.16

// replace github.com/myriadeinc/zircon/internal => ./internal

// replace github.com/myriadeinc/zircon/xmrlib => ./xmrlib

require (
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/go-redis/redis/v8 v8.11.0 // indirect
	github.com/myriadeinc/zircon_proto v0.0.0-20210704180748-04f6be9f590b // indirect
	github.com/rs/zerolog v1.26.0
	github.com/ybbus/jsonrpc v2.1.2+incompatible
	google.golang.org/grpc v1.38.0 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)
