package rpc

import (
	"fmt"

	cm "github.com/open-falcon/falcon-plus/common/model"
	cu "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/hbs/cache"
)

// Hbs TODO:
type Hbs int

// GetHosts 获取监控主机
func (t *Hbs) GetHosts(args string, reply *cm.ExpressionResponse) error {
	return nil
}

// GetExpressions 获取表达式
func (t *Hbs) GetExpressions(req cm.NullRPCRequest, reply *cm.ExpressionResponse) error {
	reply.Expressions = cache.ExpressionCache.Get()
	return nil
}

// GetStrategies 获取策略
func (t *Hbs) GetStrategies(req cm.NullRPCRequest, reply *cm.StrategiesResponse) error {
	reply.HostStrategies = []*cm.HostStrategy{}
	// 一个机器ID对应多个模板ID
	hidTids := cache.HostTemplateIDs.GetMap()
	sz := len(hidTids)
	if sz == 0 {
		return nil
	}

	// Judge需要的是hostname，此处要把HostId转换为hostname
	// 查出的hosts，是不处于维护时间内的
	hosts := cache.MonitoredHosts.Get()
	if len(hosts) == 0 {
		// 所有机器都处于维护状态，汗
		return nil
	}

	tpls := cache.TemplateCache.GetMap()
	if len(tpls) == 0 {
		return nil
	}

	strategies := cache.Strategies.GetMap()
	if len(strategies) == 0 {
		return nil
	}

	// 做个索引，给一个tplId，可以很方便的找到对应了哪些Strategy
	tpl2Strategies := Tpl2Strategies(strategies)

	hostStrategies := make([]*cm.HostStrategy, 0, sz)
	for hostId, tplIds := range hidTids {

		h, exists := hosts[hostId]
		if !exists {
			continue
		}

		// 计算当前host配置了哪些监控策略
		ss := CalcInheritStrategies(tpls, tplIds, tpl2Strategies)
		if len(ss) <= 0 {
			continue
		}

		hs := cm.HostStrategy{
			Hostname:   h.Name,
			Strategies: ss,
		}

		hostStrategies = append(hostStrategies, &hs)

	}

	reply.HostStrategies = hostStrategies
	return nil
}

// Tpl2Strategies TODO:
func Tpl2Strategies(strategies map[int]*cm.Strategy) map[int][]*cm.Strategy {
	ret := make(map[int][]*cm.Strategy)
	for _, s := range strategies {
		if s == nil || s.Template == nil {
			continue
		}
		if _, exists := ret[s.Template.ID]; exists {
			ret[s.Template.ID] = append(ret[s.Template.ID], s)
		} else {
			ret[s.Template.ID] = []*cm.Strategy{s}
		}
	}
	return ret
}

// CalcInheritStrategies TODO:
func CalcInheritStrategies(allTpls map[int]*cm.Template, tids []int, tpl2Strategies map[int][]*cm.Strategy) []cm.Strategy {
	// 根据模板的继承关系，找到每个机器对应的模板全量
	/**
	 * host_id =>
	 * |a |d |a |a |a |
	 * |  |  |b |b |f |
	 * |  |  |  |c |  |
	 * |  |  |  |  |  |
	 */
	tpl_buckets := [][]int{}
	for _, tid := range tids {
		ids := cache.ParentIds(allTpls, tid)
		if len(ids) <= 0 {
			continue
		}
		tpl_buckets = append(tpl_buckets, ids)
	}

	// 每个host 关联的模板，有继承关系的放到同一个bucket中，其他的放在各自单独的bucket中
	/**
	 * host_id =>
	 * |a |d |a |
	 * |b |  |f |
	 * |c |  |  |
	 * |  |  |  |
	 */
	count := len(tpl_buckets)
	uniq_tpl_buckets := [][]int{}
	for i := 0; i < count; i++ {
		var valid bool = true
		for j := 0; j < count; j++ {
			if i == j {
				continue
			}
			if slice_int_eq(tpl_buckets[i], tpl_buckets[j]) {
				break
			}
			if slice_int_lt(tpl_buckets[i], tpl_buckets[j]) {
				valid = false
				break
			}
		}
		if valid {
			uniq_tpl_buckets = append(uniq_tpl_buckets, tpl_buckets[i])
		}
	}

	// 继承覆盖父模板策略，得到每个模板聚合后的策略列表
	strategies := []cm.Strategy{}

	exists_by_id := make(map[int]struct{})
	for _, bucket := range uniq_tpl_buckets {

		// 开始计算一个桶，先计算老的tid，再计算新的，所以可以覆盖
		// 该桶最终结果
		bucket_stras_map := make(map[string][]*cm.Strategy)
		for _, tid := range bucket {

			// 一个tid对应的策略列表
			the_tid_stras := make(map[string][]*cm.Strategy)

			if stras, ok := tpl2Strategies[tid]; ok {
				for _, s := range stras {
					uuid := fmt.Sprintf("metric: %s/tags: %v", s.Metric, cu.SortedTags(s.Tags))
					if _, ok2 := the_tid_stras[uuid]; ok2 {
						the_tid_stras[uuid] = append(the_tid_stras[uuid], s)
					} else {
						the_tid_stras[uuid] = []*cm.Strategy{s}
					}
				}
			}

			// 覆盖父模板
			for uuid, ss := range the_tid_stras {
				bucket_stras_map[uuid] = ss
			}
		}

		last_tid := bucket[len(bucket)-1]

		// 替换所有策略的模板为最年轻的模板
		for _, ss := range bucket_stras_map {
			for _, s := range ss {
				valStrategy := *s
				// exists_by_id[s.Id] 是根据策略ID去重，不太确定是否真的需要，不过加上肯定没问题
				if _, exist := exists_by_id[valStrategy.ID]; !exist {
					if valStrategy.Template.ID != last_tid {
						valStrategy.Template = allTpls[last_tid]
					}
					strategies = append(strategies, valStrategy)
					exists_by_id[valStrategy.ID] = struct{}{}
				}
			}
		}
	}

	return strategies
}

func slice_int_contains(list []int, target int) bool {
	for _, b := range list {
		if b == target {
			return true
		}
	}
	return false
}

func slice_int_eq(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, av := range a {
		if av != b[i] {
			return false
		}
	}
	return true
}

func slice_int_lt(a []int, b []int) bool {
	for _, i := range a {
		if !slice_int_contains(b, i) {
			return false
		}
	}
	return true
}
