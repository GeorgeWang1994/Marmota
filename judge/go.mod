module marmota/judge

go 1.17

require (
	github.com/garyburd/redigo v1.6.3
	github.com/hashicorp/net-rpc-msgpackrpc v0.0.0-20151116020338-a14192a58a69
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.3.0
)

require (
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	marmota/pkg v0.0.1
)

replace marmota/pkg v0.0.1 => ../pkg
