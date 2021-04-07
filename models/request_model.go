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

import (
	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/polynetwork/poly-nft-bridge/utils/decimal"
	"math/big"
)

type PolyBridgeInfoResp struct {
	Version string
	URL     string
}

type ErrorRsp struct {
	Message string
}

func MakeErrorRsp(messgae string) *ErrorRsp {
	errorRsp := &ErrorRsp{
		Message: messgae,
	}
	return errorRsp
}

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

type PriceMarketRsp struct {
	TokenBasicName string
	MarketName     string
	Name           string
	Price          string
	Ind            uint64
	Time           int64
	TokenBasic     *TokenBasicRsp
}

func MakePriceMarketRsp(priceMarket *PriceMarket) *PriceMarketRsp {
	price := new(big.Float).Quo(new(big.Float).SetInt64(priceMarket.Price), new(big.Float).SetInt64(basedef.PRICE_PRECISION))
	priceMarketRsp := &PriceMarketRsp{
		TokenBasicName: priceMarket.TokenBasicName,
		MarketName:     priceMarket.MarketName,
		Name:           priceMarket.Name,
		Price:          price.String(),
		Ind:            priceMarket.Ind,
		Time:           priceMarket.Time,
	}
	if priceMarket.TokenBasic != nil {
		priceMarketRsp.TokenBasic = MakeTokenBasicRsp(priceMarket.TokenBasic)
	}
	return priceMarketRsp
}

type TokenBasicReq struct {
	Name string
}

type TokenBasicRsp struct {
	Name         string
	Precision    uint64
	Price        string
	Ind          uint64
	Time         int64
	Property     int64
	PriceMarkets []*PriceMarketRsp
	Tokens       []*TokenRsp
}

func MakeTokenBasicRsp(tokenBasic *TokenBasic) *TokenBasicRsp {
	price := new(big.Float).Quo(new(big.Float).SetInt64(tokenBasic.Price), new(big.Float).SetInt64(basedef.PRICE_PRECISION))
	tokenBasicRsp := &TokenBasicRsp{
		Name:      tokenBasic.Name,
		Time:      tokenBasic.Time,
		Precision: tokenBasic.Precision,
		Price:     price.String(),
		Ind:       tokenBasic.Ind,
		Property:  tokenBasic.Property,
		Tokens:    nil,
	}
	if tokenBasic.Tokens != nil {
		for _, token := range tokenBasic.Tokens {
			tokenBasicRsp.Tokens = append(tokenBasicRsp.Tokens, MakeTokenRsp(token))
		}
	}
	if tokenBasic.PriceMarkets != nil {
		for _, priceMarket := range tokenBasic.PriceMarkets {
			tokenBasicRsp.PriceMarkets = append(tokenBasicRsp.PriceMarkets, MakePriceMarketRsp(priceMarket))
		}
	}
	return tokenBasicRsp
}

type TokenBasicsReq struct {
}

type TokenBasicsRsp struct {
	TotalCount  uint64
	TokenBasics []*TokenBasicRsp
}

func MakeTokenBasicsRsp(tokenBasics []*TokenBasic) *TokenBasicsRsp {
	tokenBasicsRsp := &TokenBasicsRsp{
		TotalCount: uint64(len(tokenBasics)),
	}
	for _, tokenBasic := range tokenBasics {
		if tokenBasic.Property == 1 {
			tokenBasicsRsp.TokenBasics = append(tokenBasicsRsp.TokenBasics, MakeTokenBasicRsp(tokenBasic))
		}
	}
	return tokenBasicsRsp
}

type TokenReq struct {
	ChainId uint64
	Hash    string
}

type TokenRsp struct {
	Hash           string
	ChainId        uint64
	Name           string
	Property       int64
	TokenBasicName string
	TokenBasic     *TokenBasicRsp
	TokenMaps      []*TokenMapRsp
}

