package requests

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/toolkits/net/httplib"
)

func CurlPlus(uri, method, token_name, token_sig string, headers, params map[string]string) (req *httplib.BeegoHttpRequest, err error) {
	if method == "GET" {
		req = httplib.Get(uri)
	} else if method == "POST" {
		req = httplib.Post(uri)
	} else if method == "PUT" {
		req = httplib.Put(uri)
	} else if method == "DELETE" {
		req = httplib.Delete(uri)
	} else if method == "HEAD" {
		req = httplib.Head(uri)
	} else {
		err = errors.New("invalid http method")
		return
	}

	req = req.SetTimeout(1*time.Second, 5*time.Second)

	token, _ := json.Marshal(map[string]string{
		"name": token_name,
		"sig":  token_sig,
	})
	req.Header("X-Falcon-Token", string(token))

	for hk, hv := range headers {
		req.Header(hk, hv)
	}

	for pk, pv := range params {
		req.Param(pk, pv)
	}

	return
}
