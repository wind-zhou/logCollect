package taillog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/wind-zhou/logagent/kafka"
)


//var (
//	tailObj *tail.Tail
//	//LogChan chan string
//)


//一个日志收集任务，对多个tailobj进行管理
type TailTask struct {
	path string  //日志路径
	topic string
	instance *tail.Tail
	ctx context.Context//哦用于之后停止协程
	cancelFunc context.CancelFunc

}





func NewTailTask(path ,topic string)(tailObj *TailTask){
	ctx,cancel:=context.WithCancel(context.Background())

	tailObj=&TailTask{
		path:path,
		topic: topic,
		ctx:ctx,
		cancelFunc:cancel,

	}
	tailObj.Init()//根据路径打开对应文件
	return

}

func (t TailTask)Init(){
	config := tail.Config{
		// File-specifc
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件那个位置开始读
		ReOpen:    true,                                 //是否重新打开
		MustExist: false,                                // Fail early if the file does not exist
		Poll:      true,                                 // Poll for file changes instead of using inotify
		Follow:    true,                                 // Continue looking for new lines (tail -f)

	}
	var err error
	t.instance,err = tail.TailFile(t.path, config) //TailFile(filename, config)

	if err != nil {
		fmt.Println("tail file err=", err)

	}


	go t.run()//直接去采集日志发往kafka
}


func (t *TailTask)run(){

	for {
		select {
		case <-t.ctx.Done():
			fmt.Printf("tail task %s_%s 结束了\n",t.path,t.topic)
			return
		case line:=	<- t.instance.Lines://从tailobj同道中读数据
			//kafka.SendToKafka(t.topic,line.Text)//这里相当于取日志如发送到kafka的速度一致，这里可以优化
			//先把日志发往一个通道
			//kafka包中有一个单的groutine去处理
			kafka.SendToChan(t.topic,line.Text)
		}
	}
}


