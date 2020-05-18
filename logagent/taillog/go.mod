module github.com/wind-zhou/logagent/taillog

go 1.13

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/hpcloud/tail v1.0.0
	golang.org/x/sys v0.0.0-20200501145240-bc7a7d42d5c3 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	github.com/wind-zhou/logagent/kafka v0.0.0
)

replace github.com/wind-zhou/logagent/kafka => ./kafka
