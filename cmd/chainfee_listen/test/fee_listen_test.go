package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/polynetwork/poly-nft-bridge/conf"
	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/polynetwork/poly-nft-bridge/dao/chainfeedao"
	"github.com/polynetwork/poly-nft-bridge/logic/fee"
)

func TestListenFee(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("current directory: %s\n", dir)
	config := conf.NewConfig("./../../conf/config_testnet.json")
	if config == nil {
		panic("read config failed!")
	}
	dao := chainfeedao.NewChainFeeDao(basedef.SERVER_STAKE, config.DBConfig)
	if dao == nil {
		panic("server is not valid")
	}
	feeListenCfgs := config.FeeListenConfig
	chainFees := make([]fee.ChainFee, 0)
	for _, cfg := range feeListenCfgs {
		chainFee := fee.NewChainFee(cfg, config.FeeUpdateSlot)
		chainFees = append(chainFees, chainFee)
	}
	feeListen := fee.NewFeeListen(config.FeeUpdateSlot, chainFees, dao)
	feeListen.ListenFee()
}
