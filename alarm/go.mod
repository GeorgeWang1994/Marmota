module marmota/alarm

go 1.17

require (
	github.com/astaxie/beego v1.12.3
	github.com/garyburd/redigo v1.6.3
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.3.0
	github.com/toolkits/container v0.0.0-20151219225805-ba7d73adeaca
	github.com/toolkits/net v0.0.0-20160910085801-3f39ab6fe3ce
)

require (
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	marmota/pkg v0.0.1
)

replace marmota/pkg v0.0.1 => ../pkg
