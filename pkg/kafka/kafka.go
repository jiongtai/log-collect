package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log-collect/config"
)

var (
	client  sarama.SyncProducer
	MsgChan chan *sarama.ProducerMessage
)

func Init(cfg *config.Config) (err error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll // 等待所有follower都回复ack，确保Kafka不会丢消息
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Partitioner = sarama.NewHashPartitioner // 对Key进行Hash，同样的Key每次都落到一个分区，这样消息是有序的
	// 使用同步producer，异步模式下有更高的性能，但是处理更复杂，这里建议先从简单的入手
	client, err = sarama.NewSyncProducer(cfg.Kafka.Address, saramaConfig)
	if err != nil {
		panic(err.Error())
	}
	MsgChan = make(chan *sarama.ProducerMessage, cfg.Kafka.ChanMaxSize)
	go run()
	return
}

func run() {
	for true {
		select {
		case msg := <-MsgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				return
			}
			fmt.Printf("send msg success, pid: %d, offset: %d, msg: %s", pid, offset, msg.Value)
		}
	}
}
