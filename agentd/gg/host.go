package gg

import (
	"marmota/agentd/cc"
	"os"
)

var hostname string

// Hostname todo: 这里可能存在重复hostname的问题
func Hostname() (string, error) {
	if hostname != "" {
		return hostname, nil
	}

	hostname = cc.Config().Hostname
	if hostname != "" {
		return hostname, nil
	}

	if os.Getenv("AGENT_NAME") != "" {
		hostname = os.Getenv("AGENT_NAME")
		return hostname, nil
	}

	var err error
	hostname, err = os.Hostname()
	if err != nil {
		return "", nil
	}
	return hostname, err
}
