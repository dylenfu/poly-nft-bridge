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
