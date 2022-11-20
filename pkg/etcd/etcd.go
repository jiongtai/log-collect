package etcd

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"log-collect/common"
	"log-collect/config"
	"time"
)

var client *clientv3.Client

func Init(cfg *config.Config) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{cfg.Etcd.Address},
		DialTimeout: cfg.Etcd.Timeout * time.Second,
	})
	if err != nil {
		return err
	}
	return
}

func GetLogPathCfg(cfg *config.Config) (logInfoList []common.LogInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	getResponse, err := client.Get(ctx, cfg.Etcd.LogCollectKey)
	if err != nil {
		return
	}
	if len(getResponse.Kvs) == 0 {
		return
	}
	res := getResponse.Kvs[0]
	err = json.Unmarshal(res.Value, &logInfoList)
	if err != nil {
		return nil, err
	}
	return
}
