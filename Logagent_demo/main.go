package main

import (
	"fmt"
	"github.com/wind-zhou/Logagent_demo/conf"
	"github.com/wind-zhou/Logagent_demo/etcd"
	"github.com/wind-zhou/Logagent_demo/kafka"
	"github.com/wind-zhou/Logagent_demo/tail"
	"gopkg.in/ini.v1"
	"sync"
)

var (
	cfg = new(conf.AppConf)
	LogEtcdConf =make([]*etcd.LogEntry,1000) //用于接受etcd拉取的配置

)


func main(){

	//0. 解析配置文件的内容
	err := ini.MapTo(cfg, "./conf/cfg.ini")
	if err != nil {
		fmt.Printf("解析文件出错：err=",err)
		return

	}
	fmt.Println("解析文件成功")
	fmt.Printf("%v\n",cfg)


	//1.初始化etcd
	err=etcd.Init([]string{cfg.EtcdConf.Address},cfg.EtcdConf.Timeout)
	if err != nil {
		fmt.Printf("初始化etcd失败，err=%v\n",err)
		return
	}
	fmt.Println("connect to etcd success")

	LogEtcdConf,err=etcd.GetConf(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Printf("拉取配置失败：err=%v\n",err)
		return
	}
	fmt.Println("拉取配置成功")

	for _,value:=range LogEtcdConf{
		fmt.Printf("%v\n",value)
	}

	//2.初始化kafka
	//就是连上kafka

	err=kafka.Init([]string{cfg.KafkaConf.Address})
	if err!=nil{
		fmt.Println("初始化失败：, err:", err)
		return
	}

	fmt.Println("kafka init success")

	//3. 初始化tail
	tail.Init(LogEtcdConf)

	fmt.Println("tail init success")

	//初始化信道，用于接受watch的变化
	newConfChan:=tail.NewConfChan()

	var wg sync.WaitGroup
	wg.Add(1)

	go etcd.WatchConf(cfg.EtcdConf.Key,newConfChan)//将配置项放入了通道，在tail模块的goroutine会读取信息
	wg.Wait()
}

