package cron

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/common/model"
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
	var strategiesResponse model.StrategiesResponse
	err := g.HbsClient.Call("Hbs.GetStrategies", model.NullRpcRequest{}, &strategiesResponse)
	if err != nil {
		log.Println("[ERROR] Hbs.GetStrategies:", err)
		return
	}

	rebuildStrategyMap(&strategiesResponse)
}

func rebuildStrategyMap(strategiesResponse *model.StrategiesResponse) {
	// endpoint:metric => [strategy1, strategy2 ...]
	m := make(map[string][]model.Strategy)
	for _, hs := range strategiesResponse.HostStrategies {
		hostname := hs.Hostname
		if hostname == g.Config().DebugHost {
			log.Debugln(hostname, "strategies:")
			bs, _ := json.Marshal(hs.Strategies)
			log.Debugln(string(bs))
		}
		for _, strategy := range hs.Strategies {
			key := fmt.Sprintf("%s/%s", hostname, strategy.Metric)
			if _, exists := m[key]; exists {
				m[key] = append(m[key], strategy)
			} else {
				m[key] = []model.Strategy{strategy}
			}
		}
	}

	g.StrategyMap.ReInit(m)
}

func syncExpression() {
	var expressionResponse model.ExpressionResponse
	err := g.HbsClient.Call("Hbs.GetExpressions", model.NullRpcRequest{}, &expressionResponse)
	if err != nil {
		log.Println("[ERROR] Hbs.GetExpressions:", err)
		return
	}

	rebuildExpressionMap(&expressionResponse)
}

func rebuildExpressionMap(expressionResponse *model.ExpressionResponse) {
	m := make(map[string][]*model.Expression)
	for _, exp := range expressionResponse.Expressions {
		for k, v := range exp.Tags {
			key := fmt.Sprintf("%s/%s=%s", exp.Metric, k, v)
			if _, exists := m[key]; exists {
				m[key] = append(m[key], exp)
			} else {
				m[key] = []*model.Expression{exp}
			}
		}
	}

	g.ExpressionMap.ReInit(m)
}

func syncFilter() {
	m := make(map[string]string)

	//M map[string][]model.Strategy
	strategyMap := g.StrategyMap.Get()
	for _, strategies := range strategyMap {
		for _, strategy := range strategies {
			m[strategy.Metric] = strategy.Metric
		}
	}

	//M map[string][]*model.Expression
	expressionMap := g.ExpressionMap.Get()
	for _, expressions := range expressionMap {
		for _, expression := range expressions {
			m[expression.Metric] = expression.Metric
		}
	}

	g.FilterMap.ReInit(m)
}
