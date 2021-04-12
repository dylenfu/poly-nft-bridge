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

func (dao *SwapDao) AddAssets(basics []*models.TokenBasic) error {
	if basics != nil && len(basics) > 0 {
		res := dao.db.Save(basics)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("add basics failed!")
		}
	}

	maps := getAssetMapsFromAsset(basics)
	if maps != nil && len(maps) > 0 {
		res := dao.db.Save(maps)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("add basics map failed!")
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

func (dao *SwapDao) RemoveAsset(name string) error {
	basic := new(models.TokenBasic)
	res := dao.db.Model(&models.TokenBasic{}).Where("name = ?", name).Preload("Tokens").First(basic)
	if res.Error != nil {
		return res.Error
	}

	basics := []*models.TokenBasic{basic}
	maps := getAssetMapsFromAsset(basics)
	for _, mp := range maps {
		dao.db.Where("src_chain_id = ? and src_token_hash = ? and dst_chain_id = ? and dst_token_hash = ?",
			mp.SrcChainId,
			strings.ToLower(mp.SrcTokenHash),
			mp.DstChainId,
			strings.ToLower(mp.DstTokenHash),
		).Delete(&models.TokenMap{})
	}
	for _, asset := range basic.Tokens {
		dao.db.Where("hash = ? and chain_id = ?", asset.Hash, asset.ChainId).Delete(&models.Token{})
	}
	dao.db.Where("name = ?", basic.Name).Delete(&models.NFTToken{})
	return nil
}

func (dao *SwapDao) RemoveAssetMaps(maps []*models.TokenMap) error {
	for _, mp := range maps {
		dao.db.Model(&models.TokenMap{}).
			Where("src_chain_id = ? and src_token_hash = ? and dst_chain_id = ? and dst_token_hash = ?",
				mp.SrcChainId,
				strings.ToLower(mp.SrcTokenHash),
				mp.DstChainId,
				strings.ToLower(mp.DstTokenHash),
			).Update("property", 0)
	}
	return nil
}

func getAssetMapsFromAsset(basics []*models.TokenBasic) []*models.TokenMap {
	maps := make([]*models.TokenMap, 0)
	for _, basic := range basics {
		for _, src := range basic.Tokens {
			for _, dst := range basic.Tokens {
				if dst.ChainId != src.ChainId {
					maps = append(maps, &models.TokenMap{
						SrcChainId:   src.ChainId,
						SrcTokenHash: src.Hash,
						DstChainId:   dst.ChainId,
						DstTokenHash: dst.Hash,
						Property:     1,
					})
				}
			}
		}
	}
	return maps
}
