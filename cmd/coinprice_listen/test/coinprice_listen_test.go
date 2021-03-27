package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/polynetwork/poly-bridge/conf"
	basedef "github.com/polynetwork/poly-bridge/const"
	"github.com/polynetwork/poly-bridge/dao/coinpricedao"
	priceListen "github.com/polynetwork/poly-bridge/logic/price"
)

func TestListenCoinPrice(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("current directory: %s\n", dir)
	config := conf.NewConfig("./../../conf/config_testnet.json")
	if config == nil {
		panic("read config failed!")
	}
	dao := coinpricedao.NewCoinPriceDao(basedef.SERVER_STAKE, config.DBConfig)
	if dao == nil {
		panic("server is not valid")
	}
	priceListenConfig := config.CoinPriceListenConfig
	priceMarkets := make([]priceListen.PriceMarket, 0)
	for _, cfg := range priceListenConfig {
		priceMarket := priceListen.NewPriceMarket(cfg)
		priceMarkets = append(priceMarkets, priceMarket)
	}
	cpListen := priceListen.NewCoinPriceListen(config.CoinPriceUpdateSlot, priceMarkets, dao)
	cpListen.ListenPrice()
}
