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
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"runtime"

	log "github.com/astaxie/beego/logs"
	"github.com/ethereum/go-ethereum/common"
	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/polynetwork/poly-nft-bridge/sdk/eth_sdk"
	"github.com/polynetwork/poly-nft-bridge/sdk/poly_sdk"
	xecdsa "github.com/polynetwork/poly-nft-bridge/utils/ecdsa"
	"github.com/polynetwork/poly-nft-bridge/utils/files"
	"github.com/polynetwork/poly-nft-bridge/utils/leveldb"
	"github.com/polynetwork/poly-nft-bridge/utils/math"
	"github.com/polynetwork/poly-nft-bridge/utils/wallet"
	"github.com/urfave/cli"
)

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

var (
	cfgPath string
	cfg     = new(Config)
	cc      *ChainConfig
	storage *leveldb.LevelDBImpl
	sdk     *eth_sdk.EthereumSdk
	adm     *ecdsa.PrivateKey
)

const defaultAccPwd = "111111"

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "poly nftbridge deploy tool"
	app.Action = startServer
	app.Version = "1.0.0"
	app.Copyright = "Copyright in 2020 The Ontology Authors"
	app.Flags = []cli.Flag{
		LogLevelFlag,
		LogDirFlag,
		ConfigPathFlag,
		ChainIDFlag,
		NFTNameFlag,
		NFTSymbolFlag,
		DstChainFlag,
		AssetFlag,
		DstAssetFlag,
		SrcAccountFlag,
		DstAccountFlag,
		FeeTokenFlag,
		ERC20TokenFlag,
		AmountFlag,
		TokenIdFlag,
	}
	app.Commands = []cli.Command{
		CmdDeployECCDContract,
		CmdDeployECCMContract,
		CmdDeployCCMPContract,
		CmdDeployNFTContract,
		CmdDeployERC20Contract,
		CmdDeployLockProxyContract,
		CmdDeployNFTWrapContract,
		CmdLockProxySetCCMP,
		CmdBindLockProxy,
		CmdGetBoundLockProxy,
		CmdBindNFTAsset,
		CmdTransferECCDOwnership,
		CmdTransferECCMOwnership,
		CmdSyncSideChainGenesis2Poly,
		CmdSyncPolyGenesis2SideChain,
		CmdNFTWrapSetFeeCollector,
		CmdNFTWrapSetLockProxy,
		CmdNFTMint,
		CmdNFTWrapLock,
		CmdERC20Mint,
		CmdERC20Transfer,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startServer(ctx *cli.Context) {
	var err error
	// load config instance
	cfgPath = ctx.GlobalString(getFlagName(ConfigPathFlag))
	if err = files.ReadJsonFile(cfgPath, cfg); err != nil {
		panic(fmt.Sprintf("read config json file, err: %v", err))
	}

	logDir := ctx.GlobalString(getFlagName(LogDirFlag))
	if err = log.SetLogger(log.AdapterFile, fmt.Sprintf(`{"filename":"%d/fee_listen.log"}`, logDir)); err != nil {
		panic(fmt.Sprintf("set logger failed, err: %v", err))
	}

	// prepare storage for persist account passphrase
	storage = leveldb.NewLevelDBInstance(cfg.LevelDB)

	// select src chainID and prepare config and accounts
	chainID := ctx.GlobalUint64(getFlagName(ChainIDFlag))
	selectChainConfig(chainID)

	if sdk, err = eth_sdk.NewEthereumSdk(cc.RPC); err != nil {
		panic(fmt.Sprintf("generate sdk for chain %d faild, err: %v", cc.ChainID, err))
	}

	if adm, err = wallet.LoadEthAccount(storage, cc.Keystore, cc.Admin, defaultAccPwd); err != nil {
		panic(fmt.Sprintf("load eth account for chain %d faild, err: %v", cc.ChainID, err))
	}
}

func handleCmdDeployECCDContract(ctx *cli.Context) error {
	addr, err := sdk.DeployECCDContract(adm)
	if err != nil {
		return fmt.Errorf("deploy eccd for chain %d failed, err: %v", cc.ChainID, err)
	}

	cc.ECCD = addr
	log.Info("deploy eccd for chain %d success %s", cc.ChainID, addr.Hex())
	return updateConfig()
}

func handleCmdDeployECCMContract(ctx *cli.Context) error {
	addr, err := sdk.DeployECCMContract(adm, cc.ECCD, cc.ChainID)
	if err != nil {
		return fmt.Errorf("deploy eccm for chain %d failed, err: %v", cc.ChainID, err)
	}
	cc.ECCM = addr
	log.Info("deploy eccm for chain %d success %s", cc.ChainID, addr.Hex())
	return updateConfig()
}

func handleCmdDeployCCMPContract(ctx *cli.Context) error {
	addr, err := sdk.DeployECCMPContract(adm, cc.ECCM)
	if err != nil {
		return fmt.Errorf("deploy ccmp for chain %d failed, err: %v", cc.ChainID, err)
	}
	cc.CCMP = addr
	log.Info("deploy ccmp for chain %d success %s", cc.ChainID, addr.Hex())
	return updateConfig()
}

func handleCmdDeployNFTContract(ctx *cli.Context) error {
	name := ctx.GlobalString(getFlagName(NFTNameFlag))
	symbol := ctx.GlobalString(getFlagName(NFTSymbolFlag))
	owner := xecdsa.Key2address(adm)
	if addr, err := sdk.DeployNFT(adm, cc.NFTLockProxy, name, symbol); err != nil {
		return fmt.Errorf("deploy nft contract for owner %s on chain %d failed, err: %v", owner.Hex(), cc.ChainID, err)
	} else {
		log.Info("deploy nft contract %s for user %s on chain %d success!", addr.Hex(), owner.Hex(), cc.ChainID)
	}
	return nil
}

func handleCmdDeployERC20Contract(ctx *cli.Context) error {
	addr, err := sdk.DeployERC20(adm)
	if err != nil {
		return fmt.Errorf("deploy erc20 token failed")
	}

	log.Info("deploy erc20 %s success", addr.Hex())
	isFeeToken := ctx.GlobalBool(getFlagName(FeeTokenFlag))
	if !isFeeToken {
		return nil
	}

	cc.FeeToken = addr
	return updateConfig()
}

func handleCmdDeployLockProxyContract(ctx *cli.Context) error {
	addr, err := sdk.DeployNFTLockProxy(adm)
	if err != nil {
		return fmt.Errorf("deploy nft lock proxy for chain %d failed, err: %v", cc.ChainID, err)
	}
	cc.NFTLockProxy = addr
	log.Info("deploy nft lock proxy for chain %d success %s", cc.ChainID, addr.Hex())
	return updateConfig()
}

func handleCmdDeployNFTWrapContract(ctx *cli.Context) error {
	addr, err := sdk.DeployWrapContract(adm, cc.ChainID)
	if err != nil {
		return err
	}

	cc.NFTWrap = addr
	log.Info("deploy wrap contract %s success!", addr.Hex())
	return updateConfig()
}

func handleCmdLockProxySetCCMP(ctx *cli.Context) error {
	hash, err := sdk.NFTLockProxySetCCMP(adm, cc.NFTLockProxy, cc.CCMP)
	if err != nil {
		return fmt.Errorf("nft lock proxy set ccmp for chain %d failed, err: %v", cc.ChainID, err)
	}
	log.Info("nft lock proxy set ccmp for chain %d success! hash %s", cc.ChainID, hash.Hex())
	return nil
}

func handleCmdBindLockProxy(ctx *cli.Context) error {
	dstChainId := ctx.GlobalUint64(getFlagName(DstChainFlag))
	dstChainCfg := customSelectChainConfig(dstChainId)

	hash, err := sdk.BindLockProxy(adm, cc.NFTLockProxy, dstChainCfg.NFTLockProxy, dstChainId)
	if err != nil {
		return fmt.Errorf("bind lock proxy (src proxy %s, dst proxy %s, src chain id %d, dst chain id %d)",
			cc.NFTLockProxy.Hex(), dstChainCfg.NFTLockProxy.Hex(), cc.ChainID, dstChainId)
	}

	log.Info("bind lock proxy (src proxy %s, dst proxy %s, src chain id %d, dst chain id %d), txhash %s",
		cc.NFTLockProxy.Hex(), dstChainCfg.NFTLockProxy.Hex(), cc.ChainID, dstChainId, hash.Hex())
	return nil
}

func handleCmdGetBoundLockProxy(ctx *cli.Context) error {
	dstChainId := ctx.GlobalUint64(getFlagName(DstChainFlag))

	if addr, err := sdk.GetBoundNFTProxy(cc.NFTLockProxy, dstChainId); err != nil {
		return fmt.Errorf("check bound nft lock proxy err: %v", err)
	} else {
		log.Info("proxy %s bound to %s in with target chain id of %d", cc.NFTLockProxy.Hex(), addr.Hex(), dstChainId)
	}

	return nil
}

func handleCmdBindNFTAsset(ctx *cli.Context) error {
	dstChainId := ctx.GlobalUint64(getFlagName(DstChainFlag))
	srcAsset := common.HexToAddress(ctx.GlobalString(getFlagName(AssetFlag)))
	dstAsset := common.HexToAddress(ctx.GlobalString(getFlagName(DstAssetFlag)))
	dstChainCfg := customSelectChainConfig(dstChainId)

	owner := xecdsa.Key2address(adm)
	hash, err := sdk.BindNFTAsset(
		adm,
		cc.NFTLockProxy,
		srcAsset,
		dstAsset,
		dstChainId,
	)
	if err != nil {
		return fmt.Errorf("bind nft asset (src chain id %d, src asset %s, src proxy %s) - "+
			"(dst chain id %d, dst asset %s, dst proxy %s)"+
			" for user %s failed, err: %v",
			cc.ChainID, srcAsset.Hex(), cc.NFTLockProxy.Hex(),
			dstChainId, dstAsset.Hex(), dstChainCfg.NFTLockProxy.Hex(),
			owner.Hex(), err)
	}

	log.Info("bind nft asset (src chain id %d, src asset %s, src proxy %s) - "+
		"(dst chain id %d, dst asset %s, dst proxy %s)"+
		" for user %s success! txhash %s",
		cc.ChainID, srcAsset.Hex(), cc.NFTLockProxy.Hex(),
		dstChainId, dstAsset.Hex(), dstChainCfg.NFTLockProxy.Hex(),
		owner.Hex(), hash.Hex())
	return nil
}

func handleCmdTransferECCDOwnership(ctx *cli.Context) error {
	if hash, err := sdk.TransferECCDOwnership(adm, cc.ECCD, cc.ECCM); err != nil {
		return fmt.Errorf("transfer eccd %s ownership to eccm %s on chain %d failed, err: %v",
			cc.ECCD.Hex(), cc.ECCM.Hex(), cc.ChainID, err)
	} else {
		log.Info("transfer eccd %s ownership to eccm %s on chain %d success, txhash: %s",
			cc.ECCD.Hex(), cc.ECCM.Hex(), cc.ChainID, hash.Hex())
	}
	return nil
}

func handleCmdTransferECCMOwnership(ctx *cli.Context) error {
	if hash, err := sdk.TransferECCMOwnership(adm, cc.ECCM, cc.CCMP); err != nil {
		return fmt.Errorf("transfer eccm %s ownership to ccmp %s on chain %d failed, err: %v",
			cc.ECCM.Hex(), cc.CCMP.Hex(), cc.ChainID, err)
	} else {
		log.Info("transfer eccm %s ownership to ccmp %s on chain %d success, txhash: %s",
			cc.ECCM.Hex(), cc.CCMP.Hex(), cc.ChainID, hash.Hex())
	}

	return nil
}

func handleCmdSyncSideChainGenesis2Poly(ctx *cli.Context) error {
	polySdk := poly_sdk.NewPolySDK(cfg.Poly.RPC)
	validators := wallet.LoadPolyAccountList(cfg.Poly.Keystore, cfg.Poly.Passphrase)
	if err := SyncSideChainGenesisHeaderToPolyChain(
		cc.ChainID,
		sdk,
		polySdk,
		validators,
	); err != nil {
		return fmt.Errorf("sync side chain %d genesis header to poly failed, err: %v", cc.ChainID, err)
	} else {
		log.Info("sync side chain %d genesis header to poly success!", cc.ChainID)
	}
	return nil
}

func handleCmdSyncPolyGenesis2SideChain(ctx *cli.Context) error {
	polySdk := poly_sdk.NewPolySDK(cfg.Poly.RPC)
	err := SyncPolyChainGenesisHeader2SideChain(
		polySdk,
		adm,
		sdk,
		cc.ECCM,
	)
	if err != nil {
		return fmt.Errorf("sync poly chain genesis header to side chain %d failed, err: %v", cc.ChainID, err)
	}
	log.Info("sync poly chain genesis header to side chain %d success!", cc.ChainID)
	return nil
}

func handleCmdNFTWrapSetFeeCollector(ctx *cli.Context) error {
	if tx, err := sdk.SetWrapFeeCollector(adm, cc.NFTWrap, cc.FeeCollector); err != nil {
		return fmt.Errorf("set fee collector failed, err: %v", err)
	} else {
		log.Info("set fee collector success, hash %s", tx.Hex())
	}
	return nil
}

func handleCmdNFTWrapSetLockProxy(ctx *cli.Context) error {
	if tx, err := sdk.SetWrapLockProxy(adm, cc.NFTWrap, cc.NFTLockProxy); err != nil {
		return fmt.Errorf("set lock proxy for wrap failed, err: %v", err)
	} else {
		log.Info("set wrap lock proxy success, hash %s", tx.Hex())
	}
	return nil
}

func handleCmdNFTMint(ctx *cli.Context) error {
	asset := common.HexToAddress(ctx.GlobalString(getFlagName(AssetFlag)))
	to := common.HexToAddress(ctx.GlobalString(getFlagName(DstAccountFlag)))
	tokenID := new(big.Int).SetUint64(ctx.GlobalUint64(getFlagName(TokenIdFlag)))
	uri := cfg.OSS + tokenID.String()
	tx, err := sdk.MintNFT(adm, asset, to, tokenID, uri)
	if err != nil {
		return err
	}
	log.Info("mint nft %d to %s success! txhash %s", tokenID.Uint64(), to.Hex(), tx.Hex())
	return nil
}

func handleCmdNFTWrapLock(ctx *cli.Context) error {
	from := common.HexToAddress(ctx.GlobalString(getFlagName(SrcAccountFlag)))
	key, err := wallet.LoadEthAccount(storage, cc.Keystore, from.Hex(), defaultAccPwd)
	if err != nil {
		return err
	}

	asset := common.HexToAddress(ctx.GlobalString(getFlagName(AssetFlag)))
	to := common.HexToAddress(ctx.GlobalString(getFlagName(DstAccountFlag)))
	dstChainID := ctx.GlobalUint64(getFlagName(DstChainFlag))
	tokenID := new(big.Int).SetUint64(ctx.GlobalUint64(getFlagName(TokenIdFlag)))
	feeAmount := math.String2BigInt(ctx.GlobalString(getFlagName(AmountFlag)))
	id := new(big.Int).SetUint64(ctx.GlobalUint64(getFlagName(LockIdFlag)))

	tx, err := sdk.WrapLock(key, cc.NFTWrap, asset, to, dstChainID, tokenID, cc.FeeToken, feeAmount, id)
	if err != nil {
		return err
	}
	log.Info("wrap lock nft %d success! txhash %s", tokenID, tx.Hex())
	return nil
}

func handleCmdERC20Mint(ctx *cli.Context) error {
	var asset common.Address
	isFeeToken := ctx.GlobalBool(getFlagName(FeeTokenFlag))
	if isFeeToken {
		asset = cc.FeeToken
	} else {
		asset = common.HexToAddress(ctx.GlobalString(getFlagName(ERC20TokenFlag)))
	}
	to := common.HexToAddress(ctx.GlobalString(getFlagName(DstAccountFlag)))
	param := ctx.GlobalString(getFlagName(AmountFlag))
	amount := math.String2BigInt(param)

	tx, err := sdk.MintERC20Token(adm, asset, to, amount)
	if err != nil {
		return err
	}
	log.Info("mint %s to %s success, hash %s", amount.String(), to.Hex(), tx.Hex())
	return nil
}

func handleCmdERC20Transfer(ctx *cli.Context) error {
	var asset common.Address
	isFeeToken := ctx.GlobalBool(getFlagName(FeeTokenFlag))
	if isFeeToken {
		asset = cc.FeeToken
	} else {
		asset = common.HexToAddress(ctx.GlobalString(getFlagName(ERC20TokenFlag)))
	}

	from := common.HexToAddress(ctx.GlobalString(getFlagName(SrcAccountFlag)))
	key, err := wallet.LoadEthAccount(storage, cc.Keystore, from.Hex(), defaultAccPwd)
	if err != nil {
		return err
	}

	to := common.HexToAddress(ctx.GlobalString(getFlagName(DstAccountFlag)))
	param := ctx.GlobalString(getFlagName(AmountFlag))
	amount := math.String2BigInt(param)
	tx, err := sdk.TransferERC20Token(key, asset, to, amount)
	if err != nil {
		return err
	}
	log.Info("%s transfer %s to %s success, tx %s", from.Hex(), amount.String(), to.Hex(), tx.Hex())
	return nil
}

func updateConfig() error {
	return files.WriteJsonFile(cfgPath, cfg, true)
}

func selectChainConfig(chainID uint64) {
	cc = customSelectChainConfig(chainID)
}

func customSelectChainConfig(chainID uint64) *ChainConfig {
	switch chainID {
	case basedef.ETHEREUM_CROSSCHAIN_ID:
		return cfg.Ethereum
	case basedef.BSC_CROSSCHAIN_ID:
		return cfg.Bsc
	case basedef.HECO_CROSSCHAIN_ID:
		return cfg.Heco
	}
	panic(fmt.Sprintf("invalid chain id %d", chainID))
}
