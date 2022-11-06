package logAgent

import (
	"fmt"
	"log-collect/config"
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
	// 初始化 tail
	err = tailfile.Init(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 收集日志发往 Kafka
	err = tailfile.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

}
