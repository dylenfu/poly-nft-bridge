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

package eth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/astaxie/beego/logs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/polynetwork/poly-nft-bridge/conf"
	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/polynetwork/poly-nft-bridge/dao/models"
	"github.com/polynetwork/poly-nft-bridge/go_abi/eccm_abi"
	nftlp "github.com/polynetwork/poly-nft-bridge/go_abi/nft_lock_proxy_abi"
	nftwp "github.com/polynetwork/poly-nft-bridge/go_abi/nft_wrap_abi"
	"github.com/polynetwork/poly-nft-bridge/sdk/eth_sdk"
)

const (
	_eth_crosschainlock   = "CrossChainLockEvent"
	_eth_crosschainunlock = "CrossChainUnlockEvent"
	_eth_lock             = "LockEvent"
	_eth_unlock           = "UnlockEvent"
)

type EthereumChainListen struct {
	ethCfg *conf.ChainListenConfig
	ethSdk *eth_sdk.EthereumSdkPro
}

func NewEthereumChainListen(cfg *conf.ChainListenConfig) *EthereumChainListen {
	ethListen := &EthereumChainListen{}
	ethListen.ethCfg = cfg
	//
	urls := cfg.GetNodesUrl()
	sdk := eth_sdk.NewEthereumSdkPro(urls, cfg.ListenSlot, cfg.ChainId)
	ethListen.ethSdk = sdk
	return ethListen
}

func (e *EthereumChainListen) WrapperAddress() common.Address {
	return common.HexToAddress(e.ethCfg.WrapperContract)
}

func (e *EthereumChainListen) ECCMAddress() common.Address {
	return common.HexToAddress(e.ethCfg.CCMContract)
}

func (e *EthereumChainListen) ProxyAddress() common.Address {
	return common.HexToAddress(e.ethCfg.ProxyContract)
}

func (e *EthereumChainListen) GetLatestHeight() (uint64, error) {
	return e.ethSdk.GetLatestHeight()
}

func (e *EthereumChainListen) GetChainListenSlot() uint64 {
	return e.ethCfg.ListenSlot
}

func (e *EthereumChainListen) GetChainId() uint64 {
	return e.ethCfg.ChainId
}

func (e *EthereumChainListen) GetChainName() string {
	return e.ethCfg.ChainName
}

func (e *EthereumChainListen) GetDefer() uint64 {
	return e.ethCfg.Defer
}

func (e *EthereumChainListen) HandleNewBlock(height uint64) (
	[]*models.WrapperTransaction,
	[]*models.SrcTransaction,
	[]*models.PolyTransaction,
	[]*models.DstTransaction,
	error,
) {

	wrapAddr := e.WrapperAddress()
	eccmAddr := e.ECCMAddress()
	proxyAddr := e.ProxyAddress()
	chainName := e.GetChainName()
	chainID := e.GetChainId()

	blockHeader, err := e.ethSdk.GetHeaderByNumber(height)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if blockHeader == nil {
		return nil, nil, nil, nil, fmt.Errorf("there is no ethereum block!")
	}
	tt := blockHeader.Time

	wrapperTransactions, err := e.getWrapperEventByBlockNumber(wrapAddr, height, height)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	for _, wtx := range wrapperTransactions {
		logs.Info("(wrapper) from chain: %s, txhash: %s", chainName, wtx.Hash)
		wtx.Time = tt
		wtx.SrcChainId = e.GetChainId()
		wtx.Status = basedef.STATE_SOURCE_DONE
	}
	eccmLockEvents, eccmUnLockEvents, err := e.getECCMEventByBlockNumber(eccmAddr, height, height)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	proxyLockEvents, proxyUnlockEvents, err := e.getProxyEventByBlockNumber(proxyAddr, height, height)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	srcTransactions := make([]*models.SrcTransaction, 0)
	dstTransactions := make([]*models.DstTransaction, 0)
	for _, lockEvent := range eccmLockEvents {
		if lockEvent.Method == _eth_crosschainlock {
			logs.Info("(lock) from chain: %s, txhash: %s, txid: %s",
				chainName, lockEvent.TxHash, lockEvent.Txid)
			srcTransaction := assembleSrcTransaction(lockEvent, proxyLockEvents, chainID, tt)
			srcTransactions = append(srcTransactions, srcTransaction)
		}
	}
	// save unLockEvent to db
	for _, unLockEvent := range eccmUnLockEvents {
		if unLockEvent.Method == _eth_crosschainunlock {
			logs.Info("(unlock) to chain: %s, txhash: %s", chainName, unLockEvent.TxHash)
			dstTransaction := assembleDstTransaction(unLockEvent, proxyUnlockEvents, chainID, tt)
			dstTransactions = append(dstTransactions, dstTransaction)
		}
	}
	return wrapperTransactions, srcTransactions, nil, dstTransactions, nil
}

