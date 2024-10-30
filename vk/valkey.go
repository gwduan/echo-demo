package vk

import (
	"echo-demo/config"

	"github.com/valkey-io/valkey-go"
)

var client valkey.Client

func ClientInit() error {
	cli, err := valkey.NewClient(valkey.MustParseURL(config.ValkeyURL()))
	if err != nil {
		cli.Close()
		return err
	}

	client = cli

	return nil
}

func Client() valkey.Client {
	return client
}
