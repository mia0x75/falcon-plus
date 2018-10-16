package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type RespJson struct {
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
		wcode int
		msg   interface{}
	)

	if len(arg) == 1 {
		wcode = http.StatusOK
		msg = arg[0]
	} else {
		wcode = arg[0].(int)
		msg = arg[1]
	}

	var body interface{}
	if wcode == 200 {
		switch msg.(type) {
		case string:
			body = RespJson{Msg: msg.(string)}
			c.JSON(http.StatusOK, body)
		default:
			c.JSON(http.StatusOK, msg)
			body = msg
		}
	} else {
		switch msg.(type) {
		case string:
			body = RespJson{Error: msg.(string)}
			c.JSON(wcode, body)
		case error:
			body = RespJson{Error: msg.(error).Error()}
			c.JSON(wcode, body)
		default:
			body = RespJson{Error: "system type error. please ask admin for help"}
			c.JSON(wcode, body)
		}
	}
	return
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	RenderDataJson(w, data)
}

func StdRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		w.WriteHeader(400)
		RenderMsgJson(w, err.Error())
		return
	}
	RenderJson(w, data)
}

func postByJson(rw http.ResponseWriter, req *http.Request, url string) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	s := buf.String()
	reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
	if err != nil {
		log.Println("Error =", err.Error())
	}
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		log.Println("Error =", err.Error())
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
		log.Println("Error =", err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}

func getRequest(rw http.ResponseWriter, url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error =", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error =", err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}
