package consume

import (
	"marmota/pkg/common/model"
)

func BuildCommonIMContent(event *model.Event) string {
	return ""
}

func GenerateIMContent(event *model.Event) string {
	return BuildCommonIMContent(event)
}