func MakeTokenRsp(token *Token) *TokenRsp {
	tokenRsp := &TokenRsp{
		Hash:           token.Hash,
		ChainId:        token.ChainId,
		Name:           token.Name,
		TokenBasicName: token.TokenBasicName,
		Property:       token.Property,
	}
	if token.TokenBasic != nil {
		tokenRsp.TokenBasic = MakeTokenBasicRsp(token.TokenBasic)
	}
	if token.TokenMaps != nil {
		for _, tokenmap := range token.TokenMaps {
			tokenRsp.TokenMaps = append(tokenRsp.TokenMaps, MakeTokenMapRsp(tokenmap))
		}
	}
	return tokenRsp
}

type TokenMapReq struct {
	ChainId uint64
	Hash    string
}

type TokenMapRsp struct {
	SrcTokenHash string
	SrcToken     *TokenRsp
	DstTokenHash string
	DstToken     *TokenRsp
	Property     int64
}

func MakeTokenMapRsp(tokenMap *TokenMap) *TokenMapRsp {
	tokenMapRsp := &TokenMapRsp{
		SrcTokenHash: tokenMap.SrcTokenHash,
		DstTokenHash: tokenMap.DstTokenHash,
		Property:     tokenMap.Property,
	}
	if tokenMap.SrcToken != nil {
		tokenMapRsp.SrcToken = MakeTokenRsp(tokenMap.SrcToken)
	}
	if tokenMap.DstToken != nil {
		tokenMapRsp.DstToken = MakeTokenRsp(tokenMap.DstToken)
	}
	return tokenMapRsp
}

type TokenMapsReq struct {
	ChainId uint64
	Hash    string
}

type TokenMapsRsp struct {
	TotalCount uint64
	TokenMaps  []*TokenMapRsp
}

func MakeTokenMapsRsp(tokenMaps []*TokenMap) *TokenMapsRsp {
	tokenMapsRsp := &TokenMapsRsp{
		TotalCount: uint64(len(tokenMaps)),
	}
	for _, tokenMap := range tokenMaps {
		tokenMapsRsp.TokenMaps = append(tokenMapsRsp.TokenMaps, MakeTokenMapRsp(tokenMap))
	}
	return tokenMapsRsp
}

//type GetFeeReq struct {
//	SrcChainId uint64
//	Hash       string
//	DstChainId uint64
//}
//
//type GetFeeRsp struct {
//	SrcChainId               uint64
//	Hash                     string
//	DstChainId               uint64
//	UsdtAmount               string
//	TokenAmount              string
//	TokenAmountWithPrecision string
//}
//
//func MakeGetFeeRsp(srcChainId uint64, hash string, dstChainId uint64, usdtAmount *big.Float, tokenAmount *big.Float, tokenAmountWithPrecision *big.Float) *GetFeeRsp {
//	getFeeRsp := &GetFeeRsp{
//		SrcChainId:               srcChainId,
//		Hash:                     hash,
//		DstChainId:               dstChainId,
//		UsdtAmount:               usdtAmount.String(),
//		TokenAmount:              tokenAmount.String(),
//		TokenAmountWithPrecision: tokenAmountWithPrecision.String(),
//	}
//	{
//		aaa, _ := usdtAmount.Float64()
//		usdtAmount := decimal.NewFromFloat(aaa)
//		getFeeRsp.UsdtAmount = usdtAmount.String()
//	}
//	{
//		precision := decimal.NewFromInt(basedef.PRICE_PRECISION)
//		aaa := new(big.Float).Mul(tokenAmount, new(big.Float).SetInt64(basedef.PRICE_PRECISION))
//		bbb, _ := aaa.Int64()
//		ccc := decimal.NewFromInt(bbb + 1)
//		tokenAmount := ccc.Div(precision)
//		getFeeRsp.TokenAmount = tokenAmount.String()
//	}
//	{
//		aaa, _ := tokenAmountWithPrecision.Float64()
//		tokenAmountWithPrecision := decimal.NewFromFloat(aaa)
//		getFeeRsp.TokenAmountWithPrecision = tokenAmountWithPrecision.String()
//	}
//	return getFeeRsp
//}
//
//type CheckFeeReq struct {
//	Hash    string
//	ChainId uint64
//}
//
//type CheckFeeRsp struct {
//	ChainId     uint64
//	Hash        string
//	PayState    int
//	Amount      string
//	MinProxyFee string
//}
//
//type CheckFeesReq struct {
//	Checks []*CheckFeeReq
//}
//
//type CheckFeesRsp struct {
//	TotalCount uint64
//	CheckFees  []*CheckFeeRsp
//}
//
//func MakeCheckFeesRsp(checkFees []*CheckFee) *CheckFeesRsp {
//	checkFeesRsp := &CheckFeesRsp{
//		TotalCount: uint64(len(checkFees)),
//	}
//	for _, checkFee := range checkFees {
//		checkFeesRsp.CheckFees = append(checkFeesRsp.CheckFees, MakeCheckFeeRsp(checkFee))
//	}
//	return checkFeesRsp
//}
//
//func MakeCheckFeeRsp(checkFee *CheckFee) *CheckFeeRsp {
//	checkFeeRsp := &CheckFeeRsp{
//		ChainId:     checkFee.ChainId,
//		Hash:        checkFee.Hash,
//		PayState:    checkFee.PayState,
//		Amount:      checkFee.Amount.String(),
//		MinProxyFee: checkFee.MinProxyFee.String(),
//	}
//	{
//		aaa, _ := checkFee.Amount.Float64()
//		bbb := decimal.NewFromFloat(aaa)
//		checkFeeRsp.Amount = bbb.String()
//	}
//	{
//		aaa, _ := checkFee.MinProxyFee.Float64()
//		bbb := decimal.NewFromFloat(aaa)
//		checkFeeRsp.MinProxyFee = bbb.String()
//	}
//	return checkFeeRsp
//}

