package main

import (
	"10_xm_logAgent/tailf"
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
)

var (
	appConfig *Config //保存配置文件信息的实例对象
)

type Config struct { //保存配置文件的信息
	logLevel  string
	logPath   string
	chan_size int

	kafkaAddr string

	//etcdAddr string

	collectConf []tailf.CollectConf
}

func loadCollectConf(conf config.Configer) error {

	var cc tailf.CollectConf
	cc.LogPath = conf.String("collect::log_path")
	if len(cc.LogPath) == 0 {
		err := errors.New("invaild collect::log_path")
		return err
	}

	cc.Topic = conf.String("collect::topic")
	if len(cc.LogPath) == 0 {
		err := errors.New("invaild collect::topic")
		return err
	}

	appConfig.collectConf = append(appConfig.collectConf, cc)
	return nil

}

func loadConf(confType, filename string) error { //读取配置文件相关内容

	conf, errConf := config.NewConfig(confType, filename) //生成一个文件读取实例conf(文件格式，文件路径)
	if errConf != nil {
		fmt.Println("new config failed, err", errConf)
		return errConf
	}

	appConfig = &Config{}
	appConfig.logLevel = conf.String("logs::log_level")
	if len(appConfig.logLevel) == 0 { //容错处理
		appConfig.logLevel = "debug"
	}

	appConfig.logPath = conf.String("logs::log_path")
	if len(appConfig.logPath) == 0 {
		appConfig.logPath = "10_xm_logAgent/logs"
	}

	appConfig.chan_size, errConf = conf.Int("collect::chan_size")
	if errConf != nil {
		appConfig.chan_size = 100
		return errConf
	}

	appConfig.kafkaAddr = conf.String("kafka::server_addr")
	if len(appConfig.kafkaAddr) == 0 {
		errConf = fmt.Errorf("invaild kafka addr")
	}

	/*appConfig.etcdAddr = conf.String("etcd::addr")
	if len(appConfig.kafkaAddr) == 0 {
		errConf = fmt.Errorf("invaild etcd addr")
	} */

	//加载收集日志信息（加载collectConf字段信息）
	errConf = loadCollectConf(conf)
	if errConf != nil {
		fmt.Println("load collect conf faild, err:%v\n", errConf)
		return errConf
	}

	return nil

}
