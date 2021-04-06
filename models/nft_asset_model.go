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

type NFTAssetBasic struct {
	Name    string      `gorm:"primaryKey;size:64;not null"`
	Time    int64       `gorm:"type:bigint(20);not null"`
	Disable int64       `gorm:"type:int;not null"`
	Assets  []*NFTAsset `gorm:"foreignKey:AssetBasicName;references:Name"`
}

type NFTAsset struct {
	Hash           string         `gorm:"primaryKey;size:66;not null"`
	ChainId        uint64         `gorm:"primaryKey;type:bigint(20);not null"`
	Name           string         `gorm:"size:64;not null"`
	BaseUri        string         `gorm:"type:varchar(128);not null"`
	AssetBasicName string         `gorm:"size:64;not null"`
	Disable        int64          `gorm:"type:int;not null"`
	AssetBasic     *NFTAssetBasic `gorm:"foreignKey:AssetBasicName;references:Name"`
	AssetMaps      []*NFTAssetMap `gorm:"foreignKey:SrcAssetHash,SrcChainId;references:Hash,ChainId"`
}

type NFTAssetMap struct {
	SrcChainId   uint64    `gorm:"primaryKey;type:bigint(20);not null"`
	SrcAssetHash string    `gorm:"primaryKey;size:66;not null"`
	SrcAsset     *NFTAsset `gorm:"foreignKey:SrcAssetHash,SrcChainId;references:Hash,ChainId"`
	DstChainId   uint64    `gorm:"primaryKey;type:bigint(20);not null"`
	DstAssetHash string    `gorm:"primaryKey;size:66;not null"`
	DstAsset     *NFTAsset `gorm:"foreignKey:DstAssetHash,DstChainId;references:Hash,ChainId"`
	Disable      int64     `gorm:"type:int;not null"`
}

type NFTToken struct {
	Hash       string         `gorm:"primaryKey;size:66;not null"`
	ChainId    uint64         `gorm:"primaryKey;type:bigint(20);not null"`
	TokenId    string         `gorm:"type:varchar(66);not null"`
	Url        string         `gorm:"type:varchar(128);not null"`
	TokenBasic *NFTAssetBasic `gorm:"foreignKey:AssetBasicName;references:Name"`
	TokenMaps  []*TokenMap    `gorm:"foreignKey:SrcAssetHash,SrcChainId;references:Hash,ChainId"`
}

type WrapperTransactionWithNFTToken struct {
	Hash         string  `gorm:"primaryKey;size:66;not null"`
	TokenID      string  `gorm:"type:varchar(66);not null"`
	User         string  `gorm:"size:64"`
	SrcChainId   uint64  `gorm:"type:bigint(20);not null"`
	BlockHeight  uint64  `gorm:"type:bigint(20);not null"`
	Time         uint64  `gorm:"type:bigint(20);not null"`
	DstChainId   uint64  `gorm:"type:bigint(20);not null"`
	DstUser      string  `gorm:"type:varchar(66);not null"`
	ServerId     uint64  `gorm:"type:bigint(20);not null"`
	FeeTokenHash string  `gorm:"size:66;not null"`
	FeeToken     *Token  `gorm:"foreignKey:FeeTokenHash,SrcChainId;references:Hash,ChainId"`
	FeeAmount    *BigInt `gorm:"type:varchar(64);not null"`
	Status       uint64  `gorm:"type:bigint(20);not null"`
}
