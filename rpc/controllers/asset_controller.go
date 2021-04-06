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
	"fmt"

	"github.com/astaxie/beego"
	"github.com/polynetwork/poly-nft-bridge/models"
)

type AssetController struct {
	beego.Controller
}

func (c *AssetController) Assets() {
	var req models.NFTAssetsReq
	if err := input(&c.Controller, &req); err != nil {
		return
	}

	assets := make([]*models.NFTAsset, 0)
	db.Where("chain_id = ?", req.ChainId).
		Preload("AssetBasic").
		Preload("AssetMaps").
		Preload("AssetMaps.DstToken").
		Find(&assets)
	data := models.MakeNFTAssetsRsp(assets)

	output(&c.Controller, data)
}

func (c *AssetController) Asset() {
	var tokenReq models.TokenReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &tokenReq); err != nil {
		c.Data["json"] = models.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}
	token := new(models.Token)
	res := db.Where("hash = ? and chain_id = ?", tokenReq.Hash, tokenReq.ChainId).Preload("AssetBasic").Preload("AssetMaps").Preload("AssetMaps.DstToken").First(token)
	if res.RowsAffected == 0 {
		c.Data["json"] = models.MakeErrorRsp(fmt.Sprintf("token: (%s,%d) does not exist", tokenReq.Hash, tokenReq.ChainId))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
		return
	}
	c.Data["json"] = models.MakeTokenRsp(token)
	c.ServeJSON()
}

func (c *AssetController) AssetBasics() {
	var tokenBasicReq models.TokenBasicReq
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &tokenBasicReq); err != nil {
		c.Data["json"] = models.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.ServeJSON()
	}
	tokenBasics := make([]*models.TokenBasic, 0)
	db.Model(&models.TokenBasic{}).Preload("Assets").Find(&tokenBasics)
	c.Data["json"] = models.MakeTokenBasicsRsp(tokenBasics)
	c.ServeJSON()
}
