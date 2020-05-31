package api

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/net/httplib"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

// Action TODO: use api/app/model/action.go
type Action struct {
	ID                 int    `json:"id"`
	Uic                string `json:"uic"`
	URL                string `json:"url"`
	Callback           int    `json:"callback"`
	BeforeCallbackSMS  int    `json:"beforeCallbackSMS"`
	BeforeCallbackMail int    `json:"beforeCallbackMail"`
	AfterCallbackSMS   int    `json:"afterCallbackSMS"`
	AfterCallbackMail  int    `json:"afterCallbackMail"`
}

// ActionCache TODO:
type ActionCache struct {
	sync.RWMutex
	M map[int]*Action
}

// Actions TODO:
var Actions = &ActionCache{M: make(map[int]*Action)}

// Get TODO:
func (s *ActionCache) Get(id int) *Action {
	s.RLock()
	defer s.RUnlock()
	val, exists := s.M[id]
	if !exists {
		return nil
	}

	return val
}

// Set TODO:
func (s *ActionCache) Set(id int, action *Action) {
	s.Lock()
	defer s.Unlock()
	s.M[id] = action
}

// GetAction TODO:
func GetAction(id int) *Action {
	action := CurlAction(id)

	if action != nil {
		Actions.Set(id, action)
	} else {
		action = Actions.Get(id)
	}

	return action
}

// CurlAction TODO:
func CurlAction(id int) *Action {
	if id <= 0 {
		return nil
	}

	uri := fmt.Sprintf("%s/api/v1/action/%d", g.Config().API.API, id)
	req := httplib.Get(uri).SetTimeout(5*time.Second, 30*time.Second)
	token, _ := json.Marshal(map[string]string{
		"name": "falcon-alarm",
		"sig":  g.Config().API.Token,
	})
	req.Header("X-Falcon-Token", string(token))

	var act Action
	err := req.ToJson(&act)
	if err != nil {
		log.Errorf("[E] curl %s fail: %v", uri, err)
		return nil
	}

	return &act
}
