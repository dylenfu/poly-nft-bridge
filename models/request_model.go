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

package models

type NFTAssetBasicReq struct {
	Name string
}

type NFTAssetBasicRsp struct {
	Name    string
	Time    int64
	Disable int64
	Assets  []*NFTAssetRsp
}

func MakeAssetBasicRsp(assetBasic *NFTAssetBasic) *NFTAssetBasicRsp {
	rsp := &NFTAssetBasicRsp{
		Name:    assetBasic.Name,
		Time:    assetBasic.Time,
		Disable: assetBasic.Disable,
		Assets:  nil,
	}
	if assetBasic.Assets != nil {
		for _, asset := range assetBasic.Assets {
			rsp.Assets = append(rsp.Assets, MakeNFTAssetRsp(asset))
		}
	}
	return rsp
}

type NFTAssetBasicsReq struct {
}

type NFTAssetBasicsRsp struct {
	TotalCount  uint64
	AssetBasics []*NFTAssetBasicRsp
}

func MakeNFTAssetBasicsRsp(tokenBasics []*NFTAssetBasic) *NFTAssetBasicsRsp {
	rsp := &NFTAssetBasicsRsp{
		TotalCount: uint64(len(tokenBasics)),
	}
	for _, assetBasic := range tokenBasics {
		if assetBasic.Disable == 0 {
			rsp.AssetBasics = append(rsp.AssetBasics, MakeAssetBasicRsp(assetBasic))
		}
	}
	return rsp
}

type NFTAssetReq struct {
	ChainId uint64
	Hash    string
}

type NFTAssetRsp struct {
	Hash           string
	ChainId        uint64
	Name           string
	Disable        int64
	BaseUri        string
	AssetBasicName string
	AssetBasic     *NFTAssetBasicRsp
	AssetMaps      []*NFTAssetMapRsp
}

func MakeNFTAssetRsp(asset *NFTAsset) *NFTAssetRsp {
	rsp := &NFTAssetRsp{
		Hash:           asset.Hash,
		ChainId:        asset.ChainId,
		Name:           asset.Name,
		BaseUri:        asset.BaseUri,
		AssetBasicName: asset.AssetBasicName,
		Disable:        asset.Disable,
	}
	if asset.AssetBasic != nil {
		rsp.AssetBasic = MakeAssetBasicRsp(asset.AssetBasic)
	}
	if asset.AssetMaps != nil {
		for _, m := range asset.AssetMaps {
			rsp.AssetMaps = append(rsp.AssetMaps, MakeNFTAssetMapRsp(m))
		}
	}
	return rsp
}

type NFTAssetsReq struct {
	ChainId uint64
}

type NFTAssetsRsp struct {
	TotalCount uint64
	Assets     []*NFTAssetRsp
}

func MakeNFTAssetsRsp(assets []*NFTAsset) *NFTAssetsRsp {
	tokensRsp := &NFTAssetsRsp{
		TotalCount: uint64(len(assets)),
	}
	for _, asset := range assets {
		tokensRsp.Assets = append(tokensRsp.Assets, MakeNFTAssetRsp(asset))
	}
	return tokensRsp
}

type NFTAssetMapReq struct {
	ChainId uint64
	Hash    string
}

type NFTAssetMapRsp struct {
	SrcTokenHash string
	SrcToken     *NFTAssetRsp
	DstTokenHash string
	DstToken     *NFTAssetRsp
	Disable      int64
}

func MakeNFTAssetMapRsp(assetMap *NFTAssetMap) *NFTAssetMapRsp {
	rsp := &NFTAssetMapRsp{
		SrcTokenHash: assetMap.SrcAssetHash,
		DstTokenHash: assetMap.DstAssetHash,
		Disable:      assetMap.Disable,
	}
	if assetMap.SrcAsset != nil {
		rsp.SrcToken = MakeNFTAssetRsp(assetMap.SrcAsset)
	}
	if assetMap.DstAsset != nil {
		rsp.DstToken = MakeNFTAssetRsp(assetMap.DstAsset)
	}
	return rsp
}

type NFTAssetMapsReq struct {
	ChainId uint64
	Hash    string
}

type NFTAssetMapsRsp struct {
	TotalCount uint64
	AssetMaps  []*NFTAssetMapRsp
}

func MakeNFTAssetMapsRsp(assetMaps []*NFTAssetMap) *NFTAssetMapsRsp {
	rsp := &NFTAssetMapsRsp{
		TotalCount: uint64(len(assetMaps)),
	}
	for _, assetMap := range assetMaps {
		rsp.AssetMaps = append(rsp.AssetMaps, MakeNFTAssetMapRsp(assetMap))
	}
	return rsp
}
