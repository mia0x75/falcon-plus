package rpc

import (
	"fmt"

	cm "github.com/open-falcon/falcon-plus/common/model"
	grh "github.com/open-falcon/falcon-plus/modules/api/graph"
)

type Graph int

func (grprpc *Graph) QueryOne(para cm.GraphQueryParam, resp *cm.GraphQueryResponse) error {
	r, _ := grh.QueryOne(para)
	if r != nil {
		resp.Values = r.Values
		resp.Counter = r.Counter
		resp.DsType = r.DsType
		resp.Endpoint = r.Endpoint
		resp.Step = r.Step
		fmt.Println(resp)
	}
	return nil
}
