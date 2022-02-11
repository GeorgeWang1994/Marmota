package gg

import (
	"bytes"
	"fmt"
	"marmota/agentd/cc"
	"marmota/pkg/utils/file"
	"os/exec"
	"strings"
)

const Version = "0.0.1"

func GetCurrPluginVersion() string {
	if !cc.Config().Plugin.Enabled {
		return "plugin not enabled"
	}

	pluginDir := cc.Config().Plugin.Dir
	if !file.IsExist(pluginDir) {
		return "plugin dir not existent"
	}

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = pluginDir

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error:%s", err.Error())
	}

	return strings.TrimSpace(out.String())
}
