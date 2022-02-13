package gg

import "marmota/agentd/cc"

var LocalIp string

func IP() string {
	if len(LocalIp) > 0 {
		return LocalIp
	}

	LocalIp = cc.Config().IP
	if LocalIp != "" {
		// use ip in pivas
		return LocalIp
	}
	
	return LocalIp
}
