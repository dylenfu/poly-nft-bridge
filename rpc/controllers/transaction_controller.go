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
		Select("tx_hash as hash, asset as asset").
		Joins("inner join wrapper_transactions on src_transfers.tx_hash = wrapper_transactions.hash").
		Joins("left join tokens on tokens.hash = wrapper_transactions.fee_token_hash and tokens.chain_id=wrapper_transactions.src_chain_id").
		Where("`from` in ? or src_transfers.dst_user in ?", req.Addresses, req.Addresses)).
		Select("src_transactions.hash as src_hash, poly_transactions.hash as poly_hash, dst_transactions.hash as dst_hash, src_transactions.chain_id as chain_id, u.asset as token_hash").
		Joins("left join src_transactions on u.hash = src_transactions.hash").
		Joins("left join poly_transactions on src_transactions.hash = poly_transactions.src_hash").
		Joins("left join dst_transactions on poly_transactions.hash = dst_transactions.poly_hash").
		Preload("WrapperTransaction").
		Preload("Token").
		Preload("SrcTransaction").
		Preload("SrcTransaction.SrcTransfer").
		Preload("PolyTransaction").
		Preload("DstTransaction").
		Preload("DstTransaction.DstTransfer").
		Preload("Asset").
		Preload("Asset.AssetBasic").
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
		chainsMap[*chain.ChainId] = chain
	}

	totalPage := (int(transactionNum) + req.PageSize - 1) / req.PageSize
	totalCnt := int(transactionNum)
	data := models.MakeTransactionsOfUserRsp(req.PageSize, req.PageNo, totalPage, totalCnt, srcPolyDstRelations, chainsMap)
	output(&c.Controller, data)
}

//func (c *TransactionController) TransactionOfHash() {
//	var transactionOfHashReq models.TransactionOfHashReq
//	var err error
//	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &transactionOfHashReq); err != nil {
//		c.Data["json"] = models.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
//		c.Ctx.ResponseWriter.WriteHeader(400)
//		c.ServeJSON()
//	}
//	srcPolyDstRelation := new(models.SrcPolyDstRelation)
//	res := db.Table("src_transactions").
//		Select("src_transactions.hash as src_hash, poly_transactions.hash as poly_hash, dst_transactions.hash as dst_hash, src_transactions.chain_id as chain_id, src_transfers.asset as token_hash").
//		Where("src_transactions.hash = ?", transactionOfHashReq.Hash).
//		Joins("inner join wrapper_transactions on src_transactions.hash = wrapper_transactions.hash").
//		Joins("left join src_transfers on src_transactions.hash = src_transfers.tx_hash").
//		Joins("left join poly_transactions on src_transactions.hash = poly_transactions.src_hash").
//		Joins("left join dst_transactions on poly_transactions.hash = dst_transactions.poly_hash").
//		Preload("WrapperTransaction").
//		Preload("SrcTransaction").
//		Preload("SrcTransaction.SrcTransfer").
//		Preload("PolyTransaction").
//		Preload("DstTransaction").
//		Preload("DstTransaction.DstTransfer").
//		Preload("Asset").
//		Preload("Asset.AssetBasic").
//		Order("src_transactions.time desc").
//		Find(srcPolyDstRelation)
//	if res.RowsAffected == 0 {
//		c.Data["json"] = models.MakeErrorRsp(fmt.Sprintf("transacion: %s does not exist", transactionOfHashReq.Hash))
//		c.Ctx.ResponseWriter.WriteHeader(400)
//		c.ServeJSON()
//		return
//	}
//	chains := make([]*models.Chain, 0)
//	db.Model(&models.Chain{}).Find(&chains)
//	chainsMap := make(map[uint64]*models.Chain)
//	for _, chain := range chains {
//		chainsMap[*chain.ChainId] = chain
//	}
//	c.Data["json"] = models.MakeTransactionRsp(srcPolyDstRelation, chainsMap)
//	c.ServeJSON()
//}
//
//func (c *TransactionController) TransactionsOfState() {
//	var transactionsOfStateReq models.TransactionsOfStateReq
//	var err error
//	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &transactionsOfStateReq); err != nil {
//		c.Data["json"] = models.MakeErrorRsp(fmt.Sprintf("request parameter is invalid!"))
//		c.Ctx.ResponseWriter.WriteHeader(400)
//		c.ServeJSON()
//	}
//	transactions := make([]*models.WrapperTransaction, 0)
//	db.Where("status = ?", transactionsOfStateReq.State).Limit(transactionsOfStateReq.PageSize).Offset(transactionsOfStateReq.PageSize * transactionsOfStateReq.PageNo).Order("time asc").Find(&transactions)
//	var transactionNum int64
//	db.Model(&models.WrapperTransaction{}).Where("status = ?", transactionsOfStateReq.State).Count(&transactionNum)
//	c.Data["json"] = models.MakeTransactionsOfStateRsp(transactionsOfStateReq.PageSize, transactionsOfStateReq.PageNo,
//		(int(transactionNum)+transactionsOfStateReq.PageSize-1)/transactionsOfStateReq.PageSize, int(transactionNum), transactions)
//	c.ServeJSON()
//}
