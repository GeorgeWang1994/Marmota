package sender

import (
	"marmota/transfer/sender/node"
	"marmota/transfer/sender/pool"
	"marmota/transfer/sender/queue"
	"marmota/transfer/sender/task"
)

func Start() {
	node.InitNodeRings()
	queue.InitSendQueues()
	pool.InitConnPools()
	task.StartSendTasks()
}
