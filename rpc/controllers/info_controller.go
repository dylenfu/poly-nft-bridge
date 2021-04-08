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
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/astaxie/beego"
	"github.com/polynetwork/poly-nft-bridge/models"
	"github.com/polynetwork/poly-nft-bridge/utils/net"
)

var (
	mode string
	port int
)

type InfoController struct {
	beego.Controller
}

func (c *InfoController) Get() {
	url, err := captureUrl()
	if err != nil {
		c.Data["json"] = models.MakeErrorRsp(err.Error())
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}
	explorer := &models.PolyBridgeInfoResp{
		Version: "v1",
		URL:     url,
	}
	c.Data["json"] = explorer
	c.ServeJSON()
}

func (c *InfoController) Home() {
	var req models.HomeReq
	if !input(&c.Controller, &req) {
		return
	}

	sdk := selectNode(req.ChainId)
	if sdk == nil {
		customInput(&c.Controller, ErrCodeRequest, "chain id not exist")
		return
	}

	assets := make([]*models.NFTAsset, 0)
	db.Where("chain_id = ?", req.ChainId).Find(&assets)
	totalCnt := len(assets)

	assetItems := make([]*models.AssetItems, 0)
	for _, v := range assets {
		addr := common.HexToAddress(v.Hash)
		tokenIds, _ := sdk.GetAssetNFTs(addr, 0, req.Size)
		items := make([]*models.Item, 0)
		if len(tokenIds) > 0 {
			urlmap, _ := sdk.GetNFTURLs(addr, tokenIds)
			if len(urlmap) > 0 {
				for tkid, url := range urlmap {
					items = append(items, &models.Item{
						TokenId: tkid,
						Url:     url,
					})
				}
			}
		}
		assetItem := &models.AssetItems{
			Asset: v,
			Items: items,
		}
		assetItems = append(assetItems, assetItem)
	}
	data := models.MakeHomeRsp(totalCnt, assetItems)
	output(&c.Controller, data)
}

func SetBaseInfo(_mode string, _port int) {
	mode = _mode
	port = _port
}

func captureUrl() (string, error) {
	switch mode {
	case "dev", "test":
		ips, err := net.GetLocalIPv4s()
		if err != nil {
			return "", err
		}
		if len(ips) == 0 {
			return "", fmt.Errorf("local IPv4s reading error")
		}
		return fmt.Sprintf("http://%s:%d/nft", ips[0], port), nil
	case "prod":
		return "http://bridge.poly.network/nft", nil
	}
	return "", fmt.Errorf("beego running mode invalid")
}
