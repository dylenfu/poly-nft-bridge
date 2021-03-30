package flags

import (
	"fmt"
	"strings"

	basedef "github.com/polynetwork/poly-nft-bridge/const"
	"github.com/urfave/cli"
)

const (
	CmdDeployECCDContract        = "deploy-eccd"
	CmdDeployECCMContract        = "deploy-eccm"
	CmdDeployCCMPContract        = "deploy-ccmp"
	CmdDeployNFTContract         = "deploy-nft"
	CmdDeployERC20Contract       = "deploy-erc20"
	CmdDeployLockProxyContract   = "deploy-proxy"
	CmdDeployNFTWrapContract     = "deploy-nft-wrap"
	CmdLockProxySetCCMP          = "proxy-set-ccmp"
	CmdBindNFTAsset              = "bind-nft"
	CmdBindLockProxy             = "bind-proxy"
	CmdGetBoundLockProxy         = "get-bound-proxy"
	CmdTransferECCDOwnership     = "transfer-eccd-ownership"
	CmdTransferECCMOwnership     = "transfer-eccm-ownership"
	CmdTransferCCMPOwnership     = "transfer-ccmp-ownership"
	CmdTransferProxyOwnership    = "transfer-proxy-ownership"
	CmdSyncSideChainGenesis2Poly = "sync-side-genesis"
	CmdSyncPolyGenesis2SideChain = "sync-poly-genesis"
	CmdNFTWrapSetFeeCollector    = "set-fee-collector"
	CmdNFTWrapSetLockProxy       = "set-wrap-lock-proxy"
	CmdNFTMint                   = "mint-nft" // todo
	CmdNFTWrapLock               = "lock-nft" // todo
	CmdERC20Mint                 = "mint-erc20"
	CmdERC20Transfer             = "transfer-erc20"
)

var (
	cmdCases = cli.StringSlice{
		CmdDeployECCDContract,
		CmdDeployECCMContract,
		CmdDeployCCMPContract,
		CmdDeployNFTContract,
		CmdDeployLockProxyContract,
		CmdLockProxySetCCMP,
		CmdBindLockProxy,
		CmdGetBoundLockProxy,
		CmdBindNFTAsset,
		CmdTransferECCDOwnership,
		CmdTransferECCMOwnership,
		//CmdTransferCCMPOwnership,
		//CmdTransferProxyOwnership,
		CmdSyncSideChainGenesis2Poly,
		CmdSyncPolyGenesis2SideChain,
		CmdNFTWrapSetFeeCollector,
		CmdNFTWrapSetLockProxy,
		CmdNFTMint,
		CmdNFTWrapLock,
		CmdERC20Mint,
		CmdERC20Transfer,
	}
)

