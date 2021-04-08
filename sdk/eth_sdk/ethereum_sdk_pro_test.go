package eth_sdk

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestEthereumSdkPro_GetNFTs(t *testing.T) {
	asset := common.HexToAddress("03d84da9432f7cb5364a8b99286f97c59f738001")
	owner := common.HexToAddress("5fb03eb21303d39967a1a119b32dd744a0fa8986")
	start, end := 0, 100
	data, err := pro.GetNFTs(asset, owner, start, end)
	assert.NoError(t, err)
	for _, tokenid := range data {
		t.Logf("tokenid %d", tokenid.Uint64())
	}
}

func TestEthereumSdkPro_GetNFTURLs(t *testing.T) {
	asset := common.HexToAddress("03d84da9432f7cb5364a8b99286f97c59f738001")
	tokens := []*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(2)}
	data, err := pro.GetNFTURLs(asset, tokens)
	assert.NoError(t, err)
	for tokenid, url := range data {
		t.Logf("tokenid %d, url %s", tokenid, url)
	}
}
