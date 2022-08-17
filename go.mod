module github.com/weaveworks-liquidmetal/fl

go 1.17

replace (
	github.com/weaveworks-liquidmetal/flintlock/api => ../flintlock/api
	github.com/weaveworks-liquidmetal/flintlock/client => ../flintlock/client
)

require (
	github.com/moby/moby v20.10.14+incompatible
	github.com/urfave/cli/v2 v2.4.0
	github.com/weaveworks-liquidmetal/flintlock/api v0.0.0-20220722132608-982d429ba641
	github.com/weaveworks-liquidmetal/flintlock/client v0.0.0-20220304105853-8fcb8aa2bafb
	github.com/yitsushi/macpot v1.0.2
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.45.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	google.golang.org/genproto v0.0.0-20220407144326-9054f6ed7bac // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
