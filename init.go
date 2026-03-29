package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/filecoin-project/lotus/api/v0api"
	lcli "github.com/filecoin-project/lotus/cli"
	"github.com/urfave/cli/v2"
)

var LotusApi = "wss://node.glif.io/space06/lotus/"

var lapi v0api.FullNode
var ctx = context.Background()

var bootstrapTime = int64(1598306400)

var dateFormat = "2006-01-02"

func init() {
	// 禁用 glog 的标志解析
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	if api := os.Getenv("FULLNODE_API_INFO"); api == "" {
		err := os.Setenv("FULLNODE_API_INFO", LotusApi)
		if err != nil {
			log.Panicln(err)
		}
	}
	if e := os.Getenv("DATE_FORMAT"); e != "" {
		dateFormat = e
	}

	var err error
	lapi, _, err = lcli.GetFullNodeAPI(cli.NewContext(&cli.App{}, nil, nil))
	if err != nil {
		log.Panicln(err)
	}
}