func (e *EthereumChainListen) HandleNewBlockBatch(startHeight, endHeight uint64) (
	[]*models.WrapperTransaction,
	[]*models.SrcTransaction,
	[]*models.PolyTransaction,
	[]*models.DstTransaction,
	error,
) {

	wrapAddr := e.WrapperAddress()
	eccmAddr := e.ECCMAddress()
	proxyAddr := e.ProxyAddress()
	chainName := e.GetChainName()
	chainID := e.GetChainId()

	wrapperTransactions, err := e.getWrapperEventByBlockNumber(wrapAddr, startHeight, endHeight)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	for _, wtx := range wrapperTransactions {
		logs.Info("(wrapper) from chain: %s, txhash: %s", chainName, wtx.Hash)
		wtx.SrcChainId = e.GetChainId()
		wtx.Status = basedef.STATE_SOURCE_DONE
	}
	eccmLockEvents, eccmUnLockEvents, err := e.getECCMEventByBlockNumber(eccmAddr, startHeight, endHeight)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	proxyLockEvents, proxyUnlockEvents, err := e.getProxyEventByBlockNumber(proxyAddr, startHeight, endHeight)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	//

	srcTransactions := make([]*models.SrcTransaction, 0)
	dstTransactions := make([]*models.DstTransaction, 0)
	for _, lockEvent := range eccmLockEvents {
		if lockEvent.Method == _eth_crosschainlock {
			logs.Info("(lock) from chain: %s, txhash: %s, txid: %s", chainName, lockEvent.TxHash, lockEvent.Txid)
			srcTransaction := assembleSrcTransaction(lockEvent, proxyLockEvents, chainID, 0)
			srcTransactions = append(srcTransactions, srcTransaction)
		}
	}
	// save unLockEvent to db
	for _, unLockEvent := range eccmUnLockEvents {
		if unLockEvent.Method == _eth_crosschainunlock {
			logs.Info("(unlock) to chain: %s, txhash: %s", chainName, unLockEvent.TxHash)
			dstTransaction := assembleDstTransaction(unLockEvent, proxyUnlockEvents, chainID, 0)
			dstTransactions = append(dstTransactions, dstTransaction)
		}
	}
	return wrapperTransactions, srcTransactions, nil, dstTransactions, nil
}

func (e *EthereumChainListen) getWrapperEventByBlockNumber(
	wrapAddr common.Address,
	startHeight, endHeight uint64) (
	[]*models.WrapperTransaction,
	error,
) {

	// todo: newPolyWrapper change to IPolyNFTWrapper
	wrapperContract, err := nftwp.NewPolyNFTWrapper(wrapAddr, e.ethSdk.GetClient())
	if err != nil {
		return nil, fmt.Errorf("GetSmartContractEventByBlock, error: %s", err.Error())
	}
	opt := &bind.FilterOpts{
		Start:   startHeight,
		End:     &endHeight,
		Context: context.Background(),
	}

	// get ethereum lock events from given block
	wrapperTransactions := make([]*models.WrapperTransaction, 0)
	lockEvents, err := wrapperContract.FilterPolyWrapperLock(opt, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("GetSmartContractEventByBlock, filter lock events :%s", err.Error())
	}
	for lockEvents.Next() {
		evt := lockEvents.Event
		wtx := wrapLockEvent2WrapTx(evt)
		wrapperTransactions = append(wrapperTransactions, wtx)
	}
	speedupEvents, err := wrapperContract.FilterPolyWrapperSpeedUp(opt, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("GetSmartContractEventByBlock, filter lock events :%s", err.Error())
	}
	for speedupEvents.Next() {
		evt := speedupEvents.Event
		wtx := wrapSpeedUpEvent2WrapTx(evt)
		wrapperTransactions = append(wrapperTransactions, wtx)
	}
	return wrapperTransactions, nil
}

