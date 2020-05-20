package main

import (
   "github.com/wind-zhou/log_transfer/conf"
   "github.com/wind-zhou/log_transfer/es"
   "gopkg.in/ini.v1"
   "fmt"

   "github.com/wind-zhou/log_transfer/kafka"

)

//1.将日志数据取出
   //2.发往ES

var (
   cfg = new(conf.LogTransfer)
)
   func main(){

   	//0. 加载配置文件
      err := ini.MapTo(cfg, "./conf/cfg.ini")
      if err != nil {
         fmt.Println("load ini failed")
         return
      }
      fmt.Println(cfg)


      //初始化es
      err=es.Init(cfg.EsCfg.Address)
      if err != nil {
         fmt.Printf("init es failed err=%v\n",err)
      }

      fmt.Println("init es success")
       //1.初始化kafka
       //1.1 链接kafka
       //1.2 每个分区读数据发往es

   	err=kafka.Init([]string{cfg.KafkaCfg.Address},cfg.KafkaCfg.Topic)
      if err != nil {
         fmt.Printf("init kafka failed err=%v\n",err)
         return
      }

      select {

      }
   }

