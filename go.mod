module ecsDemo

go 1.15

require (
	github.com/golang/protobuf v1.5.2
	github.com/gonutz/prototype v1.0.7
	github.com/spf13/cast v1.4.1
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	github.com/wwj31/dogactor v1.0.3
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.26.0
)

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
