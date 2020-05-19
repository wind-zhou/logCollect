package taillog

import (
	"fmt"
	"github.com/wind-zhou/logagent/etcd"
	"time"
)

//管理taillog例所有的对象


var tskMgr *taillogMgr


//tailTask 管理者
type taillogMgr struct {
	logEntry []*etcd.LogEntry
	tskMap map[string]*TailTask//存储每次的tsk
	newConfChan chan []*etcd.LogEntry

}


func Init(logConf []*etcd.LogEntry){//从etcd中拿到配置文件，然后根据配置文件初始化
	tskMgr=&taillogMgr{
		logEntry: logConf,//把当前日志收集项配置信息保存起来
		tskMap:make(map[string]*TailTask,16),//用来存储各个task
		newConfChan: make(chan []*etcd.LogEntry),//无缓冲区通道，用来接收热更改配置

	}

	for _,logEntry:=range logConf{
		//conf：*etcd.logEntry
		tailObj:=NewTailTask(logEntry.Path,logEntry.Topic)

		mk:=fmt.Sprintf("%s_%s",logEntry.Path,logEntry.Topic)
		tskMgr.tskMap[mk]=tailObj//记录起了多少个tail实例


	}

	go tskMgr.run()//负责从chan读取更新的配置
}

//监听自己的newConfChan,有了新配置就处理
func (t *taillogMgr)run(){

	for {
		select {
		case newConf:= <- t.newConfChan:
			//
			//fmt.Println("--------------------")
			//fmt.Println("tail模块收到的配置为：")
			//for _,value:=range newConf{
			//	fmt.Printf("%v\n",value)
			//}
			//fmt.Println("--------------------")


			for _,conf:=range newConf{
				mk:=fmt.Sprintf("%s_%s",conf.Path,conf.Topic)
				_,ok:=t.tskMap[mk]//判断是否该项是否为原来的配置项
				if ok{
					//原来就有
					continue
				}else {
					//新增的
					tailObj:=NewTailTask(conf.Path,conf.Topic)//NewTailTask会根据配置文件建立和日志的联系，并读取

					fmt.Printf("tail task %s_%s 启动了了\n",tailObj.path,tailObj.topic)
					t.tskMap[mk]=tailObj
				}
			}


			//找出t.logEntry有但newconf没有的，删除掉
			for _,c1:=range t.logEntry{  //从原来配置中一次拿出配置项，去新的配置中逐一比较
				isDelete:=true
				for _,c2:=range newConf {
					if c2.Path == c1.Path && c2.Topic == c1.Topic {
						isDelete = false
						continue
					}
				}
				if isDelete{
					//把c1对应的这个tailObj停掉,怎么停掉这个之前的协程呢？用context
					mk:=fmt.Sprintf("%s_%s",c1.Path,c1.Topic)
					t.tskMap[mk].cancelFunc()
					delete(t.tskMap, mk) //删除记录

					//fmt.Println("*********************************")
					//fmt.Println("删除某个日志项后，系统存储的日志项为：")
					//
					//for _,v:=range t.tskMap{
					//	fmt.Printf("%v\n",v)
					//}
					//fmt.Println("*********************************")


				}
			}

			//2.配置删除
			fmt.Println("新配置来了",newConf)
		default:
			time.Sleep(time.Second)

		}
	}

}

//定义一个函数，暴露newConfChan

func NewConfChan()  chan<- []*etcd.LogEntry{
	return tskMgr.newConfChan

}