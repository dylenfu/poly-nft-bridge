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
	app.Version = "1.0.0"
	app.Copyright = "Copyright in 2020 The Ontology Authors"
	app.Flags = []cli.Flag{
		LogLevelFlag,
		//LogDirFlag,
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
		CmdSample,
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
		CmdRegisterSideChain,
		CmdApproveSideChain,
		CmdSyncSideChainGenesis2Poly,
		CmdSyncPolyGenesis2SideChain,
		CmdNFTWrapSetFeeCollector,
		CmdNFTWrapSetLockProxy,
		CmdNFTMint,
		CmdNFTWrapLock,
		CmdERC20Mint,
		CmdERC20Transfer,
		CmdGetERC20Balance,
	}

	app.Before = beforeCommands
	return app
}

func main() {

	app := setupApp()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// action execute after commands
func beforeCommands(ctx *cli.Context) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// load config instance
	cfgPath = ctx.GlobalString(getFlagName(ConfigPathFlag))
	if err = files.ReadJsonFile(cfgPath, cfg); err != nil {
		return fmt.Errorf("read config json file, err: %v", err)
	}

	//logDir := ctx.GlobalString(getFlagName(LogDirFlag))
	//logFormat := fmt.Sprintf(`{"filename":"%s/deploy.log", "perm": "0777"}`, logDir)
	loglevel := ctx.GlobalUint64(getFlagName(LogLevelFlag))
	logFormat := fmt.Sprintf(`{"level:":"%d"}`, loglevel)
	if err := log.SetLogger("console", logFormat); err != nil {
		return fmt.Errorf("set logger failed, err: %v", err)
	}

	// prepare storage for persist account passphrase
	storage = leveldb.NewLevelDBInstance(cfg.LevelDB)

	// select src chainID and prepare config and accounts
	chainID := ctx.GlobalUint64(getFlagName(ChainIDFlag))
	selectChainConfig(chainID)

	if sdk, err = eth_sdk.NewEthereumSdk(cc.RPC); err != nil {
		return fmt.Errorf("generate sdk for chain %d faild, err: %v", cc.ChainID, err)
	}

	if adm, err = wallet.LoadEthAccount(storage, cc.Keystore, cc.Admin, defaultAccPwd); err != nil {
		return fmt.Errorf("load eth account for chain %d faild, err: %v", cc.ChainID, err)
	}

	return nil
}

func handleSample(ctx *cli.Context) error {
	log.Info("start to debug sample...")
	feeToken := ctx.BoolT(getFlagName(FeeTokenFlag))
	log.Info("feeToken %v", feeToken)
	return nil
}

func handleCmdDeployECCDContract(ctx *cli.Context) error {
	log.Info("start to deploy eccd contract...")

	addr, err := sdk.DeployECCDContract(adm)
	if err != nil {
		return fmt.Errorf("deploy eccd for chain %d failed, err: %v", cc.ChainID, err)
	}

	cc.ECCD = addr.Hex()
	log.Info("deploy eccd for chain %d success %s", cc.ChainID, addr.Hex())
	return updateConfig()
}

func handleCmdDeployECCMContract(ctx *cli.Context) error {
	log.Info("start to deploy eccm contract...")

	eccd := common.HexToAddress(cc.ECCD)
	addr, err := sdk.DeployECCMContract(adm, eccd, cc.ChainID)
	if err != nil {
		return fmt.Errorf("deploy eccm for chain %d failed, err: %v", cc.ChainID, err)
	}
	cc.ECCM = addr.Hex()
	log.Info("deploy eccm for chain %d success %s", cc.ChainID, addr.Hex())
	return updateConfig()
}

func handleCmdDeployCCMPContract(ctx *cli.Context) error {
	log.Info("start to deploy ccmp contract...")

	eccm := common.HexToAddress(cc.ECCM)
	addr, err := sdk.DeployECCMPContract(adm, eccm)
	if err != nil {
		return fmt.Errorf("deploy ccmp for chain %d failed, err: %v", cc.ChainID, err)
	}
	cc.CCMP = addr.Hex()
	log.Info("deploy ccmp for chain %d success %s", cc.ChainID, addr.Hex())
	return updateConfig()
}

