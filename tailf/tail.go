package tailf

//本文件功能是从日志文件内不断读取数据，并组成完整信息结构发送到chan内。并提供输出chan内的数据的接口

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"time"
)

type CollectConf struct { //保存收集的日志信息
	LogPath string
	Topic   string //头，既谁的日志（从配置文件内获得）
}

type TailObj struct { //一个tail对象信息
	tail *tail.Tail
	conf CollectConf
}

type Textmsg struct {
	Msg   string //日志内容
	Topic string //从CollectConf的Topic字段内获得
}

type TailObjMgr struct { //管理所有的tail对象
	tailsObjs []*TailObj
	msgChan   chan *Textmsg
}

var (
	tailObjMgr *TailObjMgr
)

func GetOneline() (msg *Textmsg) { //从管道获取从日志文件读取出来的内容
	msg = <-tailObjMgr.msgChan
	return
}

func InitTail(conf []CollectConf, chanSize int) (err error) {

	//参数检验
	if len(conf) == 0 {
		err = fmt.Errorf("invalid config for log collect, conf:%v", conf)
		return
	}

	//初始化chan
	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *Textmsg, chanSize),
	}

	for _, v := range conf {
		//初始化对象内自定义的conf
		obj := &TailObj{
			conf: v,
		}
		//初始化对象内的tail对象
		tails, errTial := tail.TailFile(v.LogPath, tail.Config{
			ReOpen:    true,
			Follow:    true,
			MustExist: false,
			Poll:      true,
		})
		if errTial != nil {
			fmt.Println("tail file err:", err)
			return
		}

		obj.tail = tails //赋值

		tailObjMgr.tailsObjs = append(tailObjMgr.tailsObjs, obj)

		go readFromTail(obj) //单独起一个协程收集一个日志的读
	}

	return
}

func readFromTail(tailObj *TailObj) { //读一个日志的信息
	logs.Debug("czw tailObj is %v", tailObj)

	for true {
		line, ok := <-tailObj.tail.Lines //通过网络tail读取一行
		if !ok {
			logs.Warn("tail file close reopen, filename:%s\n", tailObj.tail.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		textMsg := &Textmsg{
			Msg:   line.Text,
			Topic: tailObj.conf.Topic,
		}

		tailObjMgr.msgChan <- textMsg
	}
}
