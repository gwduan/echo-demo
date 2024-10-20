package config

import (
	"encoding/json"
	"os"
)

type DemoConfig struct {
	ServerAddr   string `json:"server_addr"`
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

	SignKey:   "secret",
	VerifyKey: "secret",

	DbName:       "mysql",
	DbUrl:        "root:root@/echo_demo?charset=utf8&parseTime=True&loc=Local",
	RecordLimit:  5,
	RecordOffset: 0,
}

func ServerAddr() string {
	return config.ServerAddr
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