func (e *EthereumChainListen) getECCMEventByBlockNumber(
	eccmAddr common.Address,
	startHeight, endHeight uint64) (
	[]*models.ECCMLockEvent,
	[]*models.ECCMUnlockEvent,
	error,
) {

	eccmContract, err := eccm_abi.NewEthCrossChainManager(eccmAddr, e.ethSdk.GetClient())
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, error: %s", err.Error())
	}
	opt := &bind.FilterOpts{
		Start:   startHeight,
		End:     &endHeight,
		Context: context.Background(),
	}
	// get ethereum lock events from given block
	eccmLockEvents := make([]*models.ECCMLockEvent, 0)
	crossChainEvents, err := eccmContract.FilterCrossChainEvent(opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, filter lock events :%s", err.Error())
	}
	for crossChainEvents.Next() {
		evt := crossChainEvents.Event
		Fee := e.GetConsumeGas(evt.Raw.TxHash)
		eccmLockEvent := crossChainEvent2ProxyLockEvent(evt, Fee)
		eccmLockEvents = append(eccmLockEvents, eccmLockEvent)
	}
	// ethereum unlock events from given block
	eccmUnlockEvents := make([]*models.ECCMUnlockEvent, 0)
	executeTxEvent, err := eccmContract.FilterVerifyHeaderAndExecuteTxEvent(opt)
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, filter unlock events :%s", err.Error())
	}

	for executeTxEvent.Next() {
		evt := executeTxEvent.Event
		Fee := e.GetConsumeGas(evt.Raw.TxHash)
		eccmUnlockEvent := verifyAndExecuteEvent2ProxyUnlockEvent(evt, Fee)
		eccmUnlockEvents = append(eccmUnlockEvents, eccmUnlockEvent)
	}
	return eccmLockEvents, eccmUnlockEvents, nil
}

func (e *EthereumChainListen) getProxyEventByBlockNumber(
	proxyAddr common.Address,
	startHeight, endHeight uint64) (
	[]*models.ProxyLockEvent,
	[]*models.ProxyUnlockEvent,
	error,
) {

	proxyContract, err := nftlp.NewNFTLockProxy(proxyAddr, e.ethSdk.GetClient())
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, error: %s", err.Error())
	}
	opt := &bind.FilterOpts{
		Start:   startHeight,
		End:     &endHeight,
		Context: context.Background(),
	}
	// get ethereum lock events from given block
	proxyLockEvents := make([]*models.ProxyLockEvent, 0)
	lockEvents, err := proxyContract.FilterLockEvent(opt)
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, filter lock events :%s", err.Error())
	}
	for lockEvents.Next() {
		proxyLockEvent := convertLockProxyEvent(lockEvents.Event)
		proxyLockEvents = append(proxyLockEvents, proxyLockEvent)
	}

	// ethereum unlock events from given block
	proxyUnlockEvents := make([]*models.ProxyUnlockEvent, 0)
	unlockEvents, err := proxyContract.FilterUnlockEvent(opt)
	if err != nil {
		return nil, nil, fmt.Errorf("GetSmartContractEventByBlock, filter unlock events :%s", err.Error())
	}
	for unlockEvents.Next() {
		proxyUnlockEvent := convertUnlockProxyEvent(unlockEvents.Event)
		proxyUnlockEvents = append(proxyUnlockEvents, proxyUnlockEvent)
	}
	return proxyLockEvents, proxyUnlockEvents, nil
}

func (e *EthereumChainListen) GetConsumeGas(hash common.Hash) uint64 {
	tx, err := e.ethSdk.GetTransactionByHash(hash)
	if err != nil {
		return 0
	}
	receipt, err := e.ethSdk.GetTransactionReceipt(hash)
	if err != nil {
		return 0
	}
	return tx.GasPrice().Uint64() * receipt.GasUsed
}

type ExtendHeightRsp struct {
	Status  uint64 `json:"status,string"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func (e *EthereumChainListen) GetExtendLatestHeight() (uint64, error) {
	if len(e.ethCfg.ExtendNodes) == 0 {
		return e.GetLatestHeight()
	}
	for i, _ := range e.ethCfg.ExtendNodes {
		height, err := e.getExtendLatestHeight(i)
		if err == nil {
			return height, nil
		}
	}
	return 0, fmt.Errorf("all extend node is not working")
}

func (e *EthereumChainListen) getExtendLatestHeight(node int) (uint64, error) {
	req, err := http.NewRequest("GET", e.ethCfg.ExtendNodes[node].Url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accepts", "application/json")
	q := url.Values{}
	q.Add("module", "proxy")
	q.Add("action", "eth_blockNumber")
	q.Add("apikey", e.ethCfg.ExtendNodes[node].Key)
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("response status code: %d", resp.StatusCode)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	extendHeight := new(ExtendHeightRsp)
	extendHeight.Status = 1
	err = json.Unmarshal(respBody, extendHeight)
	if err != nil {
		return 0, err
	}
	if extendHeight.Status == 0 {
		return 0, fmt.Errorf(extendHeight.Result)
	}
	height, err := hexutil.DecodeBig(extendHeight.Result)
	if err != nil {
		return 0, err
	}
	return height.Uint64(), nil
}
