package plugins

import (
	"io/ioutil"
	"log"
	"marmota/agentd/cc"
	"path/filepath"
	"strconv"
	"strings"
)

// ListPlugins return: dict{sys/ntp/60_ntp.py : *Plugin}
func ListPlugins(script_path string) map[string]*Plugin {
	ret := make(map[string]*Plugin)
	if script_path == "" {
		return ret
	}

	abs_path := filepath.Join(cc.Config().Plugin.Dir, script_path)
	fs, err := ioutil.ReadDir(abs_path)
	if err != nil {
		log.Println("can not list files under", abs_path)
		return ret
	}

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		filename := f.Name()
		arr := strings.Split(filename, "_")
		if len(arr) < 2 {
			continue
		}

		// filename should be: $cycle_$xx
		var cycle int
		cycle, err = strconv.Atoi(arr[0])
		if err != nil {
			continue
		}

		fpath := filepath.Join(script_path, filename)
		plugin := &Plugin{FilePath: fpath, MTime: f.ModTime().Unix(), Cycle: cycle, Args: ""}
		ret[fpath] = plugin
	}
	return ret
}
