package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/common/model"
)

func PostPush(L []*model.JsonMetaData) error {
	bs, err := json.Marshal(L)
	if err != nil {
		return err
	}

	bf := bytes.NewBuffer(bs)

	resp, err := http.Post(PostPushUrl, "application/json", bf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	content := string(body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code %d != 200, response: %s", resp.StatusCode, content)
	}

	log.Debugf("[D] response: %s", content)

	return nil
}
