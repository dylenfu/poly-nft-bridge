package wallet

import (
	"testing"

	"github.com/polynetwork/poly-nft-bridge/utils/leveldb"
	"github.com/stretchr/testify/assert"
)

func TestLoadEthAccount(t *testing.T) {
	keystoreDir := "/Users/dylen/software/nft-bridge/poly-nft-bridge/build/devnet/deploy_tool/keystore/eth/"
	storeDir := "/Users/dylen/software/nft-bridge/poly-nft-bridge/build/devnet/deploy_tool/leveldb/"
	account := "0x31c0dd87B33Dcd66f9a255Cf4CF39287F8AE593C"
	passphrase := "111111"
	storage := leveldb.NewLevelDBInstance(storeDir)
	key, err := LoadEthAccount(storage, keystoreDir, account, passphrase)
	assert.NoError(t, err)

	t.Log(key.PublicKey.X.String(), key.PublicKey.Y.String())
}