type WrapperTransactionReq struct {
	Hash string
}

type WrapperTransactionRsp struct {
	Hash         string
	User         string
	SrcChainId   uint64
	BlockHeight  uint64
	Time         uint64
	DstChainId   uint64
	DstUser      string
	ServerId     uint64
	FeeTokenHash string
	FeeAmount    string
	State        uint64
}

func MakeWrapperTransactionRsp(transaction *WrapperTransaction) *WrapperTransactionRsp {
	transactionRsp := &WrapperTransactionRsp{
		Hash:         transaction.Hash,
		User:         transaction.User,
		SrcChainId:   transaction.SrcChainId,
		BlockHeight:  transaction.BlockHeight,
		Time:         transaction.Time,
		DstChainId:   transaction.DstChainId,
		DstUser:      transaction.DstUser,
		ServerId:     transaction.ServerId,
		FeeTokenHash: transaction.FeeTokenHash,
		FeeAmount:    transaction.FeeAmount.String(),
		State:        transaction.Status,
	}
	return transactionRsp
}

type WrapperTransactionsReq struct {
	PageSize int
	PageNo   int
}

type WrapperTransactionsRsp struct {
	PageSize     int
	PageNo       int
	TotalPage    int
	TotalCount   int
	Transactions []*WrapperTransactionRsp
}

func MakeWrapperTransactionsRsp(
	pageSize, pageNo, totalPage, totalCount int,
	transactions []*WrapperTransaction,
) *WrapperTransactionsRsp {

	rsp := &WrapperTransactionsRsp{
		PageSize:   pageSize,
		PageNo:     pageNo,
		TotalPage:  totalPage,
		TotalCount: totalCount,
	}
	for _, v := range transactions {
		tx := MakeWrapperTransactionRsp(v)
		rsp.Transactions = append(rsp.Transactions, tx)
	}
	return rsp
}

type TransactionOfHashReq struct {
	Hash string
}

type TransactionStateRsp struct {
	Hash       string
	ChainId    uint64
	Blocks     uint64
	NeedBlocks uint64
	Time       uint64
}

type TransactionRsp struct {
	Hash             string
	User             string
	SrcChainId       uint64
	BlockHeight      uint64
	Time             uint64
	DstChainId       uint64
	DstUser          string
	TokenId          string
	ServerId         uint64
	FeeToken         *TokenRsp
	FeeAmount        string
	State            uint64
	Asset            *NFTAssetRsp
	TransactionState []*TransactionStateRsp
}

