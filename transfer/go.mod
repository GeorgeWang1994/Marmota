module marmota/transfer

go 1.17

require (
	github.com/hashicorp/net-rpc-msgpackrpc v0.0.0-20151116020338-a14192a58a69
	github.com/spf13/cobra v1.3.0
	github.com/toolkits/concurrent v0.0.0-20150624120057-a4371d70e3e3
	github.com/toolkits/consistent v0.0.0-20150827090850-a6f56a64d1b1
	github.com/toolkits/container v0.0.0-20151219225805-ba7d73adeaca
	github.com/toolkits/proc v0.0.0-20170520054645-8c734d0eb018
	marmota/pkg v0.0.1
)

require (
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/toolkits/conn_pool v0.0.0-20170512061817-2b758bec1177 // indirect
	github.com/toolkits/time v0.0.0-20160524122720-c274716e8d7f // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace marmota/pkg v0.0.1 => ../pkg
