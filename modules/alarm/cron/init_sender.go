package cron

import (
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

// 变量定义
var (
	IMWorkerChan   chan int
	SMSWorkerChan  chan int
	MailWorkerChan chan int
)

// InitSenderWorker 初始化告警通道
func InitSenderWorker() {
	workerConfig := g.Config().Worker
	IMWorkerChan = make(chan int, workerConfig.IM)
	SMSWorkerChan = make(chan int, workerConfig.SMS)
	MailWorkerChan = make(chan int, workerConfig.Mail)
}
