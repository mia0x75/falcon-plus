package cron

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
)

func SyncStrategies() {
	go func() {
		d := time.Duration(g.Config().HBS.Interval) * time.Second
		for range time.Tick(d) {
			syncStrategies()
			syncExpression()
			syncFilter()
		}
	}()
}

func syncStrategies() {
	var strategiesResponse cm.StrategiesResponse
	err := g.HBSClient.Call("Hbs.GetStrategies", cm.NullRPCRequest{}, &strategiesResponse)
	if err != nil {
		log.Errorf("[E] Hbs.GetStrategies: %v", err)
		return
	}

	rebuildStrategyMap(&strategiesResponse)
}

func rebuildStrategyMap(strategiesResponse *cm.StrategiesResponse) {
	// endpoint:metric => [strategy1, strategy2 ...]
	m := make(map[string][]cm.Strategy)
	for _, hs := range strategiesResponse.HostStrategies {
		hostname := hs.Hostname
		if hostname == g.Config().DebugHost {
			bs, _ := json.Marshal(hs.Strategies)
			log.Debugf("[D] %s, strategies: %s", hostname, string(bs))
		}
		for _, strategy := range hs.Strategies {
			key := fmt.Sprintf("%s/%s", hostname, strategy.Metric)
			if _, exists := m[key]; exists {
				m[key] = append(m[key], strategy)
			} else {
				m[key] = []cm.Strategy{strategy}
			}
		}
	}

	g.StrategyMap.ReInit(m)
}

func syncExpression() {
	var expressionResponse cm.ExpressionResponse
	err := g.HBSClient.Call("Hbs.GetExpressions", cm.NullRPCRequest{}, &expressionResponse)
	if err != nil {
		log.Errorf("[E] Hbs.GetExpressions: %v", err)
		return
	}

	rebuildExpressionMap(&expressionResponse)
}

func rebuildExpressionMap(expressionResponse *cm.ExpressionResponse) {
	m := make(map[string][]*cm.Expression)
	for _, exp := range expressionResponse.Expressions {
		for k, v := range exp.Tags {
			key := fmt.Sprintf("%s/%s=%s", exp.Metric, k, v)
			if _, exists := m[key]; exists {
				m[key] = append(m[key], exp)
			} else {
				m[key] = []*cm.Expression{exp}
			}
		}
	}

	g.ExpressionMap.ReInit(m)
}

func syncFilter() {
	m := make(map[string]string)

	// M map[string][]cm.Strategy
	strategyMap := g.StrategyMap.Get()
	for _, strategies := range strategyMap {
		for _, strategy := range strategies {
			m[strategy.Metric] = strategy.Metric
		}
	}

	// M map[string][]*cm.Expression
	expressionMap := g.ExpressionMap.Get()
	for _, expressions := range expressionMap {
		for _, expression := range expressions {
			m[expression.Metric] = expression.Metric
		}
	}

	g.FilterMap.ReInit(m)
}
