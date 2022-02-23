module marmota/store

go 1.17

require (
	github.com/niean/goperfcounter v0.0.0-20160108100052-24860a8d3fac
	github.com/open-falcon/rrdlite v0.0.0-20200214140804-bf5829f786ad
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.3.0
	github.com/toolkits/cache v0.0.0-20190218093630-cfb07b7585e5
	github.com/toolkits/concurrent v0.0.0-20150624120057-a4371d70e3e3
	github.com/toolkits/consistent v0.0.0-20150827090850-a6f56a64d1b1
	github.com/toolkits/container v0.0.0-20151219225805-ba7d73adeaca
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	github.com/toolkits/proc v0.0.0-20170520054645-8c734d0eb018
	github.com/toolkits/time v0.0.0-20160524122720-c274716e8d7f
)

require (
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/niean/go-metrics-lite v0.0.0-20151230091537-b5d30971b578 // indirect
	github.com/niean/gotools v0.0.0-20151221085310-ff3f51fc5c60 // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	marmota/pkg v0.0.1
)

replace marmota/pkg v0.0.1 => ../pkg
