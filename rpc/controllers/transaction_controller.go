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

type TransactionController struct {
	beego.Controller
}

func (c *TransactionController) Transactions() {
	var req models.WrapperTransactionsReq
	if !input(&c.Controller, &req) {
		return
	}

	transactions := make([]*models.WrapperTransaction, 0)
	db.Limit(req.PageSize).
		Offset(req.PageSize * req.PageNo).
		Order("time asc").
		Find(&transactions)

	var transactionNum int64
	db.Model(&models.WrapperTransaction{}).Count(&transactionNum)

	totalPage := (int(transactionNum) + req.PageSize - 1) / req.PageSize
	totalCnt := int(transactionNum)
	data := models.MakeWrapperTransactionsRsp(req.PageSize, req.PageNo, totalPage, totalCnt, transactions)
	output(&c.Controller, data)
}

func (c *TransactionController) TransactionsOfAddress() {
	var req models.TransactionsOfAddressReq

	if !input(&c.Controller, &req) {
		return
	}

	// load relations
	srcPolyDstRelations := make([]*models.SrcPolyDstRelation, 0)
	db.Table("(?) as u", db.Model(&models.SrcTransfer{}).
		Select("src_transfers.tx_hash as hash, src_transfers.asset as asset, wrapper_transactions.fee_token_hash as fee_token_hash").
		Joins("inner join wrapper_transactions on src_transfers.tx_hash = wrapper_transactions.hash").
		Where("`from` in ? or src_transfers.dst_user in ?", req.Addresses, req.Addresses)).
		Select("src_transactions.hash as src_hash, " +
			"poly_transactions.hash as poly_hash, " +
			"dst_transactions.hash as dst_hash, " +
			"src_transactions.chain_id as chain_id," +
			"u.asset as asset_hash, u.fee_token_hash as fee_token_hash").
		Joins("left join src_transactions on u.hash = src_transactions.hash").
		Joins("left join poly_transactions on src_transactions.hash = poly_transactions.src_hash").
		Joins("left join dst_transactions on poly_transactions.hash = dst_transactions.poly_hash").
		Preload("WrapperTransaction").
		Preload("Asset").
		Preload("Asset.AssetBasic").
		Preload("FeeToken").
		Preload("FeeToken.TokenBasic").
		Preload("SrcTransaction").
		Preload("SrcTransaction.SrcTransfer").
		Preload("PolyTransaction").
		Preload("DstTransaction").
		Preload("DstTransaction.DstTransfer").
		Limit(req.PageSize).Offset(req.PageSize * req.PageNo).
		Order("src_transactions.time desc").
		Find(&srcPolyDstRelations)

	// get transaction number
	var transactionNum int64
	db.Model(&models.SrcTransfer{}).
		Joins("inner join wrapper_transactions on src_transfers.tx_hash = wrapper_transactions.hash").
		Where("`from` in ? or src_transfers.dst_user in ?", req.Addresses, req.Addresses).
		Count(&transactionNum)

	// get chains
	chains := make([]*models.Chain, 0)
	db.Model(&models.Chain{}).Find(&chains)
	chainsMap := make(map[uint64]*models.Chain)
	for _, chain := range chains {
		chainsMap[chain.ChainId] = chain
	}

	totalPage := (int(transactionNum) + req.PageSize - 1) / req.PageSize
	totalCnt := int(transactionNum)
	data := models.MakeTransactionsOfUserRsp(req.PageSize, req.PageNo, totalPage, totalCnt, srcPolyDstRelations, chainsMap)
	output(&c.Controller, data)
}

func (c *TransactionController) TransactionOfHash() {
	var req models.TransactionOfHashReq

	if !input(&c.Controller, &req) {
		return
	}

	srcPolyDstRelation := new(models.SrcPolyDstRelation)
	res := db.Table("(?) as u", db.Model(&models.SrcTransfer{}).
		Select("src_transfers.tx_hash as hash, src_transfers.asset as asset, wrapper_transactions.fee_token_hash as fee_token_hash").
		Joins("inner join wrapper_transactions on src_transfers.tx_hash = wrapper_transactions.hash").
		Where("src_transfers.tx_hash =?", req.Hash)).
		Select("src_transactions.hash as src_hash, " +
			"poly_transactions.hash as poly_hash, " +
			"dst_transactions.hash as dst_hash, " +
			"src_transactions.chain_id as chain_id," +
			"u.asset as asset_hash, u.fee_token_hash as fee_token_hash").
		Joins("left join src_transactions on u.hash = src_transactions.hash").
		Joins("left join poly_transactions on src_transactions.hash = poly_transactions.src_hash").
		Joins("left join dst_transactions on poly_transactions.hash = dst_transactions.poly_hash").
		Preload("WrapperTransaction").
		Preload("Asset").
		Preload("Asset.AssetBasic").
		Preload("FeeToken").
		Preload("FeeToken.TokenBasic").
		Preload("SrcTransaction").
		Preload("SrcTransaction.SrcTransfer").
		Preload("PolyTransaction").
		Preload("DstTransaction").
		Preload("DstTransaction.DstTransfer").
		Order("src_transactions.time desc").
		Find(srcPolyDstRelation)

	if res.RowsAffected == 0 {
		notExist(&c.Controller)
		return
	}

	chains := make([]*models.Chain, 0)
	db.Model(&models.Chain{}).Find(&chains)
	chainsMap := make(map[uint64]*models.Chain)
	for _, chain := range chains {
		chainsMap[chain.ChainId] = chain
	}

	data := models.MakeTransactionRsp(srcPolyDstRelation, chainsMap)
	output(&c.Controller, data)
}

func (c *TransactionController) TransactionsOfState() {
	var req models.TransactionsOfStateReq
	if !input(&c.Controller, &req) {
		return
	}

	transactions := make([]*models.WrapperTransaction, 0)
	db.Where("status = ?", req.State).
		Limit(req.PageSize).
		Offset(req.PageSize * req.PageNo).
		Order("time asc").
		Find(&transactions)

	var transactionNum int64
	db.Model(&models.WrapperTransaction{}).
		Where("status = ?", req.State).
		Count(&transactionNum)

	totalPage := (int(transactionNum) + req.PageSize - 1) / req.PageSize
	totalCount := int(transactionNum)
	data := models.MakeTransactionsOfStateRsp(req.PageSize, req.PageNo, totalPage, totalCount, transactions)
	output(&c.Controller, data)
}