func MakeTransactionRsp(transaction *SrcPolyDstRelation, chainsMap map[uint64]*Chain) *TransactionRsp {
	transactionRsp := &TransactionRsp{
		Hash:        transaction.WrapperTransaction.Hash,
		User:        transaction.WrapperTransaction.User,
		SrcChainId:  transaction.WrapperTransaction.SrcChainId,
		BlockHeight: transaction.WrapperTransaction.BlockHeight,
		Time:        transaction.WrapperTransaction.Time,
		DstChainId:  transaction.WrapperTransaction.DstChainId,
		ServerId:    transaction.WrapperTransaction.ServerId,
		FeeAmount:   transaction.WrapperTransaction.FeeAmount.String(),
		TokenId:     transaction.SrcTransaction.SrcTransfer.Amount.String(),
		DstUser:     transaction.SrcTransaction.SrcTransfer.DstUser,
		State:       transaction.WrapperTransaction.Status,
	}
	if transaction.FeeToken != nil {
		transactionRsp.FeeToken = MakeTokenRsp(transaction.FeeToken)
		precision := decimal.NewFromInt(basedef.Int64FromFigure(int(transaction.FeeToken.TokenBasic.Precision)))
		{
			bbb := decimal.NewFromBigInt(&transaction.WrapperTransaction.FeeAmount.Int, 0)
			feeAmount := bbb.Div(precision)
			transactionRsp.FeeAmount = feeAmount.String()
		}
	}
	if transaction.SrcTransaction != nil {
		transactionRsp.TransactionState = append(transactionRsp.TransactionState, &TransactionStateRsp{
			Hash:    transaction.SrcTransaction.Hash,
			ChainId: transaction.SrcTransaction.ChainId,
			Blocks:  transaction.SrcTransaction.Height,
			Time:    transaction.SrcTransaction.Time,
		})
	} else {
		transactionRsp.TransactionState = append(transactionRsp.TransactionState, &TransactionStateRsp{
			Hash:    "",
			ChainId: transaction.WrapperTransaction.SrcChainId,
			Blocks:  0,
			Time:    0,
		})
	}
	if transaction.PolyTransaction != nil {
		transactionRsp.TransactionState = append(transactionRsp.TransactionState, &TransactionStateRsp{
			Hash:    transaction.PolyTransaction.Hash,
			ChainId: transaction.PolyTransaction.ChainId,
			Blocks:  transaction.PolyTransaction.Height,
			Time:    transaction.PolyTransaction.Time,
		})
	} else {
		transactionRsp.TransactionState = append(transactionRsp.TransactionState, &TransactionStateRsp{
			Hash:    "",
			ChainId: 0,
			Blocks:  0,
			Time:    0,
		})
	}
	if transaction.DstTransaction != nil {
		transactionRsp.TransactionState = append(transactionRsp.TransactionState, &TransactionStateRsp{
			Hash:    transaction.DstTransaction.Hash,
			ChainId: transaction.DstTransaction.ChainId,
			Blocks:  transaction.DstTransaction.Height,
			Time:    transaction.DstTransaction.Time,
		})
	} else {
		transactionRsp.TransactionState = append(transactionRsp.TransactionState, &TransactionStateRsp{
			Hash:    "",
			ChainId: transaction.WrapperTransaction.DstChainId,
			Blocks:  0,
			Time:    0,
		})
	}
	for _, state := range transactionRsp.TransactionState {
		chain, ok := chainsMap[state.ChainId]
		if ok {
			if state.ChainId == transaction.WrapperTransaction.DstChainId {
				state.NeedBlocks = 1
			} else {
				state.NeedBlocks = chain.BackwardBlockNumber
			}
			if state.Blocks <= 1 {
				continue
			}
			state.Blocks = chain.Height - state.Blocks
			if state.Blocks > state.NeedBlocks {
				state.Blocks = state.NeedBlocks
			}
		}
	}
	return transactionRsp
}

type TransactionsOfAddressReq struct {
	Addresses []string
	PageSize  int
	PageNo    int
}