func handleCmdDeployNFTContract(ctx *cli.Context) error {
	log.Info("start to deploy nft contract...")

	name := flag2string(ctx, NFTNameFlag)
	symbol := flag2string(ctx, NFTSymbolFlag)
	owner := xecdsa.Key2address(adm)
	proxy := common.HexToAddress(cc.NFTLockProxy)
	if addr, err := sdk.DeployNFT(adm, proxy, name, symbol); err != nil {
		return fmt.Errorf("deploy nft contract for owner %s on chain %d failed, err: %v", owner.Hex(), cc.ChainID, err)
	} else {
		log.Info("deploy nft contract %s for user %s on chain %d success!", addr.Hex(), owner.Hex(), cc.ChainID)
	}
	return nil
}

func handleCmdDeployERC20Contract(ctx *cli.Context) error {
	log.Info("start to deploy erc20 token......")

	addr, err := sdk.DeployERC20(adm)
	if err != nil {
		return fmt.Errorf("deploy erc20 token failed, err: %v", err)
	}

	log.Info("deploy erc20 %s success", addr.Hex())
	isFeeToken := ctx.BoolT(getFlagName(FeeTokenFlag))
	if !isFeeToken {
		return nil
	}

	cc.FeeToken = addr.Hex()
	return updateConfig()
}

func handleCmdDeployLockProxyContract(ctx *cli.Context) error {
	log.Info("start to deploy nft lock proxy contract...")

	addr, err := sdk.DeployNFTLockProxy(adm)
	if err != nil {
		return fmt.Errorf("deploy nft lock proxy for chain %d failed, err: %v", cc.ChainID, err)
	}
	cc.NFTLockProxy = addr.Hex()
	log.Info("deploy nft lock proxy for chain %d success %s", cc.ChainID, addr.Hex())
	return updateConfig()
}

func handleCmdDeployNFTWrapContract(ctx *cli.Context) error {
	log.Info("start to deploy nft wrap contract...")

	addr, err := sdk.DeployWrapContract(adm, cc.ChainID)
	if err != nil {
		return err
	}

	cc.NFTWrap = addr.Hex()
	log.Info("deploy wrap contract %s success!", addr.Hex())
	return updateConfig()
}

func handleCmdLockProxySetCCMP(ctx *cli.Context) error {
	log.Info("start to set ccmp for lock proxy contract...")

	proxy := common.HexToAddress(cc.NFTLockProxy)
	ccmp := common.HexToAddress(cc.CCMP)
	hash, err := sdk.NFTLockProxySetCCMP(adm, proxy, ccmp)
	if err != nil {
		return fmt.Errorf("nft lock proxy set ccmp for chain %d failed, err: %v", cc.ChainID, err)
	}
	log.Info("nft lock proxy set ccmp for chain %d success! hash %s", cc.ChainID, hash.Hex())
	return nil
}

func handleCmdBindLockProxy(ctx *cli.Context) error {
	log.Info("start to bind lock proxy...")

	dstChainId := flag2Uint64(ctx, DstChainFlag)
	dstChainCfg := customSelectChainConfig(dstChainId)
	proxy := common.HexToAddress(cc.NFTLockProxy)
	dstProxy := common.HexToAddress(dstChainCfg.NFTLockProxy)

	hash, err := sdk.BindLockProxy(adm, proxy, dstProxy, dstChainId)
	if err != nil {
		return fmt.Errorf("bind lock proxy (src proxy %s, dst proxy %s, src chain id %d, dst chain id %d)",
			cc.NFTLockProxy, dstChainCfg.NFTLockProxy, cc.ChainID, dstChainId)
	}

	log.Info("bind lock proxy (src proxy %s, dst proxy %s, src chain id %d, dst chain id %d), txhash %s",
		cc.NFTLockProxy, dstChainCfg.NFTLockProxy, cc.ChainID, dstChainId, hash.Hex())
	return nil
}

