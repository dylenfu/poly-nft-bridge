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

import (
	"strings"

	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/urfave/cli"
)

var (
	LogLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~6). 0:Trace 1:Debug 2:Info 3:Warn 4:Error 5:Fatal 6:MaxLevel",
		Value: 1,
	}

	LogDirFlag = cli.StringFlag{
		Name:  "logdir",
		Usage: "log directory",
		Value: "./logs",
	}

	ConfigPathFlag = cli.StringFlag{
		Name:  "cliconfig",
		Usage: "Server config file `<path>`",
		Value: "./config.json",
	}

	ChainIDFlag = cli.Uint64Flag{
		Name:  "chain",
		Usage: "select chainID",
		Value: basedef.ETHEREUM_CROSSCHAIN_ID,
	}

	NFTNameFlag = cli.StringFlag{
		Name:  "name",
		Usage: "set nft name for deploy nft contract, etc.",
		Value: "",
	}

	NFTSymbolFlag = cli.StringFlag{
		Name:  "symbol",
		Usage: "set nft symbol for deploy nft contract, etc.",
		Value: "",
	}

	DstChainFlag = cli.Uint64Flag{
		Name:  "dstChain",
		Usage: "set dest chain for cross chain",
		Value: 0,
	}

	AssetFlag = cli.Uint64Flag{
		Name:  "asset",
		Usage: "set asset for cross chain or mint nft",
		Value: 0,
	}

	DstAssetFlag = cli.Uint64Flag{
		Name:  "dstAsset",
		Usage: "set dest asset for cross chain",
		Value: 0,
	}

	SrcAccountFlag = cli.StringFlag{
		Name:  "from",
		Usage: "set from account",
	}
	DstAccountFlag = cli.StringFlag{
		Name:  "to",
		Usage: "set to account",
	}

	FeeTokenFlag = cli.BoolTFlag{
		Name:  "feeToken",
		Usage: "choose erc20 token to be fee token",
	}

	ERC20TokenFlag = cli.BoolFlag{
		Name:  "erc20Token",
		Usage: "choose erc20 token to be fee token",
	}

	AmountFlag = cli.StringFlag{
		Name:  "amount",
		Usage: "transfer amount or fee amount",
		Value: "",
	}

	TokenIdFlag = cli.Uint64Flag{
		Name:  "tokenId",
		Usage: "set token id while mint nft",
	}

	LockIdFlag = cli.Uint64Flag{
		Name:  "lockId",
		Usage: "wrap lock nft item id",
	}
)

