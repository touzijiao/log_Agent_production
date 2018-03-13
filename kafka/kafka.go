package kafka

//此包功能为把读取到的日志信息写入kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	Client sarama.SyncProducer
)

func InitKafka(addr string) (err error) { //初始化kafka

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	//client, err := sarama.NewSyncProducer([]string{"192.168.1.166:9092"}, config)
	Client, err = sarama.NewSyncProducer([]string{addr}, config) //初始化client
	if err != nil {
		logs.Error("init kafka producer failed, err", err)
		return
	}

	logs.Debug("init kafka success")
	return

}

func SendToKafka(data, topic string) (err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	_, _, err = Client.SendMessage(msg)
	if err != nil {
		logs.Error("send message failed, err:%v data:%v topic:%v", err, data, topic)
		return
	}

	//logs.Debug("send succ, pid:%v offset:%v, topic:%v\n", pid, offset, topic)
	return
}
