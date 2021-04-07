package controllers

import (
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/polynetwork/poly-nft-bridge/models"
	"github.com/polynetwork/poly-nft-bridge/sdk/eth_sdk"
)

type ItemController struct {
	beego.Controller
}

// todo: cache url and token ids
func (c *ItemController) Items() {
	var req models.ItemsOfAddressReq
	if !input(&c.Controller, &req) {
		return
	}

	start := req.PageNo * req.PageSize
	end := start + req.PageSize

	sdk := selectSDK(req.ChainId)
	if sdk == nil {
		notExist(&c.Controller)
		return
	}

	owner := common.HexToAddress(req.Address)
	asset := common.HexToAddress(req.Asset)

	items := make([]*models.Item, 0)
	list := sdk.GetOwnerNFTs(asset, owner, start, end)
	if len(list) > 0 {
		urlmap := sdk.GetOwnerNFTUrls(asset, list)
		for _, v := range list {
			items = append(items, &models.Item{
				TokenId: v.Uint64(),
				Url:     urlmap[v.Uint64()],
			})
		}
	}

	data := models.MakeItemsOfAddressRsp(req.PageSize, req.PageNo, items)
	output(&c.Controller, data)
}

func selectSDK(chainId uint64) *eth_sdk.EthereumSdk {
	var url string
	switch chainId {
	case basedef.ETHEREUM_CROSSCHAIN_ID:
		url = "http://localhost:8545"
	case basedef.BSC_CROSSCHAIN_ID:
		url = "http://localhost:8546"
	case basedef.HECO_CROSSCHAIN_ID:
		url = "http://localhost:8547"
	default:
		return nil
	}

	if sdk, err := eth_sdk.NewEthereumSdk(url); err != nil {
		return nil
	} else {
		return sdk
	}
}
