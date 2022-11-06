package tailfile

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"log-collect/config"
	"log-collect/pkg/kafka"
	"time"
)

var (
	TaskSlice []*tail.Tail
)

func Init(cfg *config.Config) error {
	tailCfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	for _, filename := range cfg.LogFilePath {
		TaskObj, err := tail.TailFile(filename, tailCfg)
		if err != nil {
			fmt.Println("init tailFile OBJ failed: " + filename)
			return err
		}
		TaskSlice = append(TaskSlice, TaskObj)
	}
	return nil
}

func Run() (err error) {
	var (
		line *tail.Line
		ok   bool
	)
	for {
		// 轮询tail对象切片
		for _, tails := range TaskSlice {
			// 开始读取数据
			line, ok = <-tails.Lines
			if !ok {
				fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
				time.Sleep(time.Second)
				continue
			}
			// 组装msg
			msg := &sarama.ProducerMessage{}
			msg.Topic = "scm"
			msg.Value = sarama.StringEncoder(line.Text)
			// 将msg丢到Kafka中的channel，在Kafka那边读取channel中的数据，再发往Kafka
			kafka.MsgChan <- msg
		}
	}
}
