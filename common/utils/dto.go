package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RespJSON struct {
	Error string `json:"error,omitempty"`
	Msg   string `json:"message,omitempty"`
}

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// func JSONR(c *gin.Context, wcode int, msg interface{}) (werror error) {
func JSONR(c *gin.Context, arg ...interface{}) (werror error) {
	var (
		code int
		msg  interface{}
	)

	if len(arg) == 1 {
		code = http.StatusOK
		msg = arg[0]
	} else {
		code = arg[0].(int)
		msg = arg[1]
	}

	var body interface{}
	if code == 200 {
		switch msg.(type) {
		case string:
			body = RespJSON{Msg: msg.(string)}
			c.JSON(http.StatusOK, body)
		default:
			c.JSON(http.StatusOK, msg)
			body = msg
		}
	} else {
		switch msg.(type) {
		case string:
			body = RespJSON{
				Error: msg.(string),
			}
		case error:
			body = RespJSON{
				Error: msg.(error).Error(),
			}
		default:
			body = RespJSON{
				Error: "system type error. please ask admin for help",
			}
		}
		c.JSON(code, body)
	}
	return
}

// RenderJSON
func RenderJSON(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

// RenderDataJSON
func RenderDataJSON(w http.ResponseWriter, data interface{}) {
	RenderJSON(w, Dto{Msg: "success", Data: data})
}

// RenderMsgJSON
func RenderMsgJSON(w http.ResponseWriter, msg string) {
	RenderJSON(w, map[string]string{"msg": msg})
}

// AutoRender
func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJSON(w, err.Error())
		return
	}
	RenderDataJSON(w, data)
}

// StdRender
func StdRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		w.WriteHeader(400)
		RenderMsgJSON(w, err.Error())
		return
	}
	RenderJSON(w, data)
}

func postByJSON(rw http.ResponseWriter, req *http.Request, url string) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	s := buf.String()
	reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
	if err != nil {
		log.Errorf("[E] %v", err)
	}
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		log.Errorf("[E] %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}

func postByForm(rw http.ResponseWriter, req *http.Request, url string) {
	req.ParseForm()
	client := &http.Client{}
	resp, err := client.PostForm(url, req.PostForm)
	if err != nil {
		log.Errorf("[E] %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}

func getRequest(rw http.ResponseWriter, url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("[E] %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("[E] %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}
