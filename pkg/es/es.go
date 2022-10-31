package es

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type LogData struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

var (
	client *elastic.Client
	ch     chan *LogData
)

func Init(address string, chanSize int, worker int) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		return err
	}
	ch = make(chan *LogData, chanSize)
	// 开辟 worker 个协程读取 kafka 数据并发送至 ES
	for i := 0; i < worker; i++ {
		go SendToES()
	}
	return
}

func SendToESChan(msg *LogData) {
	ch <- msg
}

func SendToES() {
	// 链式操作
	for {
		select {
		case msg := <-ch:
			do, err := client.Index().
				Index(msg.Topic).
				BodyJson(msg).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Indexed %s to index %s, type %s", do.Id, do.Index, do.Type)
		default:
			time.Sleep(time.Second)
		}
	}
}
