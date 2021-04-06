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

package poly

import (
	"fmt"
	"math/big"

	"github.com/astaxie/beego/logs"
	pcm "github.com/polynetwork/poly-go-sdk/common"
	"github.com/polynetwork/poly-nft-bridge/conf"
	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/polynetwork/poly-nft-bridge/models"
	"github.com/polynetwork/poly-nft-bridge/sdk/poly_sdk"
)

const (
	_makeProof    = "makeProof"
	_btcTxToRelay = "btcTxToRelay"
)

type PolyChainListen struct {
	polyCfg *conf.ChainListenConfig
	polySdk *poly_sdk.PolySDKPro
}

func NewPolyChainListen(cfg *conf.ChainListenConfig) *PolyChainListen {
	polyListen := &PolyChainListen{}
	polyListen.polyCfg = cfg
	urls := cfg.GetNodesUrl()
	sdk := poly_sdk.NewPolySDKPro(urls, cfg.ListenSlot, cfg.ChainId)
	polyListen.polySdk = sdk
	return polyListen
}

func (p *PolyChainListen) GetLatestHeight() (uint64, error) {
	return p.polySdk.GetCurrentBlockHeight()
}

func (p *PolyChainListen) GetChainListenSlot() uint64 {
	return p.polyCfg.ListenSlot
}

func (p *PolyChainListen) GetChainId() uint64 {
	return p.polyCfg.ChainId
}

func (p *PolyChainListen) GetChainName() string {
	return p.polyCfg.ChainName
}

func (p *PolyChainListen) GetDefer() uint64 {
	return p.polyCfg.Defer
}

func (p *PolyChainListen) HandleNewBlock(height uint64) (
	[]*models.WrapperTransaction,
	[]*models.SrcTransaction,
	[]*models.PolyTransaction,
	[]*models.DstTransaction,
	error,
) {

	block, err := p.polySdk.GetBlockByHeight(height)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if block == nil {
		return nil, nil, nil, nil, fmt.Errorf("there is no poly block!")
	}

	events, err := p.polySdk.GetSmartContractEventByBlock(height)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tt := uint64(block.Header.Timestamp)
	chainID := p.GetChainId()
	polyTransactions := make([]*models.PolyTransaction, 0)
	for _, event := range events {
		for _, notify := range event.Notify {
			if notify.ContractAddress == p.polyCfg.ECCMContract {
				states := notify.States.([]interface{})
				contractMethod, _ := states[0].(string)
				logs.Info("chain: %s, tx hash: %s", p.GetChainName(), event.TxHash)

				if contractMethod != _makeProof && contractMethod != _btcTxToRelay {
					continue
				}
				if len(states) < 4 {
					continue
				}

				tx := assemblePolyTransaction(event, states, chainID, height, tt)
				polyTransactions = append(polyTransactions, tx)
			}
		}
	}
	return nil, nil, polyTransactions, nil, nil
}

func (p *PolyChainListen) GetExtendLatestHeight() (uint64, error) {
	if len(p.polyCfg.ExtendNodes) == 0 {
		return p.GetLatestHeight()
	}
	return p.GetLatestHeight()
}

func assemblePolyTransaction(
	event *pcm.SmartContactEvent,
	states []interface{},
	chainID, height, timestamp uint64,
) *models.PolyTransaction {

	fchainid := uint32(states[1].(float64))
	tchainid := uint32(states[2].(float64))
	mctx := &models.PolyTransaction{}
	mctx.ChainId = chainID
	mctx.Hash = event.TxHash
	mctx.State = uint64(event.State)
	mctx.Fee = &models.BigInt{*big.NewInt(0)}
	mctx.Time = timestamp
	mctx.Height = height
	mctx.SrcChainId = uint64(fchainid)
	mctx.DstChainId = uint64(tchainid)
	if uint64(fchainid) == basedef.ETHEREUM_CROSSCHAIN_ID ||
		uint64(fchainid) == basedef.BSC_CROSSCHAIN_ID ||
		uint64(fchainid) == basedef.HECO_CROSSCHAIN_ID {
		mctx.SrcHash = states[3].(string)
	} else {
		mctx.SrcHash = basedef.HexStringReverse(states[3].(string))
	}
	return mctx
}
