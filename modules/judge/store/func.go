package store

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
)

type Function interface {
	Compute(L *SafeLinkedList) (vs []*cm.HistoryData, leftValue float64, isTriggered bool, isEnough bool)
}

type MaxFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (f MaxFunction) Compute(L *SafeLinkedList) (vs []*cm.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	max := vs[0].Value
	for i := 1; i < f.Limit; i++ {
		if max < vs[i].Value {
			max = vs[i].Value
		}
	}

	leftValue = max
	isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
	return
}

type MinFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (f MinFunction) Compute(L *SafeLinkedList) (vs []*cm.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	min := vs[0].Value
	for i := 1; i < f.Limit; i++ {
		if min > vs[i].Value {
			min = vs[i].Value
		}
	}

	leftValue = min
	isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
	return
}

type AllFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (f AllFunction) Compute(L *SafeLinkedList) (vs []*cm.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	isTriggered = true
	for i := 0; i < f.Limit; i++ {
		isTriggered = checkIsTriggered(vs[i].Value, f.Operator, f.RightValue)
		if !isTriggered {
			break
		}
	}

	leftValue = vs[0].Value
	return
}

type LookupFunction struct {
	Function
	Num        int
	Limit      int
	Operator   string
	RightValue float64
}

func (f LookupFunction) Compute(L *SafeLinkedList) (vs []*cm.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	leftValue = vs[0].Value

	for n, i := 0, 0; i < f.Limit; i++ {
		if checkIsTriggered(vs[i].Value, f.Operator, f.RightValue) {
			n++
			if n == f.Num {
				isTriggered = true
				return
			}
		}
	}

	return
}

type SumFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (f SumFunction) Compute(L *SafeLinkedList) (vs []*cm.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	sum := 0.0
	for i := 0; i < f.Limit; i++ {
		sum += vs[i].Value
	}

	leftValue = sum
	isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
	return
}

type AvgFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (f AvgFunction) Compute(L *SafeLinkedList) (vs []*cm.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}

	sum := 0.0
	for i := 0; i < f.Limit; i++ {
		sum += vs[i].Value
	}

	leftValue = sum / float64(f.Limit)
	isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
	return
}

type DiffFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

// 只要有一个点的diff触发阈值，就报警
func (f DiffFunction) Compute(L *SafeLinkedList) (vs []*cm.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	// 此处f.Limit要+1，因为通常说diff(#3)，是当前点与历史的3个点相比较
	// 然而最新点已经在linkedlist的第一个位置，所以……
	vs, isEnough = L.HistoryData(f.Limit + 1)
	if !isEnough {
		return
	}

	if len(vs) == 0 {
		isEnough = false
		return
	}

	first := vs[0].Value

	isTriggered = false
	for i := 1; i < f.Limit+1; i++ {
		// diff是当前值减去历史值
		leftValue = first - vs[i].Value
		isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
		if isTriggered {
			break
		}
	}

	return
}

// pdiff(#3)
type PDiffFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

func (f PDiffFunction) Compute(L *SafeLinkedList) (vs []*cm.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit + 1)
	if !isEnough {
		return
	}

	if len(vs) == 0 {
		isEnough = false
		return
	}

	first := vs[0].Value

	isTriggered = false
	for i := 1; i < f.Limit+1; i++ {
		if vs[i].Value == 0 {
			continue
		}

		leftValue = (first - vs[i].Value) / vs[i].Value * 100.0
		isTriggered = checkIsTriggered(leftValue, f.Operator, f.RightValue)
		if isTriggered {
			break
		}
	}

	return
}

type StdDeviationFunction struct {
	Function
	Limit      int
	Operator   string
	RightValue float64
}

/*
	离群点检测函数，更多请参考3-sigma算法：https://en.wikipedia.org/wiki/68%E2%80%9395%E2%80%9399.7_rule
	stddev(#10) = 3 // 取最新 **10** 个点的数据分别计算得到他们的标准差和均值，分别计为 σ 和 μ，其中当前值计为 X，那么当 X 落在区间 [μ-3σ, μ+3σ] 之外时则报警。
*/
func (f StdDeviationFunction) Compute(L *SafeLinkedList) (vs []*cm.HistoryData, leftValue float64, isTriggered bool, isEnough bool) {
	vs, isEnough = L.HistoryData(f.Limit)
	if !isEnough {
		return
	}
	if len(vs) == 0 {
		isEnough = false
		return
	}
	leftValue = vs[0].Value
	var datas []float64
	for _, i := range vs {
		datas = append(datas, i.Value)
	}
	isTriggered = false
	std := cu.ComputeStdDeviation(datas)
	mean := cu.ComputeMean(datas)
	upperBound := mean + f.RightValue*std
	lowerBound := mean - f.RightValue*std
	if leftValue < lowerBound || leftValue > upperBound {
		isTriggered = true
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

// @str: e.g. max(#3) min(#3) all(#3) sum(#3) avg(#3) diff(#3) pdiff(#3) lookup(#2,3) stddev(#3)
func ParseFuncFromString(str string, operator string, rightValue float64) (fn Function, err error) {
	if str == "" {
		return nil, fmt.Errorf("func can not be null!")
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
	case "pdiff":
		fn = &PDiffFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	case "lookup":
		fn = &LookupFunction{Num: args[0], Limit: args[1], Operator: operator, RightValue: rightValue}
	case "stddev":
		fn = &StdDeviationFunction{Limit: args[0], Operator: operator, RightValue: rightValue}
	default:
		err = fmt.Errorf("not_supported_method")
	}

	return
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
