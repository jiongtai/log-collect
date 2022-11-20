package tailfile

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"log-collect/common"
	"log-collect/pkg/kafka"
	"time"
)

type tailTask struct {
	taskObj *tail.Tail
	path    string
	topic   string
}

func Init(logInfoList []common.LogInfo) (err error) {
	tailCfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	for _, log := range logInfoList {
		task := tailTask{
			path:  log.Path,
			topic: log.Topic,
		}
		task.taskObj, err = tail.TailFile(log.Path, tailCfg)
		if err != nil {
			fmt.Println("init tailFile OBJ failed: " + log.Path)
			continue
		}
		go run(task)
	}
	return nil
}

func run(task tailTask) (err error) {
	var (
		line *tail.Line
		ok   bool
	)
	for {
		// 开始读取数据
		line, ok = <-task.taskObj.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", task.path)
			time.Sleep(time.Second)
			continue
		}
		// 组装msg
		msg := &sarama.ProducerMessage{}
		msg.Topic = task.topic
		msg.Value = sarama.StringEncoder(line.Text)
		// 将msg丢到Kafka中的channel，在Kafka那边读取channel中的数据，再发往Kafka
		kafka.MsgChan <- msg
	}
}
