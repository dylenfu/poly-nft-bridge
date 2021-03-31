## Deploy tool


### prepare
给几条链的admin账户转native token，这里以console模式为例:
```bash
eth.accounts

personal.unlockAccount("0xa**b");

eth.getBalance("0xa**b");

eth.sendTransaction({from: "0xc**d4",to: "0x7a**6a", value: "7**0"});
```
这里需要注意，`personal.unlockAccount`需要在节点启动的时候配置`--allow-insecure-unlock`

#### cmds

1.deploy mintable erc20 token:
```shell script
./deploy-tool 
```
