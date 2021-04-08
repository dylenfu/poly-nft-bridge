package eth_sdk

import "testing"

var pro *EthereumSdkPro

var urls = []string{
	"http://localhost:8545",
}

func TestMain(m *testing.M) {
	pro = NewEthereumSdkPro(urls, 1, 2)
	m.Run()
}