func handleCmdGetBoundLockProxy(ctx *cli.Context) error {
	log.Info("start to get bound lock proxy contract...")

	dstChainId := ctx.GlobalUint64(getFlagName(DstChainFlag))
	proxy := common.HexToAddress(cc.NFTLockProxy)

	if addr, err := sdk.GetBoundNFTProxy(proxy, dstChainId); err != nil {
		return fmt.Errorf("check bound nft lock proxy err: %v", err)
	} else {
		log.Info("proxy %s bound to %s in with target chain id of %d", cc.NFTLockProxy, addr.Hex(), dstChainId)
	}

	return nil
}

func handleCmdBindNFTAsset(ctx *cli.Context) error {
	log.Info("start to bind nft asset...")

	srcAsset := flag2address(ctx, AssetFlag)
	dstAsset := flag2address(ctx, DstAssetFlag)
	dstChainId := flag2Uint64(ctx, DstChainFlag)
	dstChainCfg := customSelectChainConfig(dstChainId)
	owner := xecdsa.Key2address(adm)
	proxy := common.HexToAddress(cc.NFTLockProxy)

	hash, err := sdk.BindNFTAsset(
		adm,
		proxy,
		srcAsset,
		dstAsset,
		dstChainId,
	)
	if err != nil {
		return fmt.Errorf("bind nft asset (src chain id %d, src asset %s, src proxy %s) - "+
			"(dst chain id %d, dst asset %s, dst proxy %s)"+
			" for user %s failed, err: %v",
			cc.ChainID, srcAsset.Hex(), cc.NFTLockProxy,
			dstChainId, dstAsset.Hex(), dstChainCfg.NFTLockProxy,
			owner.Hex(), err)
	}

	log.Info("bind nft asset (src chain id %d, src asset %s, src proxy %s) - "+
		"(dst chain id %d, dst asset %s, dst proxy %s)"+
		" for user %s success! txhash %s",
		cc.ChainID, srcAsset.Hex(), cc.NFTLockProxy,
		dstChainId, dstAsset.Hex(), dstChainCfg.NFTLockProxy,
		owner.Hex(), hash.Hex())
	return nil
}

func handleCmdTransferECCDOwnership(ctx *cli.Context) error {
	log.Info("start to transfer eccd ownership...")

	eccd := common.HexToAddress(cc.ECCD)
	eccm := common.HexToAddress(cc.ECCM)

	if hash, err := sdk.TransferECCDOwnership(adm, eccd, eccm); err != nil {
		return fmt.Errorf("transfer eccd %s ownership to eccm %s on chain %d failed, err: %v",
			cc.ECCD, cc.ECCM, cc.ChainID, err)
	} else {
		log.Info("transfer eccd %s ownership to eccm %s on chain %d success, txhash: %s",
			cc.ECCD, cc.ECCM, cc.ChainID, hash.Hex())
	}
	return nil
}

func handleCmdTransferECCMOwnership(ctx *cli.Context) error {
	log.Info("start to transfer eccm ownership...")

	eccm := common.HexToAddress(cc.ECCM)
	ccmp := common.HexToAddress(cc.CCMP)

	if hash, err := sdk.TransferECCMOwnership(adm, eccm, ccmp); err != nil {
		return fmt.Errorf("transfer eccm %s ownership to ccmp %s on chain %d failed, err: %v",
			cc.ECCM, cc.CCMP, cc.ChainID, err)
	} else {
		log.Info("transfer eccm %s ownership to ccmp %s on chain %d success, txhash: %s",
			cc.ECCM, cc.CCMP, cc.ChainID, hash.Hex())
	}

	return nil
}

func handleCmdRegisterSideChain(ctx *cli.Context) error {
	validators, err := wallet.LoadPolyAccountList(cfg.Poly.Keystore, cfg.Poly.Passphrase)
	if err != nil {
		return err
	}
	polySdk, err := poly_sdk.NewPolySdkAndSetChainID(cfg.Poly.RPC)
	if err != nil {
		return err
	}

	eccd := common.HexToAddress(cc.ECCD)
	router := cc.ChainID
	if err := polySdk.RegisterSideChain(validators[0], cc.ChainID, router, eccd, cc.ChainName); err != nil {
		return fmt.Errorf("failed to register side chain, err: %s", err)
	}

	log.Info("register side chain %d eccd %s success", cc.ChainID, eccd.Hex())
	return nil
}

