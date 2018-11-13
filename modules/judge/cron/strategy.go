package cron

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
)

func SyncStrategies() {
	go func() {
		d := time.Duration(g.Config().Hbs.Interval) * time.Second
		for range time.Tick(d) {
			syncStrategies()
			syncExpression()
			syncFilter()
		}
	}()
}

func syncStrategies() {
	var strategiesResponse cmodel.StrategiesResponse
	err := g.HbsClient.Call("Hbs.GetStrategies", cmodel.NullRpcRequest{}, &strategiesResponse)
	if err != nil {
		log.Errorf("[E] Hbs.GetStrategies: %v", err)
		return
	}

	rebuildStrategyMap(&strategiesResponse)
}

func rebuildStrategyMap(strategiesResponse *cmodel.StrategiesResponse) {
	// endpoint:metric => [strategy1, strategy2 ...]
	m := make(map[string][]cmodel.Strategy)
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
				m[key] = []cmodel.Strategy{strategy}
			}
		}
	}

	g.StrategyMap.ReInit(m)
}

func syncExpression() {
	var expressionResponse cmodel.ExpressionResponse
	err := g.HbsClient.Call("Hbs.GetExpressions", cmodel.NullRpcRequest{}, &expressionResponse)
	if err != nil {
		log.Errorf("[E] Hbs.GetExpressions: %v", err)
		return
	}

	rebuildExpressionMap(&expressionResponse)
}

func rebuildExpressionMap(expressionResponse *cmodel.ExpressionResponse) {
	m := make(map[string][]*cmodel.Expression)
	for _, exp := range expressionResponse.Expressions {
		for k, v := range exp.Tags {
			key := fmt.Sprintf("%s/%s=%s", exp.Metric, k, v)
			if _, exists := m[key]; exists {
				m[key] = append(m[key], exp)
			} else {
				m[key] = []*cmodel.Expression{exp}
			}
		}
	}

	g.ExpressionMap.ReInit(m)
}

func syncFilter() {
	m := make(map[string]string)

	//M map[string][]cmodel.Strategy
	strategyMap := g.StrategyMap.Get()
	for _, strategies := range strategyMap {
		for _, strategy := range strategies {
			m[strategy.Metric] = strategy.Metric
		}
	}

	//M map[string][]*cmodel.Expression
	expressionMap := g.ExpressionMap.Get()
	for _, expressions := range expressionMap {
		for _, expression := range expressions {
			m[expression.Metric] = expression.Metric
		}
	}

	g.FilterMap.ReInit(m)
}
