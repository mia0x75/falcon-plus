package requests

import (
	"encoding/json"

	"github.com/open-falcon/falcon-plus/common/utils"
)

func PostJsonBody(url string, v interface{}) ([]byte, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	resp, err := utils.Post(url, bs)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
