package test

import (
	"fmt"
	"github.com/Shopify/sarama"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestKafkaProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有follower都回复ack，确保Kafka不会丢消息
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner // 对Key进行Hash，同样的Key每次都落到一个分区，这样消息是有序的

	// 使用同步producer，异步模式下有更高的性能，但是处理更复杂，这里建议先从简单的入手
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	defer func() {
		_ = producer.Close()
	}()
	if err != nil {
		panic(err.Error())
	}

	msgCount := 4
	// 模拟4个消息
	for i := 0; i < msgCount; i++ {
		rand.Seed(int64(time.Now().Nanosecond()))
		msg := &sarama.ProducerMessage{
			Topic: "testAutoSyncOffset",
			Value: sarama.StringEncoder("hello+" + strconv.Itoa(rand.Int())),
		}

		t1 := time.Now().Nanosecond()
		partition, offset, err := producer.SendMessage(msg)
		t2 := time.Now().Nanosecond()

		if err == nil {
			fmt.Println("produce success, partition:", partition, ",offset:", offset, ",cost:", (t2-t1)/(1000*1000), " ms")
		} else {
			fmt.Println(err.Error())
		}
	}
}

func TestKafkaConsumer(t *testing.T) {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = consumer.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	partitionConsumer, err := consumer.ConsumePartition("business_log", 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = partitionConsumer.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	for msg := range partitionConsumer.Messages() {
		fmt.Printf("partition:%d offset:%d key:%s val:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
	}
}
