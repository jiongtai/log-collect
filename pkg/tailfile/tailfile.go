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
	taskObj *tail.Tail
)

func Init(cfg *config.Config) (err error) {
	tailCfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	taskObj, err = tail.TailFile(cfg.LogFilePath, tailCfg)
	if err != nil {
		fmt.Println("init tailFile OBJ failed: " + cfg.LogFilePath)
		return err
	}
	return nil
}

func Run() (err error) {
	var (
		line *tail.Line
		ok   bool
	)
	for {
		// 开始读取数据
		line, ok = <-taskObj.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", taskObj.Filename)
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
