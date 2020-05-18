package main

import (
	"fmt"
	"github.com/wind-zhou/logagent/conf"
	"github.com/wind-zhou/logagent/etcd"
	"github.com/wind-zhou/logagent/kafka"
	"github.com/wind-zhou/logagent/taillog"
	"gopkg.in/ini.v1"
	"sync"

	"time"
)

//func run() {
//
//	//1.读取日志
//	for {
//		select {
//		case line := <-taillog.ReadChan():
//			//2.调用发送函数发给kafka
//			kafka.SendToKafka(cfg.KafkaConf.Topic, line.Text)
//		default:
//			time.Sleep(time.Second)
//
//		}
//	}
//
//}

var (
	cfg = new(conf.AppConf)
)

func main() {


	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Println("load ini failed")
		return
	}
	//1.初始化kafka连接

	err = kafka.Init([]string{cfg.KafkaConf.Address},cfg.KafkaConf.Chanmaxsize)
	if err != nil {
		fmt.Printf("init kafka ,err:%v\n", err)
		return
	}
	fmt.Println("init Kafka sucess")

	//初始化etcd

	err =etcd.Init(cfg.EtcdConf.Address,time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("init etcd sucess")


	//2.1  从etcd中拉取配置信息

	logEntryConf,err:=etcd.GetConf("xxx")
	if err!=nil{

		fmt.Printf("etcd.GetConf err:%v\n",err)
		return

	}
	fmt.Printf("get conf from etcd sucess,%v\n",logEntryConf)
	//2.2 并watch该配置


	//watch

	for index,value:=range logEntryConf{
		fmt.Printf("index:%v value:%v\n",index,value)
	}


	//3.收集日志发往kafka
    //3.1 循环每个日志收集项
	taillog.Init(logEntryConf)
	newConfChan:=taillog.NewConfChan()//从tail中获取暴露的通道
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(cfg.EtcdConf.Key,newConfChan)//哨兵发现配置变化会通知上面通道
	wg.Wait()
}

