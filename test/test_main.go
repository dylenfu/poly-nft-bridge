package test

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/polynetwork/poly-nft-bridge/conf"
)

const (
	configPath = "../config_devnet.json"
	PrivateEnc = "56b446a2de5edfccee1581fbba79e8bb5c269e28ab4c0487860afb7e2c2d2b6e"
)

var (
	config     *conf.Config
	privateKey *ecdsa.PrivateKey
)

func TestMain(m *testing.M) {
	config = conf.NewConfig(configPath)

	key, err := crypto.HexToECDSA(PrivateEnc)
	if err != nil {
		panic(err)
	}
	privateKey = key

	m.Run()
}
