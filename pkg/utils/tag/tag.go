package tag

import (
	"bytes"
	"marmota/pkg/utils/bufferPool"
	"sort"
	"strings"
)

func SortedTags(tags map[string]string) string {
	if tags == nil {
		return ""
	}

	size := len(tags)

	if size == 0 {
		return ""
	}

	ret := bufferPool.BufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.BufferPool.Put(ret)

	if size == 1 {
		for k, v := range tags {
			ret.WriteString(k)
			ret.WriteString("=")
			ret.WriteString(v)
		}
		return ret.String()
	}

	keys := make([]string, size)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for j, key := range keys {
		ret.WriteString(key)
		ret.WriteString("=")
		ret.WriteString(tags[key])
		if j != size-1 {
			ret.WriteString(",")
		}
	}

	return ret.String()
}

func TagsDict(s string) map[string]string {
	if s == "" {
		return map[string]string{}
	}

	if strings.ContainsRune(s, ' ') {
		s = strings.Replace(s, " ", "", -1)
	}

	tagDict := make(map[string]string)
	tags := strings.Split(s, ",")
	for _, tag := range tags {
		idx := strings.IndexRune(tag, '=')
		if idx != -1 {
			tagDict[tag[:idx]] = tag[idx+1:]
		}
	}
	return tagDict
}
