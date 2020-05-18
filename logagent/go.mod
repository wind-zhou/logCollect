module github.com/wind-zhou/logagent

go 1.13

require github.com/wind-zhou/logagent/kafka v0.0.0

require (
	github.com/Shopify/toxiproxy v2.1.4+incompatible // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/frankban/quicktest v1.9.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/wind-zhou/logagent/etcd v0.0.0
	github.com/wind-zhou/logagent/taillog v0.0.0
	gopkg.in/ini.v1 v1.55.0

)

replace github.com/wind-zhou/logagent/kafka => ./kafka

replace github.com/wind-zhou/logagent/taillog => ./taillog

replace github.com/wind-zhou/logagent/etcd => ./etcd

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4

replace google.golang.org/grpc v1.29.1 => google.golang.org/grpc v1.26.0
