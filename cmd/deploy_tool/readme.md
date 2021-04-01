## Deploy tool


### 准备
给几条链的admin账户转native token，这里以console模式为例:
```bash
eth.accounts

personal.unlockAccount("0xa**b");

eth.getBalance("0xa**b");

eth.sendTransaction({from: "0xc**d4",to: "0x7a**6a", value: "7**0"});
```
这里需要注意，`personal.unlockAccount`需要在节点启动的时候配置`--allow-insecure-unlock`

#### 部署erc20token

1.部署erc20合约:
```shell script
./deploy_tool --chain=2 deployERC20
./deploy_tool --chain=6 deployERC20
./deploy_tool --chain=7 deployERC20
```
这本合约，可以作为wrapper的feeToken(参数`--feeToken=true`, 因为默认为true，所以忽略), 同时继承了mintable属性，我们可以在dev环境下准备相应的token做测试.

2.mint erc20 token:
```shell script
./deploy_tool --chain=2 mintERC20 --to=0x31c0dd87B33Dcd66f9a255Cf4CF39287F8AE593C --amount=100000000000000000000000000
./deploy_tool --chain=6 mintERC20 --to=0x896fB9Dd4Bddd1C4ea2cab3df66C632AD736a9D1 --amount=100000000000000000000000000
```
这个fee后面可以用来支付wrapper手续费

3.transfer erc20 token:

本地测试数据
```dtd
eth user1 0x5Fb03EB21303D39967a1a119B32DD744a0fA8986 user2 0xB9933ff0CB5C5B42b12972C9826703E10BFDd863
bsc user1 0x8cbE1493A2894e32985E45e7e3394f3FEA15Afb2 user2 0xa252dCBF98D02218b4E5B7B00d8FE7646592394E
heco user1 0x95598C69B02925De711D4015F85b49527381aF6d user2 0xE4Ecc16675d1e0A587f1435003786afE23B71733
```

该部分测试必须在部署完feeToken, 并mint一部分给管理员后操作.每个用户10000个feeToken

```shell script
./deploy_tool --chain=2 transferERC20 --from=0x31c0dd87B33Dcd66f9a255Cf4CF39287F8AE593C --to=0x5Fb03EB21303D39967a1a119B32DD744a0fA8986 --amount=10000000000000000000000
./deploy_tool --chain=2 transferERC20 --from=0x31c0dd87B33Dcd66f9a255Cf4CF39287F8AE593C --to=0xB9933ff0CB5C5B42b12972C9826703E10BFDd863 --amount=10000000000000000000000

./deploy_tool --chain=6 transferERC20 --from=0x896fB9Dd4Bddd1C4ea2cab3df66C632AD736a9D1 --to=0x8cbE1493A2894e32985E45e7e3394f3FEA15Afb2 --amount=10000000000000000000000
./deploy_tool --chain=6 transferERC20 --from=0x896fB9Dd4Bddd1C4ea2cab3df66C632AD736a9D1 --to=0xa252dCBF98D02218b4E5B7B00d8FE7646592394E --amount=10000000000000000000000

./deploy_tool --chain=2 erc20Balance --from=0x5Fb03EB21303D39967a1a119B32DD744a0fA8986
./deploy_tool --chain=2 erc20Balance --from=0xB9933ff0CB5C5B42b12972C9826703E10BFDd863

./deploy_tool --chain=6 erc20Balance --from=0x8cbE1493A2894e32985E45e7e3394f3FEA15Afb2
./deploy_tool --chain=6 erc20Balance --from=0xa252dCBF98D02218b4E5B7B00d8FE7646592394E

./deploy_tool --chain=2 approveERC20 --from=0x5Fb03EB21303D39967a1a119B32DD744a0fA8986 --to=0xa64697B45eD4566Af42d8811B0320a1636c13BC2 --amount=1000000000000000000000000
./deploy_tool --chain=2 erc20Allowance --from=0x5Fb03EB21303D39967a1a119B32DD744a0fA8986 --to=0xa64697B45eD4566Af42d8811B0320a1636c13BC2

```
或者使用native token作为feeToken
```shell script
./deploy_tool --chain=2 transferNative --from=0x31c0dd87B33Dcd66f9a255Cf4CF39287F8AE593C --to=0x5Fb03EB21303D39967a1a119B32DD744a0fA8986 --amount=100000000000000000000
./deploy_tool --chain=2 transferNative --from=0x31c0dd87B33Dcd66f9a255Cf4CF39287F8AE593C --to=0xB9933ff0CB5C5B42b12972C9826703E10BFDd863 --amount=100000000000000000000

./deploy_tool --chain=6 transferNative --from=0x896fB9Dd4Bddd1C4ea2cab3df66C632AD736a9D1 --to=0x8cbE1493A2894e32985E45e7e3394f3FEA15Afb2 --amount=100000000000000000000
./deploy_tool --chain=6 transferNative --from=0x896fB9Dd4Bddd1C4ea2cab3df66C632AD736a9D1 --to=0xa252dCBF98D02218b4E5B7B00d8FE7646592394E --amount=100000000000000000000
```
#### 链基础合约

