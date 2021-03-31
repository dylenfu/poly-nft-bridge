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

package main

import "github.com/ethereum/go-ethereum/common"

type Config struct {
	Ethereum *ChainConfig
	Bsc      *ChainConfig
	Heco     *ChainConfig
	Poly     *PolyConfig

	// leveldb direction
	LevelDB string

	// oss
	OSS string
}

type ChainConfig struct {
	ChainID  uint64
	RPC      string
	Admin    string
	Keystore string

	ECCD common.Address
	ECCM common.Address
	CCMP common.Address

	NFTLockProxy common.Address
	NFTWrap      common.Address
	FeeToken     common.Address
	FeeCollector common.Address
}

type PolyConfig struct {
	RPC        string
	Keystore   string
	Passphrase string
}
