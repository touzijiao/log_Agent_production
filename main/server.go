package main

import (
	"10_xm_logAgent/kafka"
	"10_xm_logAgent/tailf"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

func serverRun() (err error) {
	for {
		msg := tailf.GetOneline() //msg为一个完整的带头与数据的对象
		err = sendToKafka(msg)
		if err != nil {
			logs.Error("send to kafka failed, err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}

func sendToKafka(msg *tailf.Textmsg) (err error) {
	fmt.Printf("read msg:%s, topic:%s\n", msg.Msg, msg.Topic)

	err = kafka.SendToKafka(msg.Msg, msg.Topic)

	return

}
