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

package swapdao

import (
	"fmt"
	"strings"

	"github.com/polynetwork/poly-nft-bridge/models"
)

func (dao *SwapDao) AddAssets(assetBasics []*models.NFTAssetBasic) error {
	if assetBasics != nil && len(assetBasics) > 0 {
		res := dao.db.Save(assetBasics)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("add assetBasics failed!")
		}
	}

	addAssetMaps := getAssetMapsFromAsset(assetBasics)
	if addAssetMaps != nil && len(addAssetMaps) > 0 {
		res := dao.db.Save(addAssetMaps)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("add assetBasics map failed!")
		}
	}
	return nil
}

func (dao *SwapDao) RemoveAssets(assets []string) error {
	for _, asset := range assets {
		if err := dao.RemoveAsset(asset); err != nil {
			return err
		}
	}
	return nil
}

func (dao *SwapDao) RemoveAsset(asset string) error {
	assetBasic := new(models.NFTAssetBasic)
	res := dao.db.Model(&models.NFTAssetBasic{}).Where("name = ?", asset).Preload("Assets").First(assetBasic)
	if res.Error != nil {
		return res.Error
	}

	assetBasics := []*models.NFTAssetBasic{assetBasic}
	assetMaps := getAssetMapsFromAsset(assetBasics)
	for _, assetMap := range assetMaps {
		dao.db.Where("src_chain_id = ? and src_asset_hash = ? and dst_chain_id = ? and dst_asset_hash = ?",
			assetMap.SrcChainId,
			strings.ToLower(assetMap.SrcAssetHash),
			assetMap.DstChainId,
			strings.ToLower(assetMap.DstAssetHash),
		).Delete(&models.NFTAssetMap{})
	}
	for _, asset := range assetBasic.Assets {
		dao.db.Where("hash = ? and chain_id = ?", asset.Hash, asset.ChainId).Delete(&models.NFTAsset{})
	}
	dao.db.Where("name = ?", assetBasic.Name).Delete(&models.NFTAssetBasic{})
	return nil
}

func (dao *SwapDao) RemoveAssetMaps(assetMaps []*models.NFTAssetMap) error {
	for _, assetMap := range assetMaps {
		dao.db.Model(&models.NFTAssetMap{}).
			Where("src_chain_id = ? and src_asset_hash = ? and dst_chain_id = ? and dst_asset_hash = ?",
				assetMap.SrcChainId,
				strings.ToLower(assetMap.SrcAssetHash),
				assetMap.DstChainId,
				strings.ToLower(assetMap.DstAssetHash),
			).Update("disable", 1)
	}
	return nil
}

func getAssetMapsFromAsset(assetBasics []*models.NFTAssetBasic) []*models.NFTAssetMap {
	assetMaps := make([]*models.NFTAssetMap, 0)
	for _, assetBasic := range assetBasics {
		for _, src := range assetBasic.Assets {
			for _, dst := range assetBasic.Assets {
				if dst.ChainId != src.ChainId {
					assetMaps = append(assetMaps, &models.NFTAssetMap{
						SrcChainId:   src.ChainId,
						SrcAssetHash: src.Hash,
						DstChainId:   dst.ChainId,
						DstAssetHash: dst.Hash,
						Disable:      0,
					})
				}
			}
		}
	}
	return assetMaps
}
