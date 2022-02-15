package sender

import (
	"marmota/transfer/sender/node"
	"marmota/transfer/sender/pool"
	"marmota/transfer/sender/queue"
)

func Start() {
	node.InitNodeRings()
	queue.InitSendQueues()
	pool.InitConnPools()
}