var (
	CmdDeployECCDContract = cli.Command{
		Name:   "deployECCD",
		Usage:  "admin account deploy ethereum cross chain data contract.",
		Action: handleCmdDeployECCDContract,
	}

	CmdDeployECCMContract = cli.Command{
		Name:   "deployECCM",
		Usage:  "admin account deploy ethereum cross chain manage contract.",
		Action: handleCmdDeployECCMContract,
	}

	CmdDeployCCMPContract = cli.Command{
		Name:   "deployCCMP",
		Usage:  "admin account deploy ethereum cross chain manager proxy contract.",
		Action: handleCmdDeployCCMPContract,
	}

	CmdDeployNFTContract = cli.Command{
		Name:   "deployNFT",
		Usage:  "admin account deploy new nft asset with mapping contract.",
		Action: handleCmdDeployNFTContract,
	}

	CmdDeployERC20Contract = cli.Command{
		Name:   "deployERC20",
		Usage:  "admin account deploy new mintable erc20 contract.",
		Action: handleCmdDeployERC20Contract,
		Flags: []cli.Flag{
			FeeTokenFlag,
		},
	}

	CmdDeployLockProxyContract = cli.Command{
		Name:   "deployNFTLockProxy",
		Usage:  "admin account deploy nft lock proxy contract.",
		Action: handleCmdDeployLockProxyContract,
	}

	CmdDeployNFTWrapContract = cli.Command{
		Name:   "deployNFTWrapper",
		Usage:  "admin account deploy nft wrapper contract.",
		Action: handleCmdDeployNFTWrapContract,
	}

	CmdLockProxySetCCMP = cli.Command{
		Name:   "proxySetCCMP",
		Usage:  "admin account set cross chain manager proxy address for lock proxy contract.",
		Action: handleCmdLockProxySetCCMP,
	}

	CmdBindLockProxy = cli.Command{
		Name:   "bindProxy",
		Usage:  "admin  account bind lock proxy contract with another side chain's lock proxy contract.",
		Action: handleCmdBindLockProxy,
		Flags: []cli.Flag{
			DstChainFlag,
		},
	}

	CmdGetBoundLockProxy = cli.Command{
		Name:   "getBoundProxy",
		Usage:  "get bound lock proxy contract.",
		Action: handleCmdGetBoundLockProxy,
		Flags: []cli.Flag{
			DstChainFlag,
		},
	}

	CmdBindNFTAsset = cli.Command{
		Name:   "bindNFT",
		Usage:  "admin account bind nft asset to side chain.",
		Action: handleCmdBindNFTAsset,
		Flags: []cli.Flag{
			AssetFlag,
			DstChainFlag,
			DstAssetFlag,
		},
	}

	CmdTransferECCDOwnership = cli.Command{
		Name:   "transferECCDOwnership",
		Usage:  "admin account transfer ethereum cross chain data contract ownership eccm contract.",
		Action: handleCmdTransferECCDOwnership,
	}

	CmdTransferECCMOwnership = cli.Command{
		Name:   "transferECCMOwnership",
		Usage:  "admin account transfer ethereum cross chain manager contract ownership to ccmp contract.",
		Action: handleCmdTransferECCMOwnership,
	}

	CmdSyncSideChainGenesis2Poly = cli.Command{
		Name:   "syncSideGenesis",
		Usage:  "sync side chain genesis header to poly chain.",
		Action: handleCmdSyncSideChainGenesis2Poly,
	}

	CmdSyncPolyGenesis2SideChain = cli.Command{
		Name:   "syncPolyGenesis",
		Usage:  "sync poly genesis header to side chain.",
		Action: handleCmdSyncPolyGenesis2SideChain,
	}

	CmdNFTWrapSetFeeCollector = cli.Command{
		Name:   "setFeeCollector",
		Usage:  "admin account set nft fee collecotr for wrap contract",
		Action: handleCmdNFTWrapSetFeeCollector,
	}

	CmdNFTWrapSetLockProxy = cli.Command{
		Name:   "setWrapLockProxy",
		Usage:  "admin account set nft lock proxy for wrap contract.",
		Action: handleCmdNFTWrapSetLockProxy,
	}

	CmdNFTMint = cli.Command{
		Name:   "mintNFT",
		Usage:  "admin account mint nft token.",
		Action: handleCmdNFTMint,
		Flags: []cli.Flag{
			AssetFlag,
			DstAccountFlag,
			TokenIdFlag,
		},
	}

	CmdNFTWrapLock = cli.Command{
		Name:   "lockNFT",
		Usage:  "lock nft token on wrap contract.",
		Action: handleCmdNFTWrapLock,
		Flags: []cli.Flag{
			SrcAccountFlag,
			AssetFlag,
			DstChainFlag,
			DstAccountFlag,
			TokenIdFlag,
			AmountFlag,
			LockIdFlag,
		},
	}

	CmdERC20Mint = cli.Command{
		Name:   "mintERC20",
		Usage:  "admin account mint erc20 token.",
		Action: handleCmdERC20Mint,
		Flags: []cli.Flag{
			FeeTokenFlag,
			ERC20TokenFlag,
			DstAccountFlag,
			AmountFlag,
		},
	}

	CmdERC20Transfer = cli.Command{
		Name:   "transferERC20",
		Usage:  "transfer ERC20 token.",
		Action: handleCmdERC20Transfer,
		Flags: []cli.Flag{
			FeeTokenFlag,
			ERC20TokenFlag,
			SrcAccountFlag,
			DstAccountFlag,
			AmountFlag,
		},
	}
)

//getFlagName deal with short flag, and return the flag name whether flag name have short name
func getFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}
