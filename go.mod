module github.com/open-falcon/falcon-plus

go 1.14

require (
	github.com/emirpasic/gods v1.12.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.4.1
	github.com/jinzhu/gorm v1.9.12
	github.com/mia0x75/go-metrics v0.0.0-20181010095928-0c0fea59b08f // indirect
	github.com/mia0x75/gopfc v0.0.0-20181011053331-f72552a8e1fb
	github.com/mia0x75/yaag v1.0.1
	github.com/mindprince/gonvml v0.0.0-20190828220739-9ebdce4bb989
	github.com/niean/gotools v0.0.0-20151221085310-ff3f51fc5c60 // indirect
	github.com/open-falcon/rrdlite v0.0.0-20200214140804-bf5829f786ad
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.6.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cobra v1.0.0
	github.com/toolkits/cache v0.0.0-20190218093630-cfb07b7585e5
	github.com/toolkits/concurrent v0.0.0-20150624120057-a4371d70e3e3
	github.com/toolkits/conn_pool v0.0.0-20170512061817-2b758bec1177
	github.com/toolkits/consistent v0.0.0-20150827090850-a6f56a64d1b1
	github.com/toolkits/container v0.0.0-20151219225805-ba7d73adeaca
	github.com/toolkits/core v0.0.0-20141116054942-0ebf14900fe2
	github.com/toolkits/cron v0.0.0-20150624115642-bebc2953afa6
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	github.com/toolkits/net v0.0.0-20160910085801-3f39ab6fe3ce
	github.com/toolkits/nux v0.0.0-20200401110743-debb3829764a
	github.com/toolkits/proc v0.0.0-20170520054645-8c734d0eb018
	github.com/toolkits/slice v0.0.0-20141116085117-e44a80af2484
	github.com/toolkits/sys v0.0.0-20170615103026-1f33b217ffaf
	github.com/toolkits/time v0.0.0-20160524122720-c274716e8d7f
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	github.com/go-resty/resty/v2 v2.3.0
)

replace github.com/open-falcon/rrdlite => github.com/mia0x75/rrdlite v0.0.0-20200510063900-8e4569028116
