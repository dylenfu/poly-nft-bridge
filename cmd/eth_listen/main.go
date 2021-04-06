/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/astaxie/beego/logs"
	"github.com/polynetwork/poly-nft-bridge/conf"
	"github.com/polynetwork/poly-nft-bridge/dao/crosschaindao"
	wp "github.com/polynetwork/poly-nft-bridge/wrap"
	"github.com/urfave/cli"
)

var chainListen *wp.CrossChainListen

var (
	logLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~6). 0:Trace 1:Debug 2:Info 3:Warn 4:Error 5:Fatal 6:MaxLevel",
		Value: 1,
	}

	logDirFlag = cli.StringFlag{
		Name:  "logdir",
		Usage: "log directory",
		Value: "logs",
	}

	configPathFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Server config file `<path>`",
		Value: "config.json",
	}

	chainFlag = cli.UintFlag{
		Name:  "chain",
		Usage: "Set chain. 2:Eth, 6:Bsc, 7:Heco",
		Value: 2,
	}
)

//getFlagName deal with short flag, and return the flag name whether flag name have short name
func getFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "wrapper listen Service"
	app.Action = StartServer
	app.Version = "1.0.0"
	app.Copyright = "Copyright in 2019 The Ontology Authors"
	app.Flags = []cli.Flag{
		logLevelFlag,
		configPathFlag,
		logDirFlag,
		chainFlag,
	}
	app.Commands = []cli.Command{}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func StartServer(ctx *cli.Context) {
	for true {
		startServer(ctx)
		sig := waitSignal()
		stopServer()
		if sig != syscall.SIGHUP {
			break
		} else {
			continue
		}
	}
}

func startServer(ctx *cli.Context) {
	// instance beego log
	loglevel := ctx.GlobalUint64(getFlagName(logLevelFlag))
	logFormat := fmt.Sprintf(`{"level:":"%d"}`, loglevel)
	if err := logs.SetLogger("console", logFormat); err != nil {
		panic(fmt.Errorf("set logger failed, err: %v", err))
	}

	// load configuration
	configFile := ctx.GlobalString(getFlagName(configPathFlag))
	config := conf.NewConfig(configFile)
	if config == nil {
		panic("startServer - read config failed!")
	}

	// read chain
	chain := ctx.GlobalUint64(getFlagName(chainFlag))

	// generate dao
	db := crosschaindao.NewCrossChainDao(config.Server, config.Backup, config.DBConfig)
	if db == nil {
		panic("server is invalid")
	}

	// generate and starting listen handler
	chainListenConfig := config.GetChainListenConfig(chain)
	if chainListenConfig == nil {
		panic("chain is invalid")
	} else {
		enc, _ := json.Marshal(chainListenConfig)
		logs.Info("%s\n", string(enc))
	}

	chainHandler := wp.NewChainHandle(chainListenConfig)
	chainListen = wp.NewCrossChainListen(chainHandler, db)
	chainListen.Start()
}

func waitSignal() os.Signal {
	exit := make(chan os.Signal, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sc)
	go func() {
		for sig := range sc {
			logs.Info("cross chain listen received signal:(%s).", sig.String())
			exit <- sig
			close(exit)
			break
		}
	}()
	sig := <-exit
	return sig
}

func stopServer() {
	chainListen.Stop()
}

func main() {
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