func dumpCmdCases() string {
	s := fmt.Sprintf("%s: %s\r\n", CmdDeployECCDContract, "deploy ethereum cross chain data contract.")
	s += fmt.Sprintf("%s: %s\r\n", CmdDeployECCMContract, "deploy ethereum cross chain manage contract.")
	s += fmt.Sprintf("%s: %s\r\n", CmdDeployCCMPContract, "deploy ethereum cross chain manager proxy contract.")
	s += fmt.Sprintf("%s: %s\r\n", CmdDeployNFTContract, "deploy new nft asset with mapping contract.")
	s += fmt.Sprintf("%s: %s\r\n", CmdDeployLockProxyContract, "deploy nft lock proxy contract.")
	s += fmt.Sprintf("%s: %s\r\n", CmdLockProxySetCCMP, "set cross chain manager proxy address for lock proxy contract.")
	s += fmt.Sprintf("%s: %s\r\n", CmdBindLockProxy, "bind lock proxy contract with another side chain's lock proxy contract.")
	s += fmt.Sprintf("%s: %s\r\n", CmdGetBoundLockProxy, "get bound lock proxy contract.")
	s += fmt.Sprintf("%s: %s\r\n", CmdBindNFTAsset, "bind nft asset to side chain.")
	s += fmt.Sprintf("%s: %s\r\n", CmdTransferECCDOwnership, "transfer ethereum cross chain data contract ownership to other account.")
	s += fmt.Sprintf("%s: %s\r\n", CmdTransferECCMOwnership, "transfer ethereum cross chain manager contract ownership to other account.")
	//s += fmt.Sprintf("%s: %s\r\n", CmdTransferCCMPOwnership, "transfer ethereum cross chain manager proxy contract ownership to other account.")
	//s += fmt.Sprintf("%s: %s\r\n", CmdTransferProxyOwnership, "transfer lock proxy contract ownership to other account.")
	s += fmt.Sprintf("%s: %s\r\n", CmdSyncSideChainGenesis2Poly, "sync side chain genesis header to poly chain.")
	s += fmt.Sprintf("%s: %s\r\n", CmdSyncPolyGenesis2SideChain, "sync poly genesis header to side chain.")
	s += fmt.Sprintf("%s: %s\r\n", CmdNFTWrapSetLockProxy, "set nft lock proxy for wrap contract.")
	s += fmt.Sprintf("%s: %s\r\n", CmdNFTWrapSetFeeCollector, "set nft fee collecotr for wrap contract")
	s += fmt.Sprintf("%s: %s\r\n", CmdNFTMint, "mint nft token.")
	s += fmt.Sprintf("%s: %s\r\n", CmdNFTWrapLock, "transfer nft token.")
	s += fmt.Sprintf("%s: %s\r\n", CmdERC20Mint, "mint erc20 token.")
	s += fmt.Sprintf("%s: %s\r\n", CmdERC20Transfer, "transfer ERC20 token.")
	return s
}

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
		Value: "./bridge_tools/conf/config_transactions.json",
	}

	ChainIDFlag = cli.Uint64Flag{
		Name:  "chain",
		Usage: "select chainID",
		Value: basedef.ETHEREUM_CROSSCHAIN_ID,
	}

	//OwnerFlag = cli.StringFlag{
	//	Name:  "owner",
	//	Usage: "set owner for deploy nft contract, etc.",
	//	Value: "",
	//}

	NFTNameFlag = cli.StringFlag{
		Name:  "nft-name",
		Usage: "set nft name for deploy nft contract, etc.",
		Value: "",
	}

	NFTSymbolFlag = cli.StringFlag{
		Name:  "nft-symbol",
		Usage: "set nft symbol for deploy nft contract, etc.",
		Value: "",
	}

	DstChainFlag = cli.Uint64Flag{
		Name:  "dst-chain",
		Usage: "set dest chain for cross chain",
		Value: 0,
	}

	AssetFlag = cli.Uint64Flag{
		Name:  "asset",
		Usage: "set asset for cross chain or mint nft",
		Value: 0,
	}

	DstAssetFlag = cli.Uint64Flag{
		Name:  "dst-asset",
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
		Name:  "fee-token",
		Usage: "choose erc20 token to be fee token",
	}

	ERC20TokenFlag = cli.BoolFlag{
		Name:  "erc20-token",
		Usage: "choose erc20 token to be fee token",
	}

	AmountFlag = cli.StringFlag{
		Name:  "amount",
		Usage: "transfer amount or fee amount",
		Value: "",
	}

	TokenIdFlag = cli.Uint64Flag{
		Name:  "token-id",
		Usage: "set token id while mint nft",
	}

	LockIdFlag = cli.Uint64Flag{
		Name:  "lock-id",
		Usage: "wrap lock nft item id",
	}

	CmdFlag = cli.StringSliceFlag{
		Name:  "cmd",
		Usage: "cross chain case:\r\n" + dumpCmdCases(),
		Value: &cmdCases,
	}
)

//GetFlagName deal with short flag, and return the flag name whether flag name have short name
func GetFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}
