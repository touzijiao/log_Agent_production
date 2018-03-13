package main

import (
	"10_xm_logAgent/kafka"
	"10_xm_logAgent/tailf"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

func main() {
	filename := "E:/go/src/10_xm_logAgent/conf/logagent.conf"
	err := loadConf("ini", filename) //加载(初始化)配置信息
	if err != nil {
		fmt.Println("load conf faild, err:%v\n", err)
		panic("load conf failed")
		return
	}
	logs.Debug("load conf succ, config: %v", appConfig)

	err = initLogger() //初始化日志
	if err != nil {
		fmt.Println("load logger failed, err:%v\n", err)
		panic("load logger failed")
		return
	}
	logs.Debug("init logger succ")

	err = tailf.InitTail(appConfig.collectConf, appConfig.chan_size) //初始化tail（tail读取出来的日子文件信息，我们在写到kafka内）
	if err != nil {
		logs.Error("init tail failed, err :%v", err)
		return
	}
	logs.Debug("init tailf succ")

	err = kafka.InitKafka(appConfig.kafkaAddr) //初始化kafka
	if err != nil {
		fmt.Println("InitKafka failed, err:%v\n", err)
		return
	}
	logs.Debug("init kafka succ")
	//defer Client.Close()

	logs.Debug("initialize all succ")

	//测试
	go func() {
		var count int
		for {
			count++
			logs.Debug("test for logger %d", count)
			time.Sleep(time.Millisecond * 1000)
		}
	}()

	err = serverRun() //取出数据发送到kafka
	if err != nil {
		logs.Error("serverRun failed, err:%v", err)
		return
	}
	logs.Info("program exited")
}
