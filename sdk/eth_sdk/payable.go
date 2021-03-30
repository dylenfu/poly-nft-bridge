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

// Notice: functions in this file only used for deploy_tool and test cases.

package eth_sdk

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	erc20 "github.com/polynetwork/poly-nft-bridge/go_abi/mintable_erc20_abi"
	nftmapping "github.com/polynetwork/poly-nft-bridge/go_abi/nft_mapping_abi"
	nftwrap "github.com/polynetwork/poly-nft-bridge/go_abi/nft_wrap_abi"
	xecdsa "github.com/polynetwork/poly-nft-bridge/utils/ecdsa"
	polycm "github.com/polynetwork/poly/common"
)

func (s *EthereumSdk) TransferNative(
	key *ecdsa.PrivateKey,
	to common.Address,
	amount *big.Int,
) (common.Hash, error) {

	from := xecdsa.Key2address(key)
	nonce, err := s.NonceAt(from)
	if err != nil {
		return EmptyHash, err
	}

	gasPrice, err := s.SuggestGasPrice()
	if err != nil {
		return EmptyHash, err
	}

	gasLimit, err := s.EstimateGas(ethereum.CallMsg{
		From: from, To: &to, Gas: 0, GasPrice: gasPrice,
		Value: amount, Data: []byte{},
	})
	if err != nil {
		return EmptyHash, err
	}

	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, []byte{})
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, key)
	if err != nil {
		return EmptyHash, err
	}
	if err := s.SendRawTransaction(signedTx); err != nil {
		return EmptyHash, err
	}

	if err := s.waitTxConfirm(signedTx.Hash()); err != nil {
		return EmptyHash, err
	}
	return signedTx.Hash(), nil
}

func (s *EthereumSdk) MintERC20Token(
	key *ecdsa.PrivateKey,
	asset, to common.Address,
	amount *big.Int) (common.Hash, error) {

	contract, err := erc20.NewERC20Mintable(asset, s.backend())
	if err != nil {
		return EmptyHash, err
	}

	auth, err := s.makeAuth(key)
	if err != nil {
		return EmptyHash, err
	}

	tx, err := contract.Mint(auth, to, amount)
	if err != nil {
		return EmptyHash, err
	}

	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}

	return tx.Hash(), nil
}

func (s *EthereumSdk) TransferERC20Token(
	key *ecdsa.PrivateKey,
	asset, to common.Address,
	amount *big.Int,
) (common.Hash, error) {

	contract, err := erc20.NewERC20Mintable(asset, s.backend())
	if err != nil {
		return EmptyHash, err
	}

	auth, err := s.makeAuth(key)
	if err != nil {
		return EmptyHash, err
	}
	tx, err := contract.Transfer(auth, to, amount)
	if err != nil {
		return EmptyHash, err
	}

	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}

	return tx.Hash(), nil
}

func (s *EthereumSdk) MintNFT(
	ownerKey *ecdsa.PrivateKey,
	asset,
	to common.Address,
	tokenID *big.Int,
	uri string,
) (common.Hash, error) {

	contract, err := nftmapping.NewCrossChainNFTMapping(asset, s.rawClient)
	if err != nil {
		return EmptyHash, err
	}

	auth, err := s.makeAuth(ownerKey)
	if err != nil {
		return EmptyHash, err
	}

	tx, err := contract.MintWithURI(auth, to, tokenID, uri)
	if err != nil {
		return EmptyHash, err
	}

	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}
	return tx.Hash(), nil
}

func (s *EthereumSdk) NFTSafeTransferFrom(
	nftOwnerKey *ecdsa.PrivateKey,
	asset,
	from,
	proxy common.Address,
	tokenID *big.Int,
	to common.Address,
	toChainID uint64,
) (common.Hash, error) {

	cm, err := nftmapping.NewCrossChainNFTMapping(asset, s.backend())
	if err != nil {
		return EmptyHash, err
	}

	auth, err := s.makeAuth(nftOwnerKey)
	if err != nil {
		return EmptyHash, err
	}
	data := assembleSafeTransferCallData(to, toChainID)
	tx, err := cm.SafeTransferFrom0(auth, from, proxy, tokenID, data)
	if err != nil {
		return EmptyHash, err
	}

	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}
	return tx.Hash(), nil
}

func (s *EthereumSdk) NFTApprove(key *ecdsa.PrivateKey, asset, to common.Address, token *big.Int) (common.Hash, error) {
	cm, err := nftmapping.NewCrossChainNFTMapping(asset, s.backend())
	if err != nil {
		return EmptyHash, err
	}
	auth, err := s.makeAuth(key)
	if err != nil {
		return EmptyHash, err
	}
	tx, err := cm.Approve(auth, to, token)
	if err != nil {
		return EmptyHash, err
	}
	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}
	return tx.Hash(), nil
}

func (s *EthereumSdk) GetNFTBalance(asset, owner common.Address) (*big.Int, error) {
	cm, err := nftmapping.NewCrossChainNFTMapping(asset, s.backend())
	if err != nil {
		return nil, err
	}
	return cm.BalanceOf(nil, owner)
}

func (s *EthereumSdk) GetNFTTokenUri(asset common.Address, tokenID *big.Int) (string, error) {
	cm, err := nftmapping.NewCrossChainNFTMapping(asset, s.backend())
	if err != nil {
		return "", err
	}
	return cm.TokenURI(nil, tokenID)
}

func (s *EthereumSdk) GetNFTApproved(asset common.Address, tokenID *big.Int) (common.Address, error) {
	cm, err := nftmapping.NewCrossChainNFTMapping(asset, s.backend())
	if err != nil {
		return EmptyAddress, err
	}
	return cm.GetApproved(nil, tokenID)
}

func (s *EthereumSdk) GetNFTOwner(asset common.Address, tokenID *big.Int) (common.Address, error) {
	cm, err := nftmapping.NewCrossChainNFTMapping(asset, s.backend())
	if err != nil {
		return EmptyAddress, err
	}
	return cm.OwnerOf(nil, tokenID)
}

func (s *EthereumSdk) WrapLock(
	key *ecdsa.PrivateKey,
	wrapAddr,
	fromAsset,
	toAddr common.Address,
	toChainId uint64,
	tokenID *big.Int,
	feeToken common.Address,
	feeAmount *big.Int,
	id *big.Int,
) (common.Hash, error) {

	wrapper, err := nftwrap.NewPolyNFTWrapper(wrapAddr, s.backend())
	if err != nil {
		return EmptyHash, err
	}

	auth, err := s.makeAuth(key)
	if err != nil {
		return EmptyHash, err
	}

	tx, err := wrapper.Lock(auth, fromAsset, toChainId, toAddr, tokenID, feeToken, feeAmount, id)
	if err != nil {
		return EmptyHash, err
	}

	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}

	return tx.Hash(), nil
}

func assembleSafeTransferCallData(toAddress common.Address, chainID uint64) []byte {
	sink := polycm.NewZeroCopySink(nil)
	sink.WriteVarBytes(toAddress.Bytes())
	sink.WriteUint64(chainID)
	return sink.Bytes()
}
