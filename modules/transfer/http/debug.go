package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/open-falcon/falcon-plus/modules/transfer/sender"
)

// SetupDebugRoutes 设置路由
func SetupDebugRoutes() {
	// conn pools
	http.HandleFunc("/debug/connpool/", func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Path[len("/debug/connpool/"):]
		args := strings.Split(urlParam, "/")

		argsLen := len(args)
		if argsLen < 1 {
			w.Write([]byte(fmt.Sprintf("bad args\n")))
			return
		}

		var result string
		receiver := args[0]
		switch receiver {
		case "judge":
			if sender.JudgeConnPools != nil {
				result = strings.Join(sender.JudgeConnPools.Proc(), "\n")
			}
		case "graph":
			if sender.GraphConnPools != nil {
				result = strings.Join(sender.GraphConnPools.Proc(), "\n")
			}
		case "transfer":
			if sender.TransferConnPools != nil {
				result = strings.Join(sender.TransferConnPools.Proc(), "\n")
			}
		default:
			result = fmt.Sprintf("bad args, module not exist\n")
		}
		w.Write([]byte(result))
	})
}
