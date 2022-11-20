package logAgent

import (
	"fmt"
	"log-collect/config"
	"log-collect/pkg/etcd"
	"log-collect/pkg/kafka"
	"log-collect/pkg/tailfile"
)

func Run(cfg *config.Config) {
	// 初始化 Kafka
	err := kafka.Init(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 初始化 etcd 连接
	err = etcd.Init(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	logPathList, err := etcd.GetLogPathCfg(cfg)
	if err != nil {
		return
	}
	// 根据 etcd 中获取到的日志收集路径信息，初始化 tail
	err = tailfile.Init(logPathList)
	if err != nil {
		fmt.Println(err)
		return
	}
	//select {}
}
