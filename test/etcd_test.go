package test

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"strconv"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	putResponse, err := client.Put(ctx, "age", strconv.Itoa(25))
	cancel()
	if err != nil {
		return
	}
	fmt.Println(putResponse)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	getResponse, err := client.Get(ctx, "age")
	if err != nil {
		return
	}
	fmt.Println(getResponse)
	cancel()
	for _, ev := range getResponse.Kvs {
		fmt.Println(ev)
	}
}

func TestEtcdWatch(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	defer client.Close()
	wc := client.Watch(context.Background(), "name")
	for ch := range wc {
		for _, event := range ch.Events {
			fmt.Printf("type: %s ; key: %s ; value: %s\n", event.Type, event.Kv.Key, event.Kv.Value)
		}
	}
}
