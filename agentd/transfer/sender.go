package transfer

import (
	"bytes"
	"log"
	"marmota/agentd/cc"
	"marmota/agentd/gg"
	"marmota/agentd/rpc"
	"marmota/pkg/common/model"
	"math/rand"
	"strings"
	"time"
)

func updateMetrics(c *rpc.ConnRPCClient, metrics []*model.MetricValue, resp *model.TransferResponse) bool {
	err := c.Call("Transfer.Update", metrics, resp)
	if err != nil {
		log.Println("call Transfer.Update fail:", c, err)
		return false
	}
	return true
}

func SendMetrics(metrics []*model.MetricValue, resp *model.TransferResponse) {
	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(cc.Config().Transfer.Addrs)) {
		addr := cc.Config().Transfer.Addrs[i]

		c := gg.TransferClient(addr)
		if c == nil {
			c = gg.InitTransferClient(addr)
		}

		if updateMetrics(c, metrics, resp) {
			break
		}
	}
}

func SendToTransfer(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	dt := cc.Config().DefaultTags
	if len(dt) > 0 {
		var buf bytes.Buffer
		default_tags_list := []string{}
		for k, v := range dt {
			buf.Reset()
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(v)
			default_tags_list = append(default_tags_list, buf.String())
		}
		default_tags := strings.Join(default_tags_list, ",")

		for i, x := range metrics {
			buf.Reset()
			if x.Tags == "" {
				metrics[i].Tags = default_tags
			} else {
				buf.WriteString(metrics[i].Tags)
				buf.WriteString(",")
				buf.WriteString(default_tags)
				metrics[i].Tags = buf.String()
			}
		}
	}

	debug := cc.Config().Debug

	if debug {
		log.Printf("=> <Total=%d> %v\n", len(metrics), metrics[0])
	}

	var resp model.TransferResponse
	SendMetrics(metrics, &resp)

	if debug {
		log.Println("<=", &resp)
	}
}
