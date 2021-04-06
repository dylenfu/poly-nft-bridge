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

package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/polynetwork/poly-nft-bridge/conf"
	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/polynetwork/poly-nft-bridge/dao/chainfeedao"
	"github.com/polynetwork/poly-nft-bridge/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSaveFee_SwapDao(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("current directory: %s\n", dir)
	config := conf.NewConfig("./../../conf/testnet.json")
	if config == nil {
		panic("read config failed!")
	}
	db := chainfeedao.NewChainFeeDao(basedef.SERVER_POLY_SWAP, config.DBConfig)
	if db == nil {
		panic("dao is invalid")
	}
	chainFees := make([]*models.ChainFee, 0)
	chainFeesJson := []byte(`[{"ChainId":2,"AssetBasicName":"Ethereum","AssetBasic":null,"MaxFee":1814309666000000000000,"MinFee":1814309666000000000000,"ProxyFee":2177171599200000000000,"Ind":1},{"ChainId":4,"AssetBasicName":"Neo","AssetBasic":null,"MaxFee":1000000000,"MinFee":1000000000,"ProxyFee":1000000000,"Ind":1},{"ChainId":8,"AssetBasicName":"Ethereum","AssetBasic":null,"MaxFee":0,"MinFee":0,"ProxyFee":0,"Ind":0}]`)
	err = json.Unmarshal(chainFeesJson, &chainFees)
	if err != nil {
		panic(err)
	}
	err = db.SaveFees(chainFees)
	if err != nil {
		panic(err)
	}
}

func TestQueryFees_SwapDao(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("current directory: %s\n", dir)
	config := conf.NewConfig("./../../conf/testnet.json")
	if config == nil {
		panic("read config failed!")
	}
	dbCfg := config.DBConfig
	db, err := gorm.Open(mysql.Open(dbCfg.User+":"+dbCfg.Password+"@tcp("+dbCfg.URL+")/"+
		dbCfg.Scheme+"?charset=utf8"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fees := make([]*models.ChainFee, 0)
	db.Debug().Model(&models.ChainFee{}).Find(&fees)
	json, _ := json.Marshal(fees)
	fmt.Printf("fees: %s\n", json)
}
