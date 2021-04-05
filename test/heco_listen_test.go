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

package test

import (
	"testing"

	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/polynetwork/poly-nft-bridge/dao/crosschaindao"
	wp "github.com/polynetwork/poly-nft-bridge/wrap"
	"github.com/stretchr/testify/assert"
)

func Test_HecoListen(t *testing.T) {
	dao := crosschaindao.NewCrossChainDao(basedef.SERVER_STAKE, true, config.DBConfig)
	assert.NotNil(t, dao)

	cfg := config.GetChainListenConfig(basedef.HECO_CROSSCHAIN_ID)
	assert.NotNil(t, cfg)
	chainHandle := wp.NewChainHandle(cfg)
	chainListen := wp.NewCrossChainListen(chainHandle, dao)
	chainListen.ListenChain()
}
