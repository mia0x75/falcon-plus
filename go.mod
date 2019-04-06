module github.com/open-falcon/falcon-plus

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.36.0

	golang.org/x/build => github.com/golang/build v0.0.0-20190228010158-44b79b8774a7
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190227175134-215aa809caaf
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190221220918-438050ddec5e
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190227174305-5b3e6a55c961
	golang.org/x/net => github.com/golang/net v0.0.0-20190227160552-c95aed5357e7
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190226205417-e64efc72b421
	golang.org/x/perf => github.com/golang/perf v0.0.0-20190124201629-844a5f5b46f4
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190226215855-775f8194d0f9
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20181108054448-85acf8d2951c
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190227232517-f0a709d59f0f
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.1.0
	google.golang.org/appengine => github.com/golang/appengine v1.4.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190227213309-4f5b463f9597
	google.golang.org/grpc => github.com/grpc/grpc-go v1.19.0
)

require (
	github.com/astaxie/beego v1.11.1
	github.com/betacraft/yaag v1.0.0
	github.com/emirpasic/gods v1.12.0
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.2.0
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/mia0x75/go-metrics v0.0.0-20181010095928-0c0fea59b08f // indirect
	github.com/mia0x75/gopfc v0.0.0-20181011053331-f72552a8e1fb
	github.com/mindprince/gonvml v0.0.0-20180514031326-b364b296c732
	github.com/niean/gotools v0.0.0-20151221085310-ff3f51fc5c60 // indirect
	github.com/open-falcon/common v0.0.0-20160912145637-b9ba65549217
	github.com/open-falcon/rrdlite v0.0.0-20170412122036-7d8646c85cc5
	github.com/radovskyb/watcher v1.0.6
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/toolkits/cache v0.0.0-20190218093630-cfb07b7585e5
	github.com/toolkits/concurrent v0.0.0-20150624120057-a4371d70e3e3
	github.com/toolkits/conn_pool v0.0.0-20170512061817-2b758bec1177
	github.com/toolkits/consistent v0.0.0-20150827090850-a6f56a64d1b1
	github.com/toolkits/container v0.0.0-20151219225805-ba7d73adeaca
	github.com/toolkits/core v0.0.0-20141116054942-0ebf14900fe2
	github.com/toolkits/cron v0.0.0-20150624115642-bebc2953afa6
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	github.com/toolkits/net v0.0.0-20160910085801-3f39ab6fe3ce
	github.com/toolkits/nux v0.0.0-20190312004434-44d006618852
	github.com/toolkits/proc v0.0.0-20170520054645-8c734d0eb018
	github.com/toolkits/slice v0.0.0-20141116085117-e44a80af2484
	github.com/toolkits/sys v0.0.0-20170615103026-1f33b217ffaf
	github.com/toolkits/time v0.0.0-20160524122720-c274716e8d7f
	github.com/ugorji/go/codec v0.0.0-20190320090025-2dc34c0b8780 // indirect
	golang.org/x/crypto v0.0.0-20190131182504-b8fe1690c613
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)
