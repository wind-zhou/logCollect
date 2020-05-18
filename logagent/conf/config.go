package conf

type AppConf struct {
	KafkaConf   `ini:"kafka"`
	EtcdConf    `ini:"etcd"`
}

type  EtcdConf struct{

	Address string`ini:"address"`
	Timeout  int `ini:"timeout"`
	Key string `ini:"key"`
}


type KafkaConf struct {
	Address string `ini:"address"`
	Chanmaxsize int `int:"chan_max_size"`
 	//Topic   string `ini:"topic"`
}

//-------unused-------
type TaillogConf struct {
	FileName string `ini:"filename"`
}