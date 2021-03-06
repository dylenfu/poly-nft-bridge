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

package crosschaindao

import (
	"github.com/polynetwork/poly-nft-bridge/conf"
	basedef "github.com/polynetwork/poly-nft-bridge/const"
	expd "github.com/polynetwork/poly-nft-bridge/dao/crosschaindao/explorerdao"
	stkd "github.com/polynetwork/poly-nft-bridge/dao/crosschaindao/stakedao"
	swpd "github.com/polynetwork/poly-nft-bridge/dao/crosschaindao/swapdao"
	"github.com/polynetwork/poly-nft-bridge/models"
)

type CrossChainDao interface {
	UpdateEvents(chain *models.Chain, wrapperTransactions []*models.WrapperTransaction, srcTransactions []*models.SrcTransaction, polyTransactions []*models.PolyTransaction, dstTransactions []*models.DstTransaction) error
	RemoveEvents(srcHashes []string, polyHashes []string, dstHashes []string) error
	GetChain(chainId uint64) (*models.Chain, error)
	UpdateChain(chain *models.Chain) error
	AddChains(chain []*models.Chain, chainFees []*models.ChainFee) error
	AddTokens(tokens []*models.TokenBasic, tokenMaps []*models.TokenMap) error
	RemoveTokens(tokens []string) error
	RemoveTokenMaps(tokenMaps []*models.TokenMap) error
	Name() string

	AddAssets(assetBasics []*models.NFTAssetBasic) error
	RemoveAssets(assets []string) error
	RemoveAssetMaps(assetMaps []*models.NFTAssetMap) error
}

func NewCrossChainDao(server string, backup bool, dbCfg *conf.DBConfig) CrossChainDao {
	if server == basedef.SERVER_POLY_SWAP {
		return swpd.NewSwapDao(dbCfg, backup)
	} else if server == basedef.SERVER_EXPLORER {
		return expd.NewExplorerDao(dbCfg)
	} else if server == basedef.SERVER_STAKE {
		return stkd.NewStakeDao()
	} else {
		panic("generate dao failed with invalid server name")
	}
}