type TransactionsOfAddressRsp struct {
	PageSize     int
	PageNo       int
	TotalPage    int
	TotalCount   int
	Transactions []*TransactionRsp
}

func MakeTransactionsOfUserRsp(pageSize int, pageNo int, totalPage int, totalCount int, transactions []*SrcPolyDstRelation, chainsMap map[uint64]*Chain) *TransactionsOfAddressRsp {
	transactionsRsp := &TransactionsOfAddressRsp{
		PageSize:   pageSize,
		PageNo:     pageNo,
		TotalPage:  totalPage,
		TotalCount: totalCount,
	}
	for _, transaction := range transactions {
		rsp := MakeTransactionRsp(transaction, chainsMap)
		transactionsRsp.Transactions = append(transactionsRsp.Transactions, rsp)
	}
	return transactionsRsp
}

//type TransactionsOfStateReq struct {
//	State    uint64
//	PageSize int
//	PageNo   int
//}
//
//type TransactionsOfStateRsp struct {
//	PageSize     int
//	PageNo       int
//	TotalPage    int
//	TotalCount   int
//	Transactions []*WrapperTransactionRsp
//}
//
//func MakeTransactionsOfStateRsp(pageSize int, pageNo int, totalPage int, totalCount int, transactions []*WrapperTransaction) *WrapperTransactionsRsp {
//	transactionsRsp := &WrapperTransactionsRsp{
//		PageSize:   pageSize,
//		PageNo:     pageNo,
//		TotalPage:  totalPage,
//		TotalCount: totalCount,
//	}
//	for _, transaction := range transactions {
//		transactionsRsp.Transactions = append(transactionsRsp.Transactions, MakeWrapperTransactionRsp(transaction))
//	}
//	return transactionsRsp
//}
//
//type AddressReq struct {
//	ChainId     uint64
//	AddressHash string
//}
//
//type AddressRsp struct {
//	AddressHash string
//	Address     string
//	ChainId     uint64
//}
//
//func MakeAddressRsp(addressHash string, chainId uint64, address string) *AddressRsp {
//	addressRsp := &AddressRsp{
//		AddressHash: addressHash,
//		Address:     address,
//		ChainId:     chainId,
//	}
//	return addressRsp
//}
//
//type PolyTransactionReq struct {
//	Hash string
//}
//
//type PolyTransactionRsp struct {
//	Hash       string
//	ChainId    uint64
//	State      uint64
//	Time       uint64
//	Fee        string
//	Height     uint64
//	SrcChainId uint64
//	SrcHash    string
//	DstChainId uint64
//	Key        string
//}
//
//func MakePolyTransactionRsp(transaction *PolyTransaction) *PolyTransactionRsp {
//	transactionRsp := &PolyTransactionRsp{
//		Hash:       transaction.Hash,
//		ChainId:    transaction.ChainId,
//		State:      transaction.State,
//		Time:       transaction.Time,
//		Fee:        transaction.Fee.String(),
//		Height:     transaction.Height,
//		SrcChainId: transaction.SrcChainId,
//		SrcHash:    transaction.SrcHash,
//		DstChainId: transaction.DstChainId,
//		Key:        transaction.Key,
//	}
//	return transactionRsp
//}
//
//type PolyTransactionsReq struct {
//	PageSize int
//	PageNo   int
//}
//
//type PolyTransactionsRsp struct {
//	PageSize     int
//	PageNo       int
//	TotalPage    int
//	TotalCount   int
//	Transactions []*PolyTransactionRsp
//}
//
//func MakePolyTransactionsRsp(pageSize int, pageNo int, totalPage int, totalCount int, transactions []*PolyTransaction) *PolyTransactionsRsp {
//	transactionsRsp := &PolyTransactionsRsp{
//		PageSize:   pageSize,
//		PageNo:     pageNo,
//		TotalPage:  totalPage,
//		TotalCount: totalCount,
//	}
//	for _, transaction := range transactions {
//		transactionsRsp.Transactions = append(transactionsRsp.Transactions, MakePolyTransactionRsp(transaction))
//	}
//	return transactionsRsp
//}
