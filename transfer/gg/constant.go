package gg

const (
	GAUGE       = "GAUGE"
	COUNTER     = "COUNTER"
	DERIVE      = "DERIVE"
	DefaultStep = 60
)

const (
	DefaultSendQueueMaxSize = 102400 //10.24w
	DefaultMinStep          = 30     //最小上报周期,单位sec
)