1.以太
```shell script
./deploy_tool --chain=2 deployECCD
./deploy_tool --chain=2 deployECCM
./deploy_tool --chain=2 deployCCMP

./deploy_tool --chain=2 transferECCDOwnership
./deploy_tool --chain=2 transferECCMOwnership

./deploy_tool --chain=2 registerSideChain
./deploy_tool --chain=2 approveSideChain

./deploy_tool --chain=2 syncSideGenesis
./deploy_tool --chain=2 syncPolyGenesis
```

2.bsc
```shell script
./deploy_tool --chain=6 deployECCD
./deploy_tool --chain=6 deployECCM
./deploy_tool --chain=6 deployCCMP

./deploy_tool --chain=6 transferECCDOwnership
./deploy_tool --chain=6 transferECCMOwnership

./deploy_tool --chain=6 registerSideChain
./deploy_tool --chain=6 approveSideChain

./deploy_tool --chain=6 syncSideGenesis
./deploy_tool --chain=6 syncPolyGenesis
```

3.heco
```shell script
./deploy_tool --chain=7 deployECCD
./deploy_tool --chain=7 deployECCM
./deploy_tool --chain=7 deployCCMP

./deploy_tool --chain=7 transferECCDOwnership
./deploy_tool --chain=7 transferECCMOwnership

./deploy_tool --chain=7 registerSideChain
./deploy_tool --chain=7 approveSideChain

./deploy_tool --chain=7 syncSideGenesis
./deploy_tool --chain=7 syncPolyGenesis
```

#### NFTLockProxy合约

```shell script
./deploy_tool --chain=2 deployNFTLockProxy
./deploy_tool --chain=2 proxySetCCMP

./deploy_tool --chain=6 deployNFTLockProxy
./deploy_tool --chain=6 proxySetCCMP

./deploy_tool --chain=7 deployNFTLockProxy
./deploy_tool --chain=7 proxySetCCMP

./deploy_tool --chain=2 bindProxy --dstChain=6
./deploy_tool --chain=6 bindProxy --dstChain=2

./deploy_tool --chain=2 bindProxy --dstChain=7
./deploy_tool --chain=7 bindProxy --dstChain=2

./deploy_tool --chain=6 bindProxy --dstChain=7
./deploy_tool --chain=7 bindProxy --dstChain=6

```

#### NFTWrap合约
```shell script
./deploy_tool --chain=2 deployNFTWrapper
./deploy_tool --chain=2 setWrapLockProxy
./deploy_tool --chain=2 setFeeCollector

./deploy_tool --chain=6 deployNFTWrapper
./deploy_tool --chain=6 setWrapLockProxy
./deploy_tool --chain=6 setFeeCollector

./deploy_tool --chain=7 deployNFTWrapper
./deploy_tool --chain=7 setWrapLockProxy
./deploy_tool --chain=7 setFeeCollector
```

#### NFT资产合约及绑定

这里以eth和bsc的跨链为例
```shell script

#todo: erc721合约_safeMint已修改.在mint的时候不要进入到onReceive方法，因为现在的lock proxy的onReceive方法中只接收来自proxy的行为 

./deploy_tool --chain=2 deployNFT --name=digtalCat1 --symbol=cat1
./deploy_tool --chain=6 deployNFT --name=digtalCat1 --symbol=cat1

./deploy_tool --chain=2 bindNFT --asset=0x35EFCE8D79D6Cae30B38F6dAC3fc55C62c146b4c --dstChain=6 --dstAsset=0x63F8eaCfbF43F027cca37aB90c0ce08E76D93679
./deploy_tool --chain=6 bindNFT --asset=0x63F8eaCfbF43F027cca37aB90c0ce08E76D93679 --dstChain=2 --dstAsset=0x35EFCE8D79D6Cae30B38F6dAC3fc55C62c146b4c
```

#### NFT跨链

```shell script

# 为relayer发送账户准备native token
./deploy_tool --chain=2 transferNative --from=0x31c0dd87B33Dcd66f9a255Cf4CF39287F8AE593C --to=0x69611B922c985fD793AFA56CE8Cfe7d8aFffeFDd --amount=100000000000000000000
./deploy_tool --chain=6 transferNative --from=0x896fB9Dd4Bddd1C4ea2cab3df66C632AD736a9D1 --to=0xc4A519453135aE569c582154A5E632668E6DADc4 --amount=100000000000000000000

# mint nft token
./deploy_tool --chain=2 mintNFT --asset=0x4A17a58141E9D0b85B0F9186c9dfCfc0DCD4425f --to=0x5Fb03EB21303D39967a1a119B32DD744a0fA8986 --tokenId=5

# wrapper lock nft
./deploy_tool --chain=2 lockNFT --dstChain=6 \
--asset=0x4A17a58141E9D0b85B0F9186c9dfCfc0DCD4425f \
--from=0x5Fb03EB21303D39967a1a119B32DD744a0fA8986 \
--to=0x8cbE1493A2894e32985E45e7e3394f3FEA15Afb2 \
--amount=10000000000000000 \
--nativeToken=true \
--tokenId=5 --lockId=5
```