package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wind-zhou/log_transfer/es"
)



//初始化client 消费者
//从kafka取数据发给es
func Init(address []string,topic string)  error {
	consumer, err := sarama.NewConsumer(address, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return err
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return err
	}
	fmt.Println("分区列表：",partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return err
		}
		//defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				//发送给es
				ld:=es.LogData{
					Topic: topic,
					Data: string(msg.Value),
				}
				es.SendToEsChan(&ld)//把读取的分区数据发送到一个通道
				if err != nil {
					fmt.Printf("es.SendToEs err=%v\n",err)
					return
				}
			}
		}(pc)
	}
	return  err
}
