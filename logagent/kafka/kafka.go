package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

//专门往kafka里写日志

type logData struct {
	topic string
	data string

}

var (
	client sarama.SyncProducer //声明一个全局的连接kafka的生产者client
	logDataChan chan *logData
)

//初始化client
func Init(addrs []string,maxSize int) (err error) {

	//先初始化配置文件
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka

	//1.用sarama.NewSyncProducer产生根据特定配置项指定位置写数据的一个生产者
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	logDataChan=make(chan *logData,maxSize)//初始化缓存通道
	//后台开启goroutine发往kafka
	go sendToKafka()
	return
}


//真正往kafka发送日志函数
func sendToKafka() {

	for {
		select {
		case ld :=<- logDataChan:

			//1.构造一个信息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			//2.发送到kafka
			partition, offset, err := client.SendMessage(msg)
			if err != nil {
				log.Printf("FAILED to send message: %s\n", err)
				return
			} else {
				log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
			}
		default:
			time.Sleep(time.Microsecond*50)


		}


	}



}

//该函数把日志
func SendToChan(topic, data string){

		msg:=&logData{
			topic: topic,
			data: data,
		}
		logDataChan <- msg
}