func handleCmdApproveSideChain(ctx *cli.Context) error {
	validators, err := wallet.LoadPolyAccountList(cfg.Poly.Keystore, cfg.Poly.Passphrase)
	if err != nil {
		return err
	}
	polySdk, err := poly_sdk.NewPolySdkAndSetChainID(cfg.Poly.RPC)
	if err != nil {
		return err
	}

	if err := polySdk.ApproveRegisterSideChain(cc.ChainID, validators); err != nil {
		return fmt.Errorf("failed to approve register side chain, err: %s", err)
	}

	log.Info("approve register side chain %d success", cc.ChainID)
	return nil
}

func handleCmdSyncSideChainGenesis2Poly(ctx *cli.Context) error {
	log.Info("start to sync side chain genesis header to poly chain...")

	polySdk, err := poly_sdk.NewPolySdkAndSetChainID(cfg.Poly.RPC)
	if err != nil {
		return err
	}
	validators, err := wallet.LoadPolyAccountList(cfg.Poly.Keystore, cfg.Poly.Passphrase)
	if err != nil {
		return err
	}

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
	log.Info("start to sync poly chain genesis header to side chain...")

	polySdk, err := poly_sdk.NewPolySdkAndSetChainID(cfg.Poly.RPC)
	if err != nil {
		return err
	}
	eccm := common.HexToAddress(cc.ECCM)

	if err := SyncPolyChainGenesisHeader2SideChain(
		polySdk,
		adm,
		sdk,
		eccm,
	); err != nil {
		return fmt.Errorf("sync poly chain genesis header to side chain %d failed, err: %v", cc.ChainID, err)
	}
	log.Info("sync poly chain genesis header to side chain %d success!", cc.ChainID)
	return nil
}

func handleCmdNFTWrapSetFeeCollector(ctx *cli.Context) error {
	log.Info("start to set fee collector for wrap contract...")

	wrapper := common.HexToAddress(cc.NFTWrap)
	feeCollector := common.HexToAddress(cc.FeeCollector)

	if tx, err := sdk.SetWrapFeeCollector(adm, wrapper, feeCollector); err != nil {
		return fmt.Errorf("set fee collector failed, err: %v", err)
	} else {
		log.Info("set fee collector success, hash %s", tx.Hex())
	}
	return nil
}

func handleCmdNFTWrapSetLockProxy(ctx *cli.Context) error {
	log.Info("start to set lock proxy for wrap contract...")

	wrapper := common.HexToAddress(cc.NFTWrap)
	proxy := common.HexToAddress(cc.NFTLockProxy)

	if tx, err := sdk.SetWrapLockProxy(adm, wrapper, proxy); err != nil {
		return fmt.Errorf("set lock proxy for wrap failed, err: %v", err)
	} else {
		log.Info("set wrap lock proxy success, hash %s", tx.Hex())
	}
	return nil
}

func handleCmdNFTMint(ctx *cli.Context) error {
	log.Info("start to mint nft...")

	asset := flag2address(ctx, AssetFlag)
	to := flag2address(ctx, DstAccountFlag)
	tokenID := flag2big(ctx, TokenIdFlag)
	uri := cfg.OSS + tokenID.String()
	tx, err := sdk.MintNFT(adm, asset, to, tokenID, uri)
	if err != nil {
		return err
	}
	log.Info("mint nft %d to %s success! txhash %s", tokenID.Uint64(), to.Hex(), tx.Hex())
	return nil
}

