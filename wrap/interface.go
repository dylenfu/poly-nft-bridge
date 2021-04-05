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

package wrap

import (
	"github.com/polynetwork/poly-nft-bridge/models"
)

type ChainHandle interface {

	// fetch current block height
	GetLatestHeight() (uint64, error)

	// fetch extend chain node block height, e.g: https://api.etherscan.io/api
	GetExtendLatestHeight() (uint64, error)

	// fetch block content, filter and save fixed cross chain events
	HandleNewBlock(height uint64) ([]*models.WrapperTransaction, []*models.SrcTransaction, []*models.PolyTransaction, []*models.DstTransaction, error)

	// listen slot denote an period which range of blocks
	GetChainListenSlot() uint64

	// fetch side chain id which settle in configuration
	GetChainId() uint64

	// fetch side chain name which settle in configuration
	GetChainName() string

	// `defer` is the diff result of normal chain node height and extend chain node height
	GetDefer() uint64
}
