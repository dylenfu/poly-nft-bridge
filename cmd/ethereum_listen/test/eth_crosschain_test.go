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
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	basedef "github.com/polynetwork/poly-nft-bridge/const"
	nftwp "github.com/polynetwork/poly-nft-bridge/go_abi/nft_wrap_abi"
	"github.com/polynetwork/poly-nft-bridge/go_abi/wrapper_abi"
	"github.com/polynetwork/poly-nft-bridge/sdk/eth_sdk"
	"github.com/stretchr/testify/assert"
)

func Test_EthCrossChain(t *testing.T) {
	cfg := config.GetChainListenConfig(basedef.ETHEREUM_CROSSCHAIN_ID)
	urls := cfg.GetNodesUrl()
	ethSdk := eth_sdk.NewEthereumSdkPro(urls, cfg.ListenSlot, basedef.ETHEREUM_CROSSCHAIN_ID)
	contractabi, err := abi.JSON(strings.NewReader(nftwp.PolyNFTWrapperABI))
	assert.NoError(t, err)

	assetHash := common.HexToAddress("0000000000000000000000000000000000000000")
	toAddress := common.Hex2Bytes("6e43f9988f2771f1a2b140cb3faad424767d39fc")
	txData, err := contractabi.Pack(
		"lock",
		assetHash,
		basedef.BSC_CROSSCHAIN_ID,
		toAddress,
		big.NewInt(int64(100000000000000000)),
		big.NewInt(10000000000000000),
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

// todo
func Test_EthGetFeeCollector(t *testing.T) {
	cfg := config.GetChainListenConfig(basedef.ETHEREUM_CROSSCHAIN_ID)
	address := common.HexToAddress(cfg.WrapperContract)
	urls := cfg.GetNodesUrl()
	ethSdk := eth_sdk.NewEthereumSdkPro(urls, cfg.ListenSlot, basedef.ETHEREUM_CROSSCHAIN_ID)
	instance, err := nftwp.NewPolyNFTWrapper(address, ethSdk.GetClient())
	assert.NoError(t, err)

	collector, _ := instance.FeeCollector(nil)
	t.Logf("collector: %s", collector.String())

	lockproxy, _ := instance.LockProxy(nil)
	t.Logf("lock proxy: %s", lockproxy.String())

	owner, _ := instance.Owner(nil)
	t.Logf("owner: %s", owner.String())
}

func Test_EthExtractFee(t *testing.T) {
	cfg := config.GetChainListenConfig(basedef.ETHEREUM_CROSSCHAIN_ID)
	urls := cfg.GetNodesUrl()
	ethSdk := eth_sdk.NewEthereumSdkPro(urls, cfg.ListenSlot, basedef.ETHEREUM_CROSSCHAIN_ID)
	contractabi, err := abi.JSON(strings.NewReader(nftwp.PolyNFTWrapperABI))
	assert.NoError(t, err)

	assetHash := common.HexToAddress("0000000000000000000000000000000000000000")
	txData, err := contractabi.Pack("extractFee", assetHash)
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
		Value:    big.NewInt(0),
		Data:     txData,
	}

	gasLimit, err := ethSdk.EstimateGas(callMsg)
	assert.NoError(t, err)
	t.Logf("gas limit: %d", gasLimit)

	tx := types.NewTransaction(nonce, wrapAddr, big.NewInt(0), gasLimit, gasPrice, txData)
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	assert.NoError(t, err)

	assert.NoError(t, ethSdk.SendRawTransaction(signedTx))
	ethSdk.WaitTransactionConfirm(signedTx.Hash())
}

func TestEthereumSpeedup(t *testing.T) {
	cfg := config.GetChainListenConfig(basedef.ETHEREUM_CROSSCHAIN_ID)
	urls := cfg.GetNodesUrl()
	ethSdk := eth_sdk.NewEthereumSdkPro(urls, cfg.ListenSlot, basedef.ETHEREUM_CROSSCHAIN_ID)
	contractabi, err := abi.JSON(strings.NewReader(wrapper_abi.IPolyWrapperABI))
	assert.NoError(t, err)

	assetHash := common.HexToAddress("0000000000000000000000000000000000000000")
	txData, err := contractabi.Pack("extractFee", assetHash)
	assert.NoError(t, err)
	fmt.Printf("TestInvokeContract - txdata:%s\n", hex.EncodeToString(txData))
	wrapperContractAddress := common.HexToAddress(cfg.WrapperContract)

	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("user address: %s\n", fromAddr.String())
	nonce, err := ethSdk.NonceAt(fromAddr)
	assert.NoError(t, err)
	gasPrice, err := ethSdk.SuggestGasPrice()
	assert.NoError(t, err)
	fmt.Printf("gas price: %s\n", gasPrice.String())
	callMsg := ethereum.CallMsg{
		From: fromAddr, To: &wrapperContractAddress, Gas: 0, GasPrice: gasPrice,
		Value: big.NewInt(0), Data: txData,
	}

	gasLimit, err := ethSdk.EstimateGas(callMsg)
	assert.NoError(t, err)
	fmt.Printf("gas limit: %d\n", gasLimit)

	tx := types.NewTransaction(nonce, wrapperContractAddress, big.NewInt(0), gasLimit, gasPrice, txData)
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	assert.NoError(t, err)

	assert.NoError(t, ethSdk.SendRawTransaction(signedTx))
	ethSdk.WaitTransactionConfirm(signedTx.Hash())
}