func handleCmdNFTWrapLock(ctx *cli.Context) error {
	log.Info("start to lock nft...")

	from := flag2address(ctx, SrcAccountFlag)
	key, err := wallet.LoadEthAccount(storage, cc.Keystore, from.Hex(), defaultAccPwd)
	if err != nil {
		return err
	}

	asset := flag2address(ctx, AssetFlag)
	to := flag2address(ctx, DstAccountFlag)
	dstChainID := flag2Uint64(ctx, DstChainFlag)
	tokenID := flag2big(ctx, TokenIdFlag)
	feeToken := common.HexToAddress(cc.FeeToken)
	fee := flag2big(ctx, AmountFlag)
	id := new(big.Int).SetUint64(flag2Uint64(ctx, LockIdFlag))
	wrapper := common.HexToAddress(cc.NFTWrap)

	fmt.Printf("-------assset:%s, to:%s, dstChainID:%d, tokenID:%s, feeToken:%s, fee:%s, id: %s",
		asset.Hex(), to.Hex(), dstChainID, tokenID.String(), feeToken.Hex(), fee.String(), id.String())

	tx, err := sdk.WrapLock(key, wrapper, asset, to, dstChainID, tokenID, feeToken, fee, id)
	if err != nil {
		return err
	}
	log.Info("wrap lock nft %d success! txhash %s", tokenID, tx.Hex())
	return nil
}

func handleCmdERC20Mint(ctx *cli.Context) error {
	log.Info("start to mint erc20 token...")

	var asset common.Address
	isFeeToken := ctx.GlobalBool(getFlagName(FeeTokenFlag))
	if isFeeToken {
		asset = common.HexToAddress(cc.FeeToken)
	} else {
		asset = flag2address(ctx, ERC20TokenFlag)
	}

	to := flag2address(ctx, DstAccountFlag)
	amount := flag2big(ctx, AmountFlag)
	log.Debug("mint to %s %s", to.Hex(), amount.String())

	tx, err := sdk.MintERC20Token(adm, asset, to, amount)
	if err != nil {
		return err
	}
	log.Info("mint %s to %s success, hash %s", amount.String(), to.Hex(), tx.Hex())
	return nil
}

func handleCmdERC20Transfer(ctx *cli.Context) error {
	log.Info("start to transfer erc20 token...")

	var asset common.Address
	isFeeToken := ctx.Bool(getFlagName(FeeTokenFlag))
	if isFeeToken {
		asset = common.HexToAddress(cc.FeeToken)
	} else {
		asset = flag2address(ctx, ERC20TokenFlag)
	}

	from := flag2address(ctx, SrcAccountFlag)
	key, err := wallet.LoadEthAccount(storage, cc.Keystore, from.Hex(), defaultAccPwd)
	if err != nil {
		return err
	}

	to := flag2address(ctx, DstAccountFlag)
	amount := flag2big(ctx, AmountFlag)
	tx, err := sdk.TransferERC20Token(key, asset, to, amount)
	if err != nil {
		return err
	}

	log.Info("%s transfer %s to %s success, tx %s", from.Hex(), amount.String(), to.Hex(), tx.Hex())
	return nil
}

func handleGetErc20Balance(ctx *cli.Context) error {
	owner := flag2address(ctx, SrcAccountFlag)

	var asset common.Address
	isFeeToken := ctx.Bool(getFlagName(FeeTokenFlag))
	if isFeeToken {
		asset = common.HexToAddress(cc.FeeToken)
	} else {
		asset = flag2address(ctx, ERC20TokenFlag)
	}

	balance, err := sdk.GetERC20Balance(asset, owner)
	if err != nil {
		return fmt.Errorf("get balance failed, err: %v", err)
	}
	log.Info("%s balance of asset %s is %s", owner.Hex(), asset.Hex(), balance.String())
	return nil
}

func updateConfig() error {
	if err := files.WriteJsonFile(cfgPath, cfg, true); err != nil {
		return err
	}
	log.Info("update config success!", cfgPath)
	return nil
}

func selectChainConfig(chainID uint64) {
	cc = customSelectChainConfig(chainID)
}

func flag2string(ctx *cli.Context, f cli.Flag) string {
	fn := getFlagName(f)
	data := ctx.String(fn)
	return data
}

func flag2address(ctx *cli.Context, f cli.Flag) common.Address {
	data := flag2string(ctx, f)
	return common.HexToAddress(data)
}

func flag2big(ctx *cli.Context, f cli.Flag) *big.Int {
	fn := getFlagName(f)
	data := ctx.String(fn)
	return math.String2BigInt(data)
}

func flag2Uint64(ctx *cli.Context, f cli.Flag) uint64 {
	fn := getFlagName(f)
	data := ctx.Uint64(fn)
	return data
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
