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

package test

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/polynetwork/poly-nft-bridge/go_abi/wrapper_abi"
	"github.com/polynetwork/poly-nft-bridge/sdk/eth_sdk"
	"github.com/stretchr/testify/assert"
)

func Test_BscCrossChain(t *testing.T) {
	cfg := config.GetChainListenConfig(basedef.BSC_CROSSCHAIN_ID)
	urls := cfg.GetNodesUrl()

	ethSdk := eth_sdk.NewEthereumSdkPro(urls, cfg.ListenSlot, basedef.BSC_CROSSCHAIN_ID)
	contractabi, err := abi.JSON(strings.NewReader(wrapper_abi.IPolyWrapperABI))
	assert.NoError(t, err)

	assetHash := common.HexToAddress("0000000000000000000000000000000000000000")
	toAddress := common.Hex2Bytes("6e43f9988f2771f1a2b140cb3faad424767d39fc")
	txData, err := contractabi.Pack(
		"lock",
		assetHash,
		uint64(2),
		toAddress,
		big.NewInt(100000000000000000),
		big.NewInt(10000000000000000),
		big.NewInt(0),
	)
	assert.NoError(t, err)
	t.Logf("TestInvokeContract - txdata:%s", hex.EncodeToString(txData))

	wrapAddr := common.HexToAddress(cfg.WrapperContract)
	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	t.Logf("user address: %s", fromAddr.String())

	nonce, err := ethSdk.NonceAt(fromAddr)
	assert.NoError(t, err)
	gasPrice, err := ethSdk.SuggestGasPrice()
	assert.NoError(t, err)
	t.Logf("gas price: %s", gasPrice.String())

	callMsg := ethereum.CallMsg{
		From:     fromAddr,
		To:       &wrapAddr,
		Gas:      0,
		GasPrice: gasPrice,
		Value:    big.NewInt(100000000000000000),
		Data:     txData,
	}

	gasLimit, err := ethSdk.EstimateGas(callMsg)
	assert.NoError(t, err)
	t.Logf("gas limit: %d", gasLimit)

	tx := types.NewTransaction(nonce, wrapAddr, big.NewInt(100000000000000000), gasLimit, gasPrice, txData)
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	assert.NoError(t, err)

	assert.NoError(t, ethSdk.SendRawTransaction(signedTx))
	ethSdk.WaitTransactionConfirm(signedTx.Hash())
}
