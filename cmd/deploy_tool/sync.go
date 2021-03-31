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

package main

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	polysdk "github.com/polynetwork/poly-go-sdk"
	"github.com/polynetwork/poly-nft-bridge/sdk/eth_sdk"
	xpolysdk "github.com/polynetwork/poly-nft-bridge/sdk/poly_sdk"
)

func SyncSideChainGenesisHeaderToPolyChain(
	sideChainID uint64,
	sideChainSdk *eth_sdk.EthereumSdk,
	polySdk *xpolysdk.PolySDK,
	validators []*polysdk.Account,
) error {

	curr, err := sideChainSdk.GetCurrentBlockHeight()
	if err != nil {
		return err
	}
	hdr, err := sideChainSdk.GetHeaderByNumber(curr)
	if err != nil {
		return err
	}

	headerEnc, err := hdr.MarshalJSON()
	if err != nil {
		return err
	}

	if err := polySdk.SyncGenesisBlock(sideChainID, validators, headerEnc); err != nil {
		return err
	}

	return nil
}

func SyncPolyChainGenesisHeader2SideChain(
	polySDK *xpolysdk.PolySDK,
	sideChainECCMOwnerKey *ecdsa.PrivateKey,
	sideChainSdk *eth_sdk.EthereumSdk,
	sideChainECCM common.Address,
) error {

	// `epoch` related with the poly validators changing,
	// we can set it as 0 if poly validators never changed on develop environment.
	var hasValidatorsBlockNumber uint64 = 0
	gB, err := polySDK.GetBlockByHeight(hasValidatorsBlockNumber)
	if err != nil {
		return err
	}

	bookeepers, err := xpolysdk.GetBookeeper(gB)
	if err != nil {
		return err
	}
	bookeepersEnc := xpolysdk.AssembleNoCompressBookeeper(bookeepers)
	headerEnc := gB.Header.ToArray()

	if _, err := sideChainSdk.InitGenesisBlock(
		sideChainECCMOwnerKey,
		sideChainECCM,
		headerEnc,
		bookeepersEnc,
	); err != nil {
		return err
	}

	return nil
}
