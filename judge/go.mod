module marmota/judge

go 1.17

require (
	github.com/garyburd/redigo v1.6.3
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.3.0
)

require gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	marmota/pkg v0.0.1
)

replace marmota/pkg v0.0.1 => ../pkg
