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

package conf

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	basedef "github.com/polynetwork/poly-nft-bridge/const"
)

type DBConfig struct {
	URL      string
	User     string
	Password string
	Scheme   string
	Debug    bool
}

type Restful struct {
	Url string
	Key string
}

type ChainListenConfig struct {
	ChainName       string
	ChainId         uint64
	ListenSlot      uint64
	Defer           uint64
	Nodes           []*Restful
	ExtendNodes     []*Restful
	WrapperContract string
	ECCMContract    string
	ProxyContract   string
}

func (cfg *ChainListenConfig) GetNodesUrl() []string {
	urls := make([]string, 0)
	for _, node := range cfg.Nodes {
		urls = append(urls, node.Url)
	}
	return urls
}

func (cfg *ChainListenConfig) GetNodesKey() []string {
	keys := make([]string, 0)
	for _, node := range cfg.Nodes {
		keys = append(keys, node.Key)
	}
	return keys
}

func (cfg *ChainListenConfig) GetExtendNodesUrl() []string {
	urls := make([]string, 0)
	for _, node := range cfg.ExtendNodes {
		urls = append(urls, node.Url)
	}
	return urls
}

func (cfg *ChainListenConfig) GetExtendNodesKey() []string {
	keys := make([]string, 0)
	for _, node := range cfg.ExtendNodes {
		keys = append(keys, node.Key)
	}
	return keys
}

type Config struct {
	Server            string
	Backup            bool
	ChainListenConfig []*ChainListenConfig
	DBConfig          *DBConfig
}

func (cfg *Config) GetChainListenConfig(chainId uint64) *ChainListenConfig {
	for _, chainListenConfig := range cfg.ChainListenConfig {
		if chainListenConfig.ChainId == chainId {
			return chainListenConfig
		}
	}
	return nil
}

func NewConfig(filePath string) *Config {
	fileContent, err := basedef.ReadFile(filePath)
	if err != nil {
		logs.Error("NewServiceConfig: failed, err: %s", err)
		return nil
	}
	config := &Config{}
	err = json.Unmarshal(fileContent, config)
	if err != nil {
		logs.Error("NewServiceConfig: failed, err: %s", err)
		return nil
	}
	return config
}
