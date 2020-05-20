package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"strings"

	"time"
)

type LogData struct {
	Topic string `json:"topic"`

	Data string`json:"data"`
}

var (
	client *elastic.Client
	ch =make(  chan *LogData,100000)
)

//初始化es
//准备接受Kafka的数据

func Init(address string)(err error){
	if !strings.HasPrefix(address,"http://"){
		address="http://"+address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		// Handle error
		return
	}

	fmt.Println("connect to es success")

	go SendToEs()
	return

}



//SendToEs 发送数据到es

func SendToEsChan(msg *LogData) {
	ch <-msg

}

func SendToEs(){
	for{
		select {
			case msg:= <-ch:

				put1, err := client.Index().Index(msg.Topic).Type("xxx").BodyJson(msg).Do(context.Background())
				if err != nil {
					// Handle error
					fmt.Println(err)

				}
				fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}

}