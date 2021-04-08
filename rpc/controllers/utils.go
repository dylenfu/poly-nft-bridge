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

package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/polynetwork/poly-nft-bridge/conf"
	"github.com/polynetwork/poly-nft-bridge/models"
	"github.com/polynetwork/poly-nft-bridge/sdk/eth_sdk"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   = newDB()
	sdks = make(map[uint64]*eth_sdk.EthereumSdkPro)
)

func newDB() *gorm.DB {
	user := beego.AppConfig.String("mysqluser")
	password := beego.AppConfig.String("mysqlpass")
	url := beego.AppConfig.String("mysqlurls")
	scheme := beego.AppConfig.String("mysqldb")
	mode := beego.AppConfig.String("runmode")
	Logger := logger.Default
	if mode == "dev" {
		Logger = Logger.LogMode(logger.Info)
	}
	db, err := gorm.Open(mysql.Open(user+":"+password+"@tcp("+url+")/"+scheme+"?charset=utf8"), &gorm.Config{Logger: Logger})
	if err != nil {
		panic(err)
	}
	return db
}

func Initialize(c *conf.Config) {
	//var err error
	//Logger := logger.Default
	//if c.DBConfig.Debug {
	//	Logger = Logger.LogMode(logger.Info)
	//}
	//link := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
	//	c.DBConfig.User,
	//	c.DBConfig.Password,
	//	c.DBConfig.URL,
	//	c.DBConfig.Scheme,
	//)
	//if db, err = gorm.Open(mysql.Open(link), &gorm.Config{Logger: Logger}); err != nil {
	//	panic(err)
	//}

	for _, v := range c.ChainListenConfig {
		pro := eth_sdk.NewEthereumSdkPro(v.GetNodesUrl(), v.ListenSlot, v.ChainId)
		sdks[v.ChainId] = pro
	}
}

func selectNode(chainID uint64) *eth_sdk.EthereumSdkPro {
	pro, ok := sdks[chainID]
	if !ok {
		return nil
	}
	return pro
}

const (
	ErrCodeRequest     int = 400
	ErrCodeNotExist    int = 404
	ErrCodeNodeInvalid int = 500
)

var errMap = map[int]string{
	ErrCodeRequest:     "request parameter is invalid!",
	ErrCodeNotExist:    "not found",
	ErrCodeNodeInvalid: "blockchain node exception",
}

func input(c *beego.Controller, req interface{}) bool {
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		code := ErrCodeRequest
		customInput(c, code, errMap[code])
		return false
	} else {
		return true
	}
}

func customInput(c *beego.Controller, code int, msg string) {
	c.Data["json"] = models.MakeErrorRsp(msg)
	c.Ctx.ResponseWriter.WriteHeader(code)
	c.ServeJSON()
}

func notExist(c *beego.Controller) {
	code := ErrCodeNotExist
	c.Data["json"] = models.MakeErrorRsp(errMap[code])
	c.Ctx.ResponseWriter.WriteHeader(code)
	c.ServeJSON()
}

func nodeInvalid(c *beego.Controller) {
	code := ErrCodeNodeInvalid
	c.Data["json"] = models.MakeErrorRsp(errMap[code])
	c.Ctx.ResponseWriter.WriteHeader(code)
	c.ServeJSON()
}

func output(c *beego.Controller, data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}

func getPageNo(totalNo, pageSize int) int {
	return (int(totalNo) + pageSize - 1) / pageSize
}
