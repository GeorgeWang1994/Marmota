package cron

import (
	"fmt"
	"log"
	"marmota/agentd/cc"
	"marmota/agentd/gg"
	"marmota/pkg/common/model"
	"time"
)

func ReportStatus() {
	if cc.Config().Heartbeat.Enabled && cc.Config().Heartbeat.Addr != "" {
		go reportAgentStatus(time.Duration(cc.Config().Heartbeat.Interval) * time.Second)
	}
}

// 向hbs同步探针状态
func reportAgentStatus(interval time.Duration) {
	for {
		hostname, err := gg.Hostname()
		if err != nil {
			hostname = fmt.Sprintf("error:%s", err.Error())
		}

		req := model.AgentReportRequest{
			Hostname:      hostname,
			IP:            gg.IP(),
			AgentVersion:  gg.Version,
			PluginVersion: gg.GetCurrPluginVersion(),
		}

		var resp model.RpcResponse
		err = gg.HBSRpcClient().Call("Agent.ReportStatus", req, &resp)
		if err != nil || resp.Code != 200 {
			log.Println("call Agent.ReportStatus fail:", err, "Request:", req, "Response:", resp)
		}

		time.Sleep(interval)
	}
}
