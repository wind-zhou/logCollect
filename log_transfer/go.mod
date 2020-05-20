module github.com/wind-zhou/log_transfer

go 1.13

require (
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/wind-zhou/log_transfer/es v0.0.0
	github.com/wind-zhou/log_transfer/kafka v0.0.0
	gopkg.in/ini.v1 v1.56.0
)

replace (
	github.com/wind-zhou/log_transfer/es => ./es
	github.com/wind-zhou/log_transfer/kafka => ./kafka
)
