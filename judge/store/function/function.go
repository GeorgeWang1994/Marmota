package function

import (
	"fmt"
	"marmota/judge/store"
	"marmota/pkg/common/model"
	"math"
	"strconv"
	"strings"
)

type Function interface {
	Compute(L *store.SafeLinkedList) (vs []*model.HistoryData, leftValue float64, isTriggered bool, isEnough bool)
}

func checkIsTriggered(leftValue float64, operator string, rightValue float64) (isTriggered bool) {
	switch operator {
	case "=", "==":
		isTriggered = math.Abs(leftValue-rightValue) < 0.0001
	case "!=":
		isTriggered = math.Abs(leftValue-rightValue) > 0.0001
	case "<":
		isTriggered = leftValue < rightValue
	case "<=":
		isTriggered = leftValue <= rightValue
	case ">":
		isTriggered = leftValue > rightValue
	case ">=":
		isTriggered = leftValue >= rightValue
	}

	return
}

func atois(s string) (ret []int, err error) {
	a := strings.Split(s, ",")
	ret = make([]int, len(a))
	for i, v := range a {
		ret[i], err = strconv.Atoi(v)
		if err != nil {
			return
		}
	}
	return
}

// @str: e.g. all(#3) sum(#3) avg(#10) diff(#10) stddev(#10)
func ParseFuncFromString(str string, operator string, rightValue float64) (fn Function, err error) {
	if str == "" {
		return nil, fmt.Errorf("function can not be null!")
	}
	idx := strings.Index(str, "#")
	args, err := atois(str[idx+1 : len(str)-1])
	if err != nil {
		return nil, err
	}

	switch str[:idx-1] {
	case "max":
		fn = &MaxFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "min":
		fn = &MinFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "all":
		fn = &AllFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "sum":
		fn = &SumFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "avg":
		fn = &AvgFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "diff":
		fn = &DiffFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	//case "pdiff":
	//	fn = &PDiffFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "lookup":
		fn = &LookupFunction{Num: args[0], Limit: args[1], Operator: operator, RightValue: rightValue}
	//case "stddev":
	//	fn = &StdDeviationFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	//case "kdiff":
	//	fn = &KDiffFunction{Num: args[0], Limit: args[1], Operator: operator, RightValue: rightValue}
	//case "kpdiff":
	//	fn = &KPDiffFunction{Num: args[0], Limit: args[1], Operator: operator, RightValue: rightValue}
	default:
		err = fmt.Errorf("not_supported_method")
	}

	return
}
