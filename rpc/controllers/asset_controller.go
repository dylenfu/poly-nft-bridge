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
	"github.com/astaxie/beego"
	"github.com/polynetwork/poly-nft-bridge/models"
)

type AssetController struct {
	beego.Controller
}

func (c *AssetController) Assets() {
	var req models.NFTAssetsReq
	if !input(&c.Controller, &req) {
		return
	}

	assets := make([]*models.NFTAsset, 0)
	db.Where("chain_id = ?", req.ChainId).
		Preload("AssetBasic").
		Preload("AssetMaps").
		Preload("AssetMaps.DstAsset").
		Find(&assets)
	data := models.MakeNFTAssetsRsp(assets)

	output(&c.Controller, data)
}

func (c *AssetController) Asset() {
	var req models.NFTAssetReq
	if !input(&c.Controller, &req) {
		return
	}

	asset := new(models.NFTAsset)
	res := db.Where("hash = ? and chain_id = ?", req.Hash, req.ChainId).
		Preload("AssetBasic").
		Preload("AssetMaps").
		Preload("AssetMaps.DstAsset").
		First(asset)
	if res.RowsAffected == 0 {
		notExist(&c.Controller)
		return
	}
	output(&c.Controller, asset)
}

func (c *AssetController) AssetBasics() {
	assetBasics := make([]*models.NFTAssetBasic, 0)
	db.Model(&models.NFTAssetBasic{}).Preload("Assets").Find(&assetBasics)
	data := models.MakeNFTAssetBasicsRsp(assetBasics)
	output(&c.Controller, data)
}
