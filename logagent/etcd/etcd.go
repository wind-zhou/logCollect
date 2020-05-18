package etcd

import (
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"time"
	"fmt"
	"context"
)

var(
	cli *clientv3.Client
)

//需要收集的日志的配置信息
type LogEntry struct {

	Path string`jaon:"path"`//日志存放路径
	Topic string`json:"topic"`//要发往kafka那个topic

}

//初始化etcd的函数

func Init(addr string,timeout time.Duration)(err error){
	//初始化etcd
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}

	return
}


//从etcd中根据key获取配置项

func GetConf(key string)(LogEtcdConf []*LogEntry ,err error ){
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {   //遍历切片得到结构体指针
		//fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		//对接受的数据反序列化

		err=json.Unmarshal(ev.Value,&LogEtcdConf)//返回的ev.Value也是切片类型
		if err!=nil{
			fmt.Printf("json.Unmarshal failed ,err=%V\n",err)
			return
		}
	}
    return
}



//etcd watch

func WatchConf(key string,newConfChan chan <- []*LogEntry){

	rch := cli.Watch(context.Background(), key) // <-chan WatchResponse

	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		//通知别人
		//通知tail.tskMgr
		//判断操作类型
		var newConf []*LogEntry
			if ev.Type!=clientv3.EventTypeDelete{

				err:=json.Unmarshal(ev.Kv.Value,&newConf)
				if err != nil {
					fmt.Printf("json.Unmarshal failed er:%v\n",err)
					continue
				}

			}


			fmt.Printf("get newConf %v\n",newConf)


			for _,value:=range newConf{
				fmt.Printf("%v\n",value)
			}

			newConfChan <- newConf


		}
	}

}