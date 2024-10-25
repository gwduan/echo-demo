package config

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type DemoConfig struct {
	ServerAddr   string `json:"server_addr"`
	AdminAddr    string `json:"admin_addr"`
	SignKey      string `json:"sign_key"`
	VerifyKey    string `json:"verify_key"`
	DbName       string `json:"db_name"`
	DbUrl        string `json:"db_url"`
	RecordLimit  int    `json:"record_limit"`
	RecordOffset int    `json:"record_offset"`
}

// Default values
var config = DemoConfig{
	ServerAddr: ":8080",
	AdminAddr:  ":8081",

	SignKey:   "secret",
	VerifyKey: "secret",

	DbName: "mysql",
	DbUrl:  "root:root@/echo_demo?charset=utf8&parseTime=True&loc=Local",

	RecordLimit:  5,
	RecordOffset: 0,
}

func ServerAddr() string {
	return config.ServerAddr
}

func AdminAddr() string {
	return config.AdminAddr
}

func SignKey() []byte {
	return []byte(config.SignKey)
}

func VerifyKey() []byte {
	return []byte(config.VerifyKey)
}

func DbName() string {
	return config.DbName
}

func DbUrl() string {
	return config.DbUrl
}

func RecordLimit() int {
	return config.RecordLimit
}

func RecordOffset() int {
	return config.RecordOffset
}

func Init() error {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &config); err != nil {
		return err
	}

	return nil
}

func Etcd(endpoints string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(endpoints, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()

	params := map[string]any{
		"server_addr":   &config.ServerAddr,
		"admin_addr":    &config.AdminAddr,
		"sign_key":      &config.SignKey,
		"verify_key":    &config.VerifyKey,
		"db_name":       &config.DbName,
		"db_url":        &config.DbUrl,
		"record_limit":  &config.RecordLimit,
		"record_offset": &config.RecordOffset,
	}
	for key, ptr := range params {
		val, err := getKey(cli, key)
		if err != nil {
			return err
		}
		if len(val) == 0 {
			continue
		}

		switch p := ptr.(type) {
		case *string:
			*p = val
		case *int:
			if num, err := strconv.Atoi(val); err == nil {
				*p = num
			}
		}

	}

	return nil
}

func getKey(cli *clientv3.Client, key string) (value string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		return "", err
	}

	for _, ev := range resp.Kvs {
		//fmt.Printf("%s: %s\n", ev.Key, ev.Value)
		return string(ev.Value), nil
	}

	return "", nil
}
